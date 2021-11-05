package dybmapi

import (
	"bufio"
	"compress/gzip"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"cloud.google.com/go/storage"
)

var ctx context.Context
var bucket *storage.BucketHandle

func processText(text string) string {
	// not necessary, it seems:
	// result := replaceGoogleNgramKeywords(text)
	result := text

	if !isAllLetters(result) {
		result = ""
	}
	return result
}

// prepareNgram takes the data from one line, processes it and prepares an Ngram that it returns
// it also returns true if the ngram is the same as the previous one so that the calling function can handle it
// if there was no previous ngram, it returns false
func prepareNgram(text, yearText, matchesText, volumesText string, ngram *Ngram, previousNgram *Ngram) bool {
	ngram.OriginalText = text
	isTheSame := false

	if previousNgram.OriginalText == ngram.OriginalText {
		isTheSame = true

		// yay, we don't have to process the text
		ngram.Text = previousNgram.Text
		ngram.Frequency = previousNgram.Frequency
	} else {
		ngram.Text = processText(text)
		ngram.Frequency = 0
	}

	// if the text is empty, it means it doesn't make sense and we don't want it; for example "B.B. --_."
	if ngram.Text != "" {
		year, _ := strconv.Atoi(yearText)
		matches, _ := strconv.Atoi(matchesText)
		volumes, _ := strconv.Atoi(volumesText)

		yearBonus := 1
		if year > 1980 {
			yearBonus += (year - 1980) / 10
		}

		volumesBonus := 1 + volumes/10
		if volumesBonus > 4 {
			volumesBonus = 4
		}

		frequency := matches * yearBonus * volumesBonus
		ngram.Frequency += frequency
	}

	return isTheSame
}

func processUrl(url string) {
	log.Println("Processing url", url)
	writtenNgrams := 0

	targetFileName := getNgramFilenameFromUrl(url)
	gcTestReader, err := bucket.Object(targetFileName).NewReader(ctx)
	if err == nil {
		gcTestReader.Close()
		log.Println("No need to process url", url, ", the file", targetFileName, "already exists")
		return
	}

	gcWriter := bucket.Object(targetFileName).NewWriter(ctx)
	gcWriter.ContentType = "text/csv"

	response, err := http.Get(url)
	if err != nil {
		log.Fatal("Unable to Get URL: ", err)
	}
	defer response.Body.Close()

	reader, err := gzip.NewReader(response.Body)
	if err != nil {
		log.Fatal("Unable to open gzipped file for reading: ", err)
	}
	defer reader.Close()

	var isTheSame bool

	var bufferBuilder strings.Builder
	bufferLength := 0
	bufferCapacity := 100000

	var ngram Ngram
	var previousNgram Ngram

	scanner := bufio.NewScanner(reader)
	log.Println("Starting to process ngrams for", targetFileName)
	for scanner.Scan() {
		line := scanner.Text()

		// format: ngram TAB year TAB match_count TAB volume_count NEWLINE
		fragments := strings.Split(line, "\t")
		if len(fragments) >= 4 {
			isTheSame = prepareNgram(fragments[0], fragments[1], fragments[2], fragments[3], &ngram, &previousNgram)

			if !isTheSame && previousNgram.Text != "" {
				bufferBuilder.WriteString(previousNgram.Text)
				bufferBuilder.WriteString(";")
				bufferBuilder.WriteString(strconv.Itoa(previousNgram.Frequency))
				bufferBuilder.WriteString("\n")

				bufferLength++
				if bufferLength > bufferCapacity {
					log.Println("Buffer full, flushing", bufferLength, "ngrams to", targetFileName)
					gcWriter.Write([]byte(bufferBuilder.String()))
					bufferBuilder.Reset()
					bufferLength = 0
				}

				writtenNgrams++
			}

			previousNgram.OriginalText = ngram.OriginalText
			previousNgram.Text = ngram.Text
			previousNgram.Frequency = ngram.Frequency
		} else {
			log.Println("WARNING: Found a line with not enough TABs:", line)
		}
	}

	if bufferLength > 0 {
		log.Println("Flushing the rest of the buffer -", bufferLength, "ngrams")
		gcWriter.Write([]byte(bufferBuilder.String()))
		bufferBuilder.Reset()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if err := gcWriter.Close(); err != nil {
		log.Fatalf("Unable to close Cloud Storage file %q: %v", targetFileName, err)
	}

	log.Println("Finished processing url", url, "with", writtenNgrams, "ngrams were written into", targetFileName)
}

func getAndProcessFiles(urls []string, n, letter string, maxUrls int) ImportData {
	log.Println("Processing urls with", n, "grams starting with", letter)

	var importData ImportData

	targetFilename := getNgramTargetFilename(n, letter)
	targetObject := bucket.Object(targetFilename)

	// let's check if the final file exists
	attrs, err := targetObject.Attrs(ctx)
	if err == nil {
		log.Println("No need to do anything, the final file", targetFilename, "already exists and its size is", attrs.Size)
	} else {
		var urlsToProcess []string
		urlsProcessed := 0
		for _, url := range urls {
			if isUrlForNgramAndLetter(url, n, letter) {
				urlsToProcess = append(urlsToProcess, url)

				urlsProcessed++
				if urlsProcessed == maxUrls {
					break
				}
			}
		}

		var wg sync.WaitGroup
		wg.Add(len(urlsToProcess))

		for _, url := range urlsToProcess {
			importData.ProcessedUrls = append(importData.ProcessedUrls, url)
			go func(url string) {
				processUrl(url)
				defer wg.Done()
			}(url)

		}

		wg.Wait()

		log.Println("Finished processing urls. Processed", urlsProcessed, "urls with", n, "-grams starting with", letter)
	}

	return importData
}

func prepareContext() []string {
	ctx = context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("Failed to create Cloud Storage client: ", err)
	}

	bucket = client.Bucket(bucketName)
	urlFileReader, err := bucket.Object(urlsFilename).NewReader(ctx)
	if err != nil {
		log.Fatal("Unable to open file: ", err)
	}
	defer urlFileReader.Close()

	urlsText, err := ioutil.ReadAll(urlFileReader)
	if err != nil {
		log.Fatal("Unable to read from file: ", err)
	}

	return strings.Split(string(urlsText), "\r\n")
}

// HandleImport handles the /import URL
func HandleImport(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /import")

	w.Header().Add("Content-type", "application/json")
	urls := prepareContext()

	letter := r.URL.Query().Get("letter")
	n := r.URL.Query().Get("n")
	max, err := strconv.Atoi(r.URL.Query().Get("max"))
	if err != nil {
		log.Fatal("Unable to get the max number of files: ", err)
	}

	importData := getAndProcessFiles(urls, n, letter, max)

	importData.UrlsFilename = urlsFilename
	resultJson, err := json.Marshal(importData)

	if err != nil {
		log.Fatal("Error when JSONing the result: ", err)
	}

	w.Write(resultJson)

	log.Println("Handling /import finished")
}

// HandleImport handles the /import URL
func HandleCombineImport(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /combine-import")

	w.Header().Add("Content-type", "application/json")

	urls := prepareContext()

	var combineImportData CombineImportData

	//var cLetters = [...]string{"c"}
	//for n := 2; n <= 2; n++ {
	//for _, letter := range cLetters {
	for n := 2; n <= 5; n++ {
		for _, letter := range validLetters {
			stringN := strconv.Itoa(n)
			targetFilename := getNgramTargetFilename(stringN, letter)
			targetObject := bucket.Object(targetFilename)

			// let's check if the file exists
			attrs, err := targetObject.Attrs(ctx)
			if err == nil {
				log.Println("No need to process target file", targetFilename, "as it already exists and its size is", attrs.Size)
				continue
			}

			var sourceFiles []string
			var sourceObjects []*storage.ObjectHandle

			var totalSize int64 = 0
			allFilesExist := true
			log.Println("Checking if all source files exist for", targetFilename)

			for _, url := range urls {
				if isUrlForNgramAndLetter(url, stringN, letter) {
					sourceFilename := getNgramFilenameFromUrl(url)
					sourceObject := bucket.Object(sourceFilename)

					// let's check if the file exists
					attrs, err := sourceObject.Attrs(ctx)
					if err != nil {
						log.Println("The source file", sourceFilename, "does not exist")
						allFilesExist = false
						break
					} else {
						log.Println("The source file", sourceFilename, "exists and its size is", attrs.Size)
						totalSize += attrs.Size
					}

					sourceFiles = append(sourceFiles, sourceFilename)
					sourceObjects = append(sourceObjects, sourceObject)
				}
			}

			if allFilesExist {
				log.Println("All the source files exist and their total size is", totalSize, "- let's try composing them together")
				_, err = targetObject.ComposerFrom(sourceObjects...).Run(ctx)
				if err != nil {
					log.Println("Error when composing:", err)
				} else {
					var combinedFile CombinedFile
					combinedFile.N = n
					combinedFile.Letter = letter
					combinedFile.SourceFiles = sourceFiles
					combinedFile.TargetFile = targetFilename
					combineImportData.CombinedFiles = append(combineImportData.CombinedFiles, combinedFile)

					log.Println("Composing successful for", targetFilename, " - deleting source files")
					for _, objectToDelete := range sourceObjects {
						err := objectToDelete.Delete(ctx)
						if err != nil {
							log.Println("Coudln't delete file", objectToDelete.ObjectName())
						}
					}
				}
			} else {
				log.Println("Not all source files exist for", targetFilename)
			}
		}
	}

	resultJson, err := json.Marshal(combineImportData)

	if err != nil {
		log.Fatal("Error when JSONing the result: ", err)
	}

	w.Write(resultJson)

	log.Println("Handling /combine-import finished")
}

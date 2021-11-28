package dybmimport

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/storage"
)

var DataFolder string

func Setup(dataFolder string) {
	DataFolder = dataFolder
}

// processNgram processes one ngram and either adjusts it for the final file
// or resets it to an empty string if it's not suitable for the final file
func processNgram(ngram string) string {
	// not necessary, it seems:
	// result := replaceGoogleNgramKeywords(text)
	result := ngram

	if !isAllLetters(result) {
		result = ""
	}
	return result
}

// prepareNgram takes the data from one line, processes it and prepares an Ngram that it returns
// it also returns true if the ngram is the same as the previous one so that the calling function can handle it
// if there was no previous ngram, it returns false
func prepareNgram(textBytes, yearBytes, matchesBytes, volumesBytes []byte, ngram *Ngram, previousNgram *Ngram) bool {
	text := string(textBytes)
	ngram.OriginalText = text
	isTheSame := false

	if previousNgram.OriginalText == ngram.OriginalText {
		isTheSame = true

		// yay, we don't have to process the text
		ngram.Text = previousNgram.Text
		ngram.Frequency = previousNgram.Frequency
	} else {
		ngram.Text = processNgram(text)
		ngram.Frequency = 0
	}

	// if the text is empty, it means it doesn't make sense and we don't want it; for example "B.B. --_."
	if ngram.Text != "" {
		year := int(yearBytes[0]-48)*1000 + int(yearBytes[1]-48)*100 + int(yearBytes[2]-48)*10 + int(yearBytes[3]-48)
		matches := convertAsciiNumberToInt(matchesBytes)
		volumes := convertAsciiNumberToInt(volumesBytes)

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

// processUrl processes one URL and it's meant to run concurrently.
// it can handle files the size of tens of gigabytes within one hour timeout (tested)
func processUrl(url string, id string) {
	log.Println("Url", url, id)
	writtenNgrams := 0

	ctx, bucket, _ := prepareContext()

	targetRawFileName := getNgramFilenameFromUrl(url, true)
	targetFileName := getNgramFilenameFromUrl(url, false)
	targetObject := bucket.Object(targetFileName)
	_, err := targetObject.Attrs(ctx)
	if err == nil {
		// log.Println("Url", url, "already here:", targetFileName, "size", attrs.Size, id)
		return
	}

	targetRawObject := bucket.Object(targetRawFileName)
	_, err = targetRawObject.Attrs(ctx)
	if err == nil {
		// log.Println("Gzip already unzipped:", targetRawFileName, "size", attrs.Size, id)
	} else {
		response, err := http.Get(url)
		if err != nil {
			log.Fatal("Unable to Get URL", url, id, err)
		}
		defer response.Body.Close()

		gzReader, err := gzip.NewReader(response.Body)
		if err != nil {
			log.Fatal("Unable to open gzip", url, id, err)
		}
		defer gzReader.Close()

		gcRawWriter := targetRawObject.NewWriter(ctx)
		gcRawWriter.ContentType = "text/csv"

		log.Println("Extract gzip", url, "to", targetRawFileName, id)
		if _, err = io.Copy(gcRawWriter, gzReader); err != nil {
			log.Println("Unable to extract gzip", url, "to", targetRawFileName, id, err)
			return
		}
		if err := gcRawWriter.Close(); err != nil {
			log.Println("Unable to close raw", targetRawFileName, id, err)
			return
		}
		log.Println("Extract OK", targetRawFileName, id)
	}

	gcReader, err := targetRawObject.NewReader(ctx)
	if err != nil {
		log.Println("Unable to open raw", targetRawFileName, id, err)
		return
	}
	defer gcReader.Close()

	gcWriter := targetObject.NewWriter(ctx)
	gcWriter.ContentType = "text/csv"

	var isTheSame bool

	var bufferBuilder strings.Builder
	bufferLength := 0
	bufferCapacity := 100000
	var linesRead int64 = 0
	var linesWritten int = 0

	var ngram Ngram
	var previousNgram Ngram
	var tabSeparator = [...]byte{9}

	scanner := bufio.NewScanner(gcReader)
	// log.Println("Scanning", targetFileName, id)
	for scanner.Scan() {
		line := scanner.Bytes()
		linesRead++
		if linesRead%100000000 == 0 {
			log.Println("+100M lines", targetRawFileName, "total", linesRead, id)
		}

		// format: ngram TAB year TAB match_count TAB volume_count NEWLINE
		fragments := bytes.Split(line, tabSeparator[:])
		if len(fragments) >= 4 {
			isTheSame = prepareNgram(fragments[0], fragments[1], fragments[2], fragments[3], &ngram, &previousNgram)

			if !isTheSame && previousNgram.Text != "" {
				bufferBuilder.WriteString(previousNgram.Text)
				bufferBuilder.WriteString(";")
				bufferBuilder.WriteString(strconv.Itoa(previousNgram.Frequency))
				bufferBuilder.WriteString("\n")

				bufferLength++
				if bufferLength > bufferCapacity {
					linesWritten += bufferLength

					log.Println("Buffer flush", bufferLength, "total", linesWritten, "to", targetFileName, id)
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
			log.Println("WARNING: Found a line with not enough TABs:", string(line), id)
		}
	}

	if bufferLength > 0 {
		log.Println("Final buffer flush", bufferLength, "to", targetFileName, id)
		gcWriter.Write([]byte(bufferBuilder.String()))
		bufferBuilder.Reset()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if err := gcWriter.Close(); err != nil {
		log.Fatal("Unable to close", targetFileName, id, err)
	} else {
		// successfully writen => we can remove the raw file
		targetRawObject.Delete(ctx)
		log.Println("Deleted raw", targetRawFileName, id)
	}

	log.Println("======= Finished", url, "with", writtenNgrams, "ngrams written into", targetFileName, id)
}

// getAndProcessFiles gets the source Google Books ngram files and processes them
// the resulting files are cleansed - meaning each ngram only appears once and all the ngrams featuring non-letters are omitted
func getAndProcessFiles(ctx context.Context, bucket *storage.BucketHandle, urls []string, n, letter string, maxUrls int, requestTime string) ImportData {
	log.Println("Processing urls", n, letter, "max", maxUrls, requestTime)

	var importData ImportData

	targetFilename := getNgramTargetFilename(n, letter)
	targetObject := bucket.Object(targetFilename)

	// let's check if the final file exists
	_, err := targetObject.Attrs(ctx)
	if err == nil {
		// log.Println("Final file", targetFilename, "exists, size", attrs.Size)
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

		for id, url := range urlsToProcess {
			importData.ProcessedUrls = append(importData.ProcessedUrls, url)
			go func(url string, id string) {
				processUrl(url, id)
				defer wg.Done()
			}(url, fmt.Sprintf("id%02d t%s", id+1, requestTime))

		}

		wg.Wait()

		log.Println("______________DONE_____________", urlsProcessed, "urls with", n, "-grams starting with", letter)
	}

	return importData
}

// prepareContext prepares the basic context to work with Cloud Storage
func prepareContext() (context.Context, *storage.BucketHandle, *storage.BucketHandle) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("Failed to create Cloud Storage client: ", err)
	}
	corpusBucket := client.Bucket(corpusBucketName)
	relatedWordsBucket := client.Bucket(relatedWordsBucketName)
	return ctx, corpusBucket, relatedWordsBucket
}

func readUrlFile(ctx context.Context, bucket *storage.BucketHandle) []string {
	urlFileReader, err := bucket.Object(urlsFictionFilename).NewReader(ctx)
	if err != nil {
		log.Fatal("Unable to open urls", urlsFictionFilename, err)
	}
	defer urlFileReader.Close()

	urlsText, err := ioutil.ReadAll(urlFileReader)
	if err != nil {
		log.Fatal("Unable to read urls", urlsFictionFilename, err)
	}

	return strings.Split(string(urlsText), "\r\n")
}

// HandleImport handles the /import URL
func HandleImport(w http.ResponseWriter, r *http.Request) {
	log.Println("START handling", r.URL)
	log.Println("Proto", r.Proto, "TLS", r.TLS, "Host", r.Host)

	/*if r.TLS == nil && r.Host != "localhost:8080" {
		log.Println("ERROR: Not handling non-https outside localhost", r.URL)
		return
	}*/

	currentTime := time.Now()
	requestTime := fmt.Sprintf("%02d:%02d:%02d", currentTime.Hour(), currentTime.Minute(), currentTime.Second())
	log.Println("Time", requestTime)

	w.Header().Add("Content-type", "application/json")
	ctx, bucket, _ := prepareContext()
	urls := readUrlFile(ctx, bucket)

	letter := r.URL.Query().Get("letter")
	n := r.URL.Query().Get("n")
	max, err := strconv.Atoi(r.URL.Query().Get("max"))
	if err != nil {
		log.Fatal("Unable to get the max number of files from query: ", err, requestTime)
	}

	importData := getAndProcessFiles(ctx, bucket, urls, n, letter, max, requestTime)

	importData.UrlsFilename = urlsFictionFilename
	resultJson, err := json.Marshal(importData)

	if err != nil {
		log.Fatal("Error when JSONing the result: ", err, requestTime)
	}

	w.Write(resultJson)

	log.Println("DONE handling", r.URL, requestTime)
}

// HandleCombineImport handles the /combine-import URL and does all the work
// It combines individual ngram files into one ngram file per letter using the Composer
func HandleCombineImport(w http.ResponseWriter, r *http.Request) {
	log.Println("START handling", r.URL)
	log.Println("Proto", r.Proto, "TLS", r.TLS, "Host", r.Host)

	/*	if r.TLS == nil && r.Host != "localhost:8080" {
		log.Println("ERROR: Not handling non-https outside localhost", r.URL)
		return
	}*/

	w.Header().Add("Content-type", "application/json")

	ctx, bucket, _ := prepareContext()
	urls := readUrlFile(ctx, bucket)
	var combineImportData CombineImportData

	allFinalFilesExist := true
	finalFiles := make(map[int][]*storage.ObjectHandle)

	//var cLetters = [...]string{"c"}
	//for n := 2; n <= 2; n++ {
	//for _, letter := range cLetters {
	for n := 2; n <= 5; n++ {
		var newFinalFiles []*storage.ObjectHandle
		finalFiles[n] = newFinalFiles

		for _, letter := range validLetters {
			stringN := strconv.Itoa(n)
			finalFilename := getNgramTargetFilename(stringN, letter)
			finalObject := bucket.Object(finalFilename)

			// let's check if the file exists
			_, err := finalObject.Attrs(ctx)
			if err == nil {
				// log.Println("No need to process target file", targetFilename, "as it already exists and its size is", attrs.Size)
				finalFiles[n] = append(finalFiles[n], finalObject)
				continue
			} else {
				allFinalFilesExist = false
			}

			var sourceFiles []string
			var sourceObjects []*storage.ObjectHandle

			var totalSize int64 = 0
			allFilesExist := true
			log.Println("Checking if all source files exist for", finalFilename)

			for _, url := range urls {
				if isUrlForNgramAndLetter(url, stringN, letter) {
					sourceFilename := getNgramFilenameFromUrl(url, false)
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
				composer := finalObject.ComposerFrom(sourceObjects...)
				composer.ObjectAttrs.ContentType = "text/csv"
				_, err = composer.Run(ctx)
				if err != nil {
					log.Println("Error when composing:", err)
				} else {
					var combinedFile CombinedFile
					combinedFile.N = n
					combinedFile.Letter = letter
					combinedFile.SourceFiles = sourceFiles
					combinedFile.TargetFile = finalFilename
					combineImportData.CombinedFiles = append(combineImportData.CombinedFiles, combinedFile)

					log.Println("Composing successful for", finalFilename, " - deleting source files")
					for _, objectToDelete := range sourceObjects {
						err := objectToDelete.Delete(ctx)
						if err != nil {
							log.Println("Coudln't delete file", objectToDelete.ObjectName())
						}
					}
				}
			} else {
				log.Println("Not all source files exist for", finalFilename)
			}
		}
	}

	if allFinalFilesExist {
		log.Println("All target files exist, let's make a FINAL file")

		var allFinalNObjects []*storage.ObjectHandle
		allFinalNFilesOK := true
		for n := 2; n <= 5; n++ {
			finalNFilename := fmt.Sprintf("Final_%dgram.csv", n)
			FinalNObject := bucket.Object(finalNFilename)

			log.Printf("Composing %dgram: %s", n, finalNFilename)
			finalNComposer := FinalNObject.ComposerFrom(finalFiles[n]...)
			finalNComposer.ObjectAttrs.ContentType = "text/csv"
			_, err := finalNComposer.Run(ctx)
			if err != nil {
				log.Println("Error when composing the final N file:", err)
				allFinalNFilesOK = false
				break
			} else {
				allFinalNObjects = append(allFinalNObjects, FinalNObject)
			}
		}

		if allFinalNFilesOK {
			bigQueryImportFilename := "_ngram_BigQueryImport.csv"
			bigQueryImportObject := bucket.Object(bigQueryImportFilename)

			log.Printf("Composing FINAL file: %s", bigQueryImportFilename)

			bigQueryComposer := bigQueryImportObject.ComposerFrom(allFinalNObjects...)
			bigQueryComposer.ObjectAttrs.ContentType = "text/csv"
			_, err := bigQueryComposer.Run(ctx)
			if err != nil {
				log.Println("Error when composing the BigQuery import file:", err)
			} else {
				log.Printf("FINAL file composed!")

				var sourceFiles []string
				sourceFiles = append(sourceFiles, "(all final files)")
				var combinedFile CombinedFile
				combinedFile.N = 25
				combinedFile.Letter = "a-z"
				combinedFile.SourceFiles = sourceFiles
				combinedFile.TargetFile = bigQueryImportFilename
				combineImportData.CombinedFiles = append(combineImportData.CombinedFiles, combinedFile)
			}
		}
	}

	resultJson, err := json.Marshal(combineImportData)

	if err != nil {
		log.Fatal("Error when JSONing the result: ", err)
	}

	w.Write(resultJson)

	log.Println("DONE handling", r.URL)
}

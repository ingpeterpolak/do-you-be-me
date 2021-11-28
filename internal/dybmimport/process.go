package dybmimport

import (
	"bufio"
	"bytes"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// HandleProcess handles the /process URL and does all the work
// It takes all the imported ngrams and tries to guess:
//	- the number of syllables
//	- th pronunciation
func HandleProcess(w http.ResponseWriter, r *http.Request) {
	log.Println("START handling", r.URL)
	log.Println("Proto", r.Proto, "TLS", r.TLS, "Host", r.Host)

	w.Header().Add("Content-type", "application/json")

	/*
		pronFilename := DataFolder + "slavic-pronunciations.csv"
		pronFile, err := os.Open(pronFilename)
		if err != nil {
			log.Fatal("Couldn't open", pronFilename, err)
		}

		var pronWords []string
		var objectsToDelete []string
		pronScanner := bufio.NewScanner(pronFile)
		for pronScanner.Scan() {
			line := pronScanner.Text()
			fragments := strings.Split(line, ";")
			pronWords = append(pronWords, fragments[0])
		}
		pronFile.Close()

		notFoundCount := 0
		allCount := 0
		ctx, _, bucket := prepareContext()
		it := bucket.Objects(ctx, nil)
		for {
			attrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatal("Error when accessing object", err)
			}

			allCount++

			word := attrs.Name[0 : len(attrs.Name)-4] // drop the .csv

			wordFound := false
			for _, pronWord := range pronWords {
				if pronWord == word {
					wordFound = true
				}
			}

			if !wordFound {
				log.Println("Word", word, "not found")
				objectsToDelete = append(objectsToDelete, word+".csv")
				notFoundCount++
			}
		}

		log.Println("Deleting", notFoundCount, "objects:")
		for _, filename := range objectsToDelete {
			objectToDelete := bucket.Object(filename)
			objectToDelete.Delete(ctx)
			log.Print(" .")
		}
		log.Println(" Deleted.")

		log.Println("Done.", notFoundCount, "words not found out of", allCount)
	*/

	// log.Println(extractRhyme("ädörÍz"))

	/*
		cmuDictFilename := DataFolder + "cmudict-0.7b.csv"
		cmuDictFile, err := os.Open(cmuDictFilename)
		if err != nil {
			log.Fatal("CMU Dict data file not present", cmuDictFilename, err)
		}

		var semicolonSeparator = [...]byte{59}
		var newLineSeparator = [...]byte{10}
		pronFilename := DataFolder + "slavic-pronunciations.csv"
		pronFile, err := os.Create(pronFilename)
		if err != nil {
			log.Fatal("Couldn't create", pronFilename, err)
		}

		cmuDict := make(map[string]string)
		cmuDictScanner := bufio.NewScanner(cmuDictFile)
		for cmuDictScanner.Scan() {
			line := cmuDictScanner.Text()
			fragments := strings.Split(line, ";")

			if containsParenthesis(fragments[0]) {
				continue
			}

			cmuDict[fragments[0]] = fragments[1]

			p := getPronunciation(fragments[1])
			r := extractRhyme(p)

			pronFile.Write([]byte(fragments[0]))
			pronFile.Write(semicolonSeparator[:])
			pronFile.Write([]byte(p))
			pronFile.Write(semicolonSeparator[:])
			pronFile.Write([]byte(r.Strong))
			pronFile.Write(semicolonSeparator[:])
			pronFile.Write([]byte(r.Average))
			pronFile.Write(semicolonSeparator[:])
			pronFile.Write([]byte(r.Weak))
			pronFile.Write(newLineSeparator[:])
		}
		pronFile.Close()
		cmuDictFile.Close()
	*/

	pronDict := make(map[string]string)
	pronFilename := DataFolder + "slavic-pronunciations.csv"
	pronFile, err := os.Open(pronFilename)
	if err != nil {
		log.Fatal("Slavic pronunciation data file not present", pronFilename, err)
	}
	pronScanner := bufio.NewScanner(pronFile)
	for pronScanner.Scan() {
		line := pronScanner.Text()
		fragments := strings.Split(line, ";")
		pronDict[fragments[0]] = fragments[1]
	}
	pronFile.Close()

	ctx, bucket, _ := prepareContext()

	n := r.URL.Query().Get("n")
	allNgramsFilename := "Final_" + n + "gram.csv"
	processedNgramsFilename := "_" + n + "grams_Processed.csv"

	ngramsFileReader, err := bucket.Object(allNgramsFilename).NewReader(ctx)
	if err != nil {
		log.Fatal("Unable to open file", allNgramsFilename, err)
	}
	defer ngramsFileReader.Close()

	gcWriter := bucket.Object(processedNgramsFilename).NewWriter(ctx)
	gcWriter.ContentType = "text/csv"

	var bufferBuilder strings.Builder
	bufferLength := 0
	bufferCapacity := 100000
	var linesWritten int = 0

	// line format:
	// word1 word2;123
	var semicolonSeparator = [...]byte{59}
	var newLineSeparator = [...]byte{10}

	i := 0
	scanner := bufio.NewScanner(ngramsFileReader)
	for scanner.Scan() {
		i++
		if i%1000000 == 0 {
			log.Println("+1.000.000 lines processed, currently at", i)
		}

		//if i > 1000000 {
		//	break
		//}

		line := scanner.Bytes()
		fragments := bytes.Split(line, semicolonSeparator[:])

		ngram := string(fragments[0])

		if !isNgramSuitableForLyrics(ngram) {
			continue
		}

		words := strings.Split(ngram, " ")

		ngramSyllables := 0
		var sb strings.Builder
		wasFound := true
		for _, word := range words {
			var pronunciation string
			pronunciation, wasFound = pronDict[strings.ToLower(word)]
			if !wasFound {
				break
			}

			sb.WriteString(pronunciation)
			syllables, _ := CountSyllables(word)
			ngramSyllables += syllables
		}

		if ngramSyllables > 16 {
			continue
		}

		if wasFound {
			frequencyRaw := convertAsciiNumberToInt(fragments[1])
			frequency := int(math.Log(float64(frequencyRaw)) * 4.9)
			if frequency < 1 {
				frequency = 1
			}
			if frequency > 100 {
				frequency = 100
			}

			rhyme := extractRhyme(sb.String())

			bufferBuilder.WriteString(ngram)
			bufferBuilder.WriteByte(semicolonSeparator[0])
			bufferBuilder.WriteString(strconv.Itoa(frequency))
			bufferBuilder.WriteByte(semicolonSeparator[0])
			bufferBuilder.WriteString(strconv.Itoa(ngramSyllables))
			bufferBuilder.WriteByte(semicolonSeparator[0])
			bufferBuilder.WriteString(rhyme.Strong)
			bufferBuilder.WriteByte(semicolonSeparator[0])
			bufferBuilder.WriteString(rhyme.Average)
			bufferBuilder.WriteByte(semicolonSeparator[0])
			bufferBuilder.WriteString(rhyme.Weak)
			bufferBuilder.WriteByte(newLineSeparator[0])
			bufferLength++

			if bufferLength > bufferCapacity {
				linesWritten += bufferLength

				log.Println("Buffer flush", bufferLength, "total", linesWritten)
				gcWriter.Write([]byte(bufferBuilder.String()))

				bufferBuilder.Reset()
				bufferLength = 0
			}
		}
	}

	if bufferLength > 0 {
		log.Println("Final buffer flush", bufferLength, "with total lines written:", linesWritten)
		gcWriter.Write([]byte(bufferBuilder.String()))
		bufferBuilder.Reset()
	}

	gcWriter.Close()

	/*
		testWords := [...]string{"counterrevolutionary", "inevitability", "irresponsibility", "undifferentiated", "responsibility"}
		for _, word := range testWords {
			syllablesCount, knownDataUsed := CountSyllables(word)
			log.Println(word, "has", syllablesCount, "syllables. Known:", knownDataUsed)
		}*/

	/*
		lyricsFilename := DataFolder + "tcc_ceds_music.csv"
		jsonOutputFilename := DataFolder + "relatedWords.json"

		inputFile, err := os.Open(lyricsFilename)
		if err != nil {
			log.Fatal("Lyrics data file not present", lyricsFilename, err)
		}

		wordsWithRelations := make(map[string]map[string]int)
		scanner := bufio.NewScanner(inputFile)

		// the first line is a header
		scanner.Scan()

		for scanner.Scan() {
			line := scanner.Text()
			fragments := strings.Split(line, ",")
			lyrics := fragments[5]
			firstWords := strings.Split(lyrics, " ")
			secondWords := firstWords[:]
			if len(firstWords) > 5 {
				for i, firstWord := range firstWords {
					if firstWord != "" {
						relatedWords, wasFound := wordsWithRelations[firstWord]
						if !wasFound {
							wordsWithRelations[firstWord] = make(map[string]int)
							relatedWords = wordsWithRelations[firstWord]
						}
						for j, secondWord := range secondWords {
							if secondWord != "" && i != j {
								relatedWords[secondWord]++
							}
						}
					}
				}
			}

		}
		inputFile.Close()

		ctx, _, bucket := prepareContext()

		for word, relatedWords := range wordsWithRelations {
			var strengths []int
			for _, strength := range relatedWords {
				strengths = append(strengths, strength)
			}

			sort.Ints(strengths)
			strengthsCount := len(strengths)
			minIndex := 0
			if strengthsCount > MaxRelatedWordsPerWord {
				minIndex = strengthsCount - 1 - MaxRelatedWordsPerWord + 1
			}

			var reducedRelatedWords []string
			for i := strengthsCount - 1; i >= minIndex; i-- {
				for relatedWord, strength := range relatedWords {
					if strength == strengths[i] {
						reducedRelatedWords = append(reducedRelatedWords, relatedWord)
						wordsWithRelations[word][relatedWord] = 0 // for next iteration, we'll ignore this word as there might be more words with the same strength
						break
					}
				}
			}

			wordWithRelationsFilename := fmt.Sprintf("%s.csv", word)
			wordWithRelationsObject := bucket.Object(wordWithRelationsFilename)
			gcWriter := wordWithRelationsObject.NewWriter(ctx)
			gcWriter.ContentType = "text/csv"

			for _, relatedWord := range reducedRelatedWords {
				wordToWrite := []byte(relatedWord)
				wordToWrite = append(wordToWrite, 10) // ASCII for Line Feed
				gcWriter.Write(wordToWrite)
			}

			gcWriter.Close()
		}

		w.Write([]byte(fmt.Sprintf("Output file %s created", jsonOutputFilename)))
	*/

	log.Println("DONE handling", r.URL)
}

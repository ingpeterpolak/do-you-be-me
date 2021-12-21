package dybmapi

import (
	"context"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/ingpeterpolak/do-you-be-me/internal/dybmpronounce"
	"github.com/ingpeterpolak/do-you-be-me/internal/dybmrhyme"

	"cloud.google.com/go/storage"
)

var relatedWords map[string]int

func filterWords(words []string) []string {
	var filteredWords []string

	for _, word := range words {
		if word != "" && word != " " && word != "a" && word != "an" {
			filteredWords = append(filteredWords, word)
		}
	}

	return filteredWords
}

func getLine(lineToRhyme, syllableCount, lyrics string) PimpedLine {
	words := filterWords(getPureWords(lyrics))
	buildRelatedWords(words)

	syllables, err := strconv.Atoi(syllableCount)
	if err != nil {
		log.Print("Syllable count is not a number. Using 6 instead")
		syllables = 6
	}

	wordsOfLineToRhyme := strings.Split(lineToRhyme, " ")
	lastWord := wordsOfLineToRhyme[len(wordsOfLineToRhyme)-1]

	pronunciation := dybmpronounce.Pronounce(lineToRhyme)
	rhyme := dybmrhyme.ExtractRhyme(pronunciation)

	ngrams := getAllMatchingNgrams(rhyme, syllables, lastWord)

	var pimpedLine PimpedLine

	pimpedLine.Line = ngrams[0]
	pimpedLine.Syllables = syllables

	return pimpedLine
}

func buildRelatedWords(lyricsWords []string) {
	relatedWords = make(map[string]int)

	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("Failed to create Cloud Storage client: ", err)
	}
	bucket := storageClient.Bucket(bucketName)

	for _, lyricsWord := range lyricsWords {
		lyricsWordObject := bucket.Object("related/" + lyricsWord + ".csv")

		reader, err := lyricsWordObject.NewReader(ctx)
		if err != nil {
			continue // the word doesn't have a related words record
		}
		defer reader.Close()

		wordsText, err := ioutil.ReadAll(reader)
		if err != nil {
			continue // let's pretend the word doesn't have a related words record
		}

		words := strings.Split(string(wordsText), "\n")

		weight := maxRelatedWordsCount
		for _, word := range words {
			relatedWords[word] += weight

			weight--
			if weight == 0 {
				break
			}
		}
	}
}

func getAllMatchingNgrams(rhyme dybmrhyme.Rhyme, syllables int, lastWord string) []string {
	var result []string
	result = append(result, "Just like Marie Antoinette")
	return result
}

package dybmimport

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

var wordsWithSyllableCount map[string]int

// initializeKnownSyllables initializes the syllable count map and reads the known words from the data file
func initializeKnownSyllables() {
	wordsWithSyllableCount = make(map[string]int)

	syllablesCountFilename := DataFolder + "syllables.csv"

	csvFile, err := os.Open(syllablesCountFilename)
	if err != nil {
		log.Println("WARNING: syllables data file not present (data/syllables.csv). Will figure out the syllable counts by guessing", err)
		wordsWithSyllableCount["a"] = 1
	} else {
		scanner := bufio.NewScanner(csvFile)
		for scanner.Scan() {
			line := scanner.Text()
			fragments := strings.Split(line, ";")
			count, _ := strconv.Atoi(fragments[1])
			wordsWithSyllableCount[fragments[0]] = count
		}
		csvFile.Close()
	}
}

// countSyllables returns the syllable count in a word
// first it tries to find the word in the list of known words, if not there, it tries to count the syllables manually
// the second return value indicates if the count was found in the known data
func CountSyllables(word string) (int, bool) {
	if wordsWithSyllableCount == nil {
		initializeKnownSyllables()
	}

	if word == "" {
		return 0, false
	}

	word = strings.ToLower(word)
	result, isKnown := wordsWithSyllableCount[word]
	if !isKnown {
		if isVowel(word[0]) {
			result++
		}
		for i := 1; i < len(word); i++ {
			if isVowel(word[i]) && !isVowel(word[i-1]) {
				result++
			}
		}
		if len(word) > 2 {
			if string(word[len(word)-1]) == "e" && (!isVowel(word[len(word)-2]) || string(word[len(word)-2]) == "u") {
				result--
			}
			if word[len(word)-2:] == "le" && !isVowel(word[len(word)-3]) {
				result++
			}
			if word[0:2] == "mc" {
				result++
			}
			if word[0:2] == "bi" && isVowel(word[2]) {
				result++
			}
		}
		if len(word) > 3 {
			if word[0:3] == "tri" && isVowel(word[3]) {
				result++
			}
		}
		if len(word) > 4 {
			if word[len(word)-3:] == "ian" && string(word[len(word)-4]) != "c" && string(word[len(word)-4]) != "t" {
				result++
			}
		}
		if len(word) > 6 {
			if word[0:2] == "co" && string(word[2]) != "u" && isVowel(word[2]) {
				result++
			}
			if word[0:3] == "pre" && isVowel(word[3]) {
				result++
			}
		}

		if result <= 0 {
			result = 1
		}

		if result >= 10 {
			result = 9
		}
	}

	return result, isKnown
}

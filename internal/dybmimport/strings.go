package dybmimport

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const AsciiNumDiff = 48 // the difference between the ascii code of a number and the number. E. g. the ascii code of the number 3 is 51

const corpusBucketName = "dybm-corpus-1"
const relatedWordsBucketName = "dybm-related-words-1"
const urlsFictionFilename = "google-ngrams-fiction-urls.txt"

var wordsWithSyllableCount map[string]int
var pronunciations map[string]string

// var urlsFilename = "google-ngrams-urls.txt"
var validLetters = [...]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

func isVowel(b byte) bool {
	r := string(b)
	return r == "a" || r == "e" || r == "i" || r == "o" || r == "u" || r == "y" || r == "ï" || r == "î" || r == "é" || r == "ê" || r == "è" || r == "à" || r == "â" || r == "ô"
}

// isAllLetters returns true if the given string is nothing but letters and spaces
func isAllLetters(s string) bool {
	result := true

	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) {
			result = false
			break
		}
	}

	return result
}

/*func replaceGoogleNgramKeywords(s string) string {
	r := strings.NewReplacer(
		"_ADJ", "",
		"_ADP", "",
		"_ADV", "",
		"_CONJ", "",
		"_DET", "",
		"_NOUN", "",
		"_NUM", "",
		"_PRON", "",
		"_PRT", "",
		"_VERB", "")
	return r.Replace(s)
}*/

// isUrlForNgramAndLetter checks if a given URL has the given ngrams and starting letter
// input url looks like this: http://storage.googleapis.com/books/ngrams/books/googlebooks-eng-all-2gram-20120701-bb.gz
// fiction has 8 more chars   http://storage.googleapis.com/books/ngrams/books/googlebooks-eng-fiction-all-2gram-20120701-a_.gz
func isUrlForNgramAndLetter(url, n, letter string) bool {
	nFromUrl := string(url[69+8])
	letterFromUrl := string(url[84+8 : 85+8])
	return n == nFromUrl && letter == letterFromUrl
}

// getNgramFilenameFromUrl returns target filename for a given url
// input url looks like this: http://storage.googleapis.com/books/ngrams/books/googlebooks-eng-all-2gram-20120701-bb.gz
// fiction has 8 more chars   http://storage.googleapis.com/books/ngrams/books/googlebooks-eng-fiction-all-2gram-20120701-a_.gz
func getNgramFilenameFromUrl(url string, raw bool) string {
	//	return fmt.Sprintf("googlebooks-eng-all-%sgram-%s.csv", n, letter)
	nFromUrl := string(url[69+8])
	lettersFromUrl := string(url[84+8 : 86+8])
	postfix := ""
	if raw {
		postfix = ".src"
	}
	//return fmt.Sprintf("googlebooks-eng-all-%sgram-%s%s.csv", nFromUrl, lettersFromUrl, postfix)
	return fmt.Sprintf("googlebooks-eng-fiction-all-%sgram-%s%s.csv", nFromUrl, lettersFromUrl, postfix)
}

// getNgramTargetFilename returns target filename for a given n and letter
// target filename example for n=2 and letter=a: googlebooks-eng-all-2gram-a.csv
func getNgramTargetFilename(n, letter string) string {
	//return fmt.Sprintf("googlebooks-eng-all-%sgram-%s.csv", n, letter)
	return fmt.Sprintf("googlebooks-eng-fiction-all-%sgram-%s.csv", n, letter)
}

// convertAsciiNumberToInt converts a string represented by []byte into a number
// it uses ascii codes of numbers so that it doesn't need to convert the whole thing into a string and then use Itoa
func convertAsciiNumberToInt(b []byte) int {
	var result int = 0
	order := len(b)
	for i := 0; i < order; i++ {
		var tenPowered int = 1
		for j := 0; j < order-i-1; j++ { // this could be done better...
			tenPowered *= 10
		}
		result += tenPowered * (int(b[i]) - AsciiNumDiff) //48 is the difference between the ascii code of a number and the number itself
	}
	return result
}

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

// isNgramSuitableForLyrics returns true if the provided ngram is suitable for lyrics
// e.g. doesn't end with " the"
func isNgramSuitableForLyrics(ngram string) bool {
	if len(ngram) > 2 {
		if ngram[len(ngram)-2:] == " a" {
			return false
		}
	}

	if len(ngram) > 3 {
		if ngram[len(ngram)-3:] == " an" || ngram[len(ngram)-3:] == " as" || ngram[len(ngram)-3:] == " de" || ngram[len(ngram)-3:] == " du" || ngram[len(ngram)-3:] == " my" {
			return false
		}
	}

	if len(ngram) > 4 {
		if ngram[len(ngram)-4:] == " the" || ngram[len(ngram)-4:] == " its" || ngram[len(ngram)-4:] == " our" || ngram[len(ngram)-4:] == " his" {
			return false
		}
	}

	if len(ngram) > 5 {
		if ngram[len(ngram)-5:] == " than" || ngram[len(ngram)-3:] == " your" || ngram[len(ngram)-3:] == " hers" {
			return false
		}
	}

	if len(ngram) > 6 {
		if ngram[len(ngram)-6:] == " their" {
			return false
		}
	}

	return true
}

func initializePronunciation() {
	pronunciations = make(map[string]string)

	pronunciations["aa0"] = "a"
	pronunciations["aa1"] = "Á"
	pronunciations["aa2"] = "Á"

	pronunciations["ae0"] = "ä"
	pronunciations["ae1"] = "Ä"
	pronunciations["ae2"] = "Ä"

	pronunciations["ah0"] = "ö"
	pronunciations["ah1"] = "A"
	pronunciations["ah2"] = "a"

	pronunciations["ao0"] = "o"
	pronunciations["ao1"] = "O"
	pronunciations["ao2"] = "O"

	pronunciations["aw0"] = "au"
	pronunciations["aw1"] = "AU"
	pronunciations["aw2"] = "AU"

	pronunciations["ay0"] = "ai"
	pronunciations["ay1"] = "AI"
	pronunciations["ay2"] = "AI"

	pronunciations["b"] = "b"
	pronunciations["ch"] = "č"
	pronunciations["d"] = "d"
	pronunciations["dh"] = "d"

	pronunciations["eh0"] = "e"
	pronunciations["eh1"] = "E"
	pronunciations["eh2"] = "E"

	pronunciations["er0"] = "ör"
	pronunciations["er1"] = "ÖR"
	pronunciations["er2"] = "ÖR"

	pronunciations["ey0"] = "ei"
	pronunciations["ey1"] = "EI"
	pronunciations["ey2"] = "EI"

	pronunciations["f"] = "f"
	pronunciations["g"] = "g"
	pronunciations["hh"] = "h"

	pronunciations["ih0"] = "i"
	pronunciations["ih1"] = "I"
	pronunciations["ih2"] = "I"

	pronunciations["iy0"] = "i"
	pronunciations["iy1"] = "Í"
	pronunciations["iy2"] = "Í"

	pronunciations["jh"] = "dž"
	pronunciations["k"] = "k"
	pronunciations["l"] = "l"
	pronunciations["m"] = "m"
	pronunciations["n"] = "n"
	pronunciations["ng"] = "n"

	pronunciations["ow0"] = "ou"
	pronunciations["ow1"] = "OU"
	pronunciations["ow2"] = "OU"

	pronunciations["oy0"] = "oi"
	pronunciations["oy1"] = "OI"
	pronunciations["oy2"] = "OI"

	pronunciations["p"] = "p"
	pronunciations["r"] = "r"
	pronunciations["s"] = "s"
	pronunciations["sh"] = "š"
	pronunciations["t"] = "t"
	pronunciations["th"] = "t"

	pronunciations["uh0"] = "u"
	pronunciations["uh1"] = "U"
	pronunciations["uh2"] = "U"

	pronunciations["uw0"] = "u"
	pronunciations["uw1"] = "U"
	pronunciations["uw2"] = "U"

	pronunciations["v"] = "v"
	pronunciations["w"] = "w"
	pronunciations["y"] = "j"
	pronunciations["z"] = "z"
	pronunciations["zh"] = "ž"
}

// getPronunciation gets the slavic pronunciation from the given CMU pronunciation
func getPronunciation(cmu string) string {
	if pronunciations == nil {
		initializePronunciation()
	}

	var builder strings.Builder
	cmuSymbols := strings.Split(cmu, " ")
	for _, symbol := range cmuSymbols {
		builder.WriteString(pronunciations[symbol])
	}

	return builder.String()
}

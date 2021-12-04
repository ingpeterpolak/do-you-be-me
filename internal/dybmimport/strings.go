package dybmimport

import (
	"fmt"
	"unicode"
)

const AsciiNumDiff = 48 // the difference between the ascii code of a number and the number. E. g. the ascii code of the number 3 is 51

const corpusBucketName = "dybm-corpus-1"
const relatedWordsBucketName = "dybm-related-words-1"
const urlsFictionFilename = "google-ngrams-fiction-urls.txt"

// var urlsFilename = "google-ngrams-urls.txt"
var validLetters = [...]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var validUcLetters = [...]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

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
	var tenPowered int = 1

	for i := order - 1; i >= 0; i-- {
		result += tenPowered * (int(b[i]) - AsciiNumDiff) //48 is the difference between the ascii code of a number and the number itself
		tenPowered *= 10
	}
	return result
}

package dybmapi

import (
	"fmt"
	"strings"
	"unicode"
)

var bucketName = "dybm-corpus-1"

// var urlsFilename = "google-ngrams-urls.txt"
var urlsFictionFilename = "google-ngrams-fiction-urls.txt"
var validLetters = [...]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

func removeSpecialChars(s string) string {
	return strings.Map(
		func(r rune) rune {
			ch := string(r)
			if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) || ch == "'" || ch == "-" || ch == "," || ch == "." || ch == "&" || ch == "!" || ch == "?" {
				return r
			}
			return -1
		},
		s,
	)
}

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

// check if a given URL has the given ngrams and starting letter
// input url looks like this: http://storage.googleapis.com/books/ngrams/books/googlebooks-eng-all-2gram-20120701-bb.gz
// fiction has 8 more chars   http://storage.googleapis.com/books/ngrams/books/googlebooks-eng-fiction-all-2gram-20120701-a_.gz
func isUrlForNgramAndLetter(url, n, letter string) bool {
	nFromUrl := string(url[69+8])
	letterFromUrl := string(url[84+8 : 85+8])
	return n == nFromUrl && letter == letterFromUrl
}

// returns target filename for a given url
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

// returns target filename for a given n and letter
// target filename example for n=2 and letter=a: googlebooks-eng-all-2gram-a.csv
func getNgramTargetFilename(n, letter string) string {
	//return fmt.Sprintf("googlebooks-eng-all-%sgram-%s.csv", n, letter)
	return fmt.Sprintf("googlebooks-eng-fiction-all-%sgram-%s.csv", n, letter)
}

func convertBytesNumberToInt(b []byte) int {
	var result int = 0
	order := len(b)
	for i := 0; i < order; i++ {
		var tenPowered int = 1
		for j := 0; j < order-i-1; j++ {
			tenPowered *= 10
		}
		result += tenPowered * (int(b[i]) - 48) //48 is the difference between the ascii code of a number and the number itself
	}
	return result
}

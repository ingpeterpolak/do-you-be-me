package dybmapi

import (
	"fmt"
	"strings"
	"unicode"
)

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
func isUrlForNgramAndLetter(url, n, letter string) bool {
	nFromUrl := string(url[69])
	letterFromUrl := string(url[84:85])
	return n == nFromUrl && letter == letterFromUrl
}

// returns target filename for a given url
// input url looks like this: http://storage.googleapis.com/books/ngrams/books/googlebooks-eng-all-2gram-20120701-bb.gz
func getNgramTargetFilename(url string) string {
	//	return fmt.Sprintf("googlebooks-eng-all-%sgram-%s.csv", n, letter)
	nFromUrl := string(url[69])
	lettersFromUrl := string(url[84:86])
	return fmt.Sprintf("googlebooks-eng-all-%sgram-%s.csv", nFromUrl, lettersFromUrl)
}

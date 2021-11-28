package dybmimport

import (
	"strings"
)

// isNgramSuitableForLyrics returns true if the provided ngram is suitable for lyrics
// e.g. doesn't end with " the", doesn't have words in UPPERCASE etc.
func isNgramSuitableForLyrics(ngramRaw string) bool {

	length := len(ngramRaw)
	for i := 0; i <= length-1-1; i++ {
		firstCh := ngramRaw[i : i+1]
		secondCh := ngramRaw[i+1 : i+1+1]

		firstIsUpper := false
		secondIsUpper := false
		for _, ch := range validUcLetters {
			if !firstIsUpper && firstCh == ch {
				firstIsUpper = true
			}
			if !secondIsUpper && secondCh == ch {
				secondIsUpper = true
			}
			if firstIsUpper && secondIsUpper {
				return false
			}
		}
	}

	ngram := strings.ToLower(ngramRaw)

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

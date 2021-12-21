package dybmapi

import (
	"strings"
	"unicode"
)

// getPureWords removes all non-alphanumeric chars and returns pure words
func getPureWords(s string) []string {
	line := strings.Map(
		func(r rune) rune {
			ch := string(r)
			if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) || ch == "'" {
				return r
			}
			return ' '
		},
		s,
	)

	var result []string
	words := strings.Split(line, " ")
	for _, word := range words {
		if word != "" {
			result = append(result, word)
		}
	}

	return result
}

// removeSpecialCharsFromLyrics removes chars that (usually) don't belong to a song lyrics
func removeSpecialCharsFromLyrics(s string) string {
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

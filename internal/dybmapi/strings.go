package dybmapi

import (
	"strings"
	"unicode"
)

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

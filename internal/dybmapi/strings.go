package dybmapi

import (
	"strings"
	"unicode"
)

func removeSpecialChars(s string) string {
	return strings.Map(
		func(r rune) rune {
			ch := string(r)
			if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) || ch == "'" || ch == "-" || ch == "," || ch == "." || ch == "&" {
				return r
			}
			return -1
		},
		s,
	)
}

package dybmapi

import (
	"strings"
	"unicode"
)

// removeNonAlphanumeric removes all non-alphanumeric chars
func removeNonAlphanumeric(s string) string {
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

	var sb strings.Builder
	words := strings.Split(line, " ")
	for _, word := range words {
		if word != "" {
			sb.WriteString(word)
			sb.WriteString(" ")
		}
	}

	return sb.String()
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

func getNextRhymeId(lastRhymeId string) string {
	if lastRhymeId == "" {
		return "A"
	}

	result := "Z"
	rhymeIds := "ABCDEFGHIJKLMNOPQRSTUVWZ"
	for i := 0; i < len(rhymeIds)-2; i++ {
		if lastRhymeId == rhymeIds[i:i+1] {
			result = rhymeIds[i+1 : i+2]
		}
	}

	return result
}

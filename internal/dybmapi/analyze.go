package dybmapi

import (
	"strings"

	"github.com/ingpeterpolak/do-you-be-me/internal/dybmpronounce"
	"github.com/ingpeterpolak/do-you-be-me/internal/dybmrhyme"
	"github.com/ingpeterpolak/do-you-be-me/internal/dybmsyllable"
)

func analyzeLines(lines []string) []PimpedLine {
	var lastRhymeId = ""
	var pimpedLines []PimpedLine
	var rhymes []Rhyme
	for i, s := range lines {

		line := removeSpecialCharsFromLyrics(s)
		pureWords := removeNonAlphanumeric(line)

		if pureWords == "" {
			continue
		}

		pronunciation := dybmpronounce.Pronounce(pureWords)
		guessedRhyme := dybmrhyme.GuessRhyme(pronunciation)

		words := strings.Split(pureWords, " ")
		syllablesCount := 0
		for _, word := range words {
			count, _ := dybmsyllable.CountSyllables(word)
			syllablesCount += count
		}

		foundRhymeId := ""
		for _, rhyme := range rhymes {
			if guessedRhyme == rhyme.Rhyme {
				foundRhymeId = rhyme.Id
				break
			}
		}

		if foundRhymeId == "" {
			lastRhymeId = dybmrhyme.GetNextRhymeId(lastRhymeId)

			var rhyme Rhyme
			rhyme.Id = lastRhymeId
			rhyme.Rhyme = guessedRhyme

			rhymes = append(rhymes, rhyme)
			foundRhymeId = rhyme.Id
		}

		var pimpedLine PimpedLine
		pimpedLine.Number = i + 1
		pimpedLine.Line = line
		pimpedLine.Syllables = syllablesCount
		pimpedLine.RhymeId = foundRhymeId

		pimpedLines = append(pimpedLines, pimpedLine)
	}

	return pimpedLines
}

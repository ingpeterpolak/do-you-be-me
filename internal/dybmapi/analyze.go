package dybmapi

import (
	"log"
	"strings"

	"github.com/ingpeterpolak/do-you-be-me/internal/dybmpronounce"
	"github.com/ingpeterpolak/do-you-be-me/internal/dybmrhyme"
	"github.com/ingpeterpolak/do-you-be-me/internal/dybmsyllable"
)

func compareRhymes(r1, r2 []string) bool {
	rhymeCount := len(r1)
	if rhymeCount != len(r2) {
		return false
	}

	areEqual := true
	for i := 0; i < rhymeCount; i++ {
		if r1[i] != r2[i] {
			areEqual = false
			break
		}
	}

	return areEqual
}

func analyzeScheme(detectedRhymes []PimpedLine) []string {
	var rhymes []string
	for _, detectedRhyme := range detectedRhymes {
		rhymes = append(rhymes, detectedRhyme.RhymeId)
	}

	lineCount := len(detectedRhymes)
	maxLinesPerRhyme := lineCount / 2
	if maxLinesPerRhyme > 8 {
		maxLinesPerRhyme = 8
	}

	var firstVerseRhymes, secondVerseRhymes []string
	introLength := 0
	verseLength := 0
	for firstLine := 0; firstLine <= 8; firstLine++ { // allow for some intro lines that might not contribute to rhyming scheme
		for i := firstLine + maxLinesPerRhyme - 1; i >= firstLine+3; i-- { // go from 8-line down to 3-line schemes

			firstVerseStart := firstLine
			firstVerseEnd := firstLine + i + 1
			secondVerseStart := firstVerseEnd
			secondVerseEnd := secondVerseStart + i + 1

			firstVerse := rhymes[firstVerseStart:firstVerseEnd]
			firstVerseRhymes = dybmrhyme.ResetRhymeIds(firstVerse)

			secondVerse := rhymes[secondVerseStart:secondVerseEnd]
			secondVerseRhymes = dybmrhyme.ResetRhymeIds(secondVerse)

			if compareRhymes(firstVerseRhymes, secondVerseRhymes) {
				verseLength = i + 1
				introLength = firstLine
				break
			}
		}
		if verseLength > 0 {
			break
		}
	}

	var result []string
	result = append(result, rhymes[0:introLength]...)
	result = append(result, firstVerseRhymes...)
	result = append(result, secondVerseRhymes...)
	result = append(result, rhymes[introLength+2*verseLength:]...)

	return result
}

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
		pimpedLine.RhymeId = foundRhymeId
		pimpedLine.Syllables = syllablesCount

		pimpedLines = append(pimpedLines, pimpedLine)
	}

	rhymeScheme := analyzeScheme(pimpedLines)
	log.Print("Scheme:", rhymeScheme)

	for i, rhyme := range rhymeScheme {
		pimpedLines[i].RhymeId = rhyme
	}

	return pimpedLines
}

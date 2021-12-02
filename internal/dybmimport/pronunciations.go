package dybmimport

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var cmuToSlavicPronunciations map[string]string
var pronunciations map[string]string

func isVowel(b byte) bool {
	r := string(b)
	return r == "a" || r == "e" || r == "i" || r == "o" || r == "u" || r == "y" || r == "ï" || r == "î" || r == "é" || r == "ê" || r == "è" || r == "à" || r == "â" || r == "ô"
}

func initializePronunciation() {
	cmuToSlavicPronunciations = make(map[string]string)

	cmuToSlavicPronunciations["aa0"] = "a"
	cmuToSlavicPronunciations["aa1"] = "Á"
	cmuToSlavicPronunciations["aa2"] = "Á"

	cmuToSlavicPronunciations["ae0"] = "ä"
	cmuToSlavicPronunciations["ae1"] = "Ä"
	cmuToSlavicPronunciations["ae2"] = "Ä"

	cmuToSlavicPronunciations["ah0"] = "ö"
	cmuToSlavicPronunciations["ah1"] = "A"
	cmuToSlavicPronunciations["ah2"] = "a"

	cmuToSlavicPronunciations["ao0"] = "o"
	cmuToSlavicPronunciations["ao1"] = "O"
	cmuToSlavicPronunciations["ao2"] = "O"

	cmuToSlavicPronunciations["aw0"] = "au"
	cmuToSlavicPronunciations["aw1"] = "AU"
	cmuToSlavicPronunciations["aw2"] = "AU"

	cmuToSlavicPronunciations["ay0"] = "aj"
	cmuToSlavicPronunciations["ay1"] = "AJ"
	cmuToSlavicPronunciations["ay2"] = "AJ"

	cmuToSlavicPronunciations["b"] = "b"
	cmuToSlavicPronunciations["ch"] = "č"
	cmuToSlavicPronunciations["d"] = "d"
	cmuToSlavicPronunciations["dh"] = "d"

	cmuToSlavicPronunciations["eh0"] = "e"
	cmuToSlavicPronunciations["eh1"] = "E"
	cmuToSlavicPronunciations["eh2"] = "E"

	cmuToSlavicPronunciations["er0"] = "ör"
	cmuToSlavicPronunciations["er1"] = "ÖR"
	cmuToSlavicPronunciations["er2"] = "ÖR"

	cmuToSlavicPronunciations["ey0"] = "ej"
	cmuToSlavicPronunciations["ey1"] = "EJ"
	cmuToSlavicPronunciations["ey2"] = "EJ"

	cmuToSlavicPronunciations["f"] = "f"
	cmuToSlavicPronunciations["g"] = "g"
	cmuToSlavicPronunciations["hh"] = "h"

	cmuToSlavicPronunciations["ih0"] = "i"
	cmuToSlavicPronunciations["ih1"] = "I"
	cmuToSlavicPronunciations["ih2"] = "I"

	cmuToSlavicPronunciations["iy0"] = "i"
	cmuToSlavicPronunciations["iy1"] = "Í"
	cmuToSlavicPronunciations["iy2"] = "Í"

	cmuToSlavicPronunciations["jh"] = "dž"
	cmuToSlavicPronunciations["k"] = "k"
	cmuToSlavicPronunciations["l"] = "l"
	cmuToSlavicPronunciations["m"] = "m"
	cmuToSlavicPronunciations["n"] = "n"
	cmuToSlavicPronunciations["ng"] = "n"

	cmuToSlavicPronunciations["ow0"] = "ou"
	cmuToSlavicPronunciations["ow1"] = "OU"
	cmuToSlavicPronunciations["ow2"] = "OU"

	cmuToSlavicPronunciations["oy0"] = "oj"
	cmuToSlavicPronunciations["oy1"] = "OJ"
	cmuToSlavicPronunciations["oy2"] = "OJ"

	cmuToSlavicPronunciations["p"] = "p"
	cmuToSlavicPronunciations["r"] = "r"
	cmuToSlavicPronunciations["s"] = "s"
	cmuToSlavicPronunciations["sh"] = "š"
	cmuToSlavicPronunciations["t"] = "t"
	cmuToSlavicPronunciations["th"] = "t"

	cmuToSlavicPronunciations["uh0"] = "u"
	cmuToSlavicPronunciations["uh1"] = "U"
	cmuToSlavicPronunciations["uh2"] = "U"

	cmuToSlavicPronunciations["uw0"] = "u"
	cmuToSlavicPronunciations["uw1"] = "U"
	cmuToSlavicPronunciations["uw2"] = "U"

	cmuToSlavicPronunciations["v"] = "v"
	cmuToSlavicPronunciations["w"] = "w"
	cmuToSlavicPronunciations["y"] = "j"
	cmuToSlavicPronunciations["z"] = "z"
	cmuToSlavicPronunciations["zh"] = "ž"
}

// getPronunciation gets the slavic pronunciation from the given CMU pronunciation
func getPronunciation(cmu string) string {
	if cmuToSlavicPronunciations == nil {
		initializePronunciation()
	}

	var builder strings.Builder
	cmuSymbols := strings.Split(cmu, " ")
	for _, symbol := range cmuSymbols {
		builder.WriteString(cmuToSlavicPronunciations[symbol])
	}

	return builder.String()
}

func isPronunciationVowel(r rune) bool {
	return r == 'a' || r == 'á' || r == 'ä' || r == 'ö' || r == 'o' || r == 'u' || r == 'i' || r == 'í' || r == 'e'
}

func Pronounce(line string) string {
	if pronunciations == nil {
		pronunciations = make(map[string]string)
		pronFilename := DataFolder + "slavic-pronunciations.csv"
		pronFile, err := os.Open(pronFilename)
		if err != nil {
			log.Fatal("Couldn't open", pronFilename, err)
		}

		pronScanner := bufio.NewScanner(pronFile)
		for pronScanner.Scan() {
			line := pronScanner.Text()
			fragments := strings.Split(line, ";")
			pronunciations[fragments[0]] = fragments[1]
		}
		pronFile.Close()
	}

	var sb strings.Builder
	words := strings.Split(line, " ")
	for _, word := range words {
		pronunciation, found := pronunciations[word]
		if found {
			sb.WriteString(pronunciation)
		} else {
			sb.WriteString(word) // TODO fix if we don't have pronunciation
		}
	}

	return sb.String()
}

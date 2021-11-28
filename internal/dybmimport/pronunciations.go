package dybmimport

import (
	"strings"
)

var pronunciations map[string]string

func isVowel(b byte) bool {
	r := string(b)
	return r == "a" || r == "e" || r == "i" || r == "o" || r == "u" || r == "y" || r == "ï" || r == "î" || r == "é" || r == "ê" || r == "è" || r == "à" || r == "â" || r == "ô"
}

func initializePronunciation() {
	pronunciations = make(map[string]string)

	pronunciations["aa0"] = "a"
	pronunciations["aa1"] = "Á"
	pronunciations["aa2"] = "Á"

	pronunciations["ae0"] = "ä"
	pronunciations["ae1"] = "Ä"
	pronunciations["ae2"] = "Ä"

	pronunciations["ah0"] = "ö"
	pronunciations["ah1"] = "A"
	pronunciations["ah2"] = "a"

	pronunciations["ao0"] = "o"
	pronunciations["ao1"] = "O"
	pronunciations["ao2"] = "O"

	pronunciations["aw0"] = "au"
	pronunciations["aw1"] = "AU"
	pronunciations["aw2"] = "AU"

	pronunciations["ay0"] = "aj"
	pronunciations["ay1"] = "AJ"
	pronunciations["ay2"] = "AJ"

	pronunciations["b"] = "b"
	pronunciations["ch"] = "č"
	pronunciations["d"] = "d"
	pronunciations["dh"] = "d"

	pronunciations["eh0"] = "e"
	pronunciations["eh1"] = "E"
	pronunciations["eh2"] = "E"

	pronunciations["er0"] = "ör"
	pronunciations["er1"] = "ÖR"
	pronunciations["er2"] = "ÖR"

	pronunciations["ey0"] = "ej"
	pronunciations["ey1"] = "EJ"
	pronunciations["ey2"] = "EJ"

	pronunciations["f"] = "f"
	pronunciations["g"] = "g"
	pronunciations["hh"] = "h"

	pronunciations["ih0"] = "i"
	pronunciations["ih1"] = "I"
	pronunciations["ih2"] = "I"

	pronunciations["iy0"] = "i"
	pronunciations["iy1"] = "Í"
	pronunciations["iy2"] = "Í"

	pronunciations["jh"] = "dž"
	pronunciations["k"] = "k"
	pronunciations["l"] = "l"
	pronunciations["m"] = "m"
	pronunciations["n"] = "n"
	pronunciations["ng"] = "n"

	pronunciations["ow0"] = "ou"
	pronunciations["ow1"] = "OU"
	pronunciations["ow2"] = "OU"

	pronunciations["oy0"] = "oj"
	pronunciations["oy1"] = "OJ"
	pronunciations["oy2"] = "OJ"

	pronunciations["p"] = "p"
	pronunciations["r"] = "r"
	pronunciations["s"] = "s"
	pronunciations["sh"] = "š"
	pronunciations["t"] = "t"
	pronunciations["th"] = "t"

	pronunciations["uh0"] = "u"
	pronunciations["uh1"] = "U"
	pronunciations["uh2"] = "U"

	pronunciations["uw0"] = "u"
	pronunciations["uw1"] = "U"
	pronunciations["uw2"] = "U"

	pronunciations["v"] = "v"
	pronunciations["w"] = "w"
	pronunciations["y"] = "j"
	pronunciations["z"] = "z"
	pronunciations["zh"] = "ž"
}

// getPronunciation gets the slavic pronunciation from the given CMU pronunciation
func getPronunciation(cmu string) string {
	if pronunciations == nil {
		initializePronunciation()
	}

	var builder strings.Builder
	cmuSymbols := strings.Split(cmu, " ")
	for _, symbol := range cmuSymbols {
		builder.WriteString(pronunciations[symbol])
	}

	return builder.String()
}

func isPronunciationVowel(r rune) bool {
	return r == 'a' || r == 'á' || r == 'ä' || r == 'ö' || r == 'o' || r == 'u' || r == 'i' || r == 'í' || r == 'e'
}

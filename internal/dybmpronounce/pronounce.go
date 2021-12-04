package dybmpronounce

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var DataFolder string

var cmuToSlavicPronunciations map[string]string
var pronouncedWords map[string]string

func Setup(dataFolder string) {
	DataFolder = dataFolder

	initWords()
	initCmuToSlavic()
}

func initWords() {
	pronouncedWords = make(map[string]string)

	filename := DataFolder + "slavic-pronunciations.csv"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Couldn't open", filename, err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fragments := strings.Split(line, ";")

		pronouncedWords[fragments[0]] = fragments[1]
	}
	file.Close()
}

func initCmuToSlavic() {
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

// cmuToSlavic gets the slavic pronunciation from the given CMU pronunciation
func cmuToSlavic(cmu string) string {
	var builder strings.Builder
	cmuSymbols := strings.Split(cmu, " ")
	for _, symbol := range cmuSymbols {
		builder.WriteString(cmuToSlavicPronunciations[symbol])
	}

	return builder.String()
}

func GetPronouncedWords() map[string]string {
	return pronouncedWords
}

func CreateSlavicPronuncationFile() {
	cmuDictFilename := DataFolder + "cmudict-0.7b.csv"
	cmuDictFile, err := os.Open(cmuDictFilename)
	if err != nil {
		log.Fatal("CMU Dict data file not present", cmuDictFilename, err)
	}

	var semicolonSeparator = [...]byte{59}
	var newLineSeparator = [...]byte{10}
	pronFilename := DataFolder + "slavic-pronunciations.csv"
	pronFile, err := os.Create(pronFilename)
	if err != nil {
		log.Fatal("Couldn't create", pronFilename, err)
	}

	cmuDict := make(map[string]string)
	cmuDictScanner := bufio.NewScanner(cmuDictFile)
	for cmuDictScanner.Scan() {
		line := cmuDictScanner.Text()
		fragments := strings.Split(line, ";")

		if containsParenthesis(fragments[0]) {
			continue
		}

		cmuDict[fragments[0]] = fragments[1]

		p := cmuToSlavic(fragments[1])

		pronFile.Write([]byte(fragments[0]))
		pronFile.Write(semicolonSeparator[:])
		pronFile.Write([]byte(p))
		pronFile.Write(newLineSeparator[:])
	}
	pronFile.Close()
	cmuDictFile.Close()
}

func Pronounce(line string) string {
	var sb strings.Builder
	words := strings.Split(line, " ")
	for _, word := range words {
		pronunciation, found := pronouncedWords[word]
		if found {
			sb.WriteString(pronunciation)
		} else {
			sb.WriteString(word) // TODO fix if we don't have pronunciation
		}
	}

	return sb.String()
}

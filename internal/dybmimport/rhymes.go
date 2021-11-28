package dybmimport

import (
	"strings"
)

// getLastSyllables returns last syllables from an expression
// count counts groups of vowels and consonants, usually there will be two groups per syllable
// so for two syllables, count = 4
func getLastSyllables(pronunciation string, count int) string {
	runes := []rune(pronunciation)

	firstCharIndex := 0
	halfSyllableCount := 1
	length := len(runes)
	wasVowel := isPronunciationVowel(runes[length-1])
	trimStartingConsonants := (wasVowel && count%2 == 0) || (!wasVowel && count%2 != 0)
	wasConsonantIn2ndGroup := false
	firstCharIndexOffset := 0
	for i := length - 2; i >= 0; i-- {
		isVowel := isPronunciationVowel(runes[i])
		if isVowel != wasVowel {
			halfSyllableCount++
			if halfSyllableCount == count+1 {
				firstCharIndex = i + 1 + firstCharIndexOffset
				if firstCharIndex > length-1 {
					firstCharIndex = length - 1
				}
				break
			}
		}
		if trimStartingConsonants && halfSyllableCount == count {
			if wasConsonantIn2ndGroup {
				firstCharIndexOffset++
			} else {
				wasConsonantIn2ndGroup = true
			}
		}
		wasVowel = isVowel
	}

	return string(runes[firstCharIndex:])
}

// reduceSoundalikes reduces phonemes that sound like each other
func reduceSoundalikes(pronunciation string) string {
	rDoubles := strings.NewReplacer(
		"aa", "a",
		"áa", "a",
		"aá", "a",
		"ee", "e",
		"ii", "i",
		"íi", "i",
		"ií", "i",
		"oo", "o",
		"uu", "u",
		"yy", "y",

		"bb", "b",
		"čč", "č",
		"dd", "d",
		"dždž", "dž",
		"ff", "f",
		"gg", "g",
		"hh", "h",
		"jj", "j",
		"kk", "k",
		"ll", "l",
		"mm", "m",
		"nn", "n",
		"pp", "p",
		"rr", "r",
		"ss", "s",
		"šš", "š",
		"tt", "t",
		"vv", "v",
		"ww", "w",
		"zz", "z",
		"žž", "ž",
		"šž", "š",
		"žš", "š")
	result := rDoubles.Replace(pronunciation)

	rUnvoiced := strings.NewReplacer(
		"bp", "p",
		"bt", "pt",
		"bč", "pč",
		"bs", "ps",
		"bš", "pš",
		"bk", "pk",
		"bf", "pf",

		"dp", "tp",
		"dt", "t",
		"dč", "č",
		"ds", "ts",
		"dš", "č",
		"dk", "tk",
		"df", "tf",

		"dzp", "tsp",
		"dzt", "tst",
		"dzč", "č",
		"dzs", "ts",
		"dzš", "č",
		"dzk", "tsk",
		"dzf", "tsf",

		"džp", "čp",
		"džt", "čt",
		"džč", "č",
		"džs", "č",
		"džš", "čš",
		"džk", "čk",
		"džf", "čf",

		"zp", "sp",
		"zt", "st",
		"zč", "sč",
		"zs", "s",
		"zš", "š",
		"zk", "sk",
		"zf", "sf",

		"žp", "šp",
		"žt", "št",
		"žč", "šč",
		"žs", "šs",
		"žš", "š",
		"žk", "šk",
		"žf", "šf",

		"gp", "kp",
		"gt", "kt",
		"gč", "kč",
		"gs", "ks",
		"gš", "kš",
		"gk", "k",
		"gf", "kf",

		"vp", "fp",
		"vt", "ft",
		"vč", "kč",
		"vs", "ks",
		"vš", "kš",
		"vk", "k",
		"vf", "kf")
	result = rUnvoiced.Replace(result)

	rVoiced := strings.NewReplacer(
		"pb", "b",
		"pd", "bd",
		"pdz", "bdz",
		"pdž", "bdž",
		"pz", "bz",
		"pž", "bž",
		"pg", "bg",
		"ph", "bh",
		"pv", "bv",
		"pm", "bm",
		"pn", "bn",
		"pl", "bl",
		"pr", "br",

		"tb", "db",
		"td", "d",
		"tdz", "dz",
		"tdž", "dž",
		"tz", "dz",
		"tž", "dž",
		"tg", "dg",
		"th", "dh",
		"tv", "dv",
		"tm", "dm",
		"tn", "dn",
		"tl", "dl",
		"tr", "dr",
		"tš", "č",

		"čb", "džb",
		"čd", "džd",
		"čdz", "dž",
		"čdž", "dž",
		"čz", "džz",
		"čž", "dž",
		"čg", "džg",
		"čh", "džh",
		"čv", "džv",
		"čm", "džm",
		"čn", "džn",
		"čl", "džl",
		"čr", "džr",

		"sb", "zb",
		"sd", "zd",
		"sdz", "dz",
		"sdž", "dž",
		"sz", "z",
		"sž", "ž",
		"sg", "zg",
		"sh", "zh",
		"sv", "zv",
		"sm", "zm",
		"sn", "zn",
		"sl", "zl",
		"sr", "zr",

		"šb", "žb",
		"šd", "žd",
		"šdz", "dz",
		"šdž", "dž",
		"šz", "ž",
		"šž", "ž",
		"šg", "žg",
		"šh", "žh",
		"šv", "žv",
		"šm", "žm",
		"šn", "žn",
		"šl", "žl",
		"šr", "žr",

		"kb", "gb",
		"kd", "gd",
		"kdz", "gdz",
		"kdž", "gdž",
		"kz", "gz",
		"kž", "gž",
		"kg", "g",
		"kh", "gh",
		"kv", "gv",
		"km", "gm",
		"kn", "gn",
		"kl", "gl",
		"kr", "gr",

		"fb", "vb",
		"fd", "vd",
		"fdz", "vdz",
		"fdž", "vdž",
		"fz", "vz",
		"fž", "vž",
		"fg", "vg",
		"fh", "vh",
		"fv", "v",
		"fm", "vm",
		"fn", "vn",
		"fl", "vl",
		"fr", "vr")
	result = rVoiced.Replace(result)

	rMisc := strings.NewReplacer(
		"rgb", "rb",
		"", "",
		"", "")
	result = rMisc.Replace(result)

	// there could be new combinations after replacements, let's run it again
	result = rUnvoiced.Replace(result)
	result = rVoiced.Replace(result)
	result = rMisc.Replace(result)

	runes := []rune(result)
	length := len(runes)
	isLastVowel := isPronunciationVowel(runes[length-1])
	if !isLastVowel {
		firstIndex := 0
		for i := length - 2; i >= 0; i-- {
			isVowel := isPronunciationVowel(runes[i])
			if isVowel {
				firstIndex = i + 1
				break
			}
		}
		head := string(runes[0:firstIndex])
		tail := string(runes[firstIndex:])
		rTail := strings.NewReplacer(
			"b", "p",
			"d", "t",
			"dž", "č",
			"tš", "č",
			"g", "k",
			"z", "s",
			"ž", "š",
			"jlt", "jt",
		)
		tail = rTail.Replace(tail)
		result = head + rTail.Replace(tail) // yes, the same replace again
	}

	return result
}

// extractRhyme extracts the syllables that are used in rhymes
// it tries to detect and return the last two syllables
// as for vowels and consonants, two situations can occur
// 1: ending with a consontant(s), e. g. heliport (returns iport)
// 2: ending with a vowel(s), e. g. factory (returns tory)
// the syllables are normalized for rhyming, e.g. brad and brat are the same rhyme, returns brat
func extractRhyme(pronunciation string) Rhyme {
	tail := getLastSyllables(strings.ToLower(pronunciation), 4)
	strongRhyme := reduceSoundalikes(tail)
	averageRhyme := getLastSyllables(strongRhyme, 3)
	weakRhyme := getLastSyllables(averageRhyme, 2)

	var rhyme Rhyme
	rhyme.Strong = strongRhyme
	rhyme.Average = averageRhyme
	rhyme.Weak = weakRhyme

	return rhyme
}

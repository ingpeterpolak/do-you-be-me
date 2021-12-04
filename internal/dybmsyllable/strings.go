package dybmsyllable

func isVowel(b byte) bool {
	r := string(b)
	return r == "a" || r == "e" || r == "i" || r == "o" || r == "u" || r == "y" || r == "ï" || r == "î" || r == "é" || r == "ê" || r == "è" || r == "à" || r == "â" || r == "ô"
}

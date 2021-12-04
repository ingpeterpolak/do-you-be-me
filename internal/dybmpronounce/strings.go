package dybmpronounce

func containsParenthesis(ngram string) bool {
	for _, r := range ngram {
		if r == '(' || r == ')' {
			return true
		}
	}
	return false
}

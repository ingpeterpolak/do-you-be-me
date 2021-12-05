package dybmrhyme

const RhymeIds = "ABCDEFGHIJKLMNOPQRSTUVWZ"

func isPronunciationVowel(r rune) bool {
	return r == 'a' || r == 'á' || r == 'ä' || r == 'ö' || r == 'o' || r == 'u' || r == 'i' || r == 'í' || r == 'e'
}

func GetNextRhymeId(lastRhymeId string) string {
	if lastRhymeId == "" {
		return "A"
	}

	result := "Z"
	for i := 0; i < len(RhymeIds)-2; i++ {
		if lastRhymeId == RhymeIds[i:i+1] {
			result = RhymeIds[i+1 : i+2]
		}
	}

	return result
}

func ResetRhymeIds(originalIds []string) []string {
	currentRhymeId := ""
	substitutes := make(map[string]string)
	for _, id := range originalIds {
		_, found := substitutes[id]
		if !found {
			currentRhymeId = GetNextRhymeId(currentRhymeId)
			substitutes[id] = currentRhymeId
		}
	}

	var resetIds []string
	for _, id := range originalIds {
		newId := substitutes[id]
		resetIds = append(resetIds, newId)
	}

	return resetIds
}

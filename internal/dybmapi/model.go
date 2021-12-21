package dybmapi

const bucketName = "dybm_words"
const maxRelatedWordsCount = 20 // there's up to 20 related words for each word

type IndexData struct {
	Message string `json:"message"`
}

type Rhyme struct {
	Id    string
	Rhyme string
}

type PimpedLine struct {
	Number    int    `json:"number"`
	Line      string `json:"line"`
	Syllables int    `json:"syllables"`
	RhymeId   string `json:"rhyme_id"`
}

type PimpedLyrics struct {
	Lines []PimpedLine `json:"lines"`
}

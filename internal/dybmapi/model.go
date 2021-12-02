package dybmapi

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

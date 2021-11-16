package dybmapi

type IndexData struct {
	Message string `json:"message"`
}

type PimpedLine struct {
	Number    int    `json:"number"`
	Line      string `json:"line"`
	Syllables byte   `json:"syllables"`
}

type PimpedLyrics struct {
	Lines []PimpedLine `json:"lines"`
}

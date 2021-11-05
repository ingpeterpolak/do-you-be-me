package dybmapi

type IndexData struct {
	Message string `json:"message"`
}

type ImportData struct {
	UrlsFilename  string   `json:"urls_filename"`
	ProcessedUrls []string `json:"processed_urls"`
}

type PimpedLine struct {
	Number int    `json:"number"`
	Line   string `json:"line"`
}

type PimpedLyrics struct {
	Lines []PimpedLine `json:"lines"`
}

type Ngram struct {
	OriginalText string
	Text         string
	Frequency    int
}

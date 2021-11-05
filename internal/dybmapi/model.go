package dybmapi

type IndexData struct {
	Message string `json:"message"`
}

type ImportData struct {
	UrlsFilename  string   `json:"urls_filename"`
	ProcessedUrls []string `json:"processed_urls"`
}

type CombinedFile struct {
	N           int      `json:"n"`
	Letter      string   `json:"letter"`
	SourceFiles []string `json:"source_files"`
	TargetFile  string   `json:"target_file"`
}

type CombineImportData struct {
	CombinedFiles []CombinedFile `json:"combined_files"`
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

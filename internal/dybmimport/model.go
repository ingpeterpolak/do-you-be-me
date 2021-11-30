package dybmimport

const MaxRelatedWordsPerWord = 20

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

type Ngram struct {
	OriginalText string
	Text         string
	Frequency    int
}

type Rhyme struct {
	Ngram        string `bigquery:"ngram"`
	LastWord     string `bigquery:"last_word"`
	Frequency    int    `bigquery:"frequency"`
	Syllables    int    `bigquery:"syllables"`
	StrongRhyme  string `bigquery:"rhyme_strong"`
	AverageRhyme string `bigquery:"rhyme_average"`
	WeakRhyme    string `bigquery:"rhyme_weak"`
}

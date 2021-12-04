package dybmrhyme

type Rhyme struct {
	Ngram        string `bigquery:"ngram"`
	LastWord     string `bigquery:"last_word"`
	Frequency    int    `bigquery:"frequency"`
	Syllables    int    `bigquery:"syllables"`
	StrongRhyme  string `bigquery:"rhyme_strong"`
	AverageRhyme string `bigquery:"rhyme_average"`
	WeakRhyme    string `bigquery:"rhyme_weak"`
}

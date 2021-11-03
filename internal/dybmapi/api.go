package dybmapi

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var TemplateFolder string

func Setup(templateFolder string) {
	TemplateFolder = templateFolder
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	indexTemplate := TemplateFolder + "index.gohtml"
	t, err := template.ParseFiles(indexTemplate)
	if err != nil {
		log.Fatal(err)
	}

	var indexData IndexData
	indexData.Message = "Welcome, songwriter!"
	err = t.Execute(w, indexData)
	if err != nil {
		log.Fatal(err)
	}
}

func HandlePimp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	lyrics := r.URL.Query().Get("lyrics")
	lines := strings.Split(lyrics, "\n")

	var pimpedLines []PimpedLine
	for i, s := range lines {

		line := removeSpecialChars(s)

		var pimpedLine PimpedLine
		pimpedLine.Number = i + 1
		pimpedLine.Line = line

		pimpedLines = append(pimpedLines, pimpedLine)
	}

	var pimpedLyrics PimpedLyrics
	pimpedLyrics.Lines = pimpedLines

	result, err := json.MarshalIndent(pimpedLyrics, "", "")

	if err != nil {
		log.Fatal("Can't!")
	}

	w.Write(result)
}

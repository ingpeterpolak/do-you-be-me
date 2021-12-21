package dybmapi

import (
	"bufio"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var TemplateFolder string
var AssetsFolder string
var DataFolder string

func Setup(templateFolder, assetsFolder, dataFolder string) {
	TemplateFolder = templateFolder
	AssetsFolder = assetsFolder
	DataFolder = dataFolder
}

func sendFavicon(w http.ResponseWriter) {
	log.Println("Returning favicon.ico")
	faviconIco, err := os.Open(AssetsFolder + "/favicon.ico")
	if err == nil {
		defer faviconIco.Close()
		stats, err := faviconIco.Stat()
		if err == nil {
			var size int64 = stats.Size()
			bytes := make([]byte, size)
			reader := bufio.NewReader(faviconIco)
			_, err = reader.Read(bytes)
			if err == nil {
				w.Header().Add("Content-type", "image/x-icon")
				w.Write(bytes)
			}
		}
	}
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	log.Println("START handling", r.URL)

	if r.URL.EscapedPath() == "/favicon.ico" {
		sendFavicon(w)
		return
	}

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

	log.Println("DONE handling", r.URL)
}

func HandlePimp(w http.ResponseWriter, r *http.Request) {
	log.Println("START handling", r.URL)

	w.Header().Add("Content-type", "application/json")

	lyrics := r.URL.Query().Get("lyrics")
	lines := analyzeLines(strings.Split(lyrics, "\n"))

	var pimpedLyrics PimpedLyrics
	pimpedLyrics.Lines = lines

	result, err := json.MarshalIndent(pimpedLyrics, "", "")

	if err != nil {
		log.Fatal("Error:", err)
	}

	w.Write(result)

	log.Println("DONE handling", r.URL)
}

func HandleLine(w http.ResponseWriter, r *http.Request) {
	log.Println("START handling", r.URL)

	w.Header().Add("Content-type", "application/json")

	lineToRhyme := r.URL.Query().Get("r")
	syllableCount := r.URL.Query().Get("s")
	words := r.URL.Query().Get("w")

	line := getLine(lineToRhyme, syllableCount, words)

	result, err := json.MarshalIndent(line, "", "")

	if err != nil {
		log.Fatal("Error:", err)
	}

	w.Write(result)

	log.Println("DONE handling", r.URL)
}

package dybmapi

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type IndexData struct {
	Message string `json:"message"`
}

type Greeting struct {
	Greeting string `json:"greeting"`
	Subject  string `json:"subject"`
	Suffix   string `json:"suffix"`
}

var templateFolder string

func SetupTemplateFolder(folder string) {
	templateFolder = folder
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	indexTemplate := templateFolder + "index.gohtml"
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

func HandleHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")

	var helloWorld Greeting
	helloWorld.Greeting = "Hello"
	helloWorld.Subject = "World"
	helloWorld.Suffix = " :)"

	result, err := json.MarshalIndent(helloWorld, "", "    ")

	if err != nil {
		log.Fatal("Can't!")
	}

	w.Write(result)
}

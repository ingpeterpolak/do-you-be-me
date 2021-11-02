package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Greeting struct {
	Greeting string `json:"greeting"`
	Subject  string `json:"subject"`
	Suffix   string `json:"suffix"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, "request from", r.Host)

	w.Header().Add("Content-type", "application/json")

	var helloWorld Greeting
	helloWorld.Greeting = "Hello"
	helloWorld.Subject = "World"
	helloWorld.Suffix = "!"

	result, err := json.MarshalIndent(helloWorld, "", "    ")

	if err != nil {
		log.Fatal("Can't!")
	}

	w.Write(result)
}

func main() {
	log.Println("Starting do/you API")

	http.HandleFunc("/", handleRequest)

	log.Println("Starting web server on port 8080")
	http.ListenAndServe(":8080", nil)
}

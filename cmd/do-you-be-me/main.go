package main

import (
	"github.com/ingpeterpolak/do-you-be-me/internal/dybmapi"

	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Starting do/you API")

	dybmapi.SetupTemplateFolder("../../web/template/")

	http.HandleFunc("/", dybmapi.HandleRoot)
	http.HandleFunc("/hello", dybmapi.HandleHello)

	// PORT environment variable is provided by Cloud Run.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Starting web server on port", port)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatal(err)
	}
}

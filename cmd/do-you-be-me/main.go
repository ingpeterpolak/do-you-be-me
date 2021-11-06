package main

import (
	"github.com/ingpeterpolak/do-you-be-me/internal/dybmapi"
	"github.com/ingpeterpolak/do-you-be-me/internal/dybmimport"

	"log"
	"net/http"
	"os"
)

var localDebug bool

// mapHandlers maps all the URLs to the correct handlers
func mapHandlers() {
	assetsFolder := "./assets"
	// look for the assets elsewhere if we're debugging locally on Windows
	if localDebug {
		assetsFolder = "../../web/assets"
	}

	fs := http.FileServer(http.Dir(assetsFolder))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/", dybmapi.HandleRoot)
	http.HandleFunc("/pimp", dybmapi.HandlePimp)

	http.HandleFunc("/import", dybmimport.HandleImport)
	http.HandleFunc("/combine-import", dybmimport.HandleCombineImport)

	http.HandleFunc("/process", dybmimport.HandleProcess)
}

// setupApi sets the correct template folder for both local debugging and production run
func setupApi() {
	templateFolder := "./"

	// look for the templates elsewhere if we're debugging locally on Windows
	if localDebug {
		templateFolder = "../../web/template/"
	}

	dybmapi.Setup(templateFolder)
}

// setupApi sets the correct template folder for both local debugging and production run
func setupImport() {
	dataFolder := "./data/"

	// look for the templates elsewhere if we're debugging locally on Windows
	if localDebug {
		dataFolder = "../../internal/dybmimport/data/"
	}

	dybmimport.Setup(dataFolder)
}

// startServer checks the correct port and starts the http server
func startServer() error {
	// PORT environment variable is provided by Cloud Run.
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Println("Starting web server on port", port)
	return http.ListenAndServe(":"+port, nil)
}

func main() {
	log.Println("Starting do/you be/me API")
	localDebug = false

	currentOs := os.Getenv("OS")
	if currentOs == "Windows_NT" {
		log.Println("Debugging locally")
		localDebug = true
	}

	setupApi()
	setupImport()
	mapHandlers()

	err := startServer()

	if err != nil {
		log.Fatal(err)
	}
}

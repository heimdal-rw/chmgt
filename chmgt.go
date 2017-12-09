package main

import (
	"log"
	"net/http"

	"github.com/mattjw79/chmgt/routing"
)

func main() {
	router := routing.NewRouter()
	// Let the user know that we're starting and on what
	// port we're listening
	log.Println("Starting server")
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

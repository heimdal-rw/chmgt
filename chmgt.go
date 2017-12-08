package main

import (
	"log"
	"net/http"

	"github.com/mattjw79/chmgt/routing"
)

func main() {
	router := routing.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

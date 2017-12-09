package main

import (
	"log"
	"net/http"
	"time"

	"github.com/mattjw79/chmgt/routing"
)

func main() {
	router := routing.NewRouter()
	// Let the user know that we're starting
	log.Println("Starting server")
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	// Let the user know where the server is running
	log.Printf("Listening on %v\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

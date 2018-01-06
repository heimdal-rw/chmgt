package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"chmgt/routing"
)

func main() {
	// Pull in config
	config := ReadConfig()
	log.Printf("config:\n%+v\n", config)

	// Let the user know tt we're starting
	log.Println("Starting server")
	router := routing.NewRouter()
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%s", config.Interface, config.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	// Let the user know where the server is running
	log.Printf("Listening on %v\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

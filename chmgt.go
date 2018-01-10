package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"chmgt/models"
	"chmgt/routing"
)

func main() {
	// Grabbing command line flags
	var (
		configFileFlag string //config file to use
	)
	flag.StringVar(&configFileFlag, "config", "./config", "Config file path to be used.")
	flag.Parse()

	// Pull in config
	config := ReadConfig(configFileFlag)
	log.Printf("config:\n%+v\n", config)

	// Create the database if it doesn't exist
	if err := models.Exists(models.DBConnection); err != nil {
		log.Printf("Creating database: %v", models.DBConnection)
		models.GenerateDatabase("./sql/sqlite.sql", models.DBConnection)
	}

	// Let the user know tt we're starting
	log.Println("Starting server")
	router := routing.NewRouter()
	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprint(config.ServerListen),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	// Let the user know where the server is running
	log.Printf("Listening on %v\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

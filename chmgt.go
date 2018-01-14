package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"chmgt/handling"
	"chmgt/models"
)

func main() {
	// Grabbing command line flags
	var (
		configFileFlag string //config file to use
	)
	flag.StringVar(&configFileFlag, "config", "", "Config file path to be used.")
	flag.Parse()

	// Pull in config
	config, err := ReadConfig(configFileFlag)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("config:\n%+v\n", config)

	// Create the database if it doesn't exist
	if config.Database != "" {
		// Overwrite default if config specifies db file
		models.DSN = config.Database
	}
	if err := models.Exists(models.DSN); err != nil {
		log.Printf("Creating database: %v", models.DSN)
		models.GenerateDatabase("./sql/sqlite.sql", models.DSN)
	}

	// Let the user know that we're starting
	log.Println("Starting server")
	router := handling.NewRouter()
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

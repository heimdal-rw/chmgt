package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"chmgt/routing"

	_ "github.com/mattn/go-sqlite3"
)

func createDatabase(dbFile string) error {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	// Read in the SQL for creating the database
	buf, err := ioutil.ReadFile("./sql/sqlite.sql")
	if err != nil {
		return err
	}
	sqlQuery := string(buf)

	// Create the schema in the database
	_, err = db.Exec(sqlQuery)
	if err != nil {
		return err
	}

	return nil
}

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
	dbFile := "./chmgt.db"
	if _, err := os.Stat(dbFile); err != nil {
		if os.IsNotExist(err) {
			log.Printf("Creating database: %v", dbFile)
			createDatabase(dbFile)
		}
	}

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

package models

import (
	"database/sql"
	"io/ioutil"
	"os"

	// Bring in the SQLite3 functionality
	_ "github.com/mattn/go-sqlite3"
)

// DBConnection provides a default connection string
var DBConnection = "./chmgt.db"

// Open opens a connection to the database
func Open(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Exists checks for database existance
func Exists(dbFile string) error {
	if _, err := os.Stat(dbFile); err != nil {
		return err
	}
	return nil
}

// GenerateDatabase reads in a sql file to create the database
func GenerateDatabase(sqlFile, dbFile string) error {
	db, err := Open(DBConnection)
	if err != nil {
		return err
	}
	defer db.Close()

	// Read in the SQL for creating the database
	buf, err := ioutil.ReadFile(sqlFile)
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

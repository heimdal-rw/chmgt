package models

import (
	"database/sql"
	"io/ioutil"
	"os"

	// Bring in the SQLite3 functionality
	_ "github.com/mattn/go-sqlite3"
)

type ChangeRequest struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	AuthorID    int    `json:"authorId"`
	RequesterID int    `json:"requesterId"`
	Description string `json:"description"`
	Reason      string `json:"reason"`
	Risk        string `json:"risk"`
	Steps       string `json:"steps"`
	Revert      string `json:"revert"`
}

var DBConnection string = "./chmgt.db"

func Open(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Exists(dbFile string) error {
	if _, err := os.Stat(dbFile); err != nil {
		return err
	}
	return nil
}

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

func GetChangeRequests(db *sql.DB) ([]ChangeRequest, error) {
	sqlQuery := `
	SELECT
		_rowid_,
		title,
		authorId,
		requesterId,
		description,
		reason,
		risk,
		steps,
		revert
	FROM changeRequest
	`
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	changeRequests := []ChangeRequest{}
	for rows.Next() {
		var cr ChangeRequest
		err = rows.Scan(
			&cr.ID,
			&cr.Title,
			&cr.AuthorID,
			&cr.RequesterID,
			&cr.Description,
			&cr.Reason,
			&cr.Risk,
			&cr.Steps,
			&cr.Revert,
		)
		if err != nil {
			return nil, err
		}
		changeRequests = append(changeRequests, cr)
	}

	return changeRequests, nil
}

func (cr *ChangeRequest) GetChange(db *sql.DB) error {
	sqlQuery := `
	SELECT
		title,
		authorId,
		requesterId,
		description,
		reason,
		risk,
		steps,
		revert
	FROM changeRequest
	WHERE _rowid_=?
	`
	rows, err := db.Query(sqlQuery, cr.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&cr.Title,
			&cr.AuthorID,
			&cr.RequesterID,
			&cr.Description,
			&cr.Reason,
			&cr.Risk,
			&cr.Steps,
			&cr.Revert,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cr *ChangeRequest) CreateChange(db *sql.DB) error {
	sqlQuery := `
	INSERT INTO changeRequest (
		title,
		authorId,
		requesterId,
		description,
		reason,
		risk,
		steps,
		revert
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := db.Exec(
		sqlQuery,
		cr.Title,
		cr.AuthorID,
		cr.RequesterID,
		cr.Description,
		cr.Reason,
		cr.Risk,
		cr.Steps,
		cr.Revert,
	)
	if err != nil {
		return err
	}
	insertedID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	cr.ID = int(insertedID)
	return nil
}

func (cr *ChangeRequest) DeleteChange(db *sql.DB) error {
	sqlQuery := "DELETE FROM changeRequest WHERE _rowid_=?"

	_, err := db.Exec(sqlQuery, cr.ID)
	if err != nil {
		return err
	}

	return nil
}

func (cr *ChangeRequest) UpdateChange(db *sql.DB) error {
	sqlQuery := `
	UPDATE changeRequest SET
		title=?,
		authorId=?,
		requesterId=?,
		description=?,
		reason=?,
		risk=?,
		steps=?,
		revert=?
	WHERE _rowid_=?
	`

	_, err := db.Exec(
		sqlQuery,
		cr.Title,
		cr.AuthorID,
		cr.RequesterID,
		cr.Description,
		cr.Reason,
		cr.Risk,
		cr.Steps,
		cr.Revert,
		cr.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

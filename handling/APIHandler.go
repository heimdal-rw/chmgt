package handling

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	// Bring in the SQLite3 functionality
	_ "github.com/mattn/go-sqlite3"
)

type changeRequest struct {
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

type changeRequests []changeRequest

type apiOption struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

type apiOptions []apiOption

// APIHandler takes care of the index page
func APIHandler(w http.ResponseWriter, r *http.Request) {
	jsonOut := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// Sample structure of how the API may work
	jsonOut.Encode(apiOptions{
		apiOption{"GET", "/api", "Get help on the API"},
		apiOption{"GET", "/api/changes", "Display all change requests"},
		apiOption{"POST", "/api/changes", "Create a new change request"},
		apiOption{"GET", "/api/changes/{id}", "Get the change request specified by ID"},
		apiOption{"PUT", "/api/changes/{id}", "Update the change request specified by ID"},
		apiOption{"DELETE", "/api/changes/{id}", "Delete the change request specified by ID"},
		apiOption{"GET", "/api/users", "Display all users"},
		apiOption{"POST", "/api/users", "Create a new user"},
	})
}

// GetChangesHandler returns change requests
func GetChangesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	dbFile := "./chmgt.db"
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	sqlQuery := "SELECT _rowid_, title, authorId, requesterId, description, reason, risk, steps, revert FROM changeRequest"
	if vars["id"] != "" {
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			log.Println(err)
			return
		}
		sqlQuery = fmt.Sprintf("%s WHERE _rowid_=%d", sqlQuery, id)
	}

	rows, err := db.Query(sqlQuery)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	jsonOut := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var crs changeRequests
	for rows.Next() {
		var cr changeRequest
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
			log.Println(err)
		} else {
			crs = append(crs, cr)
		}
	}
	jsonOut.Encode(crs)
}

// CreateChangeHandler creates a new change in the database
func CreateChangeHandler(w http.ResponseWriter, r *http.Request) {
	var cr changeRequest

	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil {
		log.Println(err)
		return
	}

	dbFile := "./chmgt.db"
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

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
		log.Println(err)
		return
	}

	jsonOut := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	insertedID, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return
	}
	jsonOut.Encode(insertedID)
}

// DeleteChangeHandler deletes the specified change request
func DeleteChangeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	dbFile := "./chmgt.db"
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	sqlQuery := "DELETE FROM changeRequest WHERE _rowid_=?"
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		return
	}
	_, err = db.Exec(sqlQuery, id)
	if err != nil {
		log.Println(err)
		return
	}
}

// UpdateChangeHandler updates the specified change request
func UpdateChangeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var cr changeRequest

	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil {
		log.Println(err)
		return
	}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		return
	}
	cr.ID = id

	dbFile := "./chmgt.db"
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

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

	_, err = db.Exec(
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
		log.Println(err)
		return
	}
}

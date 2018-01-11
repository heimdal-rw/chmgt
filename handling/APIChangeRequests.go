package handling

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"chmgt/models"

	"github.com/gorilla/mux"
)

// GetChangesHandler returns change requests
func GetChangesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	db, err := models.Open(models.DBConnection)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	var crs []models.ChangeRequest
	// sqlQuery := "SELECT _rowid_, title, authorId, requesterId, description, reason, risk, steps, revert FROM changeRequest"
	if vars["id"] != "" {
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			log.Println(err)
			return
		}
		// sqlQuery = fmt.Sprintf("%s WHERE _rowid_=%d", sqlQuery, id)
		var cr models.ChangeRequest
		cr.ID = id
		err = cr.GetChange(db)
		if err != nil {
			log.Println(err)
			return
		}
		crs = append(crs, cr)
	} else {
		crs, err = models.GetChangeRequests(db)
		if err != nil {
			log.Println(err)
			return
		}
	}

	jsonOut := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonOut.Encode(crs)
}

// CreateChangeHandler creates a new change in the database
func CreateChangeHandler(w http.ResponseWriter, r *http.Request) {
	var cr models.ChangeRequest

	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil {
		log.Println(err)
		return
	}

	db, err := models.Open(models.DBConnection)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	err = cr.CreateChange(db)
	if err != nil {
		log.Println(err)
		return
	}

	jsonOut := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonOut.Encode(cr.ID)
}

// DeleteChangeHandler deletes the specified change request
func DeleteChangeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	db, err := models.Open(models.DBConnection)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		return
	}

	var cr models.ChangeRequest
	cr.ID = id
	err = cr.DeleteChange(db)
	if err != nil {
		log.Println(err)
		return
	}
}

// UpdateChangeHandler updates the specified change request
func UpdateChangeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var cr models.ChangeRequest

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

	db, err := models.Open(models.DBConnection)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	err = cr.UpdateChange(db)
	if err != nil {
		log.Println(err)
		return
	}
}

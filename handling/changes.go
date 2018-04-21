package handling

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/heimdal-rw/chmgt/models"
	"gopkg.in/mgo.v2/bson"
)

// GetChanges returns a json response with the requested changes
func (h *Handler) GetChanges(w http.ResponseWriter, r *http.Request) {
	var (
		changes models.Changes
		err     error
	)
	vars := mux.Vars(r)

	changes, err = h.Datasource.GetChanges(vars["id"])
	if err != nil {
		switch err {
		case models.ErrNotFound:
			APIWriteFailure(w, "no changes found", http.StatusNotFound)
		case models.ErrObjID:
			APIWriteFailure(w, "invalid object id", http.StatusBadRequest)
		default:
			APIWriteFailure(w, "", http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	APIWriteSuccess(w, changes)
}

// InsertChange inserts a new change into the database
func (h *Handler) InsertChange(w http.ResponseWriter, r *http.Request) {
	var change models.Change

	if err := json.NewDecoder(r.Body).Decode(&change); err != nil {
		APIWriteFailure(w, "error parsing json", http.StatusBadRequest)
		return
	}
	id, err := h.Datasource.InsertChange(change)
	if err != nil {
		APIWriteFailure(w, "error inserting change", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	APIWriteSuccess(w, id)
}

// UpdateChange updates a change database record
func (h *Handler) UpdateChange(w http.ResponseWriter, r *http.Request) {
	var change models.Change

	vars := mux.Vars(r)
	if err := json.NewDecoder(r.Body).Decode(&change); err != nil {
		APIWriteFailure(w, "error parsing json", http.StatusBadRequest)
		return
	}
	change.ID = bson.ObjectIdHex(vars["id"])
	if err := h.Datasource.UpdateChange(change); err != nil {
		APIWriteFailure(w, "error updating change", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	APIWriteSuccess(w, nil)
}

// RemoveChange removes a change from the database
func (h *Handler) RemoveChange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if err := h.Datasource.RemoveChange(vars["id"]); err != nil {
		if err == models.ErrNotFound {
			APIWriteFailure(w, "change not found", http.StatusNotFound)
		} else {
			APIWriteFailure(w, "error removing change", http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	APIWriteSuccess(w, nil)
}

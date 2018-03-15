package handling

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/heimdal-rw/chmgt/models"

	"github.com/gorilla/mux"
)

// GetChangeRequestsHandler returns change requests
func (h *Handler) GetChangeRequestsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var err error
	crs, err := h.Datasource.GetChangeRequests(vars["id"])
	if err != nil {
		if err.Error() == "invalid object id" {
			APIWriteFailure(w, "no change request found", http.StatusNotFound)
			return
		}
		if vars["id"] == "" || err == models.ErrNoRows {
			APIWriteFailure(w, "no change request found", http.StatusNotFound)
			return
		}
		APIWriteFailure(w, "", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	APIWriteSuccess(w, crs)
}

// CreateChangeRequestHandler creates a new change request in the database
func (h *Handler) CreateChangeRequestHandler(w http.ResponseWriter, r *http.Request) {
	var cr models.Item

	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil {
		APIWriteFailure(w, "error parsing json", http.StatusBadRequest)
		return
	}

	err = h.Datasource.InsertChangeRequest(cr)
	if err != nil {
		APIWriteFailure(w, "error inserting change request", http.StatusInternalServerError)
		return
	}

	APIWriteSuccess(w, cr["_id"])
}

// DeleteChangeRequestHandler deletes the specified change request
func (h *Handler) DeleteChangeRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["id"] == "" {
		APIWriteFailure(w, "no id specified", http.StatusBadRequest)
		return
	}

	crs, err := h.Datasource.GetChangeRequests(vars["id"])
	if err != nil {
		if err == models.ErrNoRows {
			APIWriteFailure(w, "specified record was not found", http.StatusNotFound)
			return
		}
		APIWriteFailure(w, "error finding change request", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = h.Datasource.RemoveChangeRequest(crs[0])
	if err != nil {
		APIWriteFailure(w, "error deleting change request", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	APIWriteSuccess(w, vars["id"])
}

// UpdateChangeRequestHandler updates the specified change request
func (h *Handler) UpdateChangeRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var cr models.Item

	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil {
		APIWriteFailure(w, "error parsing json", http.StatusBadRequest)
		return
	}
	cr.SetID(vars["id"])

	err = h.Datasource.UpdateChangeRequest(cr)
	if err != nil {
		APIWriteFailure(w, "error updating change request", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	APIWriteSuccess(w, vars["id"])
}

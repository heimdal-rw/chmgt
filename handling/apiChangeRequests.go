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
		if vars["id"] == "" || err == models.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(APIResponseJSON{
				true,
				"no change request found",
				nil,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponseJSON{
			true,
			"unexpected error occurred",
			nil,
		})
		log.Println(err)
		return
	}

	if crs == nil {
		crs = make([]models.Item, 0)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(APIResponseJSON{
		false,
		"success",
		crs,
	})
}

// CreateChangeRequestHandler creates a new change request in the database
func (h *Handler) CreateChangeRequestHandler(w http.ResponseWriter, r *http.Request) {
	var cr models.Item

	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil {
		log.Println(err)
		return
	}

	err = h.Datasource.InsertChangeRequest(cr)
	if err != nil {
		log.Println(err)
		return
	}

	jsonOut := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonOut.Encode(cr["_id"])
}

// DeleteChangeRequestHandler deletes the specified change request
func (h *Handler) DeleteChangeRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["id"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	crs, err := h.Datasource.GetChangeRequests(vars["id"])
	if err != nil {
		if err == models.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println(err)
		return
	}

	err = h.Datasource.RemoveChangeRequest(crs[0])
	if err != nil {
		log.Println(err)
		return
	}
}

// UpdateChangeRequestHandler updates the specified change request
func (h *Handler) UpdateChangeRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var cr models.Item

	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil {
		log.Println(err)
		return
	}
	cr.SetID(vars["id"])

	err = h.Datasource.UpdateChangeRequest(cr)
	if err != nil {
		log.Println(err)
		return
	}
}

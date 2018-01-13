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
func (env *Env) GetChangesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	crs := make([]*models.ChangeRequest, 0)
	if vars["id"] != "" {
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			log.Println(err)
			return
		}
		cr := new(models.ChangeRequest)
		cr, err = env.DB.GetChangeRequest(id)
		if err != nil {
			log.Println(err)
			return
		}
		crs = append(crs, cr)
	} else {
		var err error
		crs, err = env.DB.GetChangeRequests()
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
func (env *Env) CreateChangeHandler(w http.ResponseWriter, r *http.Request) {
	cr := new(models.ChangeRequest)

	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil {
		log.Println(err)
		return
	}

	err = env.DB.CreateChangeRequest(cr)
	if err != nil {
		log.Println(err)
		return
	}

	jsonOut := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonOut.Encode(cr.ID)
}

// DeleteChangeHandler deletes the specified change request
func (env *Env) DeleteChangeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		return
	}

	err = env.DB.DeleteChangeRequest(id)
	if err != nil {
		log.Println(err)
		return
	}
}

// UpdateChangeHandler updates the specified change request
func (env *Env) UpdateChangeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cr := new(models.ChangeRequest)

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

	err = env.DB.UpdateChangeRequest(cr)
	if err != nil {
		log.Println(err)
		return
	}
}

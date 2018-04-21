package handling

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	"github.com/heimdal-rw/chmgt/models"
)

// GetUsers returns a json response with the requested users
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	var (
		users models.Users
		err   error
	)
	vars := mux.Vars(r)

	users, err = h.Datasource.GetUsers(vars["id"])
	if err != nil {
		switch err {
		case models.ErrNotFound:
			APIWriteFailure(w, "no users found", http.StatusNotFound)
		case models.ErrObjID:
			APIWriteFailure(w, "invalid object id", http.StatusBadRequest)
		default:
			APIWriteFailure(w, "", http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	// mask the passwords so that they're not sent in the response
	for idx := range users {
		users[idx].Password = "********"
	}

	APIWriteSuccess(w, users)
}

// InsertUser inserts a new user into the database
func (h *Handler) InsertUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		APIWriteFailure(w, "error parsing json", http.StatusBadRequest)
		return
	}
	id, err := h.Datasource.InsertUser(user)
	if err != nil {
		APIWriteFailure(w, "error inserting user", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	APIWriteSuccess(w, id)
}

// UpdateUser updates a user's database record
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	vars := mux.Vars(r)
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		APIWriteFailure(w, "error parsing json", http.StatusBadRequest)
		return
	}
	user.ID = bson.ObjectIdHex(vars["id"])
	if err := h.Datasource.UpdateUser(user); err != nil {
		APIWriteFailure(w, "error inserting user", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	APIWriteSuccess(w, nil)
}

// RemoveUser removes a user from the database
func (h *Handler) RemoveUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if err := h.Datasource.RemoveUser(vars["id"]); err != nil {
		if err == models.ErrNotFound {
			APIWriteFailure(w, "change not found", http.StatusNotFound)
		} else {
			APIWriteFailure(w, "error removing user", http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	APIWriteSuccess(w, nil)
}

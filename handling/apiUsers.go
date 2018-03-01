package handling

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/heimdal-rw/chmgt/models"

	"github.com/gorilla/mux"
)

// GetUsersHandler returns users
func (h *Handler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var err error
	users, err := h.Datasource.GetUsers(vars["id"])
	if err != nil {
		if vars["id"] == "" || err == models.ErrNoRows {
			APIWriteFailure(w, "no user found", http.StatusNotFound)
			return
		}
		APIWriteFailure(w, "", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	APIWriteSuccess(w, users)
}

// CreateUserHandler creates a new user in the database
func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.Item

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		APIWriteFailure(w, "error parsing json", http.StatusBadRequest)
		return
	}

	if user["username"] == nil {
		APIWriteFailure(w, "username not specified", http.StatusBadRequest)
		return
	}

	err = h.Datasource.InsertUser(user)
	if err != nil {
		if strings.HasPrefix(err.Error(), "E11000") {
			APIWriteFailure(w, "duplicate username", http.StatusBadRequest)
			return
		}
		APIWriteFailure(w, "error inserting change request", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	APIWriteSuccess(w, user["_id"])
}

// DeleteUserHandler deletes the specified user
func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["id"] == "" {
		APIWriteFailure(w, "no id specified", http.StatusBadRequest)
		return
	}

	users, err := h.Datasource.GetUsers(vars["id"])
	if err != nil {
		if err == models.ErrNoRows {
			APIWriteFailure(w, "specified record was not found", http.StatusNotFound)
			return
		}
		APIWriteFailure(w, "error finding user", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = h.Datasource.RemoveUser(users[0])
	if err != nil {
		APIWriteFailure(w, "error deleting user", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	APIWriteSuccess(w, vars["id"])
}

// UpdateUserHandler updates the specified user
func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var user models.Item

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		APIWriteFailure(w, "error parsing json", http.StatusBadRequest)
		return
	}
	user.SetID(vars["id"])

	// Username wasn't included in the update, but we MUST
	// have a username, so get it from the DB
	if user["username"] == nil {
		u, err := h.Datasource.GetUsers(vars["id"])
		if err != nil {
			if err == models.ErrNoRows {
				APIWriteFailure(w, "error finding user", http.StatusNotFound)
				return
			}
			APIWriteFailure(w, "", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		user["username"] = u[0]["username"]
	}

	err = h.Datasource.UpdateUser(user)
	if err != nil {
		APIWriteFailure(w, "error updating change request", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	APIWriteSuccess(w, vars["id"])
}

package handling

import (
	"encoding/json"
	"log"
	"net/http"

	"chmgt/models"

	"github.com/gorilla/mux"
)

// GetUsersHandler returns users
func (h *Handler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var err error
	users, err := h.Datasource.GetUsers(vars["id"])
	if err != nil {
		if vars["id"] == "" && err == models.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println(err)
		return
	}

	jsonOut := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonOut.Encode(users)
}

// CreateUserHandler creates a new user in the database
func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Println(err)
		return
	}

	err = h.Datasource.InsertUser(user)
	if err != nil {
		log.Println(err)
		return
	}

	jsonOut := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonOut.Encode(user.ID)
}

// DeleteUserHandler deletes the specified change request
func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["id"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := h.Datasource.GetUsers(vars["id"])
	if err != nil {
		if err == models.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println(err)
		return
	}

	err = h.Datasource.RemoveUser(&users[0])
	if err != nil {
		log.Println(err)
		return
	}
}

// UpdateUserHandler updates the specified change request
func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := new(models.User)

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Println(err)
		return
	}
	user.SetID(vars["id"])

	err = h.Datasource.UpdateUser(user)
	if err != nil {
		log.Println(err)
		return
	}
}

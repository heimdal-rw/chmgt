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
	var user models.Item

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("create user json decode error:", err)
		return
	}

	if user["username"] == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("username not specified"))
		return
	}

	err = h.Datasource.InsertUser(user)
	if err != nil {
		if strings.HasPrefix(err.Error(), "E11000") {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("duplicate username"))
		}
		log.Println("create user mongo insert error:", err)
		return
	}

	jsonOut := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonOut.Encode(user["_id"])
}

// DeleteUserHandler deletes the specified user
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

	err = h.Datasource.RemoveUser(users[0])
	if err != nil {
		log.Println(err)
		return
	}
}

// UpdateUserHandler updates the specified user
func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var user models.Item

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		return
	}
	user.SetID(vars["id"])

	// Username wasn't included in the update, but we MUST
	// have a username, so get it from the DB
	if user["username"] == nil {
		u, err := h.Datasource.GetUsers(vars["id"])
		if err != nil {
			if err == models.ErrNoRows {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Println(err)
			return
		}
		user["username"] = u[0]["username"]
	}

	err = h.Datasource.UpdateUser(user)
	if err != nil {
		log.Println(err)
		return
	}
}

package handling

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"chmgt/models"

	"github.com/gorilla/mux"
)

// GetUsersHandler returns users
func (env *Env) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	users := make([]*models.User, 0)
	if vars["id"] != "" {
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			log.Println(err)
			return
		}

		user := new(models.User)
		user.ID = id
		err = env.DB.GetUser(id)
		if err != nil {
			log.Println(err)
			return
		}
		users = append(users, user)
	} else {
		var err error
		users, err = env.DB.GetUsers()
		if err != nil {
			log.Println(err)
			return
		}
	}

	jsonOut := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonOut.Encode(users)
}

// CreateUserHandler creates a new user in the database
func (env *Env) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		return
	}

	err = env.DB.CreateUser(user)
	if err != nil {
		log.Println(err)
		return
	}

	jsonOut := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonOut.Encode(user.ID)
}

// DeleteUserHandler deletes the specified change request
func (env *Env) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		return
	}

	err = env.DB.DeleteUser(id)
	if err != nil {
		log.Println(err)
		return
	}
}

// UpdateUserHandler updates the specified change request
func (env *Env) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := new(models.User)

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		return
	}

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		return
	}
	user.ID = id

	err = env.DB.UpdateUser(user)
	if err != nil {
		log.Println(err)
		return
	}
}

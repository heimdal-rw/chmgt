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
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	db, err := models.Open(models.DBConnection)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	var users []models.User
	if vars["id"] != "" {
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			log.Println(err)
			return
		}

		var user models.User
		user.ID = id
		err = user.GetUser(db)
		if err != nil {
			log.Println(err)
			return
		}
		users = append(users, user)
	} else {
		users, err = models.GetUsers(db)
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
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
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

	err = user.CreateUser(db)
	if err != nil {
		log.Println(err)
		return
	}

	jsonOut := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	jsonOut.Encode(user.ID)
}

// DeleteUserHandler deletes the specified change request
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
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

	var user models.User
	user.ID = id
	err = user.DeleteUser(db)
	if err != nil {
		log.Println(err)
		return
	}
}

// UpdateUserHandler updates the specified change request
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var user models.User

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

	db, err := models.Open(models.DBConnection)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	err = user.UpdateUser(db)
	if err != nil {
		log.Println(err)
		return
	}
}

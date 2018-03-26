package handling

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/heimdal-rw/chmgt/models"
)

// APIResponse creates the structure of an API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// WriteJSON format a JSON response and write it back to the client
func (j APIResponse) WriteJSON(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(j)
}

// APIWriteSuccess builds a success response
func APIWriteSuccess(w http.ResponseWriter, data interface{}) {
	APIResponse{
		true,
		"success",
		data,
	}.WriteJSON(w, http.StatusOK)
}

// APIWriteFailure builds a failure response
func APIWriteFailure(w http.ResponseWriter, msg string, status int) {
	if msg == "" {
		msg = "unknown error occured"
	}

	APIResponse{
		false,
		msg,
		nil,
	}.WriteJSON(w, status)
}

func getVars(r *http.Request) (map[string]string, error) {
	vars := mux.Vars(r)

	// Uppercase the first character
	vars["collection"] = strings.Title(vars["collection"])
	for _, vc := range models.ValidCollections {
		if vars["collection"] == vc {
			return vars, nil
		}
	}

	return nil, models.ErrInvalidCollection
}

// GetItemsHandler returns all items in a collection
func (h *Handler) GetItemsHandler(w http.ResponseWriter, r *http.Request) {
	vars, err := getVars(r)
	if err != nil {
		APIWriteFailure(w, "invalid collection specified", http.StatusBadRequest)
		return
	}

	items, err := h.Datasource.GetItems(vars["id"], vars["collection"])
	if err != nil {
		if err != nil {
			switch err {
			case models.ErrNoRows:
				APIWriteFailure(w, "no item found", http.StatusNotFound)
			case models.ErrObjID:
				APIWriteFailure(w, "invalid object id", http.StatusBadRequest)
			default:
				APIWriteFailure(w, "", http.StatusInternalServerError)
				log.Println(err)
			}
			return
		}
	}

	APIWriteSuccess(w, items)
}

// CreateItemHandler creates a new item in the specified collection
func (h *Handler) CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	vars, err := getVars(r)
	if err != nil {
		APIWriteFailure(w, "invalid collection specified", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		APIWriteFailure(w, "error parsing json", http.StatusBadRequest)
		return
	}

	if strings.EqualFold("Users", vars["collection"]) {
		if item["username"] == nil {
			APIWriteFailure(w, "username not specified", http.StatusBadRequest)
			return
		}
	}

	err = h.Datasource.InsertItem(item, vars["collection"])
	if err != nil {
		if strings.HasPrefix(err.Error(), "E11000") {
			APIWriteFailure(w, "duplicate item", http.StatusBadRequest)
			return
		}
		APIWriteFailure(w, "error inserting item", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	APIWriteSuccess(w, item["_id"])
}

// DeleteItemHandler deletes the specified item from the specified collection
func (h *Handler) DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	vars, err := getVars(r)
	if err != nil {
		APIWriteFailure(w, "invalid collection specified", http.StatusBadRequest)
		return
	}

	if vars["id"] == "" {
		APIWriteFailure(w, "no id specified", http.StatusBadRequest)
		return
	}

	items, err := h.Datasource.GetItems(vars["id"], vars["collection"])
	if err != nil {
		if err == models.ErrNoRows {
			APIWriteFailure(w, "specified record was not found", http.StatusNotFound)
			return
		}
		APIWriteFailure(w, "error finding item", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	err = h.Datasource.RemoveItem(items[0], vars["collection"])
	if err != nil {
		APIWriteFailure(w, "error deleting item", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	APIWriteSuccess(w, vars["id"])
}

// UpdateItemHandler updates the specified user
func (h *Handler) UpdateItemHandler(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	vars, err := getVars(r)
	if err != nil {
		APIWriteFailure(w, "invalid path specified", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		APIWriteFailure(w, "error parsing json", http.StatusBadRequest)
		return
	}
	err = item.SetID(vars["id"])
	if err != nil {
		APIWriteFailure(w, "invalid object id", http.StatusBadRequest)
		return
	}

	if strings.EqualFold("Users", vars["collection"]) {
		// Username wasn't included in the update, but we MUST
		// have a username, so get it from the DB
		if item["username"] == nil {
			i, err := h.Datasource.GetItems(vars["id"], vars["collection"])
			if err != nil {
				if err == models.ErrNoRows {
					APIWriteFailure(w, "error finding user", http.StatusNotFound)
					return
				}
				APIWriteFailure(w, "", http.StatusInternalServerError)
				log.Println(err)
				return
			}
			item["username"] = i[0]["username"]
		}
	}

	err = h.Datasource.UpdateItem(item, vars["collection"])
	if err != nil {
		APIWriteFailure(w, "error updating item", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	APIWriteSuccess(w, vars["id"])
}

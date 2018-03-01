package handling

import (
	"encoding/json"
	"net/http"
)

// APIResponseJSON creates the structure of an API response
type APIResponseJSON struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (j APIResponseJSON) WriteJSON(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(j)
}

func APIWriteSuccess(w http.ResponseWriter, data interface{}) {
	APIResponseJSON{
		true,
		"success",
		data,
	}.WriteJSON(w, http.StatusOK)
}

func APIWriteFailure(w http.ResponseWriter, msg string, status int) {
	if msg == "" {
		msg = "unknown error occured"
	}

	APIResponseJSON{
		false,
		msg,
		nil,
	}.WriteJSON(w, status)
}

package handling

import (
	"encoding/json"
	"net/http"
)

type apiOption struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

type apiOptions []apiOption

// APIHandler takes care of the index page
func APIHandler(w http.ResponseWriter, r *http.Request) {
	jsonOut := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// Sample structure of how the API may work
	jsonOut.Encode(apiOptions{
		apiOption{"GET", "/api", "Get help on the API"},
		apiOption{"GET", "/api/changes", "Display all change requests"},
		apiOption{"POST", "/api/changes", "Create a new change request"},
		apiOption{"GET", "/api/changes/{id}", "Get the change request specified by ID"},
		apiOption{"PUT", "/api/changes/{id}", "Update the change request specified by ID"},
		apiOption{"DELETE", "/api/changes/{id}", "Delete the change request specified by ID"},
		apiOption{"GET", "/api/users", "Display all users"},
		apiOption{"POST", "/api/users", "Create a new user"},
		apiOption{"GET", "/api/users/{id}", "Get the user specified by ID"},
		apiOption{"PUT", "/api/users/{id}", "Update the user specified by ID"},
		apiOption{"DELETE", "/api/users/{id}", "Delete the user specified by ID"},
	})
}

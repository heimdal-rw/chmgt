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
		apiOption{"POST", "/api/change", "Create a new change request"},
		apiOption{"GET", "/api/change/{id}", "Get the change request specified by ID"},
		apiOption{"PUT", "/api/change/{id}", "Update the change request specified by ID"},
		apiOption{"DELETE", "/api/change/{id}", "Delete the change request specified by ID"},
	})
}

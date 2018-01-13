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

// APIHandler takes care of the index page
func APIHandler(w http.ResponseWriter, r *http.Request) {
	jsonOut := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// Sample structure of how the API may work

	apiOptions := make([]apiOption, 0)
	for _, route := range RouteDefs {
		opt := apiOption{
			route.Method,
			route.Pattern,
			route.Description,
		}
		apiOptions = append(apiOptions, opt)
	}

	jsonOut.Encode(apiOptions)
}

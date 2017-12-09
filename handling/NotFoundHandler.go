package handling

import (
	"net/http"
)

// NotFoundHandler takes care of any unhandled paths
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

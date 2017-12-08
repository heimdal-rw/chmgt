package handling

import (
	"fmt"
	"net/http"
)

// IndexHandler takes care of the index page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a test.")
}

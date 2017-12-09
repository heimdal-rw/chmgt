package routing

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/mattjw79/chmgt/handling"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		alice.New(handling.LogHandler).ThenFunc(handling.IndexHandler),
	},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.Handler)
	}
	router.NotFoundHandler = alice.New(handling.LogHandler).ThenFunc(handling.NotFoundHandler)
	return router
}

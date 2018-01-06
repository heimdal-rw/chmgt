package routing

import (
	"net/http"

	"chmgt/handling"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// Route provides all of the items needed to build a route
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
}

// Routes provides an array of routes
type Routes []Route

var routes = Routes{
	Route{
		"API",
		"GET",
		"/api",
		alice.New(handling.LogHandler).ThenFunc(handling.APIHandler),
	},
	Route{
		"GetChanges",
		"GET",
		"/api/changes",
		alice.New(handling.LogHandler).ThenFunc(handling.GetChangesHandler),
	},
}

// NewRouter builds the routing structure
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	// Attach each route to the router
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.Handler)
	}
	// Serve out static files
	// This MUST be after all other routes as it is a catch-all
	router.PathPrefix("/").Handler(
		alice.New(handling.LogHandler).Then(http.StripPrefix(
			"/",
			http.FileServer(http.Dir("static")),
		)),
	)
	return router
}

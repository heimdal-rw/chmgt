package routing

import (
	"net/http"

	"chmgt/handling"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// Route provides all of the items needed to build a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	Description string
	Handler     http.Handler
}

// Routes provides an array of routes
type Routes []Route

// Definitions provides the route definitions
var Definitions = Routes{
	Route{
		"API",
		"GET",
		"/api",
		"Get help on the API",
		alice.New(handling.LogHandler).ThenFunc(handling.APIHandler),
	},
	Route{
		"GetChanges",
		"GET",
		"/api/changes",
		"Display all change requests",
		alice.New(handling.LogHandler).ThenFunc(handling.GetChangesHandler),
	},
	Route{
		"GetChange",
		"GET",
		"/api/changes/{id}",
		"Display a single change request specified by ID",
		alice.New(handling.LogHandler).ThenFunc(handling.GetChangesHandler),
	},
	Route{
		"CreateChange",
		"POST",
		"/api/changes",
		"Create a new change request",
		alice.New(handling.LogHandler).ThenFunc(handling.CreateChangeHandler),
	},
	Route{
		"DeleteChange",
		"DELETE",
		"/api/changes/{id}",
		"Delete the change request specified by ID",
		alice.New(handling.LogHandler).ThenFunc(handling.DeleteChangeHandler),
	},
	Route{
		"UpdateChange",
		"PUT",
		"/api/changes/{id}",
		"Update the change request specified by ID",
		alice.New(handling.LogHandler).ThenFunc(handling.UpdateChangeHandler),
	},
	Route{
		"GetUsers",
		"GET",
		"/api/users",
		"Display all users",
		alice.New(handling.LogHandler).ThenFunc(handling.GetUsersHandler),
	},
	Route{
		"GetUser",
		"GET",
		"/api/users/{id}",
		"Display a single user specified by ID",
		alice.New(handling.LogHandler).ThenFunc(handling.GetUsersHandler),
	},
	Route{
		"CreateUser",
		"POST",
		"/api/users",
		"Create a new user",
		alice.New(handling.LogHandler).ThenFunc(handling.CreateUserHandler),
	},
	Route{
		"DeleteUser",
		"DELETE",
		"/api/users/{id}",
		"Delete the user specified by ID",
		alice.New(handling.LogHandler).ThenFunc(handling.DeleteUserHandler),
	},
	Route{
		"UpdateUser",
		"PUT",
		"/api/users/{id}",
		"Update the user specified by ID",
		alice.New(handling.LogHandler).ThenFunc(handling.UpdateUserHandler),
	},
}

// NewRouter builds the routing structure
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	// Attach each route to the router
	for _, route := range Definitions {
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

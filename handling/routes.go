package handling

import (
	"chmgt/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type Env struct {
	DB models.Datastore
}

// Route provides all of the items needed to build a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	Description string
	Handler     http.Handler
}

// RouteDefs provides an array of routes
var RouteDefs []Route

// NewRouter builds the routing structure
func NewRouter() *mux.Router {
	db, err := models.NewDB("chmgt.db")
	if err != nil {
		log.Fatal(err)
	}

	env := &Env{db}

	RouteDefs = []Route{
		Route{
			"API",
			"GET",
			"/api",
			"Get help on the API",
			alice.New(LogHandler).ThenFunc(APIHandler),
		},
		Route{
			"GetChanges",
			"GET",
			"/api/changes",
			"Display all change requests",
			alice.New(LogHandler).ThenFunc(env.GetChangesHandler),
		},
		Route{
			"GetChange",
			"GET",
			"/api/changes/{id}",
			"Display a single change request specified by ID",
			alice.New(LogHandler).ThenFunc(env.GetChangesHandler),
		},
		Route{
			"CreateChange",
			"POST",
			"/api/changes",
			"Create a new change request",
			alice.New(LogHandler).ThenFunc(env.CreateChangeHandler),
		},
		Route{
			"DeleteChange",
			"DELETE",
			"/api/changes/{id}",
			"Delete the change request specified by ID",
			alice.New(LogHandler).ThenFunc(env.DeleteChangeHandler),
		},
		Route{
			"UpdateChange",
			"PUT",
			"/api/changes/{id}",
			"Update the change request specified by ID",
			alice.New(LogHandler).ThenFunc(env.UpdateChangeHandler),
		},
		Route{
			"GetUsers",
			"GET",
			"/api/users",
			"Display all users",
			alice.New(LogHandler).ThenFunc(env.GetUsersHandler),
		},
		Route{
			"GetUser",
			"GET",
			"/api/users/{id}",
			"Display a single user specified by ID",
			alice.New(LogHandler).ThenFunc(env.GetUsersHandler),
		},
		Route{
			"CreateUser",
			"POST",
			"/api/users",
			"Create a new user",
			alice.New(LogHandler).ThenFunc(env.CreateUserHandler),
		},
		Route{
			"DeleteUser",
			"DELETE",
			"/api/users/{id}",
			"Delete the user specified by ID",
			alice.New(LogHandler).ThenFunc(env.DeleteUserHandler),
		},
		Route{
			"UpdateUser",
			"PUT",
			"/api/users/{id}",
			"Update the user specified by ID",
			alice.New(LogHandler).ThenFunc(env.UpdateUserHandler),
		},
	}

	router := mux.NewRouter().StrictSlash(true)
	// Attach each route to the router
	for _, route := range RouteDefs {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.Handler)
	}
	// Serve out static files
	// This MUST be after all other routes as it is a catch-all
	router.PathPrefix("/").Handler(
		alice.New(LogHandler).Then(http.StripPrefix(
			"/",
			http.FileServer(http.Dir("static")),
		)),
	)
	return router
}

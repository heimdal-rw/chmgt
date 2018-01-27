package handling

import (
	"net/http"

	"chmgt/config"
	"chmgt/models"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// Handler encompases all request handling
type Handler struct {
	Router     http.Handler
	Config     *config.Config
	Datasource *models.Datasource
}

// NewHandler builds the handler interface and routes
func NewHandler(config *config.Config) (*Handler, error) {
	handler := new(Handler)
	handler.Config = config

	router := mux.NewRouter()

	router.
		Methods("GET").
		Path("/api/users/{id}").
		Handler(alice.New(
			handler.SetConfig,
			handler.SetLogging,
		).ThenFunc(handler.GetUsersHandler))

	router.
		Methods("GET").
		Path("/api/users").
		Handler(alice.New(
			handler.SetConfig,
			handler.SetLogging,
		).ThenFunc(handler.GetUsersHandler))

	router.
		Methods("POST").
		Path("/api/users").
		Handler(alice.New(
			handler.SetConfig,
			handler.SetLogging,
		).ThenFunc(handler.CreateUserHandler))

	router.
		Methods("DELETE").
		Path("/api/users/{id}").
		Handler(alice.New(
			handler.SetConfig,
			handler.SetLogging,
		).ThenFunc(handler.DeleteUserHandler))

	router.
		Methods("PUT").
		Path("/api/users/{id}").
		Handler(alice.New(
			handler.SetConfig,
			handler.SetLogging,
		).ThenFunc(handler.UpdateUserHandler))

	// This is a "catch-all" that serves static files and logs
	// any 404s from bad requests
	router.
		PathPrefix("/").
		Handler(alice.New(
			handler.SetConfig,
			handler.SetLogging,
		).Then(
			http.StripPrefix("/", http.FileServer(http.Dir("static"))),
		))

	// Use gorilla's recovery handler to continue running in case of a panic
	handler.Router = handlers.RecoveryHandler()(router)

	return handler, nil
}

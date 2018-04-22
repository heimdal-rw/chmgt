package handling

import (
	"encoding/json"
	"net/http"

	"github.com/heimdal-rw/chmgt/config"
	"github.com/heimdal-rw/chmgt/models"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

var commonHandlers alice.Chain
var authenticatedHandlers alice.Chain

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

	commonHandlers = alice.New(
		handler.SetConfig,
		handler.SetLogging,
		handler.CheckHeaders,
	)

	authenticatedHandlers = alice.New(
		handler.CheckAuthentication,
	)
	authenticatedHandlers = authenticatedHandlers.Extend(commonHandlers)

	router := mux.NewRouter()

	router.
		Methods("POST").
		Path("/api/authenticate").
		Handler(commonHandlers.ThenFunc(handler.AuthenticateHandler))

	// users api
	router.
		Methods("GET").
		Path("/api/users").
		Handler(authenticatedHandlers.ThenFunc(handler.GetUsers))

	router.
		Methods("GET").
		Path("/api/users/{id}").
		Handler(authenticatedHandlers.ThenFunc(handler.GetUsers))

	router.
		Methods("POST").
		Path("/api/users").
		Handler(authenticatedHandlers.ThenFunc(handler.InsertUser))

	router.
		Methods("PUT").
		Path("/api/users/{id}").
		Handler(authenticatedHandlers.ThenFunc(handler.UpdateUser))

	router.
		Methods("DELETE").
		Path("/api/users/{id}").
		Handler(authenticatedHandlers.ThenFunc(handler.RemoveUser))

	// changes api
	router.
		Methods("GET").
		Path("/api/changes").
		Handler(authenticatedHandlers.ThenFunc(handler.GetChanges))

	router.
		Methods("GET").
		Path("/api/changes/{id}").
		Handler(authenticatedHandlers.ThenFunc(handler.GetChanges))

	router.
		Methods("POST").
		Path("/api/changes").
		Handler(authenticatedHandlers.ThenFunc(handler.InsertChange))

	router.
		Methods("PUT").
		Path("/api/changes/{id}").
		Handler(authenticatedHandlers.ThenFunc(handler.UpdateChange))

	router.
		Methods("DELETE").
		Path("/api/changes/{id}").
		Handler(authenticatedHandlers.ThenFunc(handler.RemoveChange))

	// This is a "catch-all" that serves static files and logs
	// any 404s from bad requests
	router.
		PathPrefix("/").
		Handler(commonHandlers.Then(
			http.StripPrefix("/", http.FileServer(http.Dir("static"))),
		))

	// Use gorilla's recovery handler to continue running in case of a panic
	handler.Router = handlers.RecoveryHandler()(router)

	return handler, nil
}

// APIResponse creates the structure of an API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// WriteJSON format a JSON response and write it back to the client
func (j APIResponse) WriteJSON(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(j)
}

// APIWriteSuccess builds a success response
func APIWriteSuccess(w http.ResponseWriter, data interface{}) {
	APIResponse{
		true,
		"success",
		data,
	}.WriteJSON(w, http.StatusOK)
}

// APIWriteFailure builds a failure response
func APIWriteFailure(w http.ResponseWriter, msg string, status int) {
	if msg == "" {
		msg = "unknown error occured"
	}

	APIResponse{
		false,
		msg,
		nil,
	}.WriteJSON(w, status)
}

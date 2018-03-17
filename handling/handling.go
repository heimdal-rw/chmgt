package handling

import (
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

	// addUserRoutes(router, handler)
	// addChangeRequestRoutes(router, handler)

	router.
		Methods("POST").
		Path("/api/authenticate").
		Handler(commonHandlers.ThenFunc(handler.AuthenticateHandler))

	router.
		Methods("GET").
		Path("/api/{collection}").
		Handler(commonHandlers.ThenFunc(handler.GetItemsHandler))

	router.
		Methods("GET").
		Path("/api/{collection}/{id}").
		Handler(commonHandlers.ThenFunc(handler.GetItemsHandler))

	router.
		Methods("POST").
		Path("/api/{collection}").
		Handler(commonHandlers.ThenFunc(handler.CreateItemHandler))

	router.
		Methods("DELETE").
		Path("/api/{collection}/{id}").
		Handler(commonHandlers.ThenFunc(handler.DeleteItemHandler))

	router.
		Methods("PUT").
		Path("/api/{collection}/{id}").
		Handler(commonHandlers.ThenFunc(handler.UpdateItemHandler))

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

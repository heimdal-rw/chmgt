package handling

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func addChangeRequestRoutes(router *mux.Router, handler *Handler) {
	router.
		Methods("GET").
		Path("/api/changeRequests/{id}").
		Handler(alice.New(
			handler.SetConfig,
			handler.SetLogging,
		).ThenFunc(handler.GetChangeRequestsHandler))

	router.
		Methods("GET").
		Path("/api/changeRequests").
		Handler(alice.New(
			handler.SetConfig,
			handler.SetLogging,
		).ThenFunc(handler.GetChangeRequestsHandler))

	router.
		Methods("POST").
		Path("/api/changeRequests").
		Handler(alice.New(
			handler.SetConfig,
			handler.SetLogging,
		).ThenFunc(handler.CreateChangeRequestHandler))

	router.
		Methods("DELETE").
		Path("/api/changeRequests/{id}").
		Handler(alice.New(
			handler.SetConfig,
			handler.SetLogging,
		).ThenFunc(handler.DeleteChangeRequestHandler))

	router.
		Methods("PUT").
		Path("/api/changeRequests/{id}").
		Handler(alice.New(
			handler.SetConfig,
			handler.SetLogging,
		).ThenFunc(handler.UpdateChangeRequestHandler))
}

package handling

import (
	"github.com/gorilla/mux"
)

func addChangeRequestRoutes(router *mux.Router, handler *Handler) {
	router.
		Methods("GET").
		Path("/api/changeRequests/{id}").
		Handler(authenticatedHandlers.ThenFunc(handler.GetChangeRequestsHandler))

	router.
		Methods("GET").
		Path("/api/changeRequests").
		Handler(authenticatedHandlers.ThenFunc(handler.GetChangeRequestsHandler))

	router.
		Methods("POST").
		Path("/api/changeRequests").
		Handler(authenticatedHandlers.ThenFunc(handler.CreateChangeRequestHandler))

	router.
		Methods("DELETE").
		Path("/api/changeRequests/{id}").
		Handler(authenticatedHandlers.ThenFunc(handler.DeleteChangeRequestHandler))

	router.
		Methods("PUT").
		Path("/api/changeRequests/{id}").
		Handler(authenticatedHandlers.ThenFunc(handler.UpdateChangeRequestHandler))
}

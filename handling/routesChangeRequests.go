package handling

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func addChangeRequestRoutes(router *mux.Router, handler *Handler) {
	commonHandlers := alice.New(
		handler.CheckAuthentication,
		handler.SetConfig,
		handler.SetLogging,
	)

	router.
		Methods("GET").
		Path("/api/changeRequests/{id}").
		Handler(commonHandlers.ThenFunc(handler.GetChangeRequestsHandler))

	router.
		Methods("GET").
		Path("/api/changeRequests").
		Handler(commonHandlers.ThenFunc(handler.GetChangeRequestsHandler))

	router.
		Methods("POST").
		Path("/api/changeRequests").
		Handler(commonHandlers.ThenFunc(handler.CreateChangeRequestHandler))

	router.
		Methods("DELETE").
		Path("/api/changeRequests/{id}").
		Handler(commonHandlers.ThenFunc(handler.DeleteChangeRequestHandler))

	router.
		Methods("PUT").
		Path("/api/changeRequests/{id}").
		Handler(commonHandlers.ThenFunc(handler.UpdateChangeRequestHandler))
}

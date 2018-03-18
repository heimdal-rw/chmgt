package handling

import (
	"github.com/gorilla/mux"
)

func addUserRoutes(router *mux.Router, handler *Handler) {
	router.
		Methods("GET").
		Path("/api/users/{id}").
		Handler(authenticatedHandlers.ThenFunc(handler.GetUsersHandler))

	router.
		Methods("GET").
		Path("/api/users").
		Handler(authenticatedHandlers.ThenFunc(handler.GetUsersHandler))

	router.
		Methods("POST").
		Path("/api/users").
		Handler(authenticatedHandlers.ThenFunc(handler.CreateUserHandler))

	router.
		Methods("DELETE").
		Path("/api/users/{id}").
		Handler(authenticatedHandlers.ThenFunc(handler.DeleteUserHandler))

	router.
		Methods("PUT").
		Path("/api/users/{id}").
		Handler(authenticatedHandlers.ThenFunc(handler.UpdateUserHandler))
}

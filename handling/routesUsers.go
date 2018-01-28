package handling

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func addUserRoutes(router *mux.Router, handler *Handler) {
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
}

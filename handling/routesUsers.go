package handling

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func addUserRoutes(router *mux.Router, handler *Handler) {
	commonHandlers := alice.New(
		handler.CheckAuthentication,
		handler.SetConfig,
		handler.SetLogging,
	)

	router.
		Methods("GET").
		Path("/api/users/{id}").
		Handler(commonHandlers.ThenFunc(handler.GetUsersHandler))

	router.
		Methods("GET").
		Path("/api/users").
		Handler(commonHandlers.ThenFunc(handler.GetUsersHandler))

	router.
		Methods("POST").
		Path("/api/users").
		Handler(commonHandlers.ThenFunc(handler.CreateUserHandler))

	router.
		Methods("DELETE").
		Path("/api/users/{id}").
		Handler(commonHandlers.ThenFunc(handler.DeleteUserHandler))

	router.
		Methods("PUT").
		Path("/api/users/{id}").
		Handler(commonHandlers.ThenFunc(handler.UpdateUserHandler))
}

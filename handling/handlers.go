package handling

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

// SetConfig makes sure all configuration settings are applied
func (h *Handler) SetConfig(next http.Handler) http.Handler {
	if h.Config.UseProxyHeaders {
		next = handlers.ProxyHeaders(next)
	}
	return next
}

// SetLogging enables logging of each request
func (h *Handler) SetLogging(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stderr, next)
}

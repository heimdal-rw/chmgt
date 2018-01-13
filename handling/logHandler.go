package handling

import (
	"log"
	"net/http"
	"time"
)

// Create a wrapped ResponseWriter so that we can
// store the statusCode for the LogWriter
type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Factory for the wrappedResponseWriter with a default
// statusCode of 200
func wrapResponseWriter(w http.ResponseWriter) *wrappedResponseWriter {
	return &wrappedResponseWriter{w, http.StatusOK}
}

// Override the ResponseWriter "WriteHeader" function
// and allow setting of the statusCode
func (wrw *wrappedResponseWriter) WriteHeader(code int) {
	wrw.statusCode = code
	wrw.ResponseWriter.WriteHeader(code)
}

// LogHandler adds a deferred log line to be printed
// after all other handlers have completed
func LogHandler(next http.Handler) http.Handler {
	h := func(w http.ResponseWriter, r *http.Request) {
		wrw := wrapResponseWriter(w)
		requestStart := time.Now()
		defer func() {
			log.Printf(
				"%v %v %v %v %v %v",
				r.Host,
				r.Method,
				r.URL.EscapedPath(),
				wrw.statusCode,
				r.Proto,
				time.Since(requestStart),
			)
		}()
		next.ServeHTTP(wrw, r)
	}
	return http.HandlerFunc(h)
}

package handling

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
)

// AuthenticateHandler checks the users credentials and issues a token
func (h *Handler) AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	var input map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to parse body"))
		return
	}

	username := input["username"].(string)
	password := input["password"].(string)

	if username != "" && password != "" {
		valid, err := h.Datasource.ValidateUser(username, password)
		if err != nil {
			APIWriteFailure(w, "error validating username and password", http.StatusInternalServerError)
			return
		}
		if valid {
			claims := &jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute * time.Duration(h.Config.Server.SessionTimeout)).Unix(),
				Issuer:    "chmgt",
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
			ss, err := token.SignedString([]byte(h.Config.Server.SessionSecret))
			if err != nil {
				APIWriteFailure(w, "", http.StatusInternalServerError)
				log.Println(err)
				return
			}

			APIWriteSuccess(w, ss)
			return
		}
	}

	// Authentication failed
	APIWriteFailure(w, "invalid username or password", http.StatusUnauthorized)
}

// CheckAuthentication verifies the user is logged in and has a valid token
func (h *Handler) CheckAuthentication(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// Check if we have a cookie with the token string first
		authCookie, err := r.Cookie("Authorization")
		if err != nil {
			// Fall back to the Authorization header for the token string
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		} else {
			// Found Authorization cookie, using that
			tokenString = authCookie.Value
		}

		// Parse the signed string
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.Config.Server.SessionSecret), nil
		})
		switch err.(type) {
		case nil:
			// Token may still be invalid, so check to make sure
			if token.Valid {
				// Authorized, continue on to the next handler
				next.ServeHTTP(w, r)
				return
			}
			APIWriteFailure(w, "token invalid", http.StatusUnauthorized)
			return
		case *jwt.ValidationError:
			if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
				APIWriteFailure(w, "token expired", http.StatusUnauthorized)
				return
			}
			APIWriteFailure(w, "token validation error", http.StatusInternalServerError)
			return
		default:
			APIWriteFailure(w, "token validation error", http.StatusInternalServerError)
			return
		}
	}
	return http.HandlerFunc(fn)
}

// SetConfig makes sure all configuration settings are applied
func (h *Handler) SetConfig(next http.Handler) http.Handler {
	if h.Config.Server.UseProxyHeaders {
		next = handlers.ProxyHeaders(next)
	}
	return next
}

// SetLogging enables logging of each request
func (h *Handler) SetLogging(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stderr, next)
}

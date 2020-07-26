package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/daemonfire/mister-fruits/internal/connector"
	"github.com/daemonfire/mister-fruits/internal/pkg/user"
)

type TokenGuardMiddleware struct {
	Authenticator connector.Authenticator
}

func (m *TokenGuardMiddleware) Handler() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//
			// Real world application should not log sensitive data
			// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
			tokenValue := r.Header.Get("Authorization")
			log.Println("Received Authorization Header", tokenValue)
			tokenValue = strings.TrimPrefix(tokenValue, "Bearer ")
			log.Println("Token Value", tokenValue)
			username, err := m.Authenticator.CheckToken(tokenValue)
			if err != nil {
				log.Println("Failed to authenticate token", err)
				w.WriteHeader(401)
				return
			}
			if username == "" {
				log.Println("Token contains invalid username", err)
				w.WriteHeader(400)
				return
			}
			log.Println("Authenticated token for user:", username)
			r = r.WithContext(user.ContextWithUser(r.Context(), username))
			next.ServeHTTP(w, r)
		})
	}
}

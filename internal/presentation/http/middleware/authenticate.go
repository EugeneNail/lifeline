package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
)

// Authenticate returns a middleware that restores the bearer token from the request, validates it, stores the user ID
// in the request context, and calls the next handler or returns 401 when authentication fails.
func Authenticate(provider auth.TokenProvider, identity authentication.RequestIdentity) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			if provider == nil {
				http.Error(writer, "authentication provider is not configured", http.StatusInternalServerError)
				return
			}

			rawToken, err := bearerTokenFromRequest(request)
			if err != nil {
				http.Error(writer, "unauthorized", http.StatusUnauthorized)
				return
			}

			token, err := provider.Restore(rawToken)
			if err != nil || !token.IsValid() {
				http.Error(writer, "unauthorized", http.StatusUnauthorized)
				return
			}

			request = identity.WithAccountID(request, token.UserID())

			next(writer, request)
		}
	}
}

// bearerTokenFromRequest returns the bearer token from the Authorization header or an error when the header is missing or malformed.
func bearerTokenFromRequest(request *http.Request) (string, error) {
	if request == nil {
		return "", fmt.Errorf("request is nil")
	}

	authorization := request.Header.Get("Authorization")
	if authorization == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authorization, bearerPrefix) {
		return "", fmt.Errorf("authorization header is not a bearer token")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authorization, bearerPrefix))
	if token == "" {
		return "", fmt.Errorf("bearer token is empty")
	}

	return token, nil
}

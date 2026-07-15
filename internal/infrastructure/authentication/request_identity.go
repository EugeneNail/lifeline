package authentication

import (
	"context"
	"fmt"
	"net/http"

	"github.com/EugeneNail/lifeline/internal/domain/auth"
)

type requestIdentityContextKey struct{}

// RequestIdentity stores and retrieves a user identifier in HTTP requests.
type RequestIdentity struct{}

// NewRequestIdentity returns a request identity helper.
func NewRequestIdentity() RequestIdentity {
	return RequestIdentity{}
}

// WithAccountID returns a new request that carries the provided user identifier.
func (identity RequestIdentity) WithAccountID(request *http.Request, userID auth.ID) *http.Request {
	ctx := context.WithValue(request.Context(), requestIdentityContextKey{}, userID)
	return request.WithContext(ctx)
}

// AccountID returns the user identifier stored in the request context or an error when it is missing or invalid.
func (identity RequestIdentity) AccountID(request *http.Request) (auth.ID, error) {
	if request == nil {
		return auth.NilID, fmt.Errorf("request is nil")
	}

	value := request.Context().Value(requestIdentityContextKey{})
	if value == nil {
		return auth.NilID, fmt.Errorf("user id is missing")
	}

	userID, ok := value.(auth.ID)
	if !ok {
		return auth.NilID, fmt.Errorf("invalid user id type %T", value)
	}

	return userID, nil
}

package authenticate

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/EugeneNail/lifeline/internal/application"
	authenticate "github.com/EugeneNail/lifeline/internal/application/usecases/authenticate"
)

// Handler adapts the authenticate use case to the HTTP transport.
type Handler struct {
	usecase *authenticate.Handler
}

// NewHandler returns a transport handler wired to the authenticate use case.
func NewHandler(usecase *authenticate.Handler) *Handler {
	return &Handler{usecase: usecase}
}

// Payload represents the JSON request body for user authentication.
type Payload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Response represents the JSON response body with issued tokens.
type Response struct {
	LoginToken   string `json:"loginToken"`
	RefreshToken string `json:"refreshToken"`
}

// Handle decodes the request, runs the use case, and maps the result to an HTTP response.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	var payload Payload
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, fmt.Errorf("decoding request body: %w", err)
	}

	result, err := handler.usecase.Handle(request.Context(), authenticate.Query{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		var fieldErrors application.FieldErrors
		if errors.As(err, &fieldErrors) {
			return http.StatusUnprocessableEntity, fieldErrors.Errors()
		}

		return http.StatusInternalServerError, fmt.Errorf("handling Authenticate command: %w", err)
	}

	return http.StatusOK, Response{
		LoginToken:   result.LoginToken,
		RefreshToken: result.RefreshToken,
	}
}

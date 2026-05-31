package refresh

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/EugeneNail/lifeline/internal/application"
	"github.com/EugeneNail/lifeline/internal/application/usecases/refresh"
)

// Handler adapts the refresh use case to the HTTP transport.
type Handler struct {
	usecase *refresh.Handler
}

// NewHandler returns a transport handler wired to the refresh use case.
func NewHandler(usecase *refresh.Handler) *Handler {
	return &Handler{usecase: usecase}
}

// Payload represents the JSON request body for token refresh.
type Payload struct {
	RefreshToken string `json:"refreshToken"`
}

// Handle decodes the request, runs the use case, and maps the result to an HTTP response.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	var payload Payload
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, fmt.Errorf("decoding request body: %w", err)
	}

	token, err := handler.usecase.Handle(request.Context(), refresh.Command{
		RefreshToken: payload.RefreshToken,
	})
	if err != nil {
		var fieldErrors application.FieldErrors
		if errors.As(err, &fieldErrors) {
			return http.StatusUnauthorized, fieldErrors.Errors()
		}

		return http.StatusInternalServerError, fmt.Errorf("handling Refresh command: %w", err)
	}

	return http.StatusOK, token.String()
}

package register_user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EugeneNail/lifeline/internal/application"
	"github.com/EugeneNail/lifeline/internal/application/usecases/register_user"
	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"net/http"
)

// Handler adapts the register-user use case to the HTTP transport.
type Handler struct {
	usecase *register_user.Handler
}

// NewHandler returns a transport handler wired to the register-user use case.
func NewHandler(usecase *register_user.Handler) *Handler {
	return &Handler{usecase: usecase}
}

// RegisterUserPayload represents the JSON request body for user registration.
type RegisterUserPayload struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

// Handle decodes the request, runs the use case, and maps the result to an HTTP response.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	var payload RegisterUserPayload
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, fmt.Errorf("decoding request body: %w", err)
	}

	id, err := handler.usecase.Handle(request.Context(), register_user.RegisterUserCommand{
		Email:                payload.Email,
		Password:             payload.Password,
		PasswordConfirmation: payload.PasswordConfirmation,
	})
	if err != nil {
		var fieldErrors application.FieldErrors
		if errors.As(err, &fieldErrors) {
			return http.StatusUnprocessableEntity, fieldErrors.Errors()
		}

		if errors.Is(err, auth.EmailAlreadyTaken) {
			return http.StatusConflict, auth.EmailAlreadyTaken.Error()
		}

		return http.StatusInternalServerError, fmt.Errorf("handling RegisterUser command: %w", err)
	}

	return http.StatusCreated, id.String()
}

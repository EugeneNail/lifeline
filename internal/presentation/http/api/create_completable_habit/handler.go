package create_completable_habit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/EugeneNail/lifeline/internal/application/usecases/create_completable_habit"
	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
)

// Handler adapts the create-completable-habit use case to the HTTP transport.
type Handler struct {
	usecase  *create_completable_habit.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the create-completable-habit use case.
func NewHandler(usecase *create_completable_habit.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Payload represents the JSON request body for completable habit creation.
type Payload struct {
	Label string `json:"label"`
	Icon  int    `json:"icon"`
}

// Handle decodes the request, runs the use case, and maps the result to an HTTP response.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	accountID, err := handler.identity.AccountID(request)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("extracting account id: %w", err)
	}

	var payload Payload
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, fmt.Errorf("decoding request body: %w", err)
	}

	id, err := handler.usecase.Handle(request.Context(), create_completable_habit.Command{
		Label:     payload.Label,
		Icon:      payload.Icon,
		AccountID: accountID.Uuid(),
	})
	if err != nil {
		var validationErrors domain.ValidationErrors
		if errors.As(err, &validationErrors) {
			return http.StatusUnprocessableEntity, validationErrors.Errors()
		}

		if errors.Is(err, habits.ErrHabitLimitExceeded) {
			return http.StatusConflict, nil
		}

		return http.StatusInternalServerError, fmt.Errorf("handling CreateCompletableHabit command: %w", err)
	}

	return http.StatusCreated, id.String()
}

package create_time_habit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/EugeneNail/lifeline/internal/application/usecases/create_time_habit"
	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
)

// Handler adapts the create-time-habit use case to the HTTP transport.
type Handler struct {
	usecase  *create_time_habit.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the create-time-habit use case.
func NewHandler(usecase *create_time_habit.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Payload represents the JSON request body for time habit creation.
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

	id, err := handler.usecase.Handle(request.Context(), create_time_habit.Command{
		Label:     payload.Label,
		Icon:      payload.Icon,
		AccountID: accountID.Uuid(),
	})
	if err != nil {
		var violations domain.Violations
		if errors.As(err, &violations) {
			return http.StatusUnprocessableEntity, violations.Violations()
		}

		if errors.Is(err, habits.ErrHabitLimitExceeded) {
			return http.StatusConflict, nil
		}

		return http.StatusInternalServerError, fmt.Errorf("handling CreateTimeHabit command: %w", err)
	}

	return http.StatusCreated, id.String()
}

package create_measurable_habit

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/EugeneNail/lifeline/internal/application/usecases/create_measurable_habit"
	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
)

// Handler adapts the create-measurable-habit use case to the HTTP transport.
type Handler struct {
	usecase  *create_measurable_habit.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the create-measurable-habit use case.
func NewHandler(usecase *create_measurable_habit.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Payload represents the JSON request body for measurable habit creation.
type Payload struct {
	Label string  `json:"label"`
	Icon  int     `json:"icon"`
	Step  float32 `json:"step"`
	Unit  string  `json:"unit"`
}

// Handle decodes the request, runs the use case, and maps the result to an HTTP response.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	accountID, err := handler.identity.AccountID(request)
	if err != nil {
		return http.StatusUnauthorized, "unauthorized"
	}

	var payload Payload
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, fmt.Errorf("decoding request body: %w", err)
	}

	id, err := handler.usecase.Handle(request.Context(), create_measurable_habit.Command{
		Label:     payload.Label,
		Icon:      payload.Icon,
		Step:      payload.Step,
		Unit:      payload.Unit,
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

		return http.StatusInternalServerError, fmt.Errorf("handling CreateMeasurableHabit command: %w", err)
	}

	return http.StatusCreated, id.String()
}

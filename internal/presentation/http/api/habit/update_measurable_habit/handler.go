package update_measurable_habit

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/update_measurable_habit"
	"net/http"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
	"github.com/google/uuid"
)

// Handler adapts the update-measurable-habit use case to the HTTP transport.
type Handler struct {
	usecase  *update_measurable_habit.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the update-measurable-habit use case.
func NewHandler(usecase *update_measurable_habit.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Payload represents the JSON request body for measurable habit updates.
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
		return http.StatusInternalServerError, fmt.Errorf("extracting account id: %w", err)
	}

	habitID, err := uuid.Parse(request.PathValue("uuid"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing habit id: %w", err)
	}

	var payload Payload
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, fmt.Errorf("decoding request body: %w", err)
	}

	_, err = handler.usecase.Handle(request.Context(), update_measurable_habit.Command{
		ID:        habitID,
		Label:     payload.Label,
		Icon:      payload.Icon,
		Step:      payload.Step,
		Unit:      payload.Unit,
		AccountID: accountID.Uuid(),
	})
	if err != nil {
		var violations domain.Violations
		if errors.As(err, &violations) {
			return http.StatusUnprocessableEntity, violations.Violations()
		}

		if errors.Is(err, habits.ErrHabitNotFound) {
			return http.StatusNotFound, nil
		}

		if errors.Is(err, habits.ErrHabitBelongsToAnotherUser) {
			return http.StatusForbidden, nil
		}

		if errors.Is(err, habits.ErrHabitIsArchived) {
			return http.StatusConflict, nil
		}

		if errors.Is(err, habits.ErrHabitIsDeleted) {
			return http.StatusGone, nil
		}

		return http.StatusInternalServerError, fmt.Errorf("handling UpdateMeasurableHabit command: %w", err)
	}

	return http.StatusNoContent, nil
}

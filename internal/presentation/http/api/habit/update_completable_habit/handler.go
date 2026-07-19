package update_completable_habit

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/update_completable_habit"
	"net/http"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
	"github.com/google/uuid"
)

// Handler adapts the update-completable-habit use case to the HTTP transport.
type Handler struct {
	usecase  *update_completable_habit.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the update-completable-habit use case.
func NewHandler(usecase *update_completable_habit.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Payload represents the JSON request body for completable habit updates.
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

	habitID, err := uuid.Parse(request.PathValue("uuid"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing habit id: %w", err)
	}

	var payload Payload
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, fmt.Errorf("decoding request body: %w", err)
	}

	_, err = handler.usecase.Handle(request.Context(), update_completable_habit.Command{
		ID:        habitID,
		Label:     payload.Label,
		Icon:      payload.Icon,
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

		return http.StatusInternalServerError, fmt.Errorf("handling UpdateCompletableHabit command: %w", err)
	}

	return http.StatusNoContent, nil
}

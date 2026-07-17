package save_completable_habit_record

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/EugeneNail/lifeline/internal/application/usecases/save_completable_habit_record"
	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
	"github.com/google/uuid"
)

// Handler adapts the save-completable-habit-record use case to the HTTP transport.
type Handler struct {
	usecase  *save_completable_habit_record.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the save-completable-habit-record use case.
func NewHandler(usecase *save_completable_habit_record.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Payload represents the JSON request body for a completable habit record save operation.
type Payload struct {
	Value bool `json:"value"`
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

	date, err := time.Parse(time.DateOnly, request.PathValue("date"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing date: %w", err)
	}

	var payload Payload
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, fmt.Errorf("decoding request body: %w", err)
	}

	if err := handler.usecase.Handle(request.Context(), save_completable_habit_record.Command{
		AccountId:          accountID.Uuid(),
		Value:              payload.Value,
		Date:               date,
		CompletableHabitId: habitID,
	}); err != nil {
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

		return http.StatusInternalServerError, fmt.Errorf("handling SaveCompletableHabitRecord command: %w", err)
	}

	return http.StatusNoContent, nil
}

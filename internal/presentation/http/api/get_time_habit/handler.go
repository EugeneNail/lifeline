package get_time_habit

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EugeneNail/lifeline/internal/application/usecases/get_time_habit"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
	"github.com/google/uuid"
)

// Handler adapts the get-time-habit use case to the HTTP transport.
type Handler struct {
	usecase  *get_time_habit.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the get-time-habit use case.
func NewHandler(usecase *get_time_habit.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Output represents the JSON response body for a time habit.
type Output struct {
	ID         string     `json:"id"`
	Label      string     `json:"label"`
	Icon       int        `json:"icon"`
	ArchivedAt *time.Time `json:"archivedAt"`
}

// Handle loads the habit and returns an HTTP response with the public fields or 404 when it is missing.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	accountID, err := handler.identity.AccountID(request)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("extracting account id: %w", err)
	}

	habitID, err := uuid.Parse(request.PathValue("uuid"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing habit id: %w", err)
	}

	habit, err := handler.usecase.Handle(request.Context(), get_time_habit.Query{
		ID:        habitID,
		AccountID: accountID.Uuid(),
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("handling GetTimeHabit query: %w", err)
	}

	if habit == nil {
		return http.StatusNotFound, nil
	}

	return http.StatusOK, Output{
		ID:         habit.ID().String(),
		Label:      habit.Label(),
		Icon:       int(habit.Icon()),
		ArchivedAt: habit.ArchivedAt(),
	}
}

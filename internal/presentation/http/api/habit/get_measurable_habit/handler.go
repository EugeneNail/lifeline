package get_measurable_habit

import (
	"fmt"
	"github.com/EugeneNail/lifeline/internal/application/usecases/habit/get_measurable_habit"
	"net/http"
	"time"

	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
	"github.com/google/uuid"
)

// Handler adapts the get-measurable-habit use case to the HTTP transport.
type Handler struct {
	usecase  *get_measurable_habit.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the get-measurable-habit use case.
func NewHandler(usecase *get_measurable_habit.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Output represents the JSON response body for a measurable habit.
type Output struct {
	ID         string     `json:"id"`
	Label      string     `json:"label"`
	Icon       int        `json:"icon"`
	Step       float32    `json:"step"`
	Unit       string     `json:"unit"`
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

	habit, err := handler.usecase.Handle(request.Context(), get_measurable_habit.Query{
		ID:        habitID,
		AccountID: accountID.Uuid(),
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("handling GetMeasurableHabit query: %w", err)
	}

	if habit == nil {
		return http.StatusNotFound, nil
	}

	return http.StatusOK, Output{
		ID:         habit.ID().String(),
		Label:      habit.Label(),
		Icon:       int(habit.Icon()),
		Step:       float32(habit.Step()),
		Unit:       string(habit.Unit()),
		ArchivedAt: habit.ArchivedAt(),
	}
}

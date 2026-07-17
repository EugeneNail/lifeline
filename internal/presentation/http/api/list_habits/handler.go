package list_habits

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EugeneNail/lifeline/internal/application/usecases/list_habits"
	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
)

// Handler adapts the list-habits use case to the HTTP transport.
type Handler struct {
	usecase  *list_habits.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the list-habits use case.
func NewHandler(usecase *list_habits.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Output represents the JSON response body for the habit list.
type Output struct {
	Measurable  []MeasurableHabit  `json:"measurable"`
	Time        []TimeHabit        `json:"time"`
	Completable []CompletableHabit `json:"completable"`
}

// MeasurableHabit represents the public measurable habit fields returned to the client.
type MeasurableHabit struct {
	ID         string     `json:"id"`
	Label      string     `json:"label"`
	Icon       int        `json:"icon"`
	Step       float32    `json:"step"`
	Unit       string     `json:"unit"`
	ArchivedAt *time.Time `json:"archivedAt"`
}

// TimeHabit represents the public time habit fields returned to the client.
type TimeHabit struct {
	ID         string     `json:"id"`
	Label      string     `json:"label"`
	Icon       int        `json:"icon"`
	ArchivedAt *time.Time `json:"archivedAt"`
}

// CompletableHabit represents the public completable habit fields returned to the client.
type CompletableHabit struct {
	ID         string     `json:"id"`
	Label      string     `json:"label"`
	Icon       int        `json:"icon"`
	ArchivedAt *time.Time `json:"archivedAt"`
}

// Handle loads the user's habits, maps them to the transport output, and returns an HTTP response.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	accountID, err := handler.identity.AccountID(request)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("extracting account id: %w", err)
	}

	result, err := handler.usecase.Handle(request.Context(), list_habits.Query{
		AccountId: accountID.Uuid(),
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("handling ListHabits query: %w", err)
	}

	return http.StatusOK, Output{
		Measurable:  mapMeasurableHabits(result.MeasurableHabits),
		Time:        mapTimeHabits(result.TimeHabits),
		Completable: mapCompletableHabits(result.CompletableHabits),
	}
}

func mapMeasurableHabits(habitsList []*habits.MeasurableHabit) []MeasurableHabit {
	output := make([]MeasurableHabit, 0, len(habitsList))
	for _, habit := range habitsList {
		output = append(output, MeasurableHabit{
			ID:         habit.ID().String(),
			Label:      habit.Label(),
			Icon:       int(habit.Icon()),
			Step:       float32(habit.Step()),
			Unit:       string(habit.Unit()),
			ArchivedAt: habit.ArchivedAt(),
		})
	}

	return output
}

func mapTimeHabits(habitsList []*habits.TimeHabit) []TimeHabit {
	output := make([]TimeHabit, 0, len(habitsList))
	for _, habit := range habitsList {
		output = append(output, TimeHabit{
			ID:         habit.ID().String(),
			Label:      habit.Label(),
			Icon:       int(habit.Icon()),
			ArchivedAt: habit.ArchivedAt(),
		})
	}

	return output
}

func mapCompletableHabits(habitsList []*habits.CompletableHabit) []CompletableHabit {
	output := make([]CompletableHabit, 0, len(habitsList))
	for _, habit := range habitsList {
		output = append(output, CompletableHabit{
			ID:         habit.ID().String(),
			Label:      habit.Label(),
			Icon:       int(habit.Icon()),
			ArchivedAt: habit.ArchivedAt(),
		})
	}

	return output
}

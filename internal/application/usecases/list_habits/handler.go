package list_habits

import (
	"context"
	"fmt"

	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/google/uuid"
)

// Handler executes the list-habits use case.
type Handler struct {
	completableHabits habits.CompletableHabitRepository
	timeHabits        habits.TimeHabitRepository
	measurableHabits  habits.MeasurableHabitRepository
}

// NewHandler returns a list-habits handler configured with the habit repositories or an error when a dependency is missing.
func NewHandler(
	completableHabits habits.CompletableHabitRepository,
	timeHabits habits.TimeHabitRepository,
	measurableHabits habits.MeasurableHabitRepository,
) (*Handler, error) {
	if completableHabits == nil {
		return nil, fmt.Errorf("list_habits handler requires a completable habit repository")
	}

	if timeHabits == nil {
		return nil, fmt.Errorf("list_habits handler requires a time habit repository")
	}

	if measurableHabits == nil {
		return nil, fmt.Errorf("list_habits handler requires a measurable habit repository")
	}

	return &Handler{
		completableHabits: completableHabits,
		timeHabits:        timeHabits,
		measurableHabits:  measurableHabits,
	}, nil
}

// Query carries the data required to list habits.
type Query struct {
	AccountId uuid.UUID
}

// Result groups habits by their type and returns them to the caller.
type Result struct {
	MeasurableHabits  []*habits.MeasurableHabit
	TimeHabits        []*habits.TimeHabit
	CompletableHabits []*habits.CompletableHabit
}

// Handle loads the user's habits from all repositories and returns them grouped by habit type.
func (handler *Handler) Handle(ctx context.Context, query Query) (Result, error) {
	completableFilter := habits.NewCompletableHabitFilter().
		WithAccountIds(query.AccountId).
		WithDeleted(false)

	completableHabits, err := handler.completableHabits.FindMany(ctx, completableFilter)
	if err != nil {
		return Result{}, fmt.Errorf("finding completable habits: %w", err)
	}

	timeFilter := habits.NewTimeHabitFilter().
		WithAccountIds(query.AccountId).
		WithDeleted(false)

	timeHabits, err := handler.timeHabits.FindMany(ctx, timeFilter)
	if err != nil {
		return Result{}, fmt.Errorf("finding time habits: %w", err)
	}

	measurableFilter := habits.NewMeasurableHabitFilter().
		WithAccountIds(query.AccountId).
		WithDeleted(false)

	measurableHabits, err := handler.measurableHabits.FindMany(ctx, measurableFilter)
	if err != nil {
		return Result{}, fmt.Errorf("finding measurable habits: %w", err)
	}

	return Result{
		MeasurableHabits:  measurableHabits,
		TimeHabits:        timeHabits,
		CompletableHabits: completableHabits,
	}, nil
}

package get_measurable_habit

import (
	"context"
	"fmt"

	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/google/uuid"
)

// Handler executes the get-measurable-habit use case.
type Handler struct {
	measurableHabits habits.MeasurableHabitRepository
}

// NewHandler returns a get-measurable-habit handler configured with the measurable habit repository or an error when the dependency is missing.
func NewHandler(measurableHabits habits.MeasurableHabitRepository) (*Handler, error) {
	if measurableHabits == nil {
		return nil, fmt.Errorf("get_measurable_habit handler requires a measurable habit repository")
	}

	return &Handler{measurableHabits: measurableHabits}, nil
}

// Query carries the data required to load a measurable habit.
type Query struct {
	ID        uuid.UUID
	AccountID uuid.UUID
}

// Handle returns the measurable habit matching the query or nil when no habit exists, or an error when lookup fails.
func (handler *Handler) Handle(ctx context.Context, query Query) (*habits.MeasurableHabit, error) {
	habit, err := handler.measurableHabits.Find(ctx, habits.NewMeasurableHabitFilter().
		WithAccountIds(query.AccountID).
		WithIds(query.ID).
		WithDeleted(false),
	)
	if err != nil {
		return nil, fmt.Errorf("finding a measurable habit: %w", err)
	}

	return habit, nil
}

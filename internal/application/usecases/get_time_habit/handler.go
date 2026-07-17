package get_time_habit

import (
	"context"
	"fmt"

	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/google/uuid"
)

// Handler executes the get-time-habit use case.
type Handler struct {
	timeHabits habits.TimeHabitRepository
}

// NewHandler returns a get-time-habit handler configured with the time habit repository or an error when the dependency is missing.
func NewHandler(timeHabits habits.TimeHabitRepository) (*Handler, error) {
	if timeHabits == nil {
		return nil, fmt.Errorf("get_time_habit handler requires a time habit repository")
	}

	return &Handler{timeHabits: timeHabits}, nil
}

// Query carries the data required to load a time habit.
type Query struct {
	ID        uuid.UUID
	AccountID uuid.UUID
}

// Handle returns the time habit matching the query or nil when no habit exists, or an error when lookup fails.
func (handler *Handler) Handle(ctx context.Context, query Query) (*habits.TimeHabit, error) {
	habit, err := handler.timeHabits.Find(ctx, habits.NewTimeHabitFilter().
		WithAccountIds(query.AccountID).
		WithIds(query.ID).
		WithDeleted(false),
	)
	if err != nil {
		return nil, fmt.Errorf("finding a time habit: %w", err)
	}

	return habit, nil
}

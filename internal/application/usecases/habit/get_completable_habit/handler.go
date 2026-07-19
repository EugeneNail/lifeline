package get_completable_habit

import (
	"context"
	"fmt"

	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/google/uuid"
)

// Handler executes the get-completable-habit use case.
type Handler struct {
	completableHabits habits.CompletableHabitRepository
}

// NewHandler returns a get-completable-habit handler configured with the completable habit repository or an error when the dependency is missing.
func NewHandler(completableHabits habits.CompletableHabitRepository) (*Handler, error) {
	if completableHabits == nil {
		return nil, fmt.Errorf("get_completable_habit handler requires a completable habit repository")
	}

	return &Handler{completableHabits: completableHabits}, nil
}

// Query carries the data required to load a completable habit.
type Query struct {
	ID        uuid.UUID
	AccountID uuid.UUID
}

// Handle returns the completable habit matching the query or nil when no habit exists, or an error when lookup fails.
func (handler *Handler) Handle(ctx context.Context, query Query) (*habits.CompletableHabit, error) {
	habit, err := handler.completableHabits.Find(ctx, habits.NewCompletableHabitFilter().
		WithAccountIds(query.AccountID).
		WithIds(query.ID).
		WithDeleted(false),
	)
	if err != nil {
		return nil, fmt.Errorf("finding a completable habit: %w", err)
	}

	return habit, nil
}

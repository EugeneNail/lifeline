package create_time_habit

import (
	"context"
	"errors"
	"fmt"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/google/uuid"
)

// Handler executes the create-time-habit use case.
type Handler struct {
	timeHabits          habits.TimeHabitRepository
	habitCreationPolicy *habits.HabitCreationPolicy
}

// NewHandler returns a create-time-habit handler configured with the time habit repository and creation policy or an error when a dependency is missing.
func NewHandler(timeHabits habits.TimeHabitRepository, habitCreationPolicy *habits.HabitCreationPolicy) (*Handler, error) {
	if timeHabits == nil {
		return nil, fmt.Errorf("create_time_habit handler requires a time habit repository")
	}

	if habitCreationPolicy == nil {
		return nil, fmt.Errorf("create_time_habit handler requires a habit creation policy")
	}

	return &Handler{
		timeHabits:          timeHabits,
		habitCreationPolicy: habitCreationPolicy,
	}, nil
}

// Command carries the data required to create a time habit.
type Command struct {
	Label     string
	Icon      int
	AccountID uuid.UUID
}

// Handle validates the command, stores a new time habit, and returns the habit identifier or field validation errors.
func (handler *Handler) Handle(ctx context.Context, command Command) (uuid.UUID, error) {
	habit, err := habits.NewTimeHabit(command.Label, command.Icon, command.AccountID)
	if err != nil {
		var violations domain.Violations
		if errors.As(err, &violations) {
			return uuid.Nil, violations
		}

		return uuid.Nil, fmt.Errorf("creating a time habit: %w", err)
	}

	if err := handler.habitCreationPolicy.EnsureCanAdd(ctx, command.AccountID); err != nil {
		if errors.Is(err, habits.ErrHabitLimitExceeded) {
			return uuid.Nil, err
		}

		return uuid.Nil, fmt.Errorf("ensuring a habit can be added: %w", err)
	}

	if err := handler.timeHabits.Add(ctx, habit); err != nil {
		return uuid.Nil, fmt.Errorf("adding a new time habit to the collection: %w", err)
	}

	return habit.ID(), nil
}

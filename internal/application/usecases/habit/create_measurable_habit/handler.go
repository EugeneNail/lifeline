package create_measurable_habit

import (
	"context"
	"errors"
	"fmt"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/google/uuid"
)

// Handler executes the create-measurable-habit use case.
type Handler struct {
	measurableHabits    habits.MeasurableHabitRepository
	habitCreationPolicy *habits.HabitCreationPolicy
}

// NewHandler returns a create-measurable-habit handler configured with the measurable habit repository and creation policy or an error when a dependency is missing.
func NewHandler(measurableHabits habits.MeasurableHabitRepository, habitCreationPolicy *habits.HabitCreationPolicy) (*Handler, error) {
	if measurableHabits == nil {
		return nil, fmt.Errorf("create_measurable_habit handler requires a measurable habit repository")
	}

	if habitCreationPolicy == nil {
		return nil, fmt.Errorf("create_measurable_habit handler requires a habit creation policy")
	}

	return &Handler{
		measurableHabits:    measurableHabits,
		habitCreationPolicy: habitCreationPolicy,
	}, nil
}

// Command carries the data required to create a measurable habit.
type Command struct {
	Label     string
	Icon      int
	Step      float32
	Unit      string
	AccountID uuid.UUID
}

// Handle validates the command, stores a new measurable habit, and returns the habit identifier or field validation errors.
func (handler *Handler) Handle(ctx context.Context, command Command) (uuid.UUID, error) {
	habit, err := habits.NewMeasurableHabit(command.Label, command.Icon, command.Step, command.Unit, command.AccountID)
	if err != nil {
		var violations domain.Violations
		if errors.As(err, &violations) {
			return uuid.Nil, violations
		}

		return uuid.Nil, fmt.Errorf("creating a measurable habit: %w", err)
	}

	if err := handler.habitCreationPolicy.EnsureCanAdd(ctx, command.AccountID); err != nil {
		if errors.Is(err, habits.ErrHabitLimitExceeded) {
			return uuid.Nil, err
		}

		return uuid.Nil, fmt.Errorf("ensuring a habit can be added: %w", err)
	}

	if err := handler.measurableHabits.Add(ctx, habit); err != nil {
		return uuid.Nil, fmt.Errorf("adding a new measurable habit to the collection: %w", err)
	}

	return habit.ID(), nil
}

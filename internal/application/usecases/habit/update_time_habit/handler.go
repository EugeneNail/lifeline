package update_time_habit

import (
	"context"
	"errors"
	"fmt"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/google/uuid"
)

// Handler executes the update-time-habit use case.
type Handler struct {
	timeHabits         habits.TimeHabitRepository
	modificationPolicy *habits.ModificationPolicy
}

// NewHandler returns an update-time-habit handler configured with the time habit repository or an error when the dependency is missing.
func NewHandler(timeHabits habits.TimeHabitRepository, modificationPolicy *habits.ModificationPolicy) (*Handler, error) {
	if timeHabits == nil {
		return nil, fmt.Errorf("update_time_habit handler requires a time habit repository")
	}

	if modificationPolicy == nil {
		return nil, fmt.Errorf("update_time_habit handler requires a modification policy")
	}

	return &Handler{
		timeHabits:         timeHabits,
		modificationPolicy: modificationPolicy,
	}, nil
}

// Command carries the data required to update a time habit.
type Command struct {
	ID        uuid.UUID
	Label     string
	Icon      int
	AccountID uuid.UUID
}

// Handle validates the command, updates a time habit, and returns the habit identifier or field validation errors.
func (handler *Handler) Handle(ctx context.Context, command Command) (uuid.UUID, error) {
	violations := domain.NewViolations()

	label, err := habits.NewLabel(command.Label)
	if err != nil {
		var violation domain.Violation
		if !errors.As(err, &violation) {
			return uuid.Nil, fmt.Errorf("creating a time habit label: %w", err)
		}

		violations.Add("label", violation)
	}

	icon, err := habits.NewIcon(command.Icon)
	if err != nil {
		var violation domain.Violation
		if !errors.As(err, &violation) {
			return uuid.Nil, fmt.Errorf("creating a time habit icon: %w", err)
		}

		violations.Add("icon", violation)
	}

	if violations.HasViolations() {
		return uuid.Nil, violations
	}

	habit, err := handler.timeHabits.Find(ctx, habits.NewTimeHabitFilter().
		WithAccountIds(command.AccountID).
		WithIds(command.ID).
		WithDeleted(false),
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("finding a time habit: %w", err)
	}

	if habit == nil {
		return uuid.Nil, habits.ErrHabitNotFound
	}

	if err := handler.modificationPolicy.Check(command.AccountID, habit); err != nil {
		return uuid.Nil, err
	}

	habit.ChangeLabel(label)
	habit.ChangeIcon(icon)

	if err := handler.timeHabits.Update(ctx, habit); err != nil {
		return uuid.Nil, fmt.Errorf("saving a time habit: %w", err)
	}

	return habit.ID(), nil
}

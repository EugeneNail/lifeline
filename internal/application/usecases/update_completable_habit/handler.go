package update_completable_habit

import (
	"context"
	"errors"
	"fmt"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/google/uuid"
)

// Handler executes the update-completable-habit use case.
type Handler struct {
	completableHabits  habits.CompletableHabitRepository
	modificationPolicy *habits.ModificationPolicy
}

// NewHandler returns an update-completable-habit handler configured with the completable habit repository or an error when the dependency is missing.
func NewHandler(completableHabits habits.CompletableHabitRepository, modificationPolicy *habits.ModificationPolicy) (*Handler, error) {
	if completableHabits == nil {
		return nil, fmt.Errorf("update_completable_habit handler requires a completable habit repository")
	}

	if modificationPolicy == nil {
		return nil, fmt.Errorf("update_completable_habit handler requires a modification policy")
	}

	return &Handler{
		completableHabits:  completableHabits,
		modificationPolicy: modificationPolicy,
	}, nil
}

// Command carries the data required to update a completable habit.
type Command struct {
	ID        uuid.UUID
	Label     string
	Icon      int
	AccountID uuid.UUID
}

// Handle validates the command, updates a completable habit, and returns the habit identifier or field validation errors.
func (handler *Handler) Handle(ctx context.Context, command Command) (uuid.UUID, error) {
	errs := domain.NewValidationErrors()

	label, err := habits.NewLabel(command.Label)
	if err != nil {
		var domainError domain.Error
		if !errors.As(err, &domainError) {
			return uuid.Nil, fmt.Errorf("creating a completable habit label: %w", err)
		}

		errs.Add("label", domainError)
	}

	icon, err := habits.NewIcon(command.Icon)
	if err != nil {
		var domainError domain.Error
		if !errors.As(err, &domainError) {
			return uuid.Nil, fmt.Errorf("creating a completable habit icon: %w", err)
		}

		errs.Add("icon", domainError)
	}

	if errs.HasErrors() {
		return uuid.Nil, errs
	}

	habit, err := handler.completableHabits.Find(ctx, habits.NewCompletableHabitFilter().
		WithAccountIds(command.AccountID).
		WithIds(command.ID).
		WithDeleted(false),
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("finding a completable habit: %w", err)
	}

	if habit == nil {
		return uuid.Nil, habits.ErrHabitNotFound
	}

	if violation := handler.modificationPolicy.Check(command.AccountID, habit); violation != nil {
		return uuid.Nil, violation
	}

	habit.ChangeLabel(label)
	habit.ChangeIcon(icon)

	if err := handler.completableHabits.Save(ctx, habit); err != nil {
		return uuid.Nil, fmt.Errorf("saving a completable habit: %w", err)
	}

	return habit.ID(), nil
}

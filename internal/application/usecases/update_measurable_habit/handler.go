package update_measurable_habit

import (
	"context"
	"errors"
	"fmt"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/google/uuid"
)

// Handler executes the update-measurable-habit use case.
type Handler struct {
	measurableHabits   habits.MeasurableHabitRepository
	modificationPolicy *habits.ModificationPolicy
}

// NewHandler returns an update-measurable-habit handler configured with the measurable habit repository or an error when the dependency is missing.
func NewHandler(measurableHabits habits.MeasurableHabitRepository, modificationPolicy *habits.ModificationPolicy) (*Handler, error) {
	if measurableHabits == nil {
		return nil, fmt.Errorf("update_measurable_habit handler requires a measurable habit repository")
	}

	if modificationPolicy == nil {
		return nil, fmt.Errorf("update_measurable_habit handler requires a modification policy")
	}

	return &Handler{
		measurableHabits:   measurableHabits,
		modificationPolicy: modificationPolicy,
	}, nil
}

// Command carries the data required to update a measurable habit.
type Command struct {
	ID        uuid.UUID
	Label     string
	Icon      int
	Step      float32
	Unit      string
	AccountID uuid.UUID
}

// Handle validates the command, updates a measurable habit, and returns the habit identifier or field validation errors.
func (handler *Handler) Handle(ctx context.Context, command Command) (uuid.UUID, error) {
	errs := domain.NewValidationErrors()

	label, err := habits.NewLabel(command.Label)
	if err != nil {
		var domainError domain.Error
		if !errors.As(err, &domainError) {
			return uuid.Nil, fmt.Errorf("creating a measurable habit label: %w", err)
		}

		errs.Add("label", domainError)
	}

	icon, err := habits.NewIcon(command.Icon)
	if err != nil {
		var domainError domain.Error
		if !errors.As(err, &domainError) {
			return uuid.Nil, fmt.Errorf("creating a measurable habit icon: %w", err)
		}

		errs.Add("icon", domainError)
	}

	step, err := habits.NewMeasurementStep(command.Step)
	if err != nil {
		var domainError domain.Error
		if !errors.As(err, &domainError) {
			return uuid.Nil, fmt.Errorf("creating a measurable habit step: %w", err)
		}

		errs.Add("step", domainError)
	}

	unit, err := habits.NewMeasurableUnit(command.Unit)
	if err != nil {
		var domainError domain.Error
		if !errors.As(err, &domainError) {
			return uuid.Nil, fmt.Errorf("creating a measurable habit unit: %w", err)
		}

		errs.Add("unit", domainError)
	}

	if errs.HasErrors() {
		return uuid.Nil, errs
	}

	habit, err := handler.measurableHabits.Find(ctx, habits.NewMeasurableHabitFilter().
		WithAccountIds(command.AccountID).
		WithIds(command.ID).
		WithDeleted(false),
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("finding a measurable habit: %w", err)
	}

	if habit == nil {
		return uuid.Nil, habits.ErrHabitNotFound
	}

	if err := handler.modificationPolicy.Check(command.AccountID, habit); err != nil {
		return uuid.Nil, err
	}

	habit.ChangeLabel(label)
	habit.ChangeIcon(icon)
	habit.ChangeStep(step)
	habit.ChangeUnit(unit)

	if err := handler.measurableHabits.Save(ctx, habit); err != nil {
		return uuid.Nil, fmt.Errorf("saving a measurable habit: %w", err)
	}

	return habit.ID(), nil
}

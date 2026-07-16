package habits

import (
	"errors"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/google/uuid"
	"github.com/samborkent/uuidv7"
)

// MeasurableHabit represents a habit measured by numeric progress.
type MeasurableHabit struct {
	id         uuid.UUID
	label      string
	icon       Icon
	step       MeasurementStep
	unit       MeasurableUnit
	createdAt  time.Time
	updatedAt  time.Time
	archivedAt *time.Time
	deletedAt  *time.Time
	accountId  uuid.UUID
}

// NewMeasurableHabit returns a measurable habit with validated fields or domain validation errors.
func NewMeasurableHabit(rawLabel string, rawIcon int, rawStep float32, rawUnit string, accountId uuid.UUID) (*MeasurableHabit, error) {
	errs := domain.NewValidationErrors()

	label, err := NewLabel(rawLabel)
	if err != nil {
		var domainError domain.Error
		if !errors.As(err, &domainError) {
			return nil, fmt.Errorf("creating a measurable habit label: %w", err)
		}

		errs.Add("label", domainError)
	}

	icon, err := NewIcon(rawIcon)
	if err != nil {
		var domainError domain.Error
		if !errors.As(err, &domainError) {
			return nil, fmt.Errorf("creating a measurable habit icon: %w", err)
		}

		errs.Add("icon", domainError)
	}

	step, err := NewMeasurementStep(rawStep)
	if err != nil {
		var domainError domain.Error
		if !errors.As(err, &domainError) {
			return nil, fmt.Errorf("creating a measurable habit step: %w", err)
		}

		errs.Add("step", domainError)
	}

	unit, err := NewMeasurableUnit(rawUnit)
	if err != nil {
		var domainError domain.Error
		if !errors.As(err, &domainError) {
			return nil, fmt.Errorf("creating a measurable habit unit: %w", err)
		}

		errs.Add("unit", domainError)
	}

	if errs.HasErrors() {
		return nil, errs
	}

	now := time.Now()

	return &MeasurableHabit{
		id:        uuid.UUID(uuidv7.New()),
		label:     label,
		icon:      icon,
		step:      step,
		unit:      unit,
		createdAt: now,
		updatedAt: now,
		accountId: accountId,
	}, nil
}

// RestoreMeasurableHabit returns a measurable habit reconstructed from persisted values without validating or changing them.
func RestoreMeasurableHabit(id uuid.UUID, label string, icon int, step float32, unit string, createdAt time.Time, updatedAt time.Time, archivedAt *time.Time, deletedAt *time.Time, accountId uuid.UUID) *MeasurableHabit {
	return &MeasurableHabit{
		id:         id,
		label:      label,
		icon:       Icon(icon),
		step:       MeasurementStep(step),
		unit:       MeasurableUnit(unit),
		createdAt:  createdAt,
		updatedAt:  updatedAt,
		archivedAt: archivedAt,
		deletedAt:  deletedAt,
		accountId:  accountId,
	}
}

// ID returns the measurable habit identifier.
func (habit *MeasurableHabit) ID() uuid.UUID { return habit.id }

// Label returns the measurable habit label.
func (habit *MeasurableHabit) Label() string { return habit.label }

// ChangeLabel updates the measurable habit label.
func (habit *MeasurableHabit) ChangeLabel(label string) {
	habit.label = label
	habit.updatedAt = time.Now()
}

// Icon returns the measurable habit icon.
func (habit *MeasurableHabit) Icon() Icon { return habit.icon }

// ChangeIcon updates the measurable habit icon.
func (habit *MeasurableHabit) ChangeIcon(icon Icon) {
	habit.icon = icon
	habit.updatedAt = time.Now()
}

// Step returns the measurable habit step.
func (habit *MeasurableHabit) Step() MeasurementStep { return habit.step }

// ChangeStep updates the measurable habit step.
func (habit *MeasurableHabit) ChangeStep(step MeasurementStep) {
	habit.step = step
	habit.updatedAt = time.Now()
}

// Unit returns the measurable habit unit.
func (habit *MeasurableHabit) Unit() MeasurableUnit { return habit.unit }

// ChangeUnit updates the measurable habit unit.
func (habit *MeasurableHabit) ChangeUnit(unit MeasurableUnit) {
	habit.unit = unit
	habit.updatedAt = time.Now()
}

// CreatedAt returns the time when the measurable habit was created.
func (habit *MeasurableHabit) CreatedAt() time.Time { return habit.createdAt }

// UpdatedAt returns the time when the measurable habit was last updated.
func (habit *MeasurableHabit) UpdatedAt() time.Time { return habit.updatedAt }

// ArchivedAt returns the time when the measurable habit was archived or nil when it is active.
func (habit *MeasurableHabit) ArchivedAt() *time.Time { return habit.archivedAt }

// DeletedAt returns the time when the measurable habit was deleted or nil when it is not deleted.
func (habit *MeasurableHabit) DeletedAt() *time.Time { return habit.deletedAt }

// AccountId returns the identifier of the account that owns the measurable habit.
func (habit *MeasurableHabit) AccountId() uuid.UUID { return habit.accountId }

// Archive marks the measurable habit as archived and updates the modification time.
func (habit *MeasurableHabit) Archive() {
	now := time.Now()
	habit.archivedAt = &now
	habit.updatedAt = now
}

// Unarchive marks the measurable habit as active and updates the modification time.
func (habit *MeasurableHabit) Unarchive() {
	habit.archivedAt = nil
	habit.updatedAt = time.Now()
}

// Delete marks the measurable habit as deleted and updates the modification time.
func (habit *MeasurableHabit) Delete() {
	now := time.Now()
	habit.deletedAt = &now
	habit.updatedAt = now
}

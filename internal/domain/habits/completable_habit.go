package habits

import (
	"errors"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/google/uuid"
	"github.com/samborkent/uuidv7"
)

// CompletableHabit represents a habit that is either completed or not completed for a day.
type CompletableHabit struct {
	id         uuid.UUID
	label      string
	icon       Icon
	createdAt  time.Time
	updatedAt  time.Time
	archivedAt *time.Time
	deletedAt  *time.Time
	accountId  uuid.UUID
}

// NewCompletableHabit returns a completable habit with validated fields or domain validation errors.
func NewCompletableHabit(rawLabel string, rawIcon int, accountId uuid.UUID) (*CompletableHabit, error) {
	violations := domain.NewViolations()

	label, err := NewLabel(rawLabel)
	if err != nil {
		var violation domain.Violation
		if !errors.As(err, &violation) {
			return nil, fmt.Errorf("creating a completable habit label: %w", err)
		}

		violations.Add("label", violation)
	}

	icon, err := NewIcon(rawIcon)
	if err != nil {
		var violation domain.Violation
		if !errors.As(err, &violation) {
			return nil, fmt.Errorf("creating a completable habit icon: %w", err)
		}

		violations.Add("icon", violation)
	}

	if violations.HasViolations() {
		return nil, violations
	}

	now := time.Now()

	return &CompletableHabit{
		id:        uuid.UUID(uuidv7.New()),
		label:     label,
		icon:      icon,
		createdAt: now,
		updatedAt: now,
		accountId: accountId,
	}, nil
}

// RestoreCompletableHabit returns a completable habit reconstructed from persisted values without validating or changing them.
func RestoreCompletableHabit(id uuid.UUID, label string, icon int, createdAt time.Time, updatedAt time.Time, archivedAt *time.Time, deletedAt *time.Time, accountId uuid.UUID) *CompletableHabit {
	return &CompletableHabit{
		id:         id,
		label:      label,
		icon:       Icon(icon),
		createdAt:  createdAt,
		updatedAt:  updatedAt,
		archivedAt: archivedAt,
		deletedAt:  deletedAt,
		accountId:  accountId,
	}
}

// ID returns the completable habit identifier.
func (habit *CompletableHabit) ID() uuid.UUID {
	return habit.id
}

// Label returns the completable habit label.
func (habit *CompletableHabit) Label() string {
	return habit.label
}

// ChangeLabel updates the completable habit label.
func (habit *CompletableHabit) ChangeLabel(label string) {
	habit.label = label
	habit.updatedAt = time.Now()
}

// Icon returns the completable habit icon.
func (habit *CompletableHabit) Icon() Icon {
	return habit.icon
}

// ChangeIcon updates the completable habit icon.
func (habit *CompletableHabit) ChangeIcon(icon Icon) {
	habit.icon = icon
	habit.updatedAt = time.Now()
}

// CreatedAt returns the time when the completable habit was created.
func (habit *CompletableHabit) CreatedAt() time.Time {
	return habit.createdAt
}

// UpdatedAt returns the time when the completable habit was last updated.
func (habit *CompletableHabit) UpdatedAt() time.Time {
	return habit.updatedAt
}

// ArchivedAt returns the time when the completable habit was archived or nil when it is active.
func (habit *CompletableHabit) ArchivedAt() *time.Time {
	return habit.archivedAt
}

// DeletedAt returns the time when the completable habit was deleted or nil when it is not deleted.
func (habit *CompletableHabit) DeletedAt() *time.Time {
	return habit.deletedAt
}

// AccountId returns the identifier of the account that owns the completable habit.
func (habit *CompletableHabit) AccountId() uuid.UUID {
	return habit.accountId
}

// Archive marks the completable habit as archived and updates the modification time.
func (habit *CompletableHabit) Archive() {
	now := time.Now()

	habit.archivedAt = &now
	habit.updatedAt = now
}

// Unarchive marks the completable habit as active and updates the modification time.
func (habit *CompletableHabit) Unarchive() {
	habit.archivedAt = nil
	habit.updatedAt = time.Now()
}

// Delete marks the completable habit as deleted and updates the modification time.
func (habit *CompletableHabit) Delete() {
	now := time.Now()

	habit.deletedAt = &now
	habit.updatedAt = now
}

package habits

import (
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/google/uuid"
	"github.com/samborkent/uuidv7"
)

// TimeHabit represents a habit that is associated with time of day.
type TimeHabit struct {
	id         uuid.UUID
	label      string
	icon       Icon
	createdAt  time.Time
	updatedAt  time.Time
	archivedAt *time.Time
	deletedAt  *time.Time
	accountId  uuid.UUID
}

// NewTimeHabit returns a time habit with validated fields or domain validation violations.
func NewTimeHabit(rawLabel string, rawIcon int, accountId uuid.UUID) (*TimeHabit, domain.Violations) {
	violations := domain.NewViolations()

	label, violation := NewLabel(rawLabel)
	if violation != nil {
		violations.Add("label", violation)
	}

	icon, violation := NewIcon(rawIcon)
	if violation != nil {
		violations.Add("icon", violation)
	}

	if violations.HasViolations() {
		return nil, violations
	}

	now := time.Now()

	return &TimeHabit{
		id:        uuid.UUID(uuidv7.New()),
		label:     label,
		icon:      icon,
		createdAt: now,
		updatedAt: now,
		accountId: accountId,
	}, nil
}

// RestoreTimeHabit returns a time habit reconstructed from persisted values without validating or changing them.
func RestoreTimeHabit(id uuid.UUID, label string, icon int, createdAt time.Time, updatedAt time.Time, archivedAt *time.Time, deletedAt *time.Time, accountId uuid.UUID) *TimeHabit {
	return &TimeHabit{
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

// ID returns the time habit identifier.
func (habit *TimeHabit) ID() uuid.UUID {
	return habit.id
}

// Label returns the time habit label.
func (habit *TimeHabit) Label() string {
	return habit.label
}

// ChangeLabel updates the time habit label.
func (habit *TimeHabit) ChangeLabel(label string) {
	habit.label = label
	habit.updatedAt = time.Now()
}

// Icon returns the time habit icon.
func (habit *TimeHabit) Icon() Icon {
	return habit.icon
}

// ChangeIcon updates the time habit icon.
func (habit *TimeHabit) ChangeIcon(icon Icon) {
	habit.icon = icon
	habit.updatedAt = time.Now()
}

// CreatedAt returns the time when the time habit was created.
func (habit *TimeHabit) CreatedAt() time.Time {
	return habit.createdAt
}

// UpdatedAt returns the time when the time habit was last updated.
func (habit *TimeHabit) UpdatedAt() time.Time {
	return habit.updatedAt
}

// ArchivedAt returns the time when the time habit was archived or nil when it is active.
func (habit *TimeHabit) ArchivedAt() *time.Time {
	return habit.archivedAt
}

// DeletedAt returns the time when the time habit was deleted or nil when it is not deleted.
func (habit *TimeHabit) DeletedAt() *time.Time {
	return habit.deletedAt
}

// AccountId returns the identifier of the account that owns the time habit.
func (habit *TimeHabit) AccountId() uuid.UUID {
	return habit.accountId
}

// Archive marks the time habit as archived and updates the modification time.
func (habit *TimeHabit) Archive() {
	now := time.Now()

	habit.archivedAt = &now
	habit.updatedAt = now
}

// Unarchive marks the time habit as active and updates the modification time.
func (habit *TimeHabit) Unarchive() {
	habit.archivedAt = nil
	habit.updatedAt = time.Now()
}

// Delete marks the time habit as deleted and updates the modification time.
func (habit *TimeHabit) Delete() {
	now := time.Now()

	habit.deletedAt = &now
	habit.updatedAt = now
}

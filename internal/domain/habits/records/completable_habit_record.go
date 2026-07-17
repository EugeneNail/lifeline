package records

import (
	"time"

	"github.com/google/uuid"
)

// CompletableHabitRecord represents a daily value for a completable habit.
type CompletableHabitRecord struct {
	completableHabitId uuid.UUID
	accountId          uuid.UUID
	date               Date
	value              bool
}

// NewCompletableHabitRecord returns a completable habit record with immutable habit, account, and date fields.
func NewCompletableHabitRecord(completableHabitId uuid.UUID, accountId uuid.UUID, date time.Time, value bool) *CompletableHabitRecord {
	return &CompletableHabitRecord{
		completableHabitId: completableHabitId,
		accountId:          accountId,
		date:               NewDate(date),
		value:              value,
	}
}

// RestoreCompletableHabitRecord returns a completable habit record reconstructed from persisted values without validating or changing them.
func RestoreCompletableHabitRecord(completableHabitId uuid.UUID, accountId uuid.UUID, date Date, value bool) *CompletableHabitRecord {
	return &CompletableHabitRecord{
		completableHabitId: completableHabitId,
		accountId:          accountId,
		date:               date,
		value:              value,
	}
}

// CompletableHabitId returns the identifier of the related completable habit.
func (record *CompletableHabitRecord) CompletableHabitId() uuid.UUID {
	return record.completableHabitId
}

// AccountId returns the identifier of the account that owns the related habit.
func (record *CompletableHabitRecord) AccountId() uuid.UUID {
	return record.accountId
}

// Date returns the record date.
func (record *CompletableHabitRecord) Date() Date {
	return record.date
}

// Value returns the record value.
func (record *CompletableHabitRecord) Value() bool {
	return record.value
}

// ChangeValue updates the record value.
func (record *CompletableHabitRecord) ChangeValue(value bool) {
	record.value = value
}

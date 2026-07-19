package moods

import (
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/google/uuid"
)

// Record represents the daily mood selected by a user.
type Record struct {
	date      domain.Date
	value     Mood
	createdAt time.Time
	updatedAt time.Time
	accountId uuid.UUID
}

// NewRecord returns a mood record with the provided validated values.
func NewRecord(date domain.Date, value Mood, accountId uuid.UUID) *Record {
	now := time.Now()
	return &Record{
		date:      date,
		value:     value,
		createdAt: now,
		updatedAt: now,
		accountId: accountId,
	}
}

// RestoreRecord returns a mood record reconstructed from persisted values without validating them.
func RestoreRecord(date domain.Date, value Mood, createdAt time.Time, updatedAt time.Time, accountId uuid.UUID) *Record {
	return &Record{
		date:      date,
		value:     value,
		createdAt: createdAt,
		updatedAt: updatedAt,
		accountId: accountId,
	}
}

// Date returns the record date.
func (record *Record) Date() domain.Date {
	return record.date
}

// Value returns the stored mood value.
func (record *Record) Value() Mood {
	return record.value
}

// ChangeValue updates the stored mood value and refreshes the modification timestamp.
func (record *Record) ChangeValue(value Mood) {
	record.value = value
	record.updatedAt = time.Now()
}

// CreatedAt returns the timestamp when the mood record was created.
func (record *Record) CreatedAt() time.Time {
	return record.createdAt
}

// UpdatedAt returns the timestamp when the mood record was last updated.
func (record *Record) UpdatedAt() time.Time {
	return record.updatedAt
}

// AccountId returns the account identifier that owns the mood record.
func (record *Record) AccountId() uuid.UUID {
	return record.accountId
}

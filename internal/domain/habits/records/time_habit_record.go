package records

import (
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/google/uuid"
)

// TimeHabitRecord represents a daily value for a time habit.
type TimeHabitRecord struct {
	timeHabitId uuid.UUID
	accountId   uuid.UUID
	date        Date
	value       TimeValue
}

// NewTimeHabitRecord returns a time habit record with immutable habit, account, and date fields and a validated time value or domain validation violations.
func NewTimeHabitRecord(timeHabitId uuid.UUID, accountId uuid.UUID, date time.Time, rawValue int) (*TimeHabitRecord, domain.Violations) {
	violations := domain.NewViolations()

	value, violation := NewTimeValue(rawValue)
	if violation != nil {
		violations.Add("value", violation)
	}

	if violations.HasViolations() {
		return nil, violations
	}

	return &TimeHabitRecord{
		timeHabitId: timeHabitId,
		accountId:   accountId,
		date:        NewDate(date),
		value:       value,
	}, nil
}

// RestoreTimeHabitRecord returns a time habit record reconstructed from persisted values without validating or changing them.
func RestoreTimeHabitRecord(timeHabitId uuid.UUID, accountId uuid.UUID, date Date, value TimeValue) *TimeHabitRecord {
	return &TimeHabitRecord{
		timeHabitId: timeHabitId,
		accountId:   accountId,
		date:        date,
		value:       value,
	}
}

// TimeHabitId returns the identifier of the related time habit.
func (record *TimeHabitRecord) TimeHabitId() uuid.UUID {
	return record.timeHabitId
}

// AccountId returns the identifier of the account that owns the related habit.
func (record *TimeHabitRecord) AccountId() uuid.UUID {
	return record.accountId
}

// Date returns the record date.
func (record *TimeHabitRecord) Date() Date {
	return record.date
}

// Value returns the record value.
func (record *TimeHabitRecord) Value() TimeValue {
	return record.value
}

// ChangeValue updates the record value after validating the provided raw minute count.
func (record *TimeHabitRecord) ChangeValue(rawValue int) domain.Violation {
	value, violation := NewTimeValue(rawValue)
	if violation != nil {
		return violation
	}

	record.value = value

	return nil
}

package records

import (
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/google/uuid"
)

// MeasurableHabitRecord represents a daily value for a measurable habit.
type MeasurableHabitRecord struct {
	measurableHabitId uuid.UUID
	accountId         uuid.UUID
	date              Date
	value             MeasurableValue
}

// NewMeasurableHabitRecord returns a measurable habit record with immutable habit, account, and date fields and a validated numeric value or domain validation violations.
func NewMeasurableHabitRecord(measurableHabitId uuid.UUID, accountId uuid.UUID, date time.Time, rawValue float32, step habits.MeasurementStep) (*MeasurableHabitRecord, domain.Violations) {
	violations := domain.NewViolations()

	value, violation := NewMeasurableValue(rawValue, step)
	if violation != nil {
		violations.Add("value", violation)
	}

	if violations.HasViolations() {
		return nil, violations
	}

	return &MeasurableHabitRecord{
		measurableHabitId: measurableHabitId,
		accountId:         accountId,
		date:              NewDate(date),
		value:             value,
	}, nil
}

// RestoreMeasurableHabitRecord returns a measurable habit record reconstructed from persisted values without validating or changing them.
func RestoreMeasurableHabitRecord(measurableHabitId uuid.UUID, accountId uuid.UUID, date Date, value MeasurableValue) *MeasurableHabitRecord {
	return &MeasurableHabitRecord{
		measurableHabitId: measurableHabitId,
		accountId:         accountId,
		date:              date,
		value:             value,
	}
}

// MeasurableHabitId returns the identifier of the related measurable habit.
func (record *MeasurableHabitRecord) MeasurableHabitId() uuid.UUID {
	return record.measurableHabitId
}

// AccountId returns the identifier of the account that owns the related habit.
func (record *MeasurableHabitRecord) AccountId() uuid.UUID {
	return record.accountId
}

// Date returns the record date.
func (record *MeasurableHabitRecord) Date() Date {
	return record.date
}

// Value returns the record value.
func (record *MeasurableHabitRecord) Value() MeasurableValue {
	return record.value
}

// ChangeValue updates the record value after validating the provided raw numeric value against the step.
func (record *MeasurableHabitRecord) ChangeValue(rawValue float32, step habits.MeasurementStep) domain.Violation {
	value, violation := NewMeasurableValue(rawValue, step)
	if violation != nil {
		return violation
	}

	record.value = value

	return nil
}

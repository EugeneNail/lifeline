package records

import (
	"errors"
	"fmt"
	"github.com/EugeneNail/lifeline/internal/domain"
	"time"

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

// NewMeasurableHabitRecord returns a measurable habit record with immutable habit, account, and date fields and a validated numeric value.
func NewMeasurableHabitRecord(measurableHabitId uuid.UUID, accountId uuid.UUID, date time.Time, rawValue float32, step habits.MeasurementStep) (*MeasurableHabitRecord, error) {
	errs := domain.NewValidationErrors()

	value, err := NewMeasurableValue(rawValue, step)
	if err != nil {
		var domainError domain.Error
		if !errors.As(err, &domainError) {
			return nil, fmt.Errorf("creating a measurable value: %w", err)
		}

		errs.Add("value", domainError)
	}

	if errs.HasErrors() {
		return nil, errs
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
func (record *MeasurableHabitRecord) ChangeValue(rawValue float32, step habits.MeasurementStep) domain.Error {
	value, err := NewMeasurableValue(rawValue, step)
	if err != nil {
		return err
	}

	record.value = value

	return nil
}

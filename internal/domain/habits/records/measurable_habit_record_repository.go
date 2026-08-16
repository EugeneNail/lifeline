package records

import "context"

// MeasurableHabitRecordRepository stores and retrieves measurable habit records.
type MeasurableHabitRecordRepository interface {
	// FindMany returns all measurable habit records matching the filter or an error when lookup fails.
	FindMany(ctx context.Context, filter MeasurableHabitRecordFilter) ([]*MeasurableHabitRecord, error)

	// Find returns the first measurable habit record matching the filter, nil when none exists, or an error when lookup fails.
	Find(ctx context.Context, filter MeasurableHabitRecordFilter) (*MeasurableHabitRecord, error)

	// Add stores a measurable habit record or returns an error when persistence fails.
	Add(ctx context.Context, record *MeasurableHabitRecord) error

	// Save updates a measurable habit record or returns an error when persistence fails.
	Save(ctx context.Context, record *MeasurableHabitRecord) error
}

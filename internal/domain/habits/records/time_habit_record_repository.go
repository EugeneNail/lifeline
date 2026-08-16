package records

import "context"

// TimeHabitRecordRepository stores and retrieves time habit records.
type TimeHabitRecordRepository interface {
	// FindMany returns all time habit records matching the filter or an error when lookup fails.
	FindMany(ctx context.Context, filter TimeHabitRecordFilter) ([]*TimeHabitRecord, error)

	// Find returns the first time habit record matching the filter, nil when none exists, or an error when lookup fails.
	Find(ctx context.Context, filter TimeHabitRecordFilter) (*TimeHabitRecord, error)

	// Add stores a time habit record or returns an error when persistence fails.
	Add(ctx context.Context, record *TimeHabitRecord) error

	// Save updates a time habit record or returns an error when persistence fails.
	Save(ctx context.Context, record *TimeHabitRecord) error
}

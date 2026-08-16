package records

import "context"

// CompletableHabitRecordRepository stores and retrieves completable habit records.
type CompletableHabitRecordRepository interface {
	// FindMany returns all completable habit records matching the filter or an error when lookup fails.
	FindMany(ctx context.Context, filter CompletableHabitRecordFilter) ([]*CompletableHabitRecord, error)

	// Find returns the first completable habit record matching the filter, nil when none exists, or an error when lookup fails.
	Find(ctx context.Context, filter CompletableHabitRecordFilter) (*CompletableHabitRecord, error)

	// Add stores a completable habit record or returns an error when persistence fails.
	Add(ctx context.Context, record *CompletableHabitRecord) error

	// Save updates a completable habit record or returns an error when persistence fails.
	Save(ctx context.Context, record *CompletableHabitRecord) error
}

package records

import (
	"context"
	"time"
)

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

func truncateRecordDates(dates ...time.Time) []time.Time {
	truncated := make([]time.Time, 0, len(dates))
	for _, date := range dates {
		truncated = append(truncated, date.Truncate(time.Hour*24))
	}

	return truncated
}

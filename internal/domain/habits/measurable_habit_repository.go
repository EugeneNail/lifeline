package habits

import "context"

// MeasurableHabitRepository stores and retrieves measurable habits.
type MeasurableHabitRepository interface {
	// Add stores a measurable habit or returns an error when persistence fails.
	Add(ctx context.Context, habit *MeasurableHabit) error

	// Count returns the number of measurable habits matching the filter or an error when counting fails.
	Count(ctx context.Context, filter MeasurableHabitFilter) (int, error)

	// Find returns the first measurable habit matching the filter, nil when none exists, or an error when lookup fails.
	Find(ctx context.Context, filter MeasurableHabitFilter) (*MeasurableHabit, error)

	// FindMany returns all measurable habits matching the filter or an error when lookup fails.
	FindMany(ctx context.Context, filter MeasurableHabitFilter) ([]*MeasurableHabit, error)
}

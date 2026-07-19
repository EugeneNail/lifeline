package habits

import "context"

// TimeHabitRepository stores and retrieves time habits.
type TimeHabitRepository interface {
	// Add stores a time habit in storage or returns an error when persistence fails.
	Add(ctx context.Context, habit *TimeHabit) error

	// Count returns the number of time habits matching the filter or an error when counting fails.
	Count(ctx context.Context, filter TimeHabitFilter) (int, error)

	// Find returns the first time habit matching the filter, nil when none exists, or an error when lookup fails.
	Find(ctx context.Context, filter TimeHabitFilter) (*TimeHabit, error)

	// FindMany returns all time habits matching the filter or an error when lookup fails.
	FindMany(ctx context.Context, filter TimeHabitFilter) ([]*TimeHabit, error)

	// Update updates a time habit in storage or returns an error when persistence fails.
	Update(ctx context.Context, habit *TimeHabit) error
}

package habits

import "context"

// CompletableHabitRepository stores and retrieves completable habits.
type CompletableHabitRepository interface {
	// Add stores a completable habit in storage or returns an error when persistence fails.
	Add(ctx context.Context, habit *CompletableHabit) error
	// Count returns the number of completable habits matching the filter or an error when counting fails.
	Count(ctx context.Context, filter CompletableHabitFilter) (int, error)
	// Find returns the first completable habit matching the filter, nil when none exists, or an error when lookup fails.
	Find(ctx context.Context, filter CompletableHabitFilter) (*CompletableHabit, error)
	// FindMany returns all completable habits matching the filter or an error when lookup fails.
	FindMany(ctx context.Context, filter CompletableHabitFilter) ([]*CompletableHabit, error)
	// Save updates a completable habit in storage or returns an error when persistence fails.
	Save(ctx context.Context, habit *CompletableHabit) error
}

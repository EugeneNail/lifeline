package habits

import "context"

// CompletableHabitRepository stores and retrieves completable habits.
type CompletableHabitRepository interface {
	Add(ctx context.Context, habit *CompletableHabit) error
	Count(ctx context.Context, filter CompletableHabitFilter) (int, error)
	Find(ctx context.Context, filter CompletableHabitFilter) (*CompletableHabit, error)
	FindMany(ctx context.Context, filter CompletableHabitFilter) ([]*CompletableHabit, error)
}

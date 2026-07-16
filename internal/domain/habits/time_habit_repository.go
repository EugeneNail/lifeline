package habits

import "context"

// TimeHabitRepository stores and retrieves time habits.
type TimeHabitRepository interface {
	Add(ctx context.Context, habit *TimeHabit) error
	Count(ctx context.Context, filter TimeHabitFilter) (int, error)
	Find(ctx context.Context, filter TimeHabitFilter) (*TimeHabit, error)
	FindMany(ctx context.Context, filter TimeHabitFilter) ([]*TimeHabit, error)
}

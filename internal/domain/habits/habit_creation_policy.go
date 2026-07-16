package habits

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const activeHabitLimit = 50

// HabitCreationPolicy checks whether a habit can be created.
type HabitCreationPolicy struct {
	completableHabits CompletableHabitRepository
	measurableHabits  MeasurableHabitRepository
}

// NewHabitCreationPolicy returns a habit creation policy configured with the completable and measurable habit repositories.
func NewHabitCreationPolicy(
	completableHabits CompletableHabitRepository,
	measurableHabits MeasurableHabitRepository,
) *HabitCreationPolicy {
	return &HabitCreationPolicy{
		completableHabits: completableHabits,
		measurableHabits:  measurableHabits,
	}
}

// EnsureCanAdd returns nil when the account has not exceeded the active habit limit, ErrHabitLimitExceeded when the limit is exceeded, or an error when counting habits fails.
func (policy *HabitCreationPolicy) EnsureCanAdd(ctx context.Context, accountId uuid.UUID) error {
	completableFilter := NewCompletableHabitFilter().
		WithAccountIds(accountId).
		WithArchived(false).
		WithDeleted(false)

	completableCount, err := policy.completableHabits.Count(ctx, completableFilter)
	if err != nil {
		return fmt.Errorf("counting active completable habits for account id %q: %w", accountId, err)
	}

	measurableFilter := NewMeasurableHabitFilter().
		WithAccountIds(accountId).
		WithArchived(false).
		WithDeleted(false)

	measurableCount, err := policy.measurableHabits.Count(ctx, measurableFilter)
	if err != nil {
		return fmt.Errorf("counting active measurable habits for account id %q: %w", accountId, err)
	}

	count := completableCount + measurableCount
	if count >= activeHabitLimit {
		return ErrHabitLimitExceeded
	}

	return nil
}

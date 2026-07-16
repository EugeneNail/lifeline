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
}

// NewHabitCreationPolicy returns a habit creation policy configured with the completable habit repository.
func NewHabitCreationPolicy(completableHabits CompletableHabitRepository) *HabitCreationPolicy {
	return &HabitCreationPolicy{completableHabits: completableHabits}
}

// EnsureCanAdd returns nil when the account has not exceeded the active habit limit, ErrHabitLimitExceeded when the limit is exceeded, or an error when counting habits fails.
func (policy *HabitCreationPolicy) EnsureCanAdd(ctx context.Context, accountId uuid.UUID) error {
	filter := NewCompletableHabitFilter().
		WithAccountIds(accountId).
		WithArchived(false).
		WithDeleted(false)

	count, err := policy.completableHabits.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("counting active completable habits for account id %q: %w", accountId, err)
	}

	if count >= activeHabitLimit {
		return ErrHabitLimitExceeded
	}

	return nil
}

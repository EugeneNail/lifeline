package journal

import (
	"context"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/auth"
)

// CreationPolicy checks whether a daily journal can be created.
type CreationPolicy struct {
	journals JournalRepository
}

// NewCreationPolicy returns a journal creation policy configured with the journal repository.
func NewCreationPolicy(repository JournalRepository) *CreationPolicy {
	return &CreationPolicy{journals: repository}
}

// Check returns nil when the account has no journal for the date, ErrDateIsOccupied when one exists, or an error when the journal lookup fails.
func (policy *CreationPolicy) Check(ctx context.Context, accountID auth.ID, date time.Time) error {
	date = date.Truncate(time.Hour * 24)
	filter := NewFilter().
		WithAccountIds(accountID).
		WithDates(date)

	journalEntry, err := policy.journals.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("finding journal by account id %q and date %q: %w", accountID, date, err)
	}

	if journalEntry != nil {
		return ErrDateIsOccupied
	}

	return nil
}

package journal

import (
	"context"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/auth"
)

// JournalCreationPolicy checks whether a daily journal can be created.
type JournalCreationPolicy struct {
	journals JournalRepository
}

// NewJournalCreationPolicy returns a journal creation policy configured with the journal repository.
func NewJournalCreationPolicy(repository JournalRepository) *JournalCreationPolicy {
	return &JournalCreationPolicy{journals: repository}
}

// EnsureCanAdd returns nil when the account has no journal for the date, ErrDateIsOccupied when one exists, or an error when the journal lookup fails.
func (policy *JournalCreationPolicy) EnsureCanAdd(ctx context.Context, accountID auth.ID, date time.Time) error {
	date = date.Truncate(time.Hour * 24)
	filter := NewJournalFilter().
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

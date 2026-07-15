package entries

import (
	"context"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/auth"
)

// EntryCreationPolicy checks whether a daily entry can be created.
type EntryCreationPolicy struct {
	entries EntryRepository
}

// NewEntryCreationPolicy returns an entry creation policy configured with the entry entries.
func NewEntryCreationPolicy(repository EntryRepository) *EntryCreationPolicy {
	return &EntryCreationPolicy{entries: repository}
}

// EnsureCanAddEntry returns nil when the account has no entry for the date, ErrDateIsOccupied when one exists, or an error when the entries lookup fails.
func (policy *EntryCreationPolicy) EnsureCanAddEntry(ctx context.Context, accountID auth.ID, date time.Time) error {
	date = date.Truncate(time.Hour * 24)
	filter := NewEntryFilter().
		WithAccountIds(accountID).
		WithDates(date)

	entry, err := policy.entries.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("finding entry by account id %q and date %q: %w", accountID, date, err)
	}

	if entry != nil {
		return ErrDateIsOccupied
	}

	return nil
}

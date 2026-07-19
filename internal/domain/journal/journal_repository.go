package journal

import "context"

// JournalRepository stores and retrieves daily journals.
type JournalRepository interface {
	Add(ctx context.Context, journal *Journal) error
	Find(ctx context.Context, filter JournalFilter) (*Journal, error)
}

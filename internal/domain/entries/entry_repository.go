package entries

import "context"

// EntryRepository stores and retrieves daily entries.
type EntryRepository interface {
	Add(ctx context.Context, entry *Entry) error
	Find(ctx context.Context, filter EntryFilter) (*Entry, error)
}

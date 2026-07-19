package moods

import "context"

// RecordRepository stores and retrieves mood records.
type RecordRepository interface {
	// FindMany returns all mood records matching the filter or an error when lookup fails.
	FindMany(ctx context.Context, filter RecordFilter) ([]*Record, error)

	// Find returns the first mood record matching the filter, nil when none exists, or an error when lookup fails.
	Find(ctx context.Context, filter RecordFilter) (*Record, error)

	// Add stores a mood record or returns an error when persistence fails.
	Add(ctx context.Context, record *Record) error

	// Update updates a mood record or returns an error when persistence fails.
	Update(ctx context.Context, record *Record) error
}

package journals

import "context"

// Repository stores and retrieves daily journals.
type Repository interface {
	// Add stores a journal in storage or returns an error when persistence fails.
	Add(ctx context.Context, journal *Journal) error

	// Find returns the first journal matching the filter, nil when none exists, or an error when lookup fails.
	Find(ctx context.Context, filter Filter) (*Journal, error)

	// FindMany returns all journals matching the filter or an error when lookup fails.
	FindMany(ctx context.Context, filter Filter) ([]*Journal, error)

	// Update stores the current journal state in storage or returns an error when persistence fails.
	Update(ctx context.Context, journal *Journal) error
}

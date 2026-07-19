package journals

import "context"

// Repository stores and retrieves daily journals.
type Repository interface {
	Add(ctx context.Context, journal *Journal) error
	Find(ctx context.Context, filter Filter) (*Journal, error)
	Update(ctx context.Context, journal *Journal) error
}

package transactions

import "context"

// Repository stores and retrieves transactions.
type Repository interface {
	// Add stores a transaction in storage or returns an error when persistence fails.
	Add(ctx context.Context, transaction *Transaction) error

	// Find returns the first transaction matching the filter, nil when none exists, or an error when lookup fails.
	Find(ctx context.Context, filter TransactionFilter) (*Transaction, error)

	// FindMany returns all transactions matching the filter or an error when lookup fails.
	FindMany(ctx context.Context, filter TransactionFilter) ([]*Transaction, error)
}

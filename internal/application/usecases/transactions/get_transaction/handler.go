package get_transaction

import (
	"context"
	"fmt"

	"github.com/EugeneNail/lifeline/internal/domain/transactions"
	"github.com/google/uuid"
)

// Handler executes the get-transaction use case.
type Handler struct {
	transactions transactions.Repository
}

// NewHandler returns a get-transaction handler configured with the transaction repository or an error when the dependency is missing.
func NewHandler(transactionsRepository transactions.Repository) (*Handler, error) {
	if transactionsRepository == nil {
		return nil, fmt.Errorf("get_transaction handler requires a transaction repository")
	}

	return &Handler{transactions: transactionsRepository}, nil
}

// Query carries the data required to load a transaction by identifier.
type Query struct {
	AccountID uuid.UUID
	ID        uuid.UUID
}

// Handle returns the transaction matching the query or nil when no transaction exists, or an error when lookup fails.
func (handler *Handler) Handle(ctx context.Context, query Query) (*transactions.Transaction, error) {
	transaction, err := handler.transactions.Find(ctx, transactions.NewTransactionFilter().
		WithAccountIds(query.AccountID).
		WithIds(query.ID),
	)
	if err != nil {
		return nil, fmt.Errorf("finding a transaction by account id %q and id %q: %w", query.AccountID, query.ID, err)
	}

	return transaction, nil
}

package list_transactions

import (
	"context"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/transactions"
	"github.com/google/uuid"
)

// Handler executes the list-transactions use case.
type Handler struct {
	transactions transactions.Repository
}

// NewHandler returns a list-transactions handler configured with the transaction repository or an error when the dependency is missing.
func NewHandler(transactionsRepository transactions.Repository) (*Handler, error) {
	if transactionsRepository == nil {
		return nil, fmt.Errorf("list_transactions handler requires a transaction repository")
	}

	return &Handler{transactions: transactionsRepository}, nil
}

// Query carries the data required to list transactions for a specific account.
type Query struct {
	AccountID uuid.UUID
	From      *time.Time
	To        *time.Time
}

// Result groups the user's transactions and returns them to the caller.
type Result struct {
	Transactions []*transactions.Transaction
}

// Handle loads the user's transactions and returns them to the caller or an error when lookup fails.
func (handler *Handler) Handle(ctx context.Context, query Query) (Result, error) {
	filter := transactions.NewTransactionFilter().
		WithAccountIds(query.AccountID)

	if query.From != nil {
		filter = filter.WithFrom(*query.From)
	}

	if query.To != nil {
		filter = filter.WithTo(*query.To)
	}

	foundTransactions, err := handler.transactions.FindMany(ctx, filter)
	if err != nil {
		return Result{}, fmt.Errorf("finding transactions: %w", err)
	}

	return Result{Transactions: foundTransactions}, nil
}

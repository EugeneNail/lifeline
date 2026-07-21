package create_transaction

import (
	"context"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/transactions"
	"github.com/google/uuid"
)

// Handler executes the create-transaction use case.
type Handler struct {
	transactions transactions.Repository
}

// NewHandler returns a create-transaction handler configured with the transaction repository or an error when the dependency is missing.
func NewHandler(transactionsRepository transactions.Repository) (*Handler, error) {
	if transactionsRepository == nil {
		return nil, fmt.Errorf("create_transaction handler requires a transaction repository")
	}

	return &Handler{transactions: transactionsRepository}, nil
}

// Command carries the data required to create a transaction.
type Command struct {
	Money       float32
	Date        time.Time
	Category    int
	Description string
	AccountID   uuid.UUID
}

// Handle validates the command, stores a new transaction, and returns the transaction identifier or field validation violations.
func (handler *Handler) Handle(ctx context.Context, command Command) (uuid.UUID, error) {
	transaction, violations := transactions.NewFromRaw(command.Money, command.Date, command.Category, command.Description, command.AccountID)
	if violations != nil {
		return uuid.Nil, violations
	}

	if err := handler.transactions.Add(ctx, transaction); err != nil {
		return uuid.Nil, fmt.Errorf("adding a new transaction to the collection: %w", err)
	}

	return transaction.ID(), nil
}

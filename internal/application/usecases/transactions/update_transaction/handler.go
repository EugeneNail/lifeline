package update_transaction

import (
	"context"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/transactions"
	"github.com/google/uuid"
)

// Handler executes the update-transaction use case.
type Handler struct {
	transactions transactions.Repository
}

// NewHandler returns an update-transaction handler configured with the transaction repository or an error when the dependency is missing.
func NewHandler(transactionsRepository transactions.Repository) (*Handler, error) {
	if transactionsRepository == nil {
		return nil, fmt.Errorf("update_transaction handler requires a transaction repository")
	}

	return &Handler{transactions: transactionsRepository}, nil
}

// Command carries the data required to update a transaction.
type Command struct {
	ID          uuid.UUID
	Money       float32
	Date        time.Time
	Direction   int
	Category    int
	Description string
	AccountID   uuid.UUID
}

// Handle validates the command, updates a transaction, and returns nil when the update succeeds or an error when it fails.
func (handler *Handler) Handle(ctx context.Context, command Command) error {
	transaction, err := handler.transactions.Find(ctx, transactions.NewTransactionFilter().
		WithAccountIds(command.AccountID).
		WithIds(command.ID),
	)
	if err != nil {
		return fmt.Errorf("finding a transaction by account id %q and id %q: %w", command.AccountID, command.ID, err)
	}

	if transaction == nil {
		return transactions.ErrTransactionNotFound
	}

	violations := domain.NewViolations()

	money, violation := transactions.NewMoney(command.Money)
	if violation != nil {
		violations.Add("money", violation)
	}

	date, violation := domain.NewDate(command.Date)
	if violation != nil {
		violations.Add("date", violation)
	}

	direction, violation := transactions.NewDirection(command.Direction)
	if violation != nil {
		violations.Add("direction", violation)
	}

	category, violation := transactions.NewCategory(command.Category)
	if violation != nil {
		violations.Add("category", violation)
	}

	description, violation := transactions.NewDescription(command.Description)
	if violation != nil {
		violations.Add("description", violation)
	}

	if violations.HasViolations() {
		return violations
	}

	transaction.ChangeMoney(money)
	transaction.ChangeDate(date)
	transaction.ChangeDirection(direction)
	transaction.ChangeCategory(category)
	transaction.ChangeDescription(description)

	if err := handler.transactions.Update(ctx, transaction); err != nil {
		return fmt.Errorf("updating a transaction %q: %w", transaction.ID(), err)
	}

	return nil
}

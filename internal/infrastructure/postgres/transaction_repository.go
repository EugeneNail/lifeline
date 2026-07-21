package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/transactions"
	"github.com/google/uuid"
)

// TransactionRepository stores transactions in PostgreSQL.
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository returns a PostgreSQL transaction repository.
func NewTransactionRepository(db *sql.DB) (*TransactionRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("TransactionRepository requires an sql.DB instance")
	}

	return &TransactionRepository{db: db}, nil
}

// Add stores the provided transaction in PostgreSQL.
func (repository *TransactionRepository) Add(ctx context.Context, transaction *transactions.Transaction) error {
	_, err := repository.db.ExecContext(
		ctx,
		`INSERT INTO transactions (id, money, date, category, description, account_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		transaction.ID(),
		float32(transaction.Money()),
		time.Time(transaction.Date()),
		int(transaction.Category()),
		string(transaction.Description()),
		transaction.AccountId(),
		transaction.CreatedAt(),
		transaction.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("executing an INSERT sql query for transaction %s: %w", transaction.ID(), err)
	}

	return nil
}

// Find returns the first transaction matching the provided filter or nil when no row exists.
func (repository *TransactionRepository) Find(ctx context.Context, filter transactions.TransactionFilter) (*transactions.Transaction, error) {
	foundTransactions, err := repository.FindMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("finding transactions: %w", err)
	}

	if len(foundTransactions) == 0 {
		return nil, nil
	}

	return foundTransactions[0], nil
}

// FindMany returns all transactions matching the provided filter.
func (repository *TransactionRepository) FindMany(ctx context.Context, filter transactions.TransactionFilter) ([]*transactions.Transaction, error) {
	query := `SELECT id, money, date, category, description, account_id, created_at, updated_at FROM transactions`
	conditions, args := repository.buildConditions(filter)

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	rows, err := repository.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing a SELECT sql query for transactions: %w", err)
	}
	defer rows.Close()

	foundTransactions := make([]*transactions.Transaction, 0)
	for rows.Next() {
		transaction, err := repository.scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning a transaction row: %w", err)
		}

		foundTransactions = append(foundTransactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating transaction rows: %w", err)
	}

	return foundTransactions, nil
}

// buildConditions converts the provided transaction filter into SQL WHERE fragments and arguments.
func (repository *TransactionRepository) buildConditions(filter transactions.TransactionFilter) ([]string, []any) {
	conditions := make([]string, 0, 3)
	args := make([]any, 0)

	if len(filter.AccountIds) > 0 {
		placeholders := make([]string, 0, len(filter.AccountIds))
		for _, accountID := range filter.AccountIds {
			args = append(args, accountID)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}

		conditions = append(conditions, fmt.Sprintf("account_id IN (%s)", strings.Join(placeholders, ", ")))
	}

	if len(filter.TransactionIds) > 0 {
		placeholders := make([]string, 0, len(filter.TransactionIds))
		for _, id := range filter.TransactionIds {
			args = append(args, id)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}

		conditions = append(conditions, fmt.Sprintf("id IN (%s)", strings.Join(placeholders, ", ")))
	}

	if len(filter.Dates) > 0 {
		placeholders := make([]string, 0, len(filter.Dates))
		for _, date := range filter.Dates {
			args = append(args, date)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}

		conditions = append(conditions, fmt.Sprintf("date IN (%s)", strings.Join(placeholders, ", ")))
	}

	return conditions, args
}

// scan converts the current SQL row into a transaction model or returns an error when reconstruction fails.
func (repository *TransactionRepository) scan(rows *sql.Rows) (*transactions.Transaction, error) {
	var (
		id          uuid.UUID
		money       float32
		rawDate     time.Time
		category    int
		description string
		accountID   uuid.UUID
		createdAt   time.Time
		updatedAt   time.Time
	)

	if err := rows.Scan(&id, &money, &rawDate, &category, &description, &accountID, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("scanning a SELECT sql result: %w", err)
	}

	return transactions.Restore(id, money, rawDate, category, description, accountID, createdAt, updatedAt), nil
}

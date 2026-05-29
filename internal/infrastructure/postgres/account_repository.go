package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"github.com/google/uuid"
)

// AccountRepository stores auth accounts in PostgreSQL.
type AccountRepository struct {
	db *sql.DB
}

// NewAccountRepository returns a PostgreSQL-backed account repository.
func NewAccountRepository(db *sql.DB) (*AccountRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("AccountRepository requires an sql.DB instance")
	}

	return &AccountRepository{db: db}, nil
}

// Add stores the provided account in PostgreSQL.
func (repository *AccountRepository) Add(ctx context.Context, account *auth.Account) error {
	_, err := repository.db.ExecContext(
		ctx,
		`INSERT INTO accounts (id, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		account.ID().Uuid(),
		account.Email().String(),
		account.Password().String(),
		account.CreatedAt,
		account.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("executing an INSERT sql query: %w", err)
	}

	return nil
}

// FindByEmail returns the account with the provided email or nil when no row exists.
func (repository *AccountRepository) FindByEmail(ctx context.Context, email auth.Email) (*auth.Account, error) {
	row := repository.db.QueryRowContext(
		ctx,
		`SELECT id, email, password, created_at, updated_at FROM accounts WHERE email = $1`,
		email.String(),
	)

	var (
		id          uuid.UUID
		storedEmail string
		password    string
		createdAt   time.Time
		updatedAt   time.Time
	)

	if err := row.Scan(&id, &storedEmail, &password, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("executing a SELECT sql query: %w", err)
	}

	account := auth.RestoreAccount(id, storedEmail, password, createdAt, updatedAt)

	return account, nil
}

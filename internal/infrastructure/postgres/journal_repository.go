package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/journals"
	"github.com/google/uuid"
)

// JournalRepository stores journals in PostgreSQL.
type JournalRepository struct {
	db *sql.DB
}

// NewJournalRepository returns a PostgreSQL journal repository.
func NewJournalRepository(db *sql.DB) (*JournalRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("Repository requires an sql.DB instance")
	}

	return &JournalRepository{db: db}, nil
}

// Add stores the provided journal in PostgreSQL.
func (repository *JournalRepository) Add(ctx context.Context, journalEntry *journals.Journal) error {
	_, err := repository.db.ExecContext(
		ctx,
		`INSERT INTO journals (date, note, created_at, updated_at, account_id) VALUES ($1, $2, $3, $4, $5)`,
		time.Time(journalEntry.Date()),
		string(journalEntry.Note()),
		journalEntry.CreatedAt(),
		journalEntry.UpdatedAt(),
		journalEntry.AccountId().Uuid(),
	)
	if err != nil {
		return fmt.Errorf("executing an INSERT sql query: %w", err)
	}

	return nil
}

// Update changes the mutable journal fields in PostgreSQL using the account and date identity.
func (repository *JournalRepository) Update(ctx context.Context, journalEntry *journals.Journal) error {
	result, err := repository.db.ExecContext(
		ctx,
		`UPDATE journals SET note = $1, updated_at = $2 WHERE account_id = $3 AND date = $4`,
		string(journalEntry.Note()),
		journalEntry.UpdatedAt(),
		journalEntry.AccountId().Uuid(),
		time.Time(journalEntry.Date()),
	)
	if err != nil {
		return fmt.Errorf("executing an UPDATE sql query for journal account %s and date %s: %w", journalEntry.AccountId().Uuid(), time.Time(journalEntry.Date()), err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking updated rows for journal account %s and date %s: %w", journalEntry.AccountId().Uuid(), time.Time(journalEntry.Date()), err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("saving journal account %s and date %s: no rows updated", journalEntry.AccountId().Uuid(), time.Time(journalEntry.Date()))
	}

	return nil
}

// Find returns the first journal matching the provided filter or nil when no row exists.
func (repository *JournalRepository) Find(ctx context.Context, filter journals.Filter) (*journals.Journal, error) {
	query := `SELECT date, note, created_at, updated_at, account_id FROM journals`
	conditions := make([]string, 0, 3)
	args := make([]any, 0)

	if len(filter.AccountIds) > 0 {
		placeholders := make([]string, 0, len(filter.AccountIds))
		for _, accountId := range filter.AccountIds {
			args = append(args, accountId.Uuid())
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}

		conditions = append(conditions, fmt.Sprintf("account_id IN (%s)", strings.Join(placeholders, ", ")))
	}

	if len(filter.Dates) > 0 {
		placeholders := make([]string, 0, len(filter.Dates))
		for _, date := range filter.Dates {
			args = append(args, date)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}

		conditions = append(conditions, fmt.Sprintf("date IN (%s)", strings.Join(placeholders, ", ")))
	}

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	query = fmt.Sprintf("%s LIMIT 1", query)

	row := repository.db.QueryRowContext(ctx, query, args...)

	var (
		date      time.Time
		note      string
		createdAt time.Time
		updatedAt time.Time
		accountId uuid.UUID
	)

	if err := row.Scan(&date, &note, &createdAt, &updatedAt, &accountId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("executing a SELECT sql query: %w", err)
	}

	journalEntry := journals.Restore(date, note, createdAt, updatedAt, accountId)

	return journalEntry, nil
}

package postgres

import (
	"context"
	"database/sql"
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
	foundJournals, err := repository.FindMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("finding journals: %w", err)
	}

	if len(foundJournals) == 0 {
		return nil, nil
	}

	return foundJournals[0], nil
}

// FindMany returns all journals matching the provided filter.
func (repository *JournalRepository) FindMany(ctx context.Context, filter journals.Filter) ([]*journals.Journal, error) {
	query := `SELECT date, note, created_at, updated_at, account_id FROM journals`
	conditions, args := buildJournalConditions(filter)

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	query = fmt.Sprintf("%s ORDER BY date DESC, created_at DESC", query)

	rows, err := repository.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing a SELECT sql query for journals: %w", err)
	}
	defer rows.Close()

	foundJournals := make([]*journals.Journal, 0)
	for rows.Next() {
		journalEntry, err := scanJournal(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning a journal row: %w", err)
		}

		foundJournals = append(foundJournals, journalEntry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating journal rows: %w", err)
	}

	return foundJournals, nil
}

// buildJournalConditions converts the provided journal filter into SQL WHERE fragments and arguments.
func buildJournalConditions(filter journals.Filter) ([]string, []any) {
	conditions := make([]string, 0, 2)
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

	return conditions, args
}

// scanJournal converts the current SQL row into a journal model or returns an error when reconstruction fails.
func scanJournal(rows *sql.Rows) (*journals.Journal, error) {
	var (
		date      time.Time
		note      string
		createdAt time.Time
		updatedAt time.Time
		accountId uuid.UUID
	)

	if err := rows.Scan(&date, &note, &createdAt, &updatedAt, &accountId); err != nil {
		return nil, fmt.Errorf("scanning a SELECT sql result: %w", err)
	}

	return journals.Restore(date, note, createdAt, updatedAt, accountId), nil
}

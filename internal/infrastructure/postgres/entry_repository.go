package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/entries"
	"github.com/google/uuid"
)

// EntryRepository stores entries in PostgreSQL.
type EntryRepository struct {
	db *sql.DB
}

// NewEntryRepository returns a PostgreSQL entry repository.
func NewEntryRepository(db *sql.DB) (*EntryRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("EntryRepository requires an sql.DB instance")
	}

	return &EntryRepository{db: db}, nil
}

// Add stores the provided entry in PostgreSQL.
func (repository *EntryRepository) Add(ctx context.Context, entry *entries.Entry) error {
	_, err := repository.db.ExecContext(
		ctx,
		`INSERT INTO entries (id, date, mood, note, created_at, updated_at, account_id) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		entry.ID().Uuid(),
		time.Time(entry.Date()),
		int(entry.Mood()),
		string(entry.Note()),
		entry.CreatedAt(),
		entry.UpdatedAt(),
		entry.AccountId().Uuid(),
	)
	if err != nil {
		return fmt.Errorf("executing an INSERT sql query: %w", err)
	}

	return nil
}

// Find returns the first entry matching the provided filter or nil when no row exists.
func (repository *EntryRepository) Find(ctx context.Context, filter entries.EntryFilter) (*entries.Entry, error) {
	query := `SELECT id, date, mood, note, created_at, updated_at, account_id FROM entries`
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

	if len(filter.Ids) > 0 {
		placeholders := make([]string, 0, len(filter.Ids))
		for _, id := range filter.Ids {
			args = append(args, id.Uuid())
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}

		conditions = append(conditions, fmt.Sprintf("id IN (%s)", strings.Join(placeholders, ", ")))
	}

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	query = fmt.Sprintf("%s LIMIT 1", query)

	row := repository.db.QueryRowContext(ctx, query, args...)

	var (
		id        uuid.UUID
		date      time.Time
		mood      int
		note      string
		createdAt time.Time
		updatedAt time.Time
		accountId uuid.UUID
	)

	if err := row.Scan(&id, &date, &mood, &note, &createdAt, &updatedAt, &accountId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("executing a SELECT sql query: %w", err)
	}

	entry := entries.Restore(id, date, mood, note, createdAt, updatedAt, accountId)

	return entry, nil
}

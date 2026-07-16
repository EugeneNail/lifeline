package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/google/uuid"
)

// CompletableHabitRepository stores completable habits in PostgreSQL.
type CompletableHabitRepository struct {
	db *sql.DB
}

// NewCompletableHabitRepository returns a PostgreSQL completable habit repository.
func NewCompletableHabitRepository(db *sql.DB) (*CompletableHabitRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("CompletableHabitRepository requires an sql.DB instance")
	}

	return &CompletableHabitRepository{db: db}, nil
}

// Add stores the provided completable habit in PostgreSQL.
func (repository *CompletableHabitRepository) Add(ctx context.Context, habit *habits.CompletableHabit) error {
	_, err := repository.db.ExecContext(
		ctx,
		`INSERT INTO completable_habits (id, label, icon, created_at, updated_at, archived_at, deleted_at, account_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		habit.ID(),
		habit.Label(),
		int(habit.Icon()),
		habit.CreatedAt(),
		habit.UpdatedAt(),
		habit.ArchivedAt(),
		habit.DeletedAt(),
		habit.AccountId(),
	)
	if err != nil {
		return fmt.Errorf("executing an INSERT sql query: %w", err)
	}

	return nil
}

// Count returns the number of completable habits matching the provided filter.
func (repository *CompletableHabitRepository) Count(ctx context.Context, filter habits.CompletableHabitFilter) (int, error) {
	query := `SELECT COUNT(*) FROM completable_habits`
	conditions, args := repository.buildConditions(filter)

	query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))

	row := repository.db.QueryRowContext(ctx, query, args...)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("executing a COUNT sql query: %w", err)
	}

	return count, nil
}

// Find returns the first completable habit matching the provided filter or nil when no row exists.
func (repository *CompletableHabitRepository) Find(ctx context.Context, filter habits.CompletableHabitFilter) (*habits.CompletableHabit, error) {
	foundHabits, err := repository.FindMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("finding completable habits: %w", err)
	}

	if len(foundHabits) == 0 {
		return nil, nil
	}

	return foundHabits[0], nil
}

// FindMany returns all completable habits matching the provided filter.
func (repository *CompletableHabitRepository) FindMany(ctx context.Context, filter habits.CompletableHabitFilter) ([]*habits.CompletableHabit, error) {
	query := `SELECT id, label, icon, created_at, updated_at, archived_at, deleted_at, account_id FROM completable_habits`
	conditions, args := repository.buildConditions(filter)

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	rows, err := repository.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing a SELECT sql query: %w", err)
	}
	defer rows.Close()

	foundHabits := make([]*habits.CompletableHabit, 0)
	for rows.Next() {
		habit, err := repository.scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning a completable habit row: %w", err)
		}

		foundHabits = append(foundHabits, habit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating completable habit rows: %w", err)
	}

	return foundHabits, nil
}

func (repository *CompletableHabitRepository) buildConditions(filter habits.CompletableHabitFilter) ([]string, []any) {
	conditions := make([]string, 0, 3)
	args := make([]any, 0)

	if len(filter.AccountIds) > 0 {
		placeholders := make([]string, 0, len(filter.AccountIds))
		for _, accountId := range filter.AccountIds {
			args = append(args, accountId)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}

		conditions = append(conditions, fmt.Sprintf("account_id IN (%s)", strings.Join(placeholders, ", ")))
	}

	if len(filter.CompletableHabitIds) > 0 {
		placeholders := make([]string, 0, len(filter.CompletableHabitIds))
		for _, id := range filter.CompletableHabitIds {
			args = append(args, id)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}

		conditions = append(conditions, fmt.Sprintf("id IN (%s)", strings.Join(placeholders, ", ")))
	}

	if filter.Archived {
		conditions = append(conditions, "archived_at IS NOT NULL")
	} else {
		conditions = append(conditions, "archived_at IS NULL")
	}

	if filter.Deleted {
		conditions = append(conditions, "deleted_at IS NOT NULL")
	} else {
		conditions = append(conditions, "deleted_at IS NULL")
	}

	return conditions, args
}

func (repository *CompletableHabitRepository) scan(rows *sql.Rows) (*habits.CompletableHabit, error) {
	var (
		id         uuid.UUID
		label      string
		icon       int
		createdAt  time.Time
		updatedAt  time.Time
		archivedAt sql.NullTime
		deletedAt  sql.NullTime
		accountId  uuid.UUID
	)

	if err := rows.Scan(&id, &label, &icon, &createdAt, &updatedAt, &archivedAt, &deletedAt, &accountId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("scanning a SELECT sql result: %w", err)
	}

	return habits.RestoreCompletableHabit(
		id,
		label,
		icon,
		createdAt,
		updatedAt,
		repository.nullTime(archivedAt),
		repository.nullTime(deletedAt),
		accountId,
	), nil
}

func (repository *CompletableHabitRepository) nullTime(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
}

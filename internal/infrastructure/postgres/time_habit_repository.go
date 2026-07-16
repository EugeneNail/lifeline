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

// TimeHabitRepository stores time habits in PostgreSQL.
type TimeHabitRepository struct {
	db *sql.DB
}

// NewTimeHabitRepository returns a PostgreSQL time habit repository.
func NewTimeHabitRepository(db *sql.DB) (*TimeHabitRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("TimeHabitRepository requires an sql.DB instance")
	}

	return &TimeHabitRepository{db: db}, nil
}

// Add stores the provided time habit in PostgreSQL.
func (repository *TimeHabitRepository) Add(ctx context.Context, habit *habits.TimeHabit) error {
	_, err := repository.db.ExecContext(
		ctx,
		`INSERT INTO time_habits (id, label, icon, created_at, updated_at, archived_at, deleted_at, account_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
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

// Count returns the number of time habits matching the provided filter.
func (repository *TimeHabitRepository) Count(ctx context.Context, filter habits.TimeHabitFilter) (int, error) {
	query := `SELECT COUNT(*) FROM time_habits`
	conditions, args := repository.buildConditions(filter)

	query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))

	row := repository.db.QueryRowContext(ctx, query, args...)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("executing a COUNT sql query: %w", err)
	}

	return count, nil
}

// Find returns the first time habit matching the provided filter or nil when no row exists.
func (repository *TimeHabitRepository) Find(ctx context.Context, filter habits.TimeHabitFilter) (*habits.TimeHabit, error) {
	foundHabits, err := repository.FindMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("finding time habits: %w", err)
	}

	if len(foundHabits) == 0 {
		return nil, nil
	}

	return foundHabits[0], nil
}

// FindMany returns all time habits matching the provided filter.
func (repository *TimeHabitRepository) FindMany(ctx context.Context, filter habits.TimeHabitFilter) ([]*habits.TimeHabit, error) {
	query := `SELECT id, label, icon, created_at, updated_at, archived_at, deleted_at, account_id FROM time_habits`
	conditions, args := repository.buildConditions(filter)

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	rows, err := repository.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing a SELECT sql query: %w", err)
	}
	defer rows.Close()

	foundHabits := make([]*habits.TimeHabit, 0)
	for rows.Next() {
		habit, err := repository.scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning a time habit row: %w", err)
		}

		foundHabits = append(foundHabits, habit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating time habit rows: %w", err)
	}

	return foundHabits, nil
}

func (repository *TimeHabitRepository) buildConditions(filter habits.TimeHabitFilter) ([]string, []any) {
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

	if len(filter.TimeHabitIds) > 0 {
		placeholders := make([]string, 0, len(filter.TimeHabitIds))
		for _, id := range filter.TimeHabitIds {
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

func (repository *TimeHabitRepository) scan(rows *sql.Rows) (*habits.TimeHabit, error) {
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

	return habits.RestoreTimeHabit(
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

func (repository *TimeHabitRepository) nullTime(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
}

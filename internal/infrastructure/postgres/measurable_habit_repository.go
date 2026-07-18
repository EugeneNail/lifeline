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

// MeasurableHabitRepository stores measurable habits in PostgreSQL.
type MeasurableHabitRepository struct {
	db *sql.DB
}

// NewMeasurableHabitRepository returns a PostgreSQL measurable habit repository.
func NewMeasurableHabitRepository(db *sql.DB) (*MeasurableHabitRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("MeasurableHabitRepository requires an sql.DB instance")
	}

	return &MeasurableHabitRepository{db: db}, nil
}

// Add stores the provided measurable habit in PostgreSQL.
func (repository *MeasurableHabitRepository) Add(ctx context.Context, habit *habits.MeasurableHabit) error {
	_, err := repository.db.ExecContext(
		ctx,
		`INSERT INTO measurable_habits (id, label, icon, step, unit, created_at, updated_at, archived_at, deleted_at, account_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		habit.ID(),
		habit.Label(),
		int(habit.Icon()),
		float32(habit.Step()),
		string(habit.Unit()),
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

// Save updates all fields of the provided measurable habit in PostgreSQL.
func (repository *MeasurableHabitRepository) Save(ctx context.Context, habit *habits.MeasurableHabit) error {
	result, err := repository.db.ExecContext(
		ctx,
		`UPDATE measurable_habits SET label = $2, icon = $3, step = $4, unit = $5, created_at = $6, updated_at = $7, archived_at = $8, deleted_at = $9, account_id = $10 WHERE id = $1`,
		habit.ID(),
		habit.Label(),
		int(habit.Icon()),
		float32(habit.Step()),
		string(habit.Unit()),
		habit.CreatedAt(),
		habit.UpdatedAt(),
		habit.ArchivedAt(),
		habit.DeletedAt(),
		habit.AccountId(),
	)
	if err != nil {
		return fmt.Errorf("executing an UPDATE sql query for measurable habit %s: %w", habit.ID(), err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking updated rows for measurable habit %s: %w", habit.ID(), err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated for habit id %q", habit.ID())
	}

	return nil
}

// Count returns the number of measurable habits matching the provided filter.
func (repository *MeasurableHabitRepository) Count(ctx context.Context, filter habits.MeasurableHabitFilter) (int, error) {
	query := `SELECT COUNT(*) FROM measurable_habits`
	conditions, args := repository.buildConditions(filter)

	query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))

	row := repository.db.QueryRowContext(ctx, query, args...)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("executing a COUNT sql query: %w", err)
	}

	return count, nil
}

// Find returns the first measurable habit matching the provided filter or nil when no row exists.
func (repository *MeasurableHabitRepository) Find(ctx context.Context, filter habits.MeasurableHabitFilter) (*habits.MeasurableHabit, error) {
	foundHabits, err := repository.FindMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("finding measurable habits: %w", err)
	}

	if len(foundHabits) == 0 {
		return nil, nil
	}

	return foundHabits[0], nil
}

// FindMany returns all measurable habits matching the provided filter.
func (repository *MeasurableHabitRepository) FindMany(ctx context.Context, filter habits.MeasurableHabitFilter) ([]*habits.MeasurableHabit, error) {
	query := `SELECT id, label, icon, step, unit, created_at, updated_at, archived_at, deleted_at, account_id FROM measurable_habits`
	conditions, args := repository.buildConditions(filter)

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	rows, err := repository.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing a SELECT sql query: %w", err)
	}
	defer rows.Close()

	foundHabits := make([]*habits.MeasurableHabit, 0)
	for rows.Next() {
		habit, err := repository.scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning a measurable habit row: %w", err)
		}

		foundHabits = append(foundHabits, habit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating measurable habit rows: %w", err)
	}

	return foundHabits, nil
}

func (repository *MeasurableHabitRepository) buildConditions(filter habits.MeasurableHabitFilter) ([]string, []any) {
	conditions := make([]string, 0, 4)
	args := make([]any, 0)

	if len(filter.AccountIds) > 0 {
		placeholders := make([]string, 0, len(filter.AccountIds))
		for _, accountId := range filter.AccountIds {
			args = append(args, accountId)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}

		conditions = append(conditions, fmt.Sprintf("account_id IN (%s)", strings.Join(placeholders, ", ")))
	}

	if len(filter.MeasurableHabitIds) > 0 {
		placeholders := make([]string, 0, len(filter.MeasurableHabitIds))
		for _, id := range filter.MeasurableHabitIds {
			args = append(args, id)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}

		conditions = append(conditions, fmt.Sprintf("id IN (%s)", strings.Join(placeholders, ", ")))
	}

	if filter.Archived != nil {
		if *filter.Archived {
			conditions = append(conditions, "archived_at IS NOT NULL")
		} else {
			conditions = append(conditions, "archived_at IS NULL")
		}
	}

	if filter.Deleted != nil {
		if *filter.Deleted {
			conditions = append(conditions, "deleted_at IS NOT NULL")
		} else {
			conditions = append(conditions, "deleted_at IS NULL")
		}
	}

	return conditions, args
}

func (repository *MeasurableHabitRepository) scan(rows *sql.Rows) (*habits.MeasurableHabit, error) {
	var (
		id         uuid.UUID
		label      string
		icon       int
		step       float32
		unit       string
		createdAt  time.Time
		updatedAt  time.Time
		archivedAt sql.NullTime
		deletedAt  sql.NullTime
		accountId  uuid.UUID
	)

	if err := rows.Scan(&id, &label, &icon, &step, &unit, &createdAt, &updatedAt, &archivedAt, &deletedAt, &accountId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("scanning a SELECT sql result: %w", err)
	}

	return habits.RestoreMeasurableHabit(
		id,
		label,
		icon,
		step,
		unit,
		createdAt,
		updatedAt,
		repository.nullTime(archivedAt),
		repository.nullTime(deletedAt),
		accountId,
	), nil
}

func (repository *MeasurableHabitRepository) nullTime(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/habits/records"
	"github.com/google/uuid"
)

// CompletableHabitRecordRepository stores completable habit records in PostgreSQL.
type CompletableHabitRecordRepository struct {
	db *sql.DB
}

// NewCompletableHabitRecordRepository returns a PostgreSQL completable habit record repository.
func NewCompletableHabitRecordRepository(db *sql.DB) (*CompletableHabitRecordRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("CompletableHabitRecordRepository requires an sql.DB instance")
	}

	return &CompletableHabitRecordRepository{db: db}, nil
}

// Add stores the provided completable habit record in PostgreSQL.
func (repository *CompletableHabitRecordRepository) Add(ctx context.Context, record *records.CompletableHabitRecord) error {
	date := time.Time(record.Date())
	_, err := repository.db.ExecContext(
		ctx,
		`INSERT INTO completable_habit_records (completable_habit_id, account_id, date, value) VALUES ($1, $2, $3, $4)`,
		record.CompletableHabitId(),
		record.AccountId(),
		date,
		record.Value(),
	)
	if err != nil {
		return fmt.Errorf("executing an INSERT sql query for completable habit record %s at %s: %w", record.CompletableHabitId(), date, err)
	}

	return nil
}

// Save updates only the value of the provided completable habit record in PostgreSQL.
func (repository *CompletableHabitRecordRepository) Save(ctx context.Context, record *records.CompletableHabitRecord) error {
	date := time.Time(record.Date())
	result, err := repository.db.ExecContext(
		ctx,
		`UPDATE completable_habit_records SET value = $2 WHERE completable_habit_id = $1 AND date = $3`,
		record.CompletableHabitId(),
		record.Value(),
		date,
	)
	if err != nil {
		return fmt.Errorf("executing an UPDATE sql query for completable habit record %s at %s: %w", record.CompletableHabitId(), date, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking updated rows for completable habit record %s at %s: %w", record.CompletableHabitId(), date, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("saving completable habit record %s at %s: no rows updated", record.CompletableHabitId(), date)
	}

	return nil
}

// FindMany returns all completable habit records matching the provided filter.
func (repository *CompletableHabitRecordRepository) FindMany(ctx context.Context, filter records.CompletableHabitRecordFilter) ([]*records.CompletableHabitRecord, error) {
	query := `SELECT completable_habit_id, account_id, date, value FROM completable_habit_records`
	conditions, args := repository.buildConditions(filter)

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	rows, err := repository.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing a SELECT sql query for completable habit records: %w", err)
	}
	defer rows.Close()

	foundRecords := make([]*records.CompletableHabitRecord, 0)
	for rows.Next() {
		record, err := repository.scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning a completable habit record row: %w", err)
		}

		foundRecords = append(foundRecords, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating completable habit record rows: %w", err)
	}

	return foundRecords, nil
}

// Find returns the first completable habit record matching the provided filter, nil when none exists, or an error when lookup fails.
func (repository *CompletableHabitRecordRepository) Find(ctx context.Context, filter records.CompletableHabitRecordFilter) (*records.CompletableHabitRecord, error) {
	foundRecords, err := repository.FindMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("finding completable habit records: %w", err)
	}

	if len(foundRecords) == 0 {
		return nil, nil
	}

	return foundRecords[0], nil
}

func (repository *CompletableHabitRecordRepository) buildConditions(filter records.CompletableHabitRecordFilter) ([]string, []any) {
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

	if len(filter.CompletableHabitRecordIds) > 0 {
		placeholders := make([]string, 0, len(filter.CompletableHabitRecordIds))
		for _, id := range filter.CompletableHabitRecordIds {
			args = append(args, id)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}

		conditions = append(conditions, fmt.Sprintf("completable_habit_id IN (%s)", strings.Join(placeholders, ", ")))
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

func (repository *CompletableHabitRecordRepository) scan(rows *sql.Rows) (*records.CompletableHabitRecord, error) {
	var (
		completableHabitId uuid.UUID
		accountId          uuid.UUID
		date               time.Time
		value              bool
	)

	if err := rows.Scan(&completableHabitId, &accountId, &date, &value); err != nil {
		return nil, fmt.Errorf("scanning a SELECT sql result: %w", err)
	}

	return records.RestoreCompletableHabitRecord(completableHabitId, accountId, records.NewDate(date), value), nil
}

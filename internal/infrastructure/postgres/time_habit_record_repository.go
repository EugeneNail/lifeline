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

// TimeHabitRecordRepository stores time habit records in PostgreSQL.
type TimeHabitRecordRepository struct {
	db *sql.DB
}

// NewTimeHabitRecordRepository returns a PostgreSQL time habit record repository.
func NewTimeHabitRecordRepository(db *sql.DB) (*TimeHabitRecordRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("TimeHabitRecordRepository requires an sql.DB instance")
	}

	return &TimeHabitRecordRepository{db: db}, nil
}

// Add stores the provided time habit record in PostgreSQL.
func (repository *TimeHabitRecordRepository) Add(ctx context.Context, record *records.TimeHabitRecord) error {
	date := time.Time(record.Date())
	_, err := repository.db.ExecContext(
		ctx,
		`INSERT INTO time_habit_records (time_habit_id, account_id, date, value) VALUES ($1, $2, $3, $4)`,
		record.TimeHabitId(),
		record.AccountId(),
		date,
		record.Value().Value(),
	)
	if err != nil {
		return fmt.Errorf("executing an INSERT sql query for time habit record %s at %s: %w", record.TimeHabitId(), date, err)
	}

	return nil
}

// Save updates only the value of the provided time habit record in PostgreSQL.
func (repository *TimeHabitRecordRepository) Save(ctx context.Context, record *records.TimeHabitRecord) error {
	date := time.Time(record.Date())
	result, err := repository.db.ExecContext(
		ctx,
		`UPDATE time_habit_records SET value = $2 WHERE time_habit_id = $1 AND date = $3`,
		record.TimeHabitId(),
		record.Value().Value(),
		date,
	)
	if err != nil {
		return fmt.Errorf("executing an UPDATE sql query for time habit record %s at %s: %w", record.TimeHabitId(), date, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking updated rows for time habit record %s at %s: %w", record.TimeHabitId(), date, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("saving time habit record %s at %s: no rows updated", record.TimeHabitId(), date)
	}

	return nil
}

// FindMany returns all time habit records matching the provided filter.
func (repository *TimeHabitRecordRepository) FindMany(ctx context.Context, filter records.TimeHabitRecordFilter) ([]*records.TimeHabitRecord, error) {
	query := `SELECT time_habit_id, account_id, date, value FROM time_habit_records`
	conditions, args := repository.buildConditions(filter)

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	rows, err := repository.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing a SELECT sql query for time habit records: %w", err)
	}
	defer rows.Close()

	foundRecords := make([]*records.TimeHabitRecord, 0)
	for rows.Next() {
		record, err := repository.scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning a time habit record row: %w", err)
		}

		foundRecords = append(foundRecords, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating time habit record rows: %w", err)
	}

	return foundRecords, nil
}

// Find returns the first time habit record matching the provided filter, nil when none exists, or an error when lookup fails.
func (repository *TimeHabitRecordRepository) Find(ctx context.Context, filter records.TimeHabitRecordFilter) (*records.TimeHabitRecord, error) {
	foundRecords, err := repository.FindMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("finding time habit records: %w", err)
	}

	if len(foundRecords) == 0 {
		return nil, nil
	}

	return foundRecords[0], nil
}

func (repository *TimeHabitRecordRepository) buildConditions(filter records.TimeHabitRecordFilter) ([]string, []any) {
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

	if len(filter.TimeHabitRecordIds) > 0 {
		placeholders := make([]string, 0, len(filter.TimeHabitRecordIds))
		for _, id := range filter.TimeHabitRecordIds {
			args = append(args, id)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}

		conditions = append(conditions, fmt.Sprintf("time_habit_id IN (%s)", strings.Join(placeholders, ", ")))
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

func (repository *TimeHabitRecordRepository) scan(rows *sql.Rows) (*records.TimeHabitRecord, error) {
	var (
		timeHabitId uuid.UUID
		accountId   uuid.UUID
		date        time.Time
		rawValue    int
	)

	if err := rows.Scan(&timeHabitId, &accountId, &date, &rawValue); err != nil {
		return nil, fmt.Errorf("scanning a SELECT sql result: %w", err)
	}

	value, err := records.NewTimeValue(rawValue)
	if err != nil {
		return nil, fmt.Errorf("restoring time habit record %s at %s: %w", timeHabitId, date, err)
	}

	return records.RestoreTimeHabitRecord(timeHabitId, accountId, records.NewDate(date), value), nil
}

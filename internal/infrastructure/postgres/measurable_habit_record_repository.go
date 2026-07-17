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

// MeasurableHabitRecordRepository stores measurable habit records in PostgreSQL.
type MeasurableHabitRecordRepository struct {
	db *sql.DB
}

// NewMeasurableHabitRecordRepository returns a PostgreSQL measurable habit record repository.
func NewMeasurableHabitRecordRepository(db *sql.DB) (*MeasurableHabitRecordRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("MeasurableHabitRecordRepository requires an sql.DB instance")
	}

	return &MeasurableHabitRecordRepository{db: db}, nil
}

// Add stores the provided measurable habit record in PostgreSQL.
func (repository *MeasurableHabitRecordRepository) Add(ctx context.Context, record *records.MeasurableHabitRecord) error {
	date := time.Time(record.Date())
	_, err := repository.db.ExecContext(
		ctx,
		`INSERT INTO measurable_habit_records (measurable_habit_id, account_id, date, value) VALUES ($1, $2, $3, $4)`,
		record.MeasurableHabitId(),
		record.AccountId(),
		date,
		record.Value().Value(),
	)
	if err != nil {
		return fmt.Errorf("executing an INSERT sql query for measurable habit record %s at %s: %w", record.MeasurableHabitId(), date, err)
	}

	return nil
}

// Save updates only the value of the provided measurable habit record in PostgreSQL.
func (repository *MeasurableHabitRecordRepository) Save(ctx context.Context, record *records.MeasurableHabitRecord) error {
	date := time.Time(record.Date())
	result, err := repository.db.ExecContext(
		ctx,
		`UPDATE measurable_habit_records SET value = $2 WHERE measurable_habit_id = $1 AND date = $3`,
		record.MeasurableHabitId(),
		record.Value().Value(),
		date,
	)
	if err != nil {
		return fmt.Errorf("executing an UPDATE sql query for measurable habit record %s at %s: %w", record.MeasurableHabitId(), date, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking updated rows for measurable habit record %s at %s: %w", record.MeasurableHabitId(), date, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("saving measurable habit record %s at %s: no rows updated", record.MeasurableHabitId(), date)
	}

	return nil
}

// FindMany returns all measurable habit records matching the provided filter.
func (repository *MeasurableHabitRecordRepository) FindMany(ctx context.Context, filter records.MeasurableHabitRecordFilter) ([]*records.MeasurableHabitRecord, error) {
	query := `SELECT measurable_habit_id, account_id, date, value FROM measurable_habit_records`
	conditions, args := repository.buildConditions(filter)

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	rows, err := repository.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing a SELECT sql query for measurable habit records: %w", err)
	}
	defer rows.Close()

	foundRecords := make([]*records.MeasurableHabitRecord, 0)
	for rows.Next() {
		record, err := repository.scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning a measurable habit record row: %w", err)
		}

		foundRecords = append(foundRecords, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating measurable habit record rows: %w", err)
	}

	return foundRecords, nil
}

// Find returns the first measurable habit record matching the provided filter, nil when none exists, or an error when lookup fails.
func (repository *MeasurableHabitRecordRepository) Find(ctx context.Context, filter records.MeasurableHabitRecordFilter) (*records.MeasurableHabitRecord, error) {
	foundRecords, err := repository.FindMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("finding measurable habit records: %w", err)
	}

	if len(foundRecords) == 0 {
		return nil, nil
	}

	return foundRecords[0], nil
}

func (repository *MeasurableHabitRecordRepository) buildConditions(filter records.MeasurableHabitRecordFilter) ([]string, []any) {
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

	if len(filter.MeasurableHabitRecordIds) > 0 {
		placeholders := make([]string, 0, len(filter.MeasurableHabitRecordIds))
		for _, id := range filter.MeasurableHabitRecordIds {
			args = append(args, id)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}

		conditions = append(conditions, fmt.Sprintf("measurable_habit_id IN (%s)", strings.Join(placeholders, ", ")))
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

func (repository *MeasurableHabitRecordRepository) scan(rows *sql.Rows) (*records.MeasurableHabitRecord, error) {
	var (
		measurableHabitId uuid.UUID
		accountId         uuid.UUID
		date              time.Time
		rawValue          float32
	)

	if err := rows.Scan(&measurableHabitId, &accountId, &date, &rawValue); err != nil {
		return nil, fmt.Errorf("scanning a SELECT sql result: %w", err)
	}

	return records.RestoreMeasurableHabitRecord(measurableHabitId, accountId, records.NewDate(date), records.MeasurableValue(rawValue)), nil
}

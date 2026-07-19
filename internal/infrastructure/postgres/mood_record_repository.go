package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	mooddomain "github.com/EugeneNail/lifeline/internal/domain/moods"
	"github.com/google/uuid"
)

// RecordRepository stores mood records in PostgreSQL.
type RecordRepository struct {
	db *sql.DB
}

// NewRecordRepository returns a PostgreSQL mood record repository.
func NewRecordRepository(db *sql.DB) (*RecordRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("RecordRepository requires an sql.DB instance")
	}

	return &RecordRepository{db: db}, nil
}

// Add stores the provided mood record in PostgreSQL.
func (repository *RecordRepository) Add(ctx context.Context, record *mooddomain.Record) error {
	date := time.Time(record.Date())
	_, err := repository.db.ExecContext(
		ctx,
		`INSERT INTO mood_records (date, account_id, value, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`,
		date,
		record.AccountId(),
		int(record.Value()),
		record.CreatedAt(),
		record.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("executing an INSERT sql query for mood record %s and account %s: %w", date, record.AccountId(), err)
	}

	return nil
}

// Update updates the mutable fields of the provided mood record in PostgreSQL.
func (repository *RecordRepository) Update(ctx context.Context, record *mooddomain.Record) error {
	date := time.Time(record.Date())
	result, err := repository.db.ExecContext(
		ctx,
		`UPDATE mood_records SET value = $3, updated_at = $4 WHERE date = $1 AND account_id = $2`,
		date,
		record.AccountId(),
		int(record.Value()),
		record.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("executing an UPDATE sql query for mood record %s and account %s: %w", date, record.AccountId(), err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking updated rows for mood record %s and account %s: %w", date, record.AccountId(), err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("saving mood record %s and account %s: no rows updated", date, record.AccountId())
	}

	return nil
}

// FindMany returns all mood records matching the provided filter.
func (repository *RecordRepository) FindMany(ctx context.Context, filter mooddomain.RecordFilter) ([]*mooddomain.Record, error) {
	query := `SELECT date, account_id, value, created_at, updated_at FROM mood_records`
	conditions, args := repository.buildConditions(filter)

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	rows, err := repository.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("executing a SELECT sql query for mood records: %w", err)
	}
	defer rows.Close()

	foundRecords := make([]*mooddomain.Record, 0)
	for rows.Next() {
		record, err := repository.scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning a mood record row: %w", err)
		}

		foundRecords = append(foundRecords, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating mood record rows: %w", err)
	}

	return foundRecords, nil
}

// Find returns the first mood record matching the provided filter, nil when none exists, or an error when lookup fails.
func (repository *RecordRepository) Find(ctx context.Context, filter mooddomain.RecordFilter) (*mooddomain.Record, error) {
	foundRecords, err := repository.FindMany(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("finding mood records: %w", err)
	}

	if len(foundRecords) == 0 {
		return nil, nil
	}

	return foundRecords[0], nil
}

func (repository *RecordRepository) buildConditions(filter mooddomain.RecordFilter) ([]string, []any) {
	conditions := make([]string, 0, 2)
	args := make([]any, 0)

	if len(filter.AccountIds) > 0 {
		placeholders := make([]string, 0, len(filter.AccountIds))
		for _, accountId := range filter.AccountIds {
			args = append(args, accountId)
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

func (repository *RecordRepository) scan(rows *sql.Rows) (*mooddomain.Record, error) {
	var (
		date      time.Time
		accountId uuid.UUID
		rawValue  int
		createdAt time.Time
		updatedAt time.Time
	)

	if err := rows.Scan(&date, &accountId, &rawValue, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("scanning a SELECT sql result: %w", err)
	}

	value, err := mooddomain.New(rawValue)
	if err != nil {
		return nil, fmt.Errorf("restoring a mood record for account %s and date %s: %w", accountId, date, err)
	}

	recordDate, err := domain.NewDate(date)
	if err != nil {
		return nil, fmt.Errorf("restoring a mood record for account %s and date %s: %w", accountId, date, err)
	}

	return mooddomain.RestoreRecord(recordDate, value, createdAt, updatedAt, accountId), nil
}

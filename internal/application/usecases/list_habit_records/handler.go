package list_habit_records

import (
	"context"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/habits/records"
	"github.com/google/uuid"
)

// Handler executes the list-habit-records use case.
type Handler struct {
	completableHabitRecords records.CompletableHabitRecordRepository
	timeHabitRecords        records.TimeHabitRecordRepository
	measurableHabitRecords  records.MeasurableHabitRecordRepository
}

// NewHandler returns a list-habit-records handler configured with the record repositories or an error when a dependency is missing.
func NewHandler(
	completableHabitRecords records.CompletableHabitRecordRepository,
	timeHabitRecords records.TimeHabitRecordRepository,
	measurableHabitRecords records.MeasurableHabitRecordRepository,
) (*Handler, error) {
	if completableHabitRecords == nil {
		return nil, fmt.Errorf("list_habit_records handler requires a completable habit record repository")
	}

	if timeHabitRecords == nil {
		return nil, fmt.Errorf("list_habit_records handler requires a time habit record repository")
	}

	if measurableHabitRecords == nil {
		return nil, fmt.Errorf("list_habit_records handler requires a measurable habit record repository")
	}

	return &Handler{
		completableHabitRecords: completableHabitRecords,
		timeHabitRecords:        timeHabitRecords,
		measurableHabitRecords:  measurableHabitRecords,
	}, nil
}

// Query carries the data required to list habit records for a specific day.
type Query struct {
	AccountId uuid.UUID
	Date      time.Time
}

// Result groups habit records by their type and returns them to the caller.
type Result struct {
	Measurable  []*records.MeasurableHabitRecord
	Time        []*records.TimeHabitRecord
	Completable []*records.CompletableHabitRecord
}

// Handle loads the user's habit records for the requested day and returns them grouped by record type.
func (handler *Handler) Handle(ctx context.Context, query Query) (Result, error) {
	date := records.NewDate(query.Date)

	completableRecords, err := handler.completableHabitRecords.FindMany(ctx, records.NewCompletableHabitRecordFilter().
		WithAccountIds(query.AccountId).
		WithDates(time.Time(date)),
	)
	if err != nil {
		return Result{}, fmt.Errorf("finding completable habit records: %w", err)
	}

	timeRecords, err := handler.timeHabitRecords.FindMany(ctx, records.NewTimeHabitRecordFilter().
		WithAccountIds(query.AccountId).
		WithDates(time.Time(date)),
	)
	if err != nil {
		return Result{}, fmt.Errorf("finding time habit records: %w", err)
	}

	measurableRecords, err := handler.measurableHabitRecords.FindMany(ctx, records.NewMeasurableHabitRecordFilter().
		WithAccountIds(query.AccountId).
		WithDates(time.Time(date)),
	)
	if err != nil {
		return Result{}, fmt.Errorf("finding measurable habit records: %w", err)
	}

	return Result{
		Measurable:  measurableRecords,
		Time:        timeRecords,
		Completable: completableRecords,
	}, nil
}

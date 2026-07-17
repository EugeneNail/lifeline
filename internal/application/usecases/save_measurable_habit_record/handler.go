package save_measurable_habit_record

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/EugeneNail/lifeline/internal/domain/habits/records"
	"github.com/google/uuid"
)

// Handler executes the save-measurable-habit-record use case.
type Handler struct {
	measurableHabitRecords records.MeasurableHabitRecordRepository
	measurableHabits       habits.MeasurableHabitRepository
	savingPolicy           *records.SavingPolicy
}

// NewHandler returns a save-measurable-habit-record handler configured with the measurable habit repository, the record repository, and the saving policy or an error when a dependency is missing.
func NewHandler(measurableHabitRecords records.MeasurableHabitRecordRepository, measurableHabits habits.MeasurableHabitRepository, savingPolicy *records.SavingPolicy) (*Handler, error) {
	if measurableHabitRecords == nil {
		return nil, fmt.Errorf("save_measurable_habit_record handler requires a measurable habit record repository")
	}

	if measurableHabits == nil {
		return nil, fmt.Errorf("save_measurable_habit_record handler requires a measurable habit repository")
	}

	if savingPolicy == nil {
		return nil, fmt.Errorf("save_measurable_habit_record handler requires a saving policy")
	}

	return &Handler{
		measurableHabitRecords: measurableHabitRecords,
		measurableHabits:       measurableHabits,
		savingPolicy:           savingPolicy,
	}, nil
}

// Command carries the data required to save a measurable habit record.
type Command struct {
	AccountId         uuid.UUID
	Value             float32
	Date              time.Time
	MeasurableHabitId uuid.UUID
}

// Handle updates the existing measurable habit record value or creates a new record when none exists, and returns an error when persistence fails.
func (handler *Handler) Handle(ctx context.Context, command Command) error {
	date := records.NewDate(command.Date)

	habit, err := handler.measurableHabits.Find(ctx, habits.NewMeasurableHabitFilter().
		WithAccountIds(command.AccountId).
		WithIds(command.MeasurableHabitId),
	)
	if err != nil {
		return fmt.Errorf("finding a measurable habit: %w", err)
	}

	if habit == nil {
		return habits.ErrHabitNotFound
	}

	if err := handler.savingPolicy.Check(command.AccountId, date, habit); err != nil {
		return err
	}

	record, err := handler.measurableHabitRecords.Find(ctx, records.NewMeasurableHabitRecordFilter().
		WithAccountIds(command.AccountId).
		WithIds(command.MeasurableHabitId).
		WithDates(time.Time(date)),
	)
	if err != nil {
		return fmt.Errorf("finding a measurable habit record: %w", err)
	}

	isNew := false
	if record == nil {
		isNew = true
		record, err = records.NewMeasurableHabitRecord(command.MeasurableHabitId, command.AccountId, time.Time(date), command.Value, habit.Step())
		if err != nil {
			var errs domain.ValidationErrors
			if errors.As(err, &errs) {
				return errs
			}

			return fmt.Errorf("creating a new measurable habit record: %w", err)
		}
	}

	if err := record.ChangeValue(command.Value, habit.Step()); err != nil {
		errs := domain.NewValidationErrors()
		errs.Add("value", err)

		return errs
	}

	if isNew {
		if err := handler.measurableHabitRecords.Add(ctx, record); err != nil {
			return fmt.Errorf("adding a new measurable habit record: %w", err)
		}
	} else {
		if err := handler.measurableHabitRecords.Save(ctx, record); err != nil {
			return fmt.Errorf("saving a measurable habit record: %w", err)
		}
	}

	return nil
}

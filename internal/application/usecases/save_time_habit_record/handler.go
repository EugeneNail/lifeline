package save_time_habit_record

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

// Handler executes the save-time-habit-record use case.
type Handler struct {
	timeHabitRecords records.TimeHabitRecordRepository
	timeHabits       habits.TimeHabitRepository
	savingPolicy     *records.SavingPolicy
}

// NewHandler returns a save-time-habit-record handler configured with the time habit repository, the record repository, and the saving policy or an error when a dependency is missing.
func NewHandler(timeHabitRecords records.TimeHabitRecordRepository, timeHabits habits.TimeHabitRepository, savingPolicy *records.SavingPolicy) (*Handler, error) {
	if timeHabitRecords == nil {
		return nil, fmt.Errorf("save_time_habit_record handler requires a time habit record repository")
	}

	if timeHabits == nil {
		return nil, fmt.Errorf("save_time_habit_record handler requires a time habit repository")
	}

	if savingPolicy == nil {
		return nil, fmt.Errorf("save_time_habit_record handler requires a saving policy")
	}

	return &Handler{
		timeHabitRecords: timeHabitRecords,
		timeHabits:       timeHabits,
		savingPolicy:     savingPolicy,
	}, nil
}

// Command carries the data required to save a time habit record.
type Command struct {
	AccountId   uuid.UUID
	Value       int
	Date        time.Time
	TimeHabitId uuid.UUID
}

// Handle updates the existing time habit record value or creates a new record when none exists, and returns an error when persistence fails.
func (handler *Handler) Handle(ctx context.Context, command Command) error {
	date := records.NewDate(command.Date)

	habit, err := handler.timeHabits.Find(ctx, habits.NewTimeHabitFilter().
		WithAccountIds(command.AccountId).
		WithIds(command.TimeHabitId),
	)
	if err != nil {
		return fmt.Errorf("finding a time habit: %w", err)
	}

	if habit == nil {
		return habits.ErrHabitNotFound
	}

	if err := handler.savingPolicy.Check(command.AccountId, date, habit); err != nil {
		return err
	}

	record, err := handler.timeHabitRecords.Find(ctx, records.NewTimeHabitRecordFilter().
		WithAccountIds(command.AccountId).
		WithIds(command.TimeHabitId).
		WithDates(time.Time(date)),
	)
	if err != nil {
		return fmt.Errorf("finding a time habit record: %w", err)
	}

	isNew := false
	if record == nil {
		isNew = true
		record, err = records.NewTimeHabitRecord(command.TimeHabitId, command.AccountId, time.Time(date), command.Value)
		if err != nil {
			var errs domain.ValidationErrors
			if errors.As(err, &errs) {
				return errs
			}

			return fmt.Errorf("creating a new time habit record: %w", err)
		}
	}

	if err := record.ChangeValue(command.Value); err != nil {
		errs := domain.NewValidationErrors()
		errs.Add("value", err)

		return errs
	}

	if isNew {
		if err := handler.timeHabitRecords.Add(ctx, record); err != nil {
			return fmt.Errorf("adding a new time habit record: %w", err)
		}
	} else {
		if err := handler.timeHabitRecords.Save(ctx, record); err != nil {
			return fmt.Errorf("saving a time habit record: %w", err)
		}
	}

	return nil
}

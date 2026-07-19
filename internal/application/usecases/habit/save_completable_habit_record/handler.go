package save_completable_habit_record

import (
	"context"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/habits"
	"github.com/EugeneNail/lifeline/internal/domain/habits/records"
	"github.com/google/uuid"
)

// Handler executes the save-completable-habit-record use case.
type Handler struct {
	completableHabitRecords records.CompletableHabitRecordRepository
	completableHabits       habits.CompletableHabitRepository
	savingPolicy            *records.SavingPolicy
}

// NewHandler returns a save-completable-habit-record handler configured with the completable habit repository, the record repository, and the saving policy or an error when a dependency is missing.
func NewHandler(completableHabitRecords records.CompletableHabitRecordRepository, completableHabits habits.CompletableHabitRepository, savingPolicy *records.SavingPolicy) (*Handler, error) {
	if completableHabitRecords == nil {
		return nil, fmt.Errorf("save_completable_habit_record handler requires a completable habit record repository")
	}

	if completableHabits == nil {
		return nil, fmt.Errorf("save_completable_habit_record handler requires a completable habit repository")
	}

	if savingPolicy == nil {
		return nil, fmt.Errorf("save_completable_habit_record handler requires a saving policy")
	}

	return &Handler{
		completableHabitRecords: completableHabitRecords,
		completableHabits:       completableHabits,
		savingPolicy:            savingPolicy,
	}, nil
}

// Command carries the data required to save a completable habit record.
type Command struct {
	AccountId          uuid.UUID
	Value              bool
	Date               time.Time
	CompletableHabitId uuid.UUID
}

// Handle updates the existing completable habit record value or creates a new record when none exists, and returns an error when persistence fails.
func (handler *Handler) Handle(ctx context.Context, command Command) error {
	date := records.NewDate(command.Date)

	habit, err := handler.completableHabits.Find(ctx, habits.NewCompletableHabitFilter().
		WithAccountIds(command.AccountId).
		WithIds(command.CompletableHabitId),
	)
	if err != nil {
		return fmt.Errorf("finding a completable habit: %w", err)
	}

	if habit == nil {
		return habits.ErrHabitNotFound
	}

	if err := handler.savingPolicy.Check(command.AccountId, date, habit); err != nil {
		return err
	}

	record, err := handler.completableHabitRecords.Find(ctx, records.NewCompletableHabitRecordFilter().
		WithAccountIds(command.AccountId).
		WithIds(command.CompletableHabitId).
		WithDates(time.Time(date)),
	)
	if err != nil {
		return fmt.Errorf("finding a completable habit record: %w", err)
	}

	isNew := false
	if record == nil {
		isNew = true
		record = records.NewCompletableHabitRecord(command.CompletableHabitId, command.AccountId, time.Time(date), command.Value)
	}

	record.ChangeValue(command.Value)

	if isNew {
		if err := handler.completableHabitRecords.Add(ctx, record); err != nil {
			return fmt.Errorf("adding a new completable habit record: %w", err)
		}
	} else {
		if err := handler.completableHabitRecords.Save(ctx, record); err != nil {
			return fmt.Errorf("saving a completable habit record: %w", err)
		}
	}

	return nil
}

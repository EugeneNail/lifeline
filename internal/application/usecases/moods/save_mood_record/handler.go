package save_mood_record

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/moods"
	"github.com/google/uuid"
)

// Handler executes the save-mood-record use case.
type Handler struct {
	moodRecords moods.RecordRepository
}

// NewHandler returns a save-mood-record handler configured with the mood record repository or an error when a dependency is missing.
func NewHandler(moodRecords moods.RecordRepository) (*Handler, error) {
	if moodRecords == nil {
		return nil, fmt.Errorf("save_mood_record handler requires a mood record repository")
	}

	return &Handler{moodRecords: moodRecords}, nil
}

// Command carries the data required to save a mood record.
type Command struct {
	AccountID uuid.UUID
	Mood      int
	Date      time.Time
}

// Handle updates the existing mood record value or creates a new record when none exists, and returns an error when persistence fails.
func (handler *Handler) Handle(ctx context.Context, command Command) error {
	violations := domain.NewViolations()

	date, err := domain.NewDate(command.Date)
	if err != nil {
		var violation domain.Violation
		if !errors.As(err, &violation) {
			return fmt.Errorf("creating a mood record date: %w", err)
		}

		violations.Add("date", violation)
	}

	value, err := moods.New(command.Mood)
	if err != nil {
		var violation domain.Violation
		if !errors.As(err, &violation) {
			return fmt.Errorf("creating a mood record value: %w", err)
		}

		violations.Add("mood", violation)
	}

	if violations.HasViolations() {
		return violations
	}

	record, err := handler.moodRecords.Find(ctx, moods.NewRecordFilter().
		WithAccountIds(command.AccountID).
		WithDates(date.Time()),
	)
	if err != nil {
		return fmt.Errorf("finding a mood record for account id %q and date %q: %w", command.AccountID, date.Time(), err)
	}

	isNew := false
	if record == nil {
		isNew = true
		record = moods.NewRecord(date, value, command.AccountID)
	} else {
		// TODO rewrite other upsert usecases like this exact line
		record.ChangeValue(value)
	}

	if isNew {
		if err := handler.moodRecords.Add(ctx, record); err != nil {
			return fmt.Errorf("adding a new mood record: %w", err)
		}
	} else {
		if err := handler.moodRecords.Update(ctx, record); err != nil {
			return fmt.Errorf("updating a mood record for account id %q and date %q: %w", command.AccountID, date.Time(), err)
		}
	}

	return nil
}

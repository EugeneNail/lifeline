package get_mood_record

import (
	"context"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/moods"
	"github.com/google/uuid"
)

// Handler executes the get-mood-record use case.
type Handler struct {
	moodRecords moods.RecordRepository
}

// NewHandler returns a get-mood-record handler configured with the mood record repository or an error when the dependency is missing.
func NewHandler(moodRecords moods.RecordRepository) (*Handler, error) {
	if moodRecords == nil {
		return nil, fmt.Errorf("get_mood_record handler requires a mood record repository")
	}

	return &Handler{moodRecords: moodRecords}, nil
}

// Query carries the data required to load a daily mood record.
type Query struct {
	AccountID uuid.UUID
	Date      time.Time
}

// Handle returns the mood record matching the query or nil when no record exists, or an error when lookup fails.
func (handler *Handler) Handle(ctx context.Context, query Query) (*moods.Record, error) {
	record, err := handler.moodRecords.Find(ctx, moods.NewRecordFilter().
		WithAccountIds(query.AccountID).
		WithDates(query.Date),
	)
	if err != nil {
		return nil, fmt.Errorf("finding a mood record: %w", err)
	}

	return record, nil
}

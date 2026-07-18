package list_habit_records

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EugeneNail/lifeline/internal/application/usecases/list_habit_records"
	"github.com/EugeneNail/lifeline/internal/domain/habits/records"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
)

// Handler adapts the list-habit-records use case to the HTTP transport.
type Handler struct {
	usecase  *list_habit_records.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the list-habit-records use case.
func NewHandler(usecase *list_habit_records.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Output represents the JSON response body for the habit record list.
type Output struct {
	Measurable  []MeasurableHabitRecord  `json:"measurable"`
	Time        []TimeHabitRecord        `json:"time"`
	Completable []CompletableHabitRecord `json:"completable"`
}

// MeasurableHabitRecord represents the public measurable habit record fields returned to the client.
type MeasurableHabitRecord struct {
	HabitId string  `json:"habitId"`
	Value   float32 `json:"value"`
}

// TimeHabitRecord represents the public time habit record fields returned to the client.
type TimeHabitRecord struct {
	HabitId string `json:"habitId"`
	Value   int    `json:"value"`
}

// CompletableHabitRecord represents the public completable habit record fields returned to the client.
type CompletableHabitRecord struct {
	HabitId string `json:"habitId"`
	Value   bool   `json:"value"`
}

// Handle loads the user's habit records for the requested day, maps them to the transport output, and returns an HTTP response.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	accountID, err := handler.identity.AccountID(request)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("extracting account id: %w", err)
	}

	date, err := time.Parse(time.DateOnly, request.PathValue("date"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing date: %w", err)
	}

	result, err := handler.usecase.Handle(request.Context(), list_habit_records.Query{
		AccountId: accountID.Uuid(),
		Date:      date,
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("handling ListHabitRecords query: %w", err)
	}

	return http.StatusOK, Output{
		Measurable:  mapMeasurableHabitRecords(result.Measurable),
		Time:        mapTimeHabitRecords(result.Time),
		Completable: mapCompletableHabitRecords(result.Completable),
	}
}

// mapMeasurableHabitRecords converts measurable domain records into transport records.
func mapMeasurableHabitRecords(recordsList []*records.MeasurableHabitRecord) []MeasurableHabitRecord {
	output := make([]MeasurableHabitRecord, 0, len(recordsList))
	for _, record := range recordsList {
		output = append(output, MeasurableHabitRecord{
			HabitId: record.MeasurableHabitId().String(),
			Value:   record.Value().Value(),
		})
	}

	return output
}

// mapTimeHabitRecords converts time domain records into transport records.
func mapTimeHabitRecords(recordsList []*records.TimeHabitRecord) []TimeHabitRecord {
	output := make([]TimeHabitRecord, 0, len(recordsList))
	for _, record := range recordsList {
		output = append(output, TimeHabitRecord{
			HabitId: record.TimeHabitId().String(),
			Value:   record.Value().Value(),
		})
	}

	return output
}

// mapCompletableHabitRecords converts completable domain records into transport records.
func mapCompletableHabitRecords(recordsList []*records.CompletableHabitRecord) []CompletableHabitRecord {
	output := make([]CompletableHabitRecord, 0, len(recordsList))
	for _, record := range recordsList {
		output = append(output, CompletableHabitRecord{
			HabitId: record.CompletableHabitId().String(),
			Value:   record.Value(),
		})
	}

	return output
}

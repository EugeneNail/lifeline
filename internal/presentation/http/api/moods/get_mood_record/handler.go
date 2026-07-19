package get_mood_record

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EugeneNail/lifeline/internal/application/usecases/moods/get_mood_record"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
)

// Handler adapts the get-mood-record use case to the HTTP transport.
type Handler struct {
	usecase  *get_mood_record.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the get-mood-record use case.
func NewHandler(usecase *get_mood_record.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Output represents the JSON response body for a mood record.
type Output struct {
	Date string `json:"date"`
	Mood int    `json:"mood"`
}

// Handle loads the mood record and returns an HTTP response with the public fields or 404 when it is missing.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	accountID, err := handler.identity.AccountID(request)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("extracting account id: %w", err)
	}

	date, err := time.Parse(time.DateOnly, request.PathValue("date"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing date: %w", err)
	}

	record, err := handler.usecase.Handle(request.Context(), get_mood_record.Query{
		AccountID: accountID.Uuid(),
		Date:      date,
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("handling GetMoodRecord query: %w", err)
	}

	if record == nil {
		return http.StatusNotFound, nil
	}

	return http.StatusOK, Output{
		Date: time.Time(record.Date()).Format(time.DateOnly),
		Mood: int(record.Value()),
	}
}

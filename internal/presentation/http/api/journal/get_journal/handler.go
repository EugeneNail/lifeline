package get_journal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EugeneNail/lifeline/internal/application/usecases/journal/get_journal"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
)

// Handler adapts the get-journal use case to the HTTP transport.
type Handler struct {
	usecase  *get_journal.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the get-journal use case.
func NewHandler(usecase *get_journal.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Output represents the JSON response body for a journal.
type Output struct {
	Date string `json:"date"`
	Note string `json:"note"`
}

// Handle loads the journal and returns an HTTP response with the public fields or 404 when it is missing.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	accountID, err := handler.identity.AccountID(request)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("extracting account id: %w", err)
	}

	date, err := time.Parse(time.DateOnly, request.PathValue("date"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing date: %w", err)
	}

	journalEntry, err := handler.usecase.Handle(request.Context(), get_journal.Query{
		AccountID: accountID,
		Date:      date,
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("handling GetJournal query: %w", err)
	}

	if journalEntry == nil {
		return http.StatusNotFound, nil
	}

	return http.StatusOK, Output{
		Date: time.Time(journalEntry.Date()).Format(time.DateOnly),
		Note: string(journalEntry.Note()),
	}
}

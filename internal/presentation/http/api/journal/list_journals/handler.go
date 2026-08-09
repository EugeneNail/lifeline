package list_journals

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/EugeneNail/lifeline/internal/application/usecases/journal/list_journals"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
	"github.com/EugeneNail/lifeline/internal/presentation/http/api/journal/resource"
)

// Handler adapts the list-journals use case to the HTTP transport.
type Handler struct {
	usecase  *list_journals.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the list-journals use case.
func NewHandler(usecase *list_journals.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Output represents the JSON response body for the journal list.
type Output struct {
	Journals []resource.Journal `json:"journals"`
}

// Handle loads journals for the requested date range, maps them to the transport output, and returns an HTTP response.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	accountID, err := handler.identity.AccountID(request)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("extracting account id: %w", err)
	}

	from, err := time.Parse(time.DateOnly, request.URL.Query().Get("from"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing a 'from' date: %w", err)
	}

	to, err := time.Parse(time.DateOnly, request.URL.Query().Get("to"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing a 'to' date: %w", err)
	}

	journalsList, err := handler.usecase.Handle(request.Context(), list_journals.NewQuery(accountID, from, to))
	if err != nil {
		if errors.Is(err, list_journals.ErrInvalidDateRange) {
			return http.StatusBadRequest, err
		}

		return http.StatusInternalServerError, fmt.Errorf("handling ListJournals query: %w", err)
	}

	output := make([]resource.Journal, 0, len(journalsList))
	for _, journal := range journalsList {
		output = append(output, resource.Journal{
			Date:      time.Time(journal.Date()).Format(time.DateOnly),
			Note:      string(journal.Note()),
			CreatedAt: journal.CreatedAt(),
			UpdatedAt: journal.UpdatedAt(),
		})
	}

	return http.StatusOK, Output{
		Journals: output,
	}
}

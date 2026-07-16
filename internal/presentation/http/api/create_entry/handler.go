package create_entry

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EugeneNail/lifeline/internal/domain/entries"
	"net/http"
	"time"

	create_entry "github.com/EugeneNail/lifeline/internal/application/usecases/create_entry"
	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
)

// Handler adapts the create-entry use case to the HTTP transport.
type Handler struct {
	usecase  *create_entry.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the create-entry use case.
func NewHandler(usecase *create_entry.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Payload represents the JSON request body for entry creation.
type Payload struct {
	Date string `json:"date"`
	Mood int    `json:"mood"`
	Note string `json:"note"`
}

// Handle decodes the request, runs the use case, and maps the result to an HTTP response.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	accountID, err := handler.identity.AccountID(request)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("extracting account id: %w", err)
	}

	var payload Payload
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, fmt.Errorf("decoding request body: %w", err)
	}

	date, err := time.Parse(time.DateOnly, payload.Date)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing date: %w", err)
	}

	id, err := handler.usecase.Handle(request.Context(), create_entry.Command{
		Date:      date,
		Mood:      payload.Mood,
		Note:      payload.Note,
		AccountID: accountID,
	})
	if err != nil {
		var validationErrors domain.ValidationErrors
		if errors.As(err, &validationErrors) {
			return http.StatusUnprocessableEntity, validationErrors.Errors()
		}

		if errors.Is(err, entries.ErrDateIsOccupied) {
			return http.StatusConflict, nil
		}

		return http.StatusInternalServerError, fmt.Errorf("handling CreateEntry command: %w", err)
	}

	return http.StatusCreated, id.String()
}

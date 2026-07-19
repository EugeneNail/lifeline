package save_mood_record

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/EugeneNail/lifeline/internal/application/usecases/moods/save_mood_record"
	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
)

// Handler adapts the save-mood-record use case to the HTTP transport.
type Handler struct {
	usecase  *save_mood_record.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the save-mood-record use case.
func NewHandler(usecase *save_mood_record.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Payload represents the JSON request body for a mood record save operation.
type Payload struct {
	Mood int `json:"mood"`
}

// Handle decodes the request, runs the use case, and maps the result to an HTTP response.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	accountID, err := handler.identity.AccountID(request)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("extracting account id: %w", err)
	}

	date, err := time.Parse(time.DateOnly, request.PathValue("date"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing date: %w", err)
	}

	var payload Payload
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, fmt.Errorf("decoding request body: %w", err)
	}

	if err := handler.usecase.Handle(request.Context(), save_mood_record.Command{
		AccountID: accountID.Uuid(),
		Mood:      payload.Mood,
		Date:      date,
	}); err != nil {
		var violations domain.Violations
		if errors.As(err, &violations) {
			return http.StatusUnprocessableEntity, violations.Violations()
		}

		return http.StatusInternalServerError, fmt.Errorf("handling SaveMoodRecord command: %w", err)
	}

	return http.StatusNoContent, nil
}

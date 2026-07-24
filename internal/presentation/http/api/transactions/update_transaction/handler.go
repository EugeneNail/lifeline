package update_transaction

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/EugeneNail/lifeline/internal/application/usecases/transactions/update_transaction"
	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/transactions"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
	"github.com/google/uuid"
)

// Handler adapts the update-transaction use case to the HTTP transport.
type Handler struct {
	usecase  *update_transaction.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the update-transaction use case.
func NewHandler(usecase *update_transaction.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Payload represents the JSON request body for transaction updates.
type Payload struct {
	Money       float32 `json:"money"`
	Date        string  `json:"date"`
	Direction   int     `json:"direction"`
	Category    int     `json:"category"`
	Description string  `json:"description"`
}

// Handle decodes the request, runs the use case, and maps the result to an HTTP response.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	accountID, err := handler.identity.AccountID(request)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("extracting account id: %w", err)
	}

	transactionID, err := uuid.Parse(request.PathValue("id"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing id: %w", err)
	}

	var payload Payload
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		return http.StatusBadRequest, fmt.Errorf("decoding request body: %w", err)
	}

	date, err := time.Parse(time.DateOnly, payload.Date)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing date: %w", err)
	}

	err = handler.usecase.Handle(request.Context(), update_transaction.Command{
		ID:          transactionID,
		Money:       payload.Money,
		Date:        date,
		Direction:   payload.Direction,
		Category:    payload.Category,
		Description: payload.Description,
		AccountID:   accountID.Uuid(),
	})
	if err != nil {
		var violations domain.Violations
		if errors.As(err, &violations) {
			return http.StatusUnprocessableEntity, violations.Violations()
		}

		if errors.Is(err, transactions.ErrTransactionNotFound) {
			return http.StatusNotFound, nil
		}

		return http.StatusInternalServerError, fmt.Errorf("handling UpdateTransaction command: %w", err)
	}

	return http.StatusNoContent, nil
}

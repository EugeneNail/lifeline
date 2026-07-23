package get_transaction

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EugeneNail/lifeline/internal/application/usecases/transactions/get_transaction"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
	"github.com/google/uuid"
)

// Handler adapts the get-transaction use case to the HTTP transport.
type Handler struct {
	usecase  *get_transaction.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the get-transaction use case.
func NewHandler(usecase *get_transaction.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Output represents the JSON response body for a transaction.
type Output struct {
	ID          string  `json:"id"`
	Money       float32 `json:"money"`
	Date        string  `json:"date"`
	Direction   int     `json:"direction"`
	Category    int     `json:"category"`
	Description string  `json:"description"`
}

// Handle loads the transaction and returns an HTTP response with the public fields or 404 when it is missing.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	accountID, err := handler.identity.AccountID(request)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("extracting account id: %w", err)
	}

	transactionID, err := uuid.Parse(request.PathValue("id"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing id: %w", err)
	}

	transaction, err := handler.usecase.Handle(request.Context(), get_transaction.Query{
		AccountID: accountID.Uuid(),
		ID:        transactionID,
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("handling GetTransaction query: %w", err)
	}

	if transaction == nil {
		return http.StatusNotFound, nil
	}

	return http.StatusOK, Output{
		ID:          transaction.ID().String(),
		Money:       float32(transaction.Money()),
		Date:        time.Time(transaction.Date()).Format(time.DateOnly),
		Direction:   int(transaction.Direction()),
		Category:    int(transaction.Category()),
		Description: string(transaction.Description()),
	}
}

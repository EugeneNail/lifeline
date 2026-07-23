package list_transactions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/EugeneNail/lifeline/internal/application/usecases/transactions/list_transactions"
	"github.com/EugeneNail/lifeline/internal/domain/transactions"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
)

// Handler adapts the list-transactions use case to the HTTP transport.
type Handler struct {
	usecase  *list_transactions.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the list-transactions use case.
func NewHandler(usecase *list_transactions.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Output represents the JSON response body for the transaction list.
type Output struct {
	Transactions []Transaction `json:"transactions"`
}

// Transaction represents the public transaction fields returned to the client.
type Transaction struct {
	ID          string  `json:"id"`
	Money       float32 `json:"money"`
	Date        string  `json:"date"`
	Direction   int     `json:"direction"`
	Category    int     `json:"category"`
	Description string  `json:"description"`
}

// Handle loads the user's transactions, maps them to the transport output, and returns an HTTP response.
func (handler *Handler) Handle(request *http.Request) (int, any) {
	accountID, err := handler.identity.AccountID(request)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("extracting account id: %w", err)
	}

	from, err := parseOptionalDate(request.URL.Query().Get("from"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing a 'from' date: %w", err)
	}

	to, err := parseOptionalDate(request.URL.Query().Get("to"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("parsing a 'to' date: %w", err)
	}

	result, err := handler.usecase.Handle(request.Context(), list_transactions.Query{
		AccountID: accountID.Uuid(),
		From:      from,
		To:        to,
	})
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("handling ListTransactions query: %w", err)
	}

	return http.StatusOK, Output{
		Transactions: mapTransactions(result.Transactions),
	}
}

func mapTransactions(transactionsList []*transactions.Transaction) []Transaction {
	output := make([]Transaction, 0, len(transactionsList))
	for _, transaction := range transactionsList {
		output = append(output, Transaction{
			ID:          transaction.ID().String(),
			Money:       float32(transaction.Money()),
			Date:        time.Time(transaction.Date()).Format(time.DateOnly),
			Direction:   int(transaction.Direction()),
			Category:    int(transaction.Category()),
			Description: string(transaction.Description()),
		})
	}

	return output
}

func parseOptionalDate(rawDate string) (*time.Time, error) {
	if rawDate == "" {
		return nil, nil
	}

	date, err := time.Parse(time.DateOnly, rawDate)
	if err != nil {
		return nil, err
	}

	return &date, nil
}

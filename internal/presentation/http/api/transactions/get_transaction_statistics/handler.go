package get_transaction_statistics

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	application "github.com/EugeneNail/lifeline/internal/application/usecases/transactions/get_transaction_statistics"
	"github.com/EugeneNail/lifeline/internal/domain/transactions"
	"github.com/EugeneNail/lifeline/internal/infrastructure/authentication"
	"github.com/EugeneNail/lifeline/internal/presentation/http/api/transactions/resource"
)

// Handler adapts the get-transaction-statistics use case to the HTTP transport.
type Handler struct {
	usecase  *application.Handler
	identity authentication.RequestIdentity
}

// NewHandler returns a transport handler wired to the get-transaction-statistics use case.
func NewHandler(usecase *application.Handler, identity authentication.RequestIdentity) *Handler {
	return &Handler{usecase: usecase, identity: identity}
}

// Handle loads transaction statistics for the requested date range and returns an HTTP response.
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

	result, err := handler.usecase.Handle(request.Context(), application.NewQuery(accountID.Uuid(), from, to))
	if err != nil {
		if errors.Is(err, application.ErrInvalidDateRange) {
			return http.StatusBadRequest, err
		}

		return http.StatusInternalServerError, fmt.Errorf("handling GetTransactionStatistics query: %w", err)
	}

	return http.StatusOK, resource.Statistics{
		Target:   mapTarget(result.Target),
		Baseline: mapBaseline(result.Baseline),
	}
}

// mapTarget converts the target statistics to transport output.
func mapTarget(target application.Target) resource.Target {
	return resource.Target{
		Overview:   mapOverview(target.Overview),
		TopFive:    mapTransactions(target.TopFive),
		Categories: mapCategories(target.Categories),
	}
}

// mapBaseline converts the baseline statistics to transport output.
func mapBaseline(baseline application.Baseline) resource.Baseline {
	return resource.Baseline{
		Overview: mapOverview(baseline.Overview),
	}
}

// mapOverview converts the domain overview to transport output.
func mapOverview(overview application.Overview) resource.Overview {
	return resource.Overview{
		Expenses:  overview.Expenses,
		Incomes:   overview.Incomes,
		NetChange: overview.NetChange,
	}
}

// mapTransactions converts the top-five transactions to transport output.
func mapTransactions(transactionsList []*transactions.Transaction) []resource.Transaction {
	output := make([]resource.Transaction, 0, len(transactionsList))
	for _, transaction := range transactionsList {
		if transaction == nil {
			continue
		}

		output = append(output, resource.Transaction{
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

// mapCategories converts the categorized expenses to transport output.
func mapCategories(categories []application.CategoryExpense) []resource.CategoryExpense {
	output := make([]resource.CategoryExpense, 0, len(categories))
	for _, category := range categories {
		output = append(output, resource.CategoryExpense{
			Absolute: category.Absolute,
			Category: int(category.Category),
			Percent:  category.Percent,
		})
	}

	return output
}

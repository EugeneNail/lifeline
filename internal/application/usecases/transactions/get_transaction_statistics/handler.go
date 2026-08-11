package get_transaction_statistics

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/transactions"
)

const baselineIntervals = 3

// Handler executes the get-transaction-statistics use case.
type Handler struct {
	transactions transactions.Repository
}

// NewHandler returns a get-transaction-statistics handler configured with the transaction repository or an error when the dependency is missing.
func NewHandler(transactionsRepository transactions.Repository) (*Handler, error) {
	if transactionsRepository == nil {
		return nil, fmt.Errorf("get_transaction_statistics handler requires a transaction repository")
	}

	return &Handler{transactions: transactionsRepository}, nil
}

// Handle returns the transaction statistics for the requested range or an error when the use case cannot be executed.
func (handler *Handler) Handle(ctx context.Context, query Query) (Result, error) {
	if query.From.After(query.To) {
		return Result{}, ErrInvalidDateRange
	}

	categories, violations := validateCategories(query.Categories)
	if violations.HasViolations() {
		return Result{}, violations
	}

	targetDates := buildDateRange(query.From, query.To)
	baselineDates := buildBaselineDates(query.From, len(targetDates))
	allDates := mergeDates(targetDates, baselineDates)

	filter := transactions.NewTransactionFilter().
		WithAccountIds(query.AccountID).
		WithCategories(categories...)
	filter.Dates = allDates

	foundTransactions, err := handler.transactions.FindMany(ctx, filter)
	if err != nil {
		return Result{}, fmt.Errorf("finding transactions for account id %q between %q and %q: %w", query.AccountID, query.From, query.To, err)
	}

	targetDateSet := toDateSet(targetDates)
	baselineDateSet := toDateSet(baselineDates)
	targetTransactions := make([]*transactions.Transaction, 0, len(foundTransactions))
	baselineTransactions := make([]*transactions.Transaction, 0, len(foundTransactions))

	for _, transaction := range foundTransactions {
		if transaction == nil {
			continue
		}

		transactionDate := dateKey(time.Time(transaction.Date()))
		switch {
		case targetDateSet[transactionDate]:
			targetTransactions = append(targetTransactions, transaction)
		case baselineDateSet[transactionDate]:
			baselineTransactions = append(baselineTransactions, transaction)
		}
	}

	target := calculateTarget(targetTransactions)

	return Result{
		Target:   target,
		Baseline: calculateBaseline(baselineTransactions),
	}, nil
}

// validateCategories returns validated transaction categories and any violations found in the raw category values.
func validateCategories(rawCategories []int) ([]transactions.Category, domain.Violations) {
	categories := make([]transactions.Category, 0, len(rawCategories))
	violations := domain.NewViolations()

	for _, rawCategory := range rawCategories {
		category, violation := transactions.NewCategory(rawCategory)
		if violation != nil {
			violations.Add("categories", violation)
			continue
		}

		categories = append(categories, category)
	}

	return categories, violations
}

// buildBaselineDates returns the three intervals that precede the target date range.
func buildBaselineDates(from time.Time, intervalDays int) []time.Time {
	baselineFrom := from.AddDate(0, 0, -baselineIntervals*intervalDays)
	baselineTo := from.AddDate(0, 0, -1)

	return buildDateRange(baselineFrom, baselineTo)
}

// buildDateRange returns all dates in the inclusive range normalized to midnight.
func buildDateRange(from time.Time, to time.Time) []time.Time {
	dates := make([]time.Time, 0)
	for currentDate := from; !currentDate.After(to); currentDate = currentDate.AddDate(0, 0, 1) {
		dates = append(dates, currentDate)
	}

	return dates
}

// mergeDates returns a de-duplicated copy of all provided dates.
func mergeDates(dateGroups ...[]time.Time) []time.Time {
	merged := make([]time.Time, 0)
	seen := make(map[time.Time]struct{})

	for _, dates := range dateGroups {
		for _, date := range dates {
			if _, ok := seen[date]; ok {
				continue
			}

			seen[date] = struct{}{}
			merged = append(merged, date)
		}
	}

	return merged
}

// toDateSet returns a lookup set for the provided dates.
func toDateSet(dates []time.Time) map[string]bool {
	set := make(map[string]bool, len(dates))
	for _, date := range dates {
		set[dateKey(date)] = true
	}

	return set
}

// dateKey returns the provided time normalized to a date-only lookup key.
func dateKey(date time.Time) string {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC).Format(time.DateOnly)
}

// calculateTarget returns aggregated transaction statistics and the top five expense transactions for the provided transactions.
func calculateTarget(transactionsList []*transactions.Transaction) Target {
	return Target{
		Overview:   calculateOverview(transactionsList),
		TopFive:    calculateTopFive(transactionsList),
		Categories: calculateCategories(transactionsList),
	}
}

// calculateOverview returns the aggregated expenses, incomes, and net change for the provided transactions.
func calculateOverview(transactionsList []*transactions.Transaction) Overview {
	overview := Overview{}

	for _, transaction := range transactionsList {
		if transaction == nil {
			continue
		}

		amount := float64(transaction.Money())
		switch transaction.Direction() {
		case transactions.Expense:
			overview.Expenses += amount
		case transactions.Income:
			overview.Incomes += amount
		}
	}

	overview.NetChange = overview.Incomes - overview.Expenses

	return overview
}

// calculateTopFive returns up to five expense transactions sorted by descending amount.
func calculateTopFive(transactionsList []*transactions.Transaction) []*transactions.Transaction {
	expenses := make([]*transactions.Transaction, 0)

	for _, transaction := range transactionsList {
		if transaction == nil {
			continue
		}

		if transaction.Direction() != transactions.Expense {
			continue
		}

		expenses = append(expenses, transaction)
	}

	sort.Slice(expenses, func(i, j int) bool {
		return expenses[i].Money() > expenses[j].Money()
	})

	if len(expenses) > 5 {
		expenses = expenses[:5]
	}

	return expenses
}

// calculateCategories returns expense statistics grouped by transaction category and sorted by descending absolute amount.
func calculateCategories(transactionsList []*transactions.Transaction) []CategoryExpense {
	expensesByCategory := make(map[transactions.Category]float64)
	totalExpenses := 0.0

	for _, transaction := range transactionsList {
		if transaction == nil {
			continue
		}

		if transaction.Direction() != transactions.Expense {
			continue
		}

		amount := float64(transaction.Money())
		totalExpenses += amount
		expensesByCategory[transaction.Category()] += amount
	}

	categories := make([]CategoryExpense, 0, len(expensesByCategory))
	for category, absolute := range expensesByCategory {
		percent := 0
		if totalExpenses > 0 {
			percent = int(math.Round((absolute / totalExpenses) * 100))
		}

		categories = append(categories, CategoryExpense{
			Absolute: absolute,
			Category: category,
			Percent:  percent,
		})
	}

	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Absolute > categories[j].Absolute
	})

	return categories
}

// calculateBaseline returns the average statistics across the three baseline intervals.
func calculateBaseline(transactionsList []*transactions.Transaction) Baseline {
	if baselineIntervals <= 0 {
		return Baseline{}
	}

	statistics := calculateTarget(transactionsList)

	return Baseline{
		Overview: Overview{
			Expenses:  statistics.Overview.Expenses / float64(baselineIntervals),
			Incomes:   statistics.Overview.Incomes / float64(baselineIntervals),
			NetChange: statistics.Overview.NetChange / float64(baselineIntervals),
		},
	}
}

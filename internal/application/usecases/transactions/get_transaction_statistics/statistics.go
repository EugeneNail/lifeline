package get_transaction_statistics

import "github.com/EugeneNail/lifeline/internal/domain/transactions"

// Target groups transaction statistics for the selected period.
type Target struct {
	Overview   Overview
	TopFive    []*transactions.Transaction
	Categories []CategoryExpense
}

// Baseline groups transaction statistics averaged across the three comparison intervals.
type Baseline struct {
	Overview Overview
}

// CategoryExpense represents aggregated expense statistics for a single transaction category.
type CategoryExpense struct {
	Absolute float64
	Category transactions.Category
	Percent  int
}

// Overview describes the aggregated transaction amounts for a range.
type Overview struct {
	Expenses  float64
	Incomes   float64
	NetChange float64
}

// Result carries the transaction statistics returned by the use case.
type Result struct {
	Target   Target
	Baseline Baseline
}

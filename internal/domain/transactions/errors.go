package transactions

import "github.com/EugeneNail/lifeline/internal/domain"

// ErrTransactionNotFound reports that a transaction with the requested identifier could not be found.
var ErrTransactionNotFound = domain.NewViolation("transaction not found")

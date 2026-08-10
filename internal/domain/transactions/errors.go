package transactions

import "github.com/EugeneNail/lifeline/internal/domain"

// ErrTransactionNotFound reports that a transaction with the requested identifier could not be found.
var ErrTransactionNotFound = domain.NewViolation("transaction not found")

// ErrTransactionBelongsToAnotherUser reports that the transaction belongs to a different user.
var ErrTransactionBelongsToAnotherUser = domain.NewViolation("transaction belongs to another user")

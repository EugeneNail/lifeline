package auth

import "context"

// AccountRepository stores and retrieves auth accounts.
type AccountRepository interface {
	Add(ctx context.Context, account *Account) error
	FindByEmail(ctx context.Context, email Email) (*Account, error)
}

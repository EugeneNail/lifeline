package auth

import "time"

// Account represents a registered application user.
type Account struct {
	id        ID
	email     Email
	password  HashedPassword
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewAccount returns an account with the provided identifier, email, and password.
func NewAccount(id ID, email Email, password HashedPassword, createdAt time.Time) *Account {
	return &Account{id: id, email: email, password: password, CreatedAt: createdAt}
}

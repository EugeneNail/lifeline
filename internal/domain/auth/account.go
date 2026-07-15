package auth

import (
	"github.com/google/uuid"
	"time"
)

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
	// TODO move field creation from usecases into the constructor
	return &Account{id: id, email: email, password: password, CreatedAt: createdAt}
}

// RestoreAccount returns an account reconstructed from persisted values.
func RestoreAccount(id uuid.UUID, email string, password string, createdAt time.Time, updatedAt time.Time) *Account {
	return &Account{
		id:        ID(id),
		email:     Email(email),
		password:  HashedPassword(password),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

package auth

import (
	"github.com/google/uuid"
	"github.com/samborkent/uuidv7"
)

// ID represents a user identifier stored as UUIDv7.
type ID uuid.UUID

// NilID is the zero-value user identifier.
var NilID = ID(uuid.Nil)

// NewID returns a new UUIDv7-based user identifier.
func NewID() ID {
	return ID(uuidv7.New())
}

// ID returns the user identifier.
func (account *Account) ID() ID {
	return account.id
}

// String returns the UUID string representation of the identifier.
func (id ID) String() string {
	return uuid.UUID(id).String()
}

// Uuid returns the UUID value of the identifier.
func (id ID) Uuid() uuid.UUID {
	return uuid.UUID(id)
}

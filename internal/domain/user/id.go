package user

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
func (user *User) ID() ID {
	return user.id
}

// ToString returns the UUID string representation of the identifier.
func (id ID) ToString() string {
	return uuid.UUID(id).String()
}

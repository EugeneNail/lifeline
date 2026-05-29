package encryption

import (
	"fmt"

	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"golang.org/x/crypto/bcrypt"
)

// BcryptPasswordHasher hashes passwords using bcrypt.
type BcryptPasswordHasher struct {
	cost int
}

// NewBcryptPasswordHasher returns a bcrypt password hasher with bcrypt.DefaultCost.
// TODO pass cost as an argument
func NewBcryptPasswordHasher() *BcryptPasswordHasher {
	return &BcryptPasswordHasher{cost: bcrypt.DefaultCost}
}

// Hash returns a bcrypt hash for the provided password.
func (hasher *BcryptPasswordHasher) Hash(password auth.Password) (auth.HashedPassword, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hasher.cost)
	if err != nil {
		return "", fmt.Errorf("generating bcrypt hash: %w", err)
	}

	return auth.HashedPassword(hashedPassword), nil
}

package encryption

import (
	"fmt"

	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"golang.org/x/crypto/bcrypt"
)

// BcryptPasswordVerifier compares raw passwords against bcrypt hashes.
type BcryptPasswordVerifier struct{}

// NewBcryptPasswordVerifier returns a bcrypt password verifier.
func NewBcryptPasswordVerifier() *BcryptPasswordVerifier {
	return &BcryptPasswordVerifier{}
}

// Verify returns nil when the provided password matches the stored hash.
func (verifier *BcryptPasswordVerifier) Verify(password auth.Password, hashedPassword auth.HashedPassword) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return fmt.Errorf("comparing bcrypt hash: %w", err)
	}

	return nil
}

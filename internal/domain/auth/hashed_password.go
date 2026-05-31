package auth

// HashedPassword represents a password hashed by the configured password hasher.
type HashedPassword string

// PasswordHasher hashes a raw password and returns the hashed representation or an error.
type PasswordHasher interface {
	Hash(password Password) (HashedPassword, error)
}

// PasswordVerifier compares a raw password against a stored hash.
type PasswordVerifier interface {
	Verify(password Password, hashedPassword HashedPassword) error
}

// String returns the string representation of the hashed password.
func (password HashedPassword) String() string {
	return string(password)
}

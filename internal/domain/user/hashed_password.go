package user

// HashedPassword represents a password hashed by the configured password hasher.
type HashedPassword string

// PasswordHasher hashes a raw password and returns the hashed representation or an error.
type PasswordHasher interface {
	Hash(password Password) (HashedPassword, error)
}

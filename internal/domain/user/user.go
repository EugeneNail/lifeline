package user

// User represents a registered application user.
type User struct {
	id       ID
	email    Email
	password HashedPassword
}

// NewUser returns a user with the provided identifier, email, and password.
func NewUser(id ID, email Email, password HashedPassword) *User {
	return &User{id: id, email: email, password: password}
}

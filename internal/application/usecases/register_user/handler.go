package register_user

import (
	"context"
	"fmt"
	"github.com/EugeneNail/lifeline/internal/application"
	"log"

	"github.com/EugeneNail/lifeline/internal/domain/auth"
)

// Handler executes the register-user use case.
type Handler struct {
	passwordHasher auth.PasswordHasher
}

// RegisterUserCommand carries the data required to register a user.
type RegisterUserCommand struct {
	Email                string
	Password             string
	PasswordConfirmation string
}

// NewHandler returns a registration handler configured with the password hasher or an error when the dependency is missing.
func NewHandler(passwordHasher auth.PasswordHasher) (*Handler, error) {
	if passwordHasher == nil {
		return nil, fmt.Errorf("register_user handler requires a password hasher")
	}

	return &Handler{passwordHasher: passwordHasher}, nil
}

// Handle validates the registration command, hashes the password, and returns the new user identifier or field validation errors.
func (h *Handler) Handle(ctx context.Context, command RegisterUserCommand) (auth.ID, error) {
	errs := application.NewFieldErrors()

	email, err := auth.NewEmail(command.Email)
	if err := errs.AddFromDomain("email", err); err != nil {
		return auth.NilID, fmt.Errorf("creating an email: %w", err)
	}

	password, err := auth.NewPassword(command.Password)
	if err := errs.AddFromDomain("password", err); err != nil {
		return auth.NilID, fmt.Errorf("creating a password: %w", err)
	}

	if command.Password != command.PasswordConfirmation {
		errs.Add("passwordConfirmation", "password confirmation must match the password")
	}

	if errs.HasErrors() {
		return auth.NilID, errs
	}

	hashedPassword, err := password.Hash(h.passwordHasher)
	if err != nil {
		return auth.NilID, fmt.Errorf("hashing the password: %w", err)
	}

	account := auth.NewAccount(auth.NewID(), email, hashedPassword)
	log.Printf("user: %+v", account)
	// repository.Save(ctx, account)

	return account.ID(), nil
}

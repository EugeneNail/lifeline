package register_user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/auth"
)

// Handler executes the register-user use case.
type Handler struct {
	passwordHasher auth.PasswordHasher
	accounts       auth.AccountRepository
}

// NewHandler returns a registration handler configured with the password hasher or an error when the dependency is missing.
func NewHandler(passwordHasher auth.PasswordHasher, accounts auth.AccountRepository) (*Handler, error) {
	if passwordHasher == nil {
		return nil, fmt.Errorf("register_user handler requires a password hasher")
	}

	if accounts == nil {
		return nil, fmt.Errorf("register_user handler requires an account repository")
	}

	return &Handler{passwordHasher: passwordHasher, accounts: accounts}, nil
}

// RegisterUserCommand carries the data required to register a user.
type RegisterUserCommand struct {
	Email                string
	Password             string
	PasswordConfirmation string
}

// Handle validates the registration command, hashes the password, and returns the new user identifier or field validation errors.
func (h *Handler) Handle(ctx context.Context, command RegisterUserCommand) (auth.ID, error) {
	violations := domain.NewViolations()

	email, err := auth.NewEmail(command.Email)
	if err != nil {
		var violation domain.Violation
		if !errors.As(err, &violation) {
			return auth.NilID, fmt.Errorf("creating an email: %w", err)
		}

		violations.Add("email", violation)
	}

	password, err := auth.NewPassword(command.Password)
	if err != nil {
		var violation domain.Violation
		if !errors.As(err, &violation) {
			return auth.NilID, fmt.Errorf("creating a password: %w", err)
		}

		violations.Add("password", violation)
	}

	if command.Password != command.PasswordConfirmation {
		violations.Add("passwordConfirmation", domain.NewViolation("password confirmation must match the password"))
	}

	if violations.HasViolations() {
		return auth.NilID, violations
	}

	existingUser, err := h.accounts.FindByEmail(ctx, email)
	if err != nil {
		return auth.NilID, fmt.Errorf("finding account by email: %w", err)
	}

	if existingUser != nil {
		return auth.NilID, auth.EmailAlreadyTaken
	}

	hashedPassword, err := password.Hash(h.passwordHasher)
	if err != nil {
		return auth.NilID, fmt.Errorf("hashing the password: %w", err)
	}

	account := auth.NewAccount(auth.NewID(), email, hashedPassword, time.Now())

	if err := h.accounts.Add(ctx, account); err != nil {
		return auth.NilID, fmt.Errorf("adding a new account to the database: %w", err)
	}

	return account.ID(), nil
}

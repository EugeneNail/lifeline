package authenticate

import (
	"context"
	"errors"
	"fmt"
	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"github.com/EugeneNail/lifeline/internal/infrastructure/config"
)

// Handler executes the authenticate use case.
type Handler struct {
	accounts         auth.AccountRepository
	passwordVerifier auth.PasswordVerifier
	tokenProvider    auth.TokenProvider
	environment      config.Environment
}

// Query carries the data required to authenticate a user.
type Query struct {
	Email       string
	Password    string
	Environment string
}

// Result contains the issued login and refresh tokens.
type Result struct {
	LoginToken   string
	RefreshToken string
}

// NewHandler returns an authentication handler configured with the account repository and token provider or an error when a dependency is missing.
func NewHandler(accounts auth.AccountRepository, passwordVerifier auth.PasswordVerifier, tokenProvider auth.TokenProvider, environment config.Environment) (*Handler, error) {
	if accounts == nil {
		return nil, fmt.Errorf("authenticate handler requires an account repository")
	}

	if passwordVerifier == nil {
		return nil, fmt.Errorf("authenticate handler requires a password verifier")
	}

	if tokenProvider == nil {
		return nil, fmt.Errorf("authenticate handler requires a token provider")
	}

	return &Handler{
		accounts:         accounts,
		passwordVerifier: passwordVerifier,
		tokenProvider:    tokenProvider,
		environment:      environment,
	}, nil
}

// Handle validates the credentials, checks the password against the stored hash, and returns login and refresh tokens or field validation errors.
func (h *Handler) Handle(ctx context.Context, command Query) (Result, error) {
	violations := domain.NewViolations()

	email, err := auth.NewEmail(command.Email)
	if err != nil {
		var violation domain.Violation
		if !errors.As(err, &violation) {
			return Result{}, fmt.Errorf("creating an email: %w", err)
		}

		violations.Add("email", violation)
	}

	password, err := auth.NewPassword(command.Password)
	if err != nil {
		var violation domain.Violation
		if !errors.As(err, &violation) {
			return Result{}, fmt.Errorf("creating a password: %w", err)
		}

		violations.Add("password", violation)
	}

	if violations.HasViolations() {
		return Result{}, violations
	}

	account, err := h.accounts.FindByEmail(ctx, email)
	if err != nil {
		return Result{}, fmt.Errorf("finding account by email: %w", err)
	}

	if account == nil {
		return Result{}, invalidCredentialsErrors()
	}

	if err := password.Verify(account.Password(), h.passwordVerifier); err != nil {
		return Result{}, invalidCredentialsErrors()
	}

	loginToken, err := h.tokenProvider.Provide(account, selectTokenLifecycle(h.environment))
	if err != nil {
		return Result{}, fmt.Errorf("creating login token: %w", err)
	}

	refreshToken, err := h.tokenProvider.Provide(account, auth.TokenLifecycleRefresh)
	if err != nil {
		return Result{}, fmt.Errorf("creating refresh token: %w", err)
	}

	return Result{
		LoginToken:   loginToken.String(),
		RefreshToken: refreshToken.String(),
	}, nil
}

// invalidCredentialsErrors returns the same field-level validation payload used for any authentication failure.
func invalidCredentialsErrors() domain.Violations {
	violations := domain.NewViolations()
	violations.Add("email", domain.NewViolation("Invalid email or password"))
	violations.Add("password", domain.NewViolation("Invalid email or password"))
	return violations
}

// selectTokenLifecycle returns a 1-hour login token lifecycle for development or the default login token lifecycle for production and unknown environments.
func selectTokenLifecycle(environment config.Environment) auth.TokenLifecycle {
	switch environment {
	case config.EnvironmentDevelopment:
		return 60 * 60 * 3
	case config.EnvironmentProduction:
		return auth.TokenLifecycleLogin
	default:
		return auth.TokenLifecycleLogin
	}
}

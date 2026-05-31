package authenticate

import (
	"context"
	"fmt"

	"github.com/EugeneNail/lifeline/internal/application"
	"github.com/EugeneNail/lifeline/internal/domain/auth"
)

// Handler executes the authenticate use case.
type Handler struct {
	accounts         auth.AccountRepository
	passwordVerifier auth.PasswordVerifier
	tokenProvider    auth.TokenProvider
}

// Query carries the data required to authenticate a user.
type Query struct {
	Email    string
	Password string
}

// Result contains the issued login and refresh tokens.
type Result struct {
	LoginToken   string
	RefreshToken string
}

// NewHandler returns an authentication handler configured with the account repository and token provider or an error when a dependency is missing.
func NewHandler(accounts auth.AccountRepository, passwordVerifier auth.PasswordVerifier, tokenProvider auth.TokenProvider) (*Handler, error) {
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
	}, nil
}

// Handle validates the credentials, checks the password against the stored hash, and returns login and refresh tokens or field validation errors.
func (h *Handler) Handle(ctx context.Context, command Query) (Result, error) {
	errs := application.NewFieldErrors()

	email, err := auth.NewEmail(command.Email)
	if err := errs.AddFromDomain("email", err); err != nil {
		return Result{}, fmt.Errorf("creating an email: %w", err)
	}

	password, err := auth.NewPassword(command.Password)
	if err := errs.AddFromDomain("password", err); err != nil {
		return Result{}, fmt.Errorf("creating a password: %w", err)
	}

	if errs.HasErrors() {
		return Result{}, errs
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

	loginToken, err := h.tokenProvider.Provide(account, auth.TokenLifecycleLogin)
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
func invalidCredentialsErrors() application.FieldErrors {
	errs := application.NewFieldErrors()
	errs.Add("email", "Invalid email or password")
	errs.Add("password", "Invalid email or password")
	return errs
}

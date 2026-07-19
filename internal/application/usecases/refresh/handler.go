package refresh

import (
	"context"
	"fmt"

	"github.com/EugeneNail/lifeline/internal/domain"
	"github.com/EugeneNail/lifeline/internal/domain/auth"
)

// Handler executes the refresh use case.
type Handler struct {
	accounts      auth.AccountRepository
	tokenProvider auth.TokenProvider
}

// Command carries the refresh token required to issue a new login token.
type Command struct {
	RefreshToken string
}

// NewHandler returns a refresh handler configured with the account repository and token provider or an error when a dependency is missing.
func NewHandler(accounts auth.AccountRepository, tokenProvider auth.TokenProvider) (*Handler, error) {
	if accounts == nil {
		return nil, fmt.Errorf("refresh handler requires an account repository")
	}

	if tokenProvider == nil {
		return nil, fmt.Errorf("refresh handler requires a token provider")
	}

	return &Handler{
		accounts:      accounts,
		tokenProvider: tokenProvider,
	}, nil
}

// Handle validates the refresh token, loads the account, and returns a new login token or field validation errors.
func (h *Handler) Handle(ctx context.Context, command Command) (auth.Token, error) {
	token, err := h.tokenProvider.Restore(command.RefreshToken)
	if err != nil {
		return nil, invalidRefreshTokenErrors()
	}

	if !token.IsValid() || token.Lifecycle() != auth.TokenLifecycleRefresh {
		return nil, invalidRefreshTokenErrors()
	}

	account, err := h.accounts.FindByID(ctx, token.UserID())
	if err != nil {
		return nil, fmt.Errorf("finding account by id: %w", err)
	}

	if account == nil {
		return nil, invalidRefreshTokenErrors()
	}

	loginToken, err := h.tokenProvider.Provide(account, auth.TokenLifecycleLogin)
	if err != nil {
		return nil, fmt.Errorf("creating login token: %w", err)
	}

	return loginToken, nil
}

// invalidRefreshTokenErrors returns the same field-level validation payload used for any refresh failure.
func invalidRefreshTokenErrors() domain.Violations {
	violations := domain.NewViolations()
	violations.Add("refreshToken", domain.NewViolation("Invalid refresh token"))
	return violations
}

package tokens

import (
	"fmt"
	"time"

	"github.com/EugeneNail/lifeline/internal/domain/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTProvider issues and restores JWT-backed auth tokens.
type JWTProvider struct {
	secret []byte
}

// JWT represents a signed authentication token backed by jwt/v5.
type JWT struct {
	raw       string
	userID    auth.ID
	lifecycle auth.TokenLifecycle
	valid     bool
}

type claims struct {
	UserID    string              `json:"uid"`
	Lifecycle auth.TokenLifecycle `json:"lifecycle"`
	jwt.RegisteredClaims
}

// NewJWTProvider returns a JWT provider configured with the given signing secret.
func NewJWTProvider(secret string) (*JWTProvider, error) {
	if secret == "" {
		return nil, fmt.Errorf("JWT provider requires a signing secret")
	}

	return &JWTProvider{
		secret: []byte(secret),
	}, nil
}

// Provide issues a login token for the provided account.
func (provider *JWTProvider) Provide(account *auth.Account, lifecycle auth.TokenLifecycle) (auth.Token, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		UserID:    account.ID().String(),
		Lifecycle: lifecycle,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   account.ID().String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(lifecycle))),
		},
	})

	rawToken, err := token.SignedString(provider.secret)
	if err != nil {
		return nil, fmt.Errorf("signing JWT: %w", err)
	}

	return JWT{
		raw:       rawToken,
		userID:    account.ID(),
		lifecycle: lifecycle,
		valid:     true,
	}, nil
}

// Restore parses, validates, and reconstructs a JWT token from its serialized form.
func (provider *JWTProvider) Restore(rawToken string) (auth.Token, error) {
	parsedToken, err := jwt.ParseWithClaims(
		rawToken,
		&claims{},
		func(token *jwt.Token) (any, error) {
			return provider.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, fmt.Errorf("restoring JWT: %w", err)
	}

	tokenClaims, ok := parsedToken.Claims.(*claims)
	if !ok {
		return nil, fmt.Errorf("restoring JWT: unexpected claims type")
	}

	parsedUserID, err := uuid.Parse(tokenClaims.UserID)
	if err != nil {
		return nil, fmt.Errorf("restoring JWT: parsing user id: %w", err)
	}

	return JWT{
		raw:       rawToken,
		userID:    auth.ID(parsedUserID),
		lifecycle: tokenClaims.Lifecycle,
		valid:     parsedToken.Valid,
	}, nil
}

// String returns the serialized JWT.
func (token JWT) String() string {
	return token.raw
}

// UserID returns the identifier embedded in the token.
func (token JWT) UserID() auth.ID {
	return token.userID
}

// IsValid reports whether the token was parsed and validated successfully.
func (token JWT) IsValid() bool {
	return token.valid
}

// Lifecycle returns the token lifecycle stored in the claims.
func (token JWT) Lifecycle() auth.TokenLifecycle {
	return token.lifecycle
}

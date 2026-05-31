package auth

type TokenLifecycle int

const (
	TokenLifecycleLogin   = 60 * 15
	TokenLifecycleRefresh = 60 * 60 * 24 * 2
)

// Token represents a signed authentication token.
type Token interface {
	String() string
	UserID() ID
	IsValid() bool
	Lifecycle() TokenLifecycle
}

type TokenProvider interface {
	Provide(account *Account, lifecycle TokenLifecycle) (Token, error)
	Restore(rawToken string) (Token, error)
}

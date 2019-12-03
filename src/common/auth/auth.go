package auth

import (
	"context"
	"errors"
)

// Definition error
var (
	ErrInvalidToken = errors.New("invalid token")
)

// TokenInfo - Token information
type TokenInfo interface {
	// Get access token
	GetAccessToken() string
	// Get token type
	GetTokenType() string
	// Get token expiration timestamp
	GetExpiresAt() int64
	// JSON encoding
	EncodeToJSON() ([]byte, error)
}

// Auther - Authentication interface
type Auther interface {
	// Generate token
	GenerateToken(ctx context.Context, userUUID string) (TokenInfo, error)

	// Destroy token
	DestroyToken(ctx context.Context, accessToken string) error

	// Resolve user ID
	ParseUserUUID(ctx context.Context, accessToken string) (string, error)

	// Release resources
	Release() error
}

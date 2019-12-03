package jwtauth

import (
	"context"
	"time"
)

// Storer - Token storage interface
type Storer interface {
	// Store token data and specify expiration time
	Set(ctx context.Context, tokenString string, expiration time.Duration) error
	// Check if the token exists
	Check(ctx context.Context, tokenString string) (bool, error)
	// Close storage
	Close() error
}

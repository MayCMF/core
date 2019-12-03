package model

import (
	"context"
)

// ITrans - Manage Transaction interface
type ITrans interface {
	// Start transaction
	Begin(ctx context.Context) (interface{}, error)
	// Submit transaction
	Commit(ctx context.Context, trans interface{}) error
	// Rollback transaction
	Rollback(ctx context.Context, trans interface{}) error
}

package controllers

import (
	"context"
)

// ITrans - Manage Transaction interface
type ITrans interface {
	// Execute transaction
	Exec(ctx context.Context, fn func(context.Context) error) error
}

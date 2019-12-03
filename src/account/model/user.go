package model

import (
	"context"

	"github.com/MayCMF/core/src/account/schema"
)

// IUser - User object storage interface
type IUser interface {
	// Query data
	Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error)
	// Query specified data
	Get(ctx context.Context, UUID string, opts ...schema.UserQueryOptions) (*schema.User, error)
	// Create data
	Create(ctx context.Context, item schema.User) error
	// Update data
	Update(ctx context.Context, UUID string, item schema.User) error
	// Delete data
	Delete(ctx context.Context, UUID string) error
	// Update status
	UpdateStatus(ctx context.Context, UUID string, status int) error
	// Update password
	UpdatePassword(ctx context.Context, UUID, password string) error
}

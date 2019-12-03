package controllers

import (
	"context"

	"github.com/MayCMF/core/src/account/schema"
)

// IUser - Manage User business logic interface
type IUser interface {
	// Query user
	Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error)
	// Get specified user
	Get(ctx context.Context, UUID string, opts ...schema.UserQueryOptions) (*schema.User, error)
	// Create user
	Create(ctx context.Context, item schema.User) (*schema.User, error)
	// Update user
	Update(ctx context.Context, UUID string, item schema.User) (*schema.User, error)
	// Delete user
	Delete(ctx context.Context, UUID string) error
	// Update status
	UpdateStatus(ctx context.Context, UUID string, status int) error
	// Query display item user
	QueryShow(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserShowQueryResult, error)
}

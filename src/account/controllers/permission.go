package controllers

import (
	"context"

	"github.com/MayCMF/core/src/account/schema"
)

// IPermission - Manage Permission business logic interface
type IPermission interface {
	// Query Permission
	Query(ctx context.Context, params schema.PermissionQueryParam, opts ...schema.PermissionQueryOptions) (*schema.PermissionQueryResult, error)
	// Query specified Permission
	Get(ctx context.Context, UUID string, opts ...schema.PermissionQueryOptions) (*schema.Permission, error)
	// Create Permission
	Create(ctx context.Context, item schema.Permission) (*schema.Permission, error)
	// update Permission
	Update(ctx context.Context, UUID string, item schema.Permission) (*schema.Permission, error)
	// delete Permission
	Delete(ctx context.Context, UUID string) error
}

package model

import (
	"context"

	"github.com/MayCMF/core/src/account/schema"
)

// IPermission - Manage Permission storage interface
type IPermission interface {
	// Query data
	Query(ctx context.Context, params schema.PermissionQueryParam, opts ...schema.PermissionQueryOptions) (*schema.PermissionQueryResult, error)
	// Query specified data
	Get(ctx context.Context, UUID string, opts ...schema.PermissionQueryOptions) (*schema.Permission, error)
	// Create data
	Create(ctx context.Context, item schema.Permission) error
	// Update data
	Update(ctx context.Context, UUID string, item schema.Permission) error
	// Update parent path
	UpdateParentPath(ctx context.Context, UUID, parentPath string) error
	// Delete data
	Delete(ctx context.Context, UUID string) error
}

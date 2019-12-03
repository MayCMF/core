package model

import (
	"context"

	"github.com/MayCMF/core/src/account/schema"
)

// IRole - Manage Role
type IRole interface {
	// Query data
	Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error)
	// Query specified data
	Get(ctx context.Context, UUID string, opts ...schema.RoleQueryOptions) (*schema.Role, error)
	// Create data
	Create(ctx context.Context, item schema.Role) error
	// Update data
	Update(ctx context.Context, UUID string, item schema.Role) error
	// Delete data
	Delete(ctx context.Context, UUID string) error
}

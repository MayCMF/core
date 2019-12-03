package controllers

import (
	"context"

	"github.com/MayCMF/core/src/account/schema"
)

// IRole - Manage Role business logic interface
type IRole interface {
	// Query role
	Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error)
	// Query specified role
	Get(ctx context.Context, UUID string, opts ...schema.RoleQueryOptions) (*schema.Role, error)
	// Create role
	Create(ctx context.Context, item schema.Role) (*schema.Role, error)
	// Update role
	Update(ctx context.Context, UUID string, item schema.Role) (*schema.Role, error)
	// Delete role
	Delete(ctx context.Context, UUID string) error
	// Get resource permissions
	GetPermissionResources(ctx context.Context, item *schema.Role) (schema.PermissionResources, error)
}

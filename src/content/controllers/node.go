package controllers

import (
	"context"

	"github.com/MayCMF/core/src/primitives/schema"
)

// INode - Node business logic interface
type INode interface {
	// Query data
	Query(ctx context.Context, params schema.NodeQueryParam, opts ...schema.NodeQueryOptions) (*schema.NodeQueryResult, error)
	// Get specified data
	Get(ctx context.Context, UUID string, opts ...schema.NodeQueryOptions) (*schema.Node, error)
	// Create data
	Create(ctx context.Context, item schema.Node) (*schema.Node, error)
	// Update data
	Update(ctx context.Context, UUID string, item schema.Node) (*schema.Node, error)
	// Delete data
	Delete(ctx context.Context, UUID string) error
}

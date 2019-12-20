package controllers

import (
	"context"

	"github.com/MayCMF/core/src/primitives/schema"
)

// IPrimitive - Primitive business logic interface
type IPrimitive interface {
	// Query data
	Query(ctx context.Context, params schema.PrimitiveQueryParam, opts ...schema.PrimitiveQueryOptions) (*schema.PrimitiveQueryResult, error)
	// Get specified data
	Get(ctx context.Context, UUID string, opts ...schema.PrimitiveQueryOptions) (*schema.Primitive, error)
	// Create data
	Create(ctx context.Context, item schema.Primitive) (*schema.Primitive, error)
	// Update data
	Update(ctx context.Context, UUID string, item schema.Primitive) (*schema.Primitive, error)
	// Delete data
	Delete(ctx context.Context, UUID string) error
}

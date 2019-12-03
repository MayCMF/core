package model

import (
	"context"

	"github.com/MayCMF/src/filemanager/schema"
)

// IExample Example storage interface
type IExample interface {
	// Query data
	Query(ctx context.Context, params schema.ExampleQueryParam, opts ...schema.ExampleQueryOptions) (*schema.ExampleQueryResult, error)
	// Query specified data
	Get(ctx context.Context, UUID string, opts ...schema.ExampleQueryOptions) (*schema.Example, error)
	// Create data
	Create(ctx context.Context, item schema.Example) error
	// Update data
	Update(ctx context.Context, UUID string, item schema.Example) error
	// Delete data
	Delete(ctx context.Context, UUID string) error
}

package model

import (
	"context"

	"github.com/MayCMF/core/src/filemanager/schema"
)

// IFile File storage interface
type IFile interface {
	// Query data
	Query(ctx context.Context, params schema.FileQueryParam, opts ...schema.FileQueryOptions) (*schema.FileQueryResult, error)
	// Query specified data
	Get(ctx context.Context, UUID string, opts ...schema.FileQueryOptions) (*schema.File, error)
	// Create data
	Create(ctx context.Context, item schema.File) error
	// Update data
	Update(ctx context.Context, UUID string, item schema.File) error
	// Delete data
	Delete(ctx context.Context, UUID string) error
	// Upload File
	Upload(ctx context.Context, item schema.File) error
}

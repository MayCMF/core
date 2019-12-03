package model

import (
	"context"

	"github.com/MayCMF/core/src/i18n/schema"
)

// ILanguage Language storage interface
type ILanguage interface {
	// Query data
	Query(ctx context.Context, params schema.LanguageQueryParam, opts ...schema.LanguageQueryOptions) (*schema.LanguageQueryResult, error)
	// Query specified data
	Get(ctx context.Context, Code string, opts ...schema.LanguageQueryOptions) (*schema.Language, error)
	// Create data
	Create(ctx context.Context, item schema.Language) error
	// Update data
	Update(ctx context.Context, Code string, item schema.Language) error
	// Delete data
	Delete(ctx context.Context, Code string) error
	// Update status
	UpdateStatus(ctx context.Context, Code string, status int) error
}

package model

import (
	"context"

	"github.com/MayCMF/core/src/i18n/schema"
)

// ICountry - Country storage interface
type ICountry interface {
	// Query data
	Query(ctx context.Context, params schema.CountryQueryParam, opts ...schema.CountryQueryOptions) (*schema.CountryQueryResult, error)
	// Query specified data
	Get(ctx context.Context, Code string, opts ...schema.CountryQueryOptions) (*schema.Country, error)
	// Create data
	Create(ctx context.Context, item schema.Country) error
	// Update data
	Update(ctx context.Context, Code string, item schema.Country) error
	// Delete data
	Delete(ctx context.Context, Code string) error
	// Update status
	UpdateStatus(ctx context.Context, Code string, status int) error
}

package implement

import (
	"context"

	"github.com/MayCMF/core/src/common/errors"
	comschema "github.com/MayCMF/core/src/common/schema"
	"github.com/MayCMF/core/src/i18n/model"
	"github.com/MayCMF/core/src/i18n/schema"
)

// NewCountry - Create a Country
func NewCountry(mCountry model.ICountry) *Country {
	return &Country{
		CountryModel: mCountry,
	}
}

// Country - Sample program
type Country struct {
	CountryModel model.ICountry
}

// Query - Query data
func (a *Country) Query(ctx context.Context, params schema.CountryQueryParam, opts ...schema.CountryQueryOptions) (*schema.CountryQueryResult, error) {
	return a.CountryModel.Query(ctx, params, opts...)
}

// Get - Get specified data
func (a *Country) Get(ctx context.Context, code string, opts ...schema.CountryQueryOptions) (*schema.Country, error) {
	item, err := a.CountryModel.Get(ctx, code, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

func (a *Country) checkCode(ctx context.Context, code string) error {
	result, err := a.CountryModel.Query(ctx, schema.CountryQueryParam{
		Code: code,
	}, schema.CountryQueryOptions{
		PageParam: &comschema.PaginationParam{PageSize: -1},
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("Number already exists")
	}
	return nil
}

func (a *Country) getUpdate(ctx context.Context, code string) (*schema.Country, error) {
	return a.Get(ctx, code)
}

// Create - Create Country data
func (a *Country) Create(ctx context.Context, item schema.Country) (*schema.Country, error) {
	err := a.checkCode(ctx, item.Code)
	if err != nil {
		return nil, err
	}

	err = a.CountryModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	return a.getUpdate(ctx, item.Code)
}

// Update - Update Country data
func (a *Country) Update(ctx context.Context, code string, item schema.Country) (*schema.Country, error) {
	oldItem, err := a.CountryModel.Get(ctx, code)
	if err != nil {
		return nil, err
	} else if oldItem == nil {
		return nil, errors.ErrNotFound
	} else if oldItem.Code != item.Code {
		err := a.checkCode(ctx, item.Code)
		if err != nil {
			return nil, err
		}
	}

	err = a.CountryModel.Update(ctx, code, item)
	if err != nil {
		return nil, err
	}
	return a.getUpdate(ctx, code)
}

// Delete - Delete data
func (a *Country) Delete(ctx context.Context, code string) error {
	oldItem, err := a.CountryModel.Get(ctx, code)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.CountryModel.Delete(ctx, code)
}

// UpdateStatus - Update status
func (a *Country) UpdateStatus(ctx context.Context, code string, status int) error {
	oldItem, err := a.CountryModel.Get(ctx, code)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.CountryModel.UpdateStatus(ctx, code, status)
}

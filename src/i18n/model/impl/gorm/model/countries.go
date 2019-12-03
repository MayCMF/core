package model

import (
	"context"

	"github.com/MayCMF/core/src/common/errors"
	"github.com/MayCMF/core/src/common/model"
	"github.com/MayCMF/core/src/i18n/model/impl/gorm/entity"
	"github.com/MayCMF/core/src/i18n/schema"
	"github.com/jinzhu/gorm"
)

// NewCountry - Create a country storage instance
func NewCountry(db *gorm.DB) *Country {
	return &Country{db}
}

// Country - Country storage
type Country struct {
	db *gorm.DB
}

func (a *Country) getQueryOption(opts ...schema.CountryQueryOptions) schema.CountryQueryOptions {
	var opt schema.CountryQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query - Query data
func (a *Country) Query(ctx context.Context, params schema.CountryQueryParam, opts ...schema.CountryQueryOptions) (*schema.CountryQueryResult, error) {
	db := entity.GetCountryDB(ctx, a.db)
	if v := params.Code; v != "" {
		db = db.Where("code=?", v)
	}
	if v := params.LikeCode; v != "" {
		db = db.Where("code LIKE ?", "%"+v+"%")
	}
	if v := params.LikeName; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	db = db.Order("code ASC")

	opt := a.getQueryOption(opts...)
	var list entity.Countries
	pr, err := model.WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.CountryQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaCountries(),
	}

	return qr, nil
}

// Get - Query specified data
func (a *Country) Get(ctx context.Context, code string, opts ...schema.CountryQueryOptions) (*schema.Country, error) {
	db := entity.GetCountryDB(ctx, a.db).Where("code=?", code)
	var item entity.Country
	ok, err := model.FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaCountry(), nil
}

// Create - Create data
func (a *Country) Create(ctx context.Context, item schema.Country) error {
	country := entity.SchemaCountry(item).ToCountry()
	result := entity.GetCountryDB(ctx, a.db).Create(country)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update - Update data
func (a *Country) Update(ctx context.Context, code string, item schema.Country) error {
	country := entity.SchemaCountry(item).ToCountry()
	result := entity.GetCountryDB(ctx, a.db).Where("code=?", code).Omit("record_id", "creator").Updates(country)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete - delete data
func (a *Country) Delete(ctx context.Context, code string) error {
	result := entity.GetCountryDB(ctx, a.db).Where("code=?", code).Delete(entity.Country{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus - update status
func (a *Country) UpdateStatus(ctx context.Context, code string, status int) error {
	result := entity.GetCountryDB(ctx, a.db).Where("code=?", code).Update("active", status)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

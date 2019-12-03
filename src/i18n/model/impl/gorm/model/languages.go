package model

import (
	"context"

	"github.com/MayCMF/core/src/common/errors"
	"github.com/MayCMF/core/src/common/model"
	"github.com/MayCMF/core/src/i18n/model/impl/gorm/entity"
	"github.com/MayCMF/core/src/i18n/schema"
	"github.com/jinzhu/gorm"
)

// NewLanguage - Create a language storage instance
func NewLanguage(db *gorm.DB) *Language {
	return &Language{db}
}

// Language - Language storage
type Language struct {
	db *gorm.DB
}

func (a *Language) getQueryOption(opts ...schema.LanguageQueryOptions) schema.LanguageQueryOptions {
	var opt schema.LanguageQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query - Query data
func (a *Language) Query(ctx context.Context, params schema.LanguageQueryParam, opts ...schema.LanguageQueryOptions) (*schema.LanguageQueryResult, error) {
	db := entity.GetLanguageDB(ctx, a.db)
	if v := params.Code; v != "" {
		db = db.Where("code=?", v)
	}
	if v := params.LikeCode; v != "" {
		db = db.Where("code LIKE ?", "%"+v+"%")
	}
	if v := params.LikeName; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.Status; v > 0 {
		db = db.Where("active=?", v)
	}
	db = db.Order("code DESC")

	opt := a.getQueryOption(opts...)
	var list entity.Languages
	pr, err := model.WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.LanguageQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaLanguages(),
	}

	return qr, nil
}

// Get - Query specified data
func (a *Language) Get(ctx context.Context, code string, opts ...schema.LanguageQueryOptions) (*schema.Language, error) {
	db := entity.GetLanguageDB(ctx, a.db).Where("code=?", code)
	var item entity.Language
	ok, err := model.FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaLanguage(), nil
}

// Create - Create data
func (a *Language) Create(ctx context.Context, item schema.Language) error {
	language := entity.SchemaLanguage(item).ToLanguage()
	result := entity.GetLanguageDB(ctx, a.db).Create(language)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update - Update data
func (a *Language) Update(ctx context.Context, code string, item schema.Language) error {
	language := entity.SchemaLanguage(item).ToLanguage()
	result := entity.GetLanguageDB(ctx, a.db).Where("code=?", code).Omit("record_id", "creator").Updates(language)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete - delete data
func (a *Language) Delete(ctx context.Context, code string) error {
	result := entity.GetLanguageDB(ctx, a.db).Where("code=?", code).Delete(entity.Language{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus - update status
func (a *Language) UpdateStatus(ctx context.Context, code string, status int) error {
	result := entity.GetLanguageDB(ctx, a.db).Where("code=?", code).Update("active", status)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

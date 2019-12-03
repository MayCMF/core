package implement

import (
	"context"

	"github.com/MayCMF/core/src/common/errors"
	comschema "github.com/MayCMF/core/src/common/schema"
	"github.com/MayCMF/core/src/i18n/model"
	"github.com/MayCMF/core/src/i18n/schema"
)

// NewLanguage - Create a Language
func NewLanguage(mLanguage model.ILanguage) *Language {
	return &Language{
		LanguageModel: mLanguage,
	}
}

// Language - Sample program
type Language struct {
	LanguageModel model.ILanguage
}

// Query - Query data
func (a *Language) Query(ctx context.Context, params schema.LanguageQueryParam, opts ...schema.LanguageQueryOptions) (*schema.LanguageQueryResult, error) {
	return a.LanguageModel.Query(ctx, params, opts...)
}

// Get - Get specified data
func (a *Language) Get(ctx context.Context, code string, opts ...schema.LanguageQueryOptions) (*schema.Language, error) {
	item, err := a.LanguageModel.Get(ctx, code, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

func (a *Language) checkCode(ctx context.Context, code string) error {
	result, err := a.LanguageModel.Query(ctx, schema.LanguageQueryParam{
		Code: code,
	}, schema.LanguageQueryOptions{
		PageParam: &comschema.PaginationParam{PageSize: -1},
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("Language already exists")
	}
	return nil
}

func (a *Language) getUpdate(ctx context.Context, code string) (*schema.Language, error) {
	return a.Get(ctx, code)
}

// Create - Create Language data
func (a *Language) Create(ctx context.Context, item schema.Language) (*schema.Language, error) {
	err := a.checkCode(ctx, item.Code)
	if err != nil {
		return nil, err
	}

	err = a.LanguageModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	return a.getUpdate(ctx, item.Code)
}

// Update - Update Language data
func (a *Language) Update(ctx context.Context, code string, item schema.Language) (*schema.Language, error) {
	oldItem, err := a.LanguageModel.Get(ctx, code)
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

	err = a.LanguageModel.Update(ctx, code, item)
	if err != nil {
		return nil, err
	}
	return a.getUpdate(ctx, code)
}

// Delete - Delete data
func (a *Language) Delete(ctx context.Context, code string) error {
	oldItem, err := a.LanguageModel.Get(ctx, code)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.LanguageModel.Delete(ctx, code)
}

// UpdateStatus - Update status
func (a *Language) UpdateStatus(ctx context.Context, code string, status int) error {
	oldItem, err := a.LanguageModel.Get(ctx, code)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.LanguageModel.UpdateStatus(ctx, code, status)
}

package model

import (
	"context"

	"github.com/MayCMF/core/src/common/errors"
	"github.com/MayCMF/core/src/common/model"
	"github.com/MayCMF/core/src/primitives/model/impl/gorm/entity"
	"github.com/MayCMF/core/src/primitives/schema"
	"github.com/jinzhu/gorm"
)

// NewPrimitive - Create a primitive storage instance
func NewPrimitive(db *gorm.DB) *Primitive {
	return &Primitive{db}
}

// Primitive - Primitive storage
type Primitive struct {
	db *gorm.DB
}

func (a *Primitive) getQueryOption(opts ...schema.PrimitiveQueryOptions) schema.PrimitiveQueryOptions {
	var opt schema.PrimitiveQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query - Query data
func (a *Primitive) Query(ctx context.Context, params schema.PrimitiveQueryParam, opts ...schema.PrimitiveQueryOptions) (*schema.PrimitiveQueryResult, error) {
	db := entity.GetPrimitiveDB(ctx, a.db)
	if v := params.Slug; v != "" {
		db = db.Where("slug=?", v)
	}
	if v := params.LikeSlug; v != "" {
		db = db.Where("slug LIKE ?", "%"+v+"%")
	}
	db = db.Order("id DESC")

	opt := a.getQueryOption(opts...)
	var list entity.Primitives
	pr, err := model.WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.PrimitiveQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaPrimitives(),
	}

	err = a.fillSchemaPrimitives(ctx, qr.Data, opts...)
	if err != nil {
		return nil, err
	}

	return qr, nil
}

// Populate Primitive object data
func (a *Primitive) fillSchemaPrimitives(ctx context.Context, items []*schema.Primitive, opts ...schema.PrimitiveQueryOptions) error {
	opt := a.getQueryOption(opts...)

	if opt.IncludeVariations {

		primitiveIDs := make([]string, len(items))
		for i, item := range items {
			primitiveIDs[i] = item.Slug
		}

		var bodyList entity.Variations
		if opt.IncludeVariations {
			items, err := a.queryVariations(ctx, primitiveIDs...)
			if err != nil {
				return err
			}
			bodyList = items
		}

		for i, item := range items {
			if len(bodyList) > 0 {
				items[i].Variations = bodyList.GetByPrimitiveID(item.Slug)
			}
		}
	}

	return nil
}

// Get - Query specified data
func (a *Primitive) Get(ctx context.Context, UUID string, opts ...schema.PrimitiveQueryOptions) (*schema.Primitive, error) {
	var item entity.Primitive
	ok, err := model.FindOne(ctx, entity.GetPrimitiveDB(ctx, a.db).Where("uuid=?", UUID), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	sitem := item.ToSchemaPrimitive()
	err = a.fillSchemaPrimitives(ctx, []*schema.Primitive{sitem}, opts...)
	if err != nil {
		return nil, err
	}

	return sitem, nil
}

// Create - Create data
func (a *Primitive) Create(ctx context.Context, item schema.Primitive) error {
	return model.ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaPrimitive(item)
		result := entity.GetPrimitiveDB(ctx, a.db).Create(sitem.ToPrimitive())

		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		for _, item := range sitem.ToVariations() {
			item.Slug = sitem.Slug
			result := entity.GetPrimitiveBodyDB(ctx, a.db).Create(item)
			if err := result.Error; err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	})
}

// Update - Update data
func (a *Primitive) Update(ctx context.Context, UUID string, item schema.Primitive) error {
	primitive := entity.SchemaPrimitive(item).ToPrimitive()
	result := entity.GetPrimitiveDB(ctx, a.db).Where("uuid=?", UUID).Omit("uuid", "creator").Updates(primitive)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete - delete data
func (a *Primitive) Delete(ctx context.Context, UUID string) error {
	result := entity.GetPrimitiveDB(ctx, a.db).Where("uuid=?", UUID).Delete(entity.Primitive{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *Primitive) queryVariations(ctx context.Context, primitiveIDs ...string) (entity.Variations, error) {
	var list entity.Variations
	result := entity.GetPrimitiveBodyDB(ctx, a.db).Where("slug IN(?)", primitiveIDs).Find(&list)
	if err := result.Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return list, nil
}

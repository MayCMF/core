package implement

import (
	"context"

	"github.com/MayCMF/core/src/common/errors"
	commonschema "github.com/MayCMF/core/src/common/schema"
	"github.com/MayCMF/core/src/common/util"
	"github.com/MayCMF/core/src/primitives/model"
	"github.com/MayCMF/core/src/primitives/schema"
)

// NewPrimitive - Create a Primitive
func NewPrimitive(mPrimitive model.IPrimitive) *Primitive {
	return &Primitive{
		PrimitiveModel: mPrimitive,
	}
}

// Primitive - Sample program
type Primitive struct {
	PrimitiveModel model.IPrimitive
}

// Query - Query data
func (a *Primitive) Query(ctx context.Context, params schema.PrimitiveQueryParam, opts ...schema.PrimitiveQueryOptions) (*schema.PrimitiveQueryResult, error) {
	return a.PrimitiveModel.Query(ctx, params, opts...)
}

// Get - Get specified data
func (a *Primitive) Get(ctx context.Context, UUID string, opts ...schema.PrimitiveQueryOptions) (*schema.Primitive, error) {
	item, err := a.PrimitiveModel.Get(ctx, UUID, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

func (a *Primitive) checkSlug(ctx context.Context, slug string) error {
	result, err := a.PrimitiveModel.Query(ctx, schema.PrimitiveQueryParam{
		Slug: slug,
	}, schema.PrimitiveQueryOptions{
		PageParam: &commonschema.PaginationParam{PageSize: -1},
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("Number already exists")
	}
	return nil
}

func (a *Primitive) getUpdate(ctx context.Context, UUID string) (*schema.Primitive, error) {
	return a.Get(ctx, UUID)
}

// Create - Create Primitive data
func (a *Primitive) Create(ctx context.Context, item schema.Primitive) (*schema.Primitive, error) {
	err := a.checkSlug(ctx, item.Slug)
	if err != nil {
		return nil, err
	}

	item.UUID = util.MustUUID()
	err = a.PrimitiveModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	return a.getUpdate(ctx, item.UUID)
}

// Update - Update Primitive data
func (a *Primitive) Update(ctx context.Context, UUID string, item schema.Primitive) (*schema.Primitive, error) {
	oldItem, err := a.PrimitiveModel.Get(ctx, UUID)
	if err != nil {
		return nil, err
	} else if oldItem == nil {
		return nil, errors.ErrNotFound
	} else if oldItem.Slug != item.Slug {
		err := a.checkSlug(ctx, item.Slug)
		if err != nil {
			return nil, err
		}
	}

	err = a.PrimitiveModel.Update(ctx, UUID, item)
	if err != nil {
		return nil, err
	}
	return a.getUpdate(ctx, UUID)
}

// Delete - Delete data
func (a *Primitive) Delete(ctx context.Context, UUID string) error {
	oldItem, err := a.PrimitiveModel.Get(ctx, UUID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.PrimitiveModel.Delete(ctx, UUID)
}

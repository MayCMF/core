package model

import (
	"context"

	"github.com/MayCMF/core/src/common/errors"
	"github.com/MayCMF/core/src/common/model"
	"github.com/MayCMF/src/filemanager/model/impl/gorm/entity"
	"github.com/MayCMF/src/filemanager/schema"
	"github.com/jinzhu/gorm"
)

// NewFile - Create a File storage instance
func NewFile(db *gorm.DB) *File {
	return &File{db}
}

// File - File storage
type File struct {
	db *gorm.DB
}

func (a *File) getQueryOption(opts ...schema.FileQueryOptions) schema.FileQueryOptions {
	var opt schema.FileQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query - Query data
func (a *File) Query(ctx context.Context, params schema.FileQueryParam, opts ...schema.FileQueryOptions) (*schema.FileQueryResult, error) {
	db := entity.GetFileDB(ctx, a.db)
	if v := params.Filename; v != "" {
		db = db.Where("filename=?", v)
	}
	if v := params.LikeFilename; v != "" {
		db = db.Where("filename LIKE ?", "%"+v+"%")
	}
	if v := params.Uri; v > 0 {
		db = db.Where("uri=?", v)
	}
	db = db.Order("id DESC")

	opt := a.getQueryOption(opts...)
	var list entity.Files
	pr, err := model.WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.FileQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaFiles(),
	}

	return qr, nil
}

// Get - Query specified data
func (a *File) Get(ctx context.Context, UUID string, opts ...schema.FileQueryOptions) (*schema.File, error) {
	db := entity.GetFileDB(ctx, a.db).Where("uuid=?", UUID)
	var item entity.File
	ok, err := model.FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaFile(), nil
}

// Create - Create data
func (a *File) Create(ctx context.Context, item schema.File) error {
	file := entity.SchemaFile(item).ToFile()
	result := entity.GetFileDB(ctx, a.db).Create(file)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update - Update data
func (a *File) Update(ctx context.Context, UUID string, item schema.File) error {
	file := entity.SchemaFile(item).ToFile()
	result := entity.GetFileDB(ctx, a.db).Where("uuid=?", UUID).Omit("uuid", "creator").Updates(file)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete - delete data
func (a *File) Delete(ctx context.Context, UUID string) error {
	result := entity.GetFileDB(ctx, a.db).Where("uuid=?", UUID).Delete(entity.File{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}


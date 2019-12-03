package implement

import (
	"context"

	"github.com/MayCMF/core/src/common/errors"
	commonschema "github.com/MayCMF/core/src/common/schema"
	"github.com/MayCMF/core/src/common/util"
	"github.com/MayCMF/core/src/filemanager/model"
	"github.com/MayCMF/core/src/filemanager/schema"
)

// NewFile - Create a File
func NewFile(mFile model.IFile) *File {
	return &File{
		FileModel: mFile,
	}
}

// File - Sample program
type File struct {
	FileModel model.IFile
}

// Query - Query data
func (a *File) Query(ctx context.Context, params schema.FileQueryParam, opts ...schema.FileQueryOptions) (*schema.FileQueryResult, error) {
	return a.FileModel.Query(ctx, params, opts...)
}

// Get - Get specified data
func (a *File) Get(ctx context.Context, UUID string, opts ...schema.FileQueryOptions) (*schema.File, error) {
	item, err := a.FileModel.Get(ctx, UUID, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

func (a *File) checkFilename(ctx context.Context, filename string) error {
	result, err := a.FileModel.Query(ctx, schema.FileQueryParam{
		Filename: filename,
	}, schema.FileQueryOptions{
		PageParam: &commonschema.PaginationParam{PageSize: -1},
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("Number already exists")
	}
	return nil
}

func (a *File) getUpdate(ctx context.Context, UUID string) (*schema.File, error) {
	return a.Get(ctx, UUID)
}

// Create - Create File data
func (a *File) Create(ctx context.Context, item schema.File) (*schema.File, error) {
	err := a.checkFilename(ctx, item.Filename)
	if err != nil {
		return nil, err
	}

	item.UUID = util.MustUUID()
	err = a.FileModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	return a.getUpdate(ctx, item.UUID)
}

// Update - Update File data
func (a *File) Update(ctx context.Context, UUID string, item schema.File) (*schema.File, error) {
	oldItem, err := a.FileModel.Get(ctx, UUID)
	if err != nil {
		return nil, err
	} else if oldItem == nil {
		return nil, errors.ErrNotFound
	} else if oldItem.Filename != item.Filename {
		err := a.checkFilename(ctx, item.Filename)
		if err != nil {
			return nil, err
		}
	}

	err = a.FileModel.Update(ctx, UUID, item)
	if err != nil {
		return nil, err
	}
	return a.getUpdate(ctx, UUID)
}

// Delete - Delete data
func (a *File) Delete(ctx context.Context, UUID string) error {
	oldItem, err := a.FileModel.Get(ctx, UUID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.FileModel.Delete(ctx, UUID)
}

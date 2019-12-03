package model

import (
	"context"

	icontext "github.com/MayCMF/core/src/common/context"
	"github.com/MayCMF/core/src/common/schema"
	"github.com/MayCMF/core/src/transaction/model"
	"github.com/jinzhu/gorm"
)

// TransFunc - Defining transaction execution functions
type TransFunc func(context.Context) error

// ExecTrans - Execute transaction
func ExecTrans(ctx context.Context, db *gorm.DB, fn TransFunc) error {
	if _, ok := icontext.FromTrans(ctx); ok {
		return fn(ctx)
	}

	transModel := model.NewTrans(db)
	trans, err := transModel.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = transModel.Rollback(ctx, trans)
			panic(r)
		}
	}()

	ctx = icontext.NewTrans(ctx, trans)
	err = fn(ctx)
	if err != nil {
		_ = transModel.Rollback(ctx, trans)
		return err
	}
	return transModel.Commit(ctx, trans)
}

// ExecTransWithLock - Execution transaction (lock)
func ExecTransWithLock(ctx context.Context, db *gorm.DB, fn TransFunc) error {
	if !icontext.FromTransLock(ctx) {
		ctx = icontext.NewTransLock(ctx)
	}
	return ExecTrans(ctx, db, fn)
}

// WrapPageQuery - Packaging with paginated queries
func WrapPageQuery(ctx context.Context, db *gorm.DB, pp *schema.PaginationParam, out interface{}) (*schema.PaginationResult, error) {
	if pp != nil {
		total, err := FindPage(ctx, db, pp.PageIndex, pp.PageSize, out)
		if err != nil {
			return nil, err
		}
		return &schema.PaginationResult{
			Total: total,
		}, nil
	}

	result := db.Find(out)
	return nil, result.Error
}

// FindPage - Query paging data
func FindPage(ctx context.Context, db *gorm.DB, pageIndex, pageSize int, out interface{}) (int, error) {
	var count int
	result := db.Count(&count)
	if err := result.Error; err != nil {
		return 0, err
	} else if count == 0 {
		return 0, nil
	}

	// If the page size is less than 0 or the page index is less than 0, the data is not queried
	if pageSize < 0 || pageIndex < 0 {
		return count, nil
	}

	if pageIndex > 0 && pageSize > 0 {
		db = db.Offset((pageIndex - 1) * pageSize)
	}
	if pageSize > 0 {
		db = db.Limit(pageSize)
	}
	result = db.Find(out)
	if err := result.Error; err != nil {
		return 0, err
	}

	return count, nil
}

// FindOne - Query a single piece of data
func FindOne(ctx context.Context, db *gorm.DB, out interface{}) (bool, error) {
	result := db.First(out)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Check - Check if the data exists
func Check(ctx context.Context, db *gorm.DB) (bool, error) {
	var count int
	result := db.Count(&count)
	if err := result.Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

package model

import (
	"context"

	"github.com/MayCMF/core/src/common/errors"
	"github.com/jinzhu/gorm"
)

// NewTrans - Create a transaction management instance
func NewTrans(db *gorm.DB) *Trans {
	return &Trans{db}
}

// Trans - Manage Transaction
type Trans struct {
	db *gorm.DB
}

// Begin - Open transaction
func (a *Trans) Begin(ctx context.Context) (interface{}, error) {
	result := a.db.Begin()
	if err := result.Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

// Commit - Submit transaction
func (a *Trans) Commit(ctx context.Context, trans interface{}) error {
	db, ok := trans.(*gorm.DB)
	if !ok {
		return errors.New("Unknow transaction")
	}

	result := db.Commit()
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Rollback - Rollback transaction
func (a *Trans) Rollback(ctx context.Context, trans interface{}) error {
	db, ok := trans.(*gorm.DB)
	if !ok {
		return errors.New("Unknow transaction")
	}

	result := db.Rollback()
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

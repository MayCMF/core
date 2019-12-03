package controllers

import (
	"context"

	"github.com/MayCMF/core/src/common"
	"github.com/MayCMF/core/src/transaction/model"
)

// NewTrans - Create a role management instance
func NewTrans(trans model.ITrans) *Trans {
	return &Trans{
		TransModel: trans,
	}
}

// Trans - Manage Transaction
type Trans struct {
	TransModel model.ITrans
}

// Exec - Execute transaction
func (a *Trans) Exec(ctx context.Context, fn func(context.Context) error) error {
	return common.ExecTrans(ctx, a.TransModel, fn)
}

package common

import (
	"context"

	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common/config"
	icontext "github.com/MayCMF/core/src/common/context"
	"github.com/MayCMF/core/src/common/util"
	"github.com/MayCMF/core/src/transaction/model"
)

// GetRootUser - Get root user
func GetRootUser() *schema.User {
	user := config.Global().Root
	return &schema.User{
		UUID:     user.UserName,
		UserName: user.UserName,
		RealName: user.RealName,
		Password: util.MD5HashString(user.Password),
	}
}

// CheckIsRootUser - Check if it is root user
func CheckIsRootUser(ctx context.Context, userUUID string) bool {
	return GetRootUser().UUID == userUUID
}

// TransFunc - Defining transaction execution functions
type TransFunc func(context.Context) error

// ExecTrans - Executive transaction
func ExecTrans(ctx context.Context, transModel model.ITrans, fn TransFunc) error {
	if _, ok := icontext.FromTrans(ctx); ok {
		return fn(ctx)
	}
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

	err = fn(icontext.NewTrans(ctx, trans))
	if err != nil {
		_ = transModel.Rollback(ctx, trans)
		return err
	}
	return transModel.Commit(ctx, trans)
}

// ExecTransWithLock - Execution transaction (lock)
func ExecTransWithLock(ctx context.Context, transModel model.ITrans, fn TransFunc) error {
	if !icontext.FromTransLock(ctx) {
		ctx = icontext.NewTransLock(ctx)
	}
	return ExecTrans(ctx, transModel, fn)
}

package implement

import (
	"context"
	"net/http"

	"github.com/MayCMF/core/src/account/model"
	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common"
	"github.com/MayCMF/core/src/common/auth"
	"github.com/MayCMF/core/src/common/errors"
	"github.com/MayCMF/core/src/common/util"
	"github.com/LyricTian/captcha"
)

// NewLogin - Create a login management instance
func NewLogin(
	a auth.Auther,
	mUser model.IUser,
	mRole model.IRole,
	mPermission model.IPermission,
) *Login {
	return &Login{
		Auth:            a,
		UserModel:       mUser,
		RoleModel:       mRole,
		PermissionModel: mPermission,
	}
}

// Login - Login management
type Login struct {
	UserModel       model.IUser
	RoleModel       model.IRole
	PermissionModel model.IPermission
	Auth            auth.Auther
}

// GetCaptcha - Get graphic verification code information
func (a *Login) GetCaptcha(ctx context.Context, length int) (*schema.LoginCaptcha, error) {
	captchaID := captcha.NewLen(length)
	item := &schema.LoginCaptcha{
		CaptchaID: captchaID,
	}
	return item, nil
}

// ResCaptcha - Generate and respond to a captcha
func (a *Login) ResCaptcha(ctx context.Context, w http.ResponseWriter, captchaID string, width, height int) error {
	err := captcha.WriteImage(w, captchaID, width, height)
	if err != nil {
		if err == captcha.ErrNotFound {
			return errors.ErrNotFound
		}
		return errors.WithStack(err)
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")
	return nil
}

// Verify - Login authentication
func (a *Login) Verify(ctx context.Context, userName, password string) (*schema.User, error) {
	// Check if it is a superuser
	root := common.GetRootUser()
	if userName == root.UserName && root.Password == password {
		return root, nil
	}

	result, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		UserName: userName,
	})
	if err != nil {
		return nil, err
	} else if len(result.Data) == 0 {
		return nil, errors.ErrInvalidUserName
	}

	item := result.Data[0]
	if item.Password != util.SHA1HashString(password) {
		return nil, errors.ErrInvalidPassword
	} else if item.Status != 1 {
		return nil, errors.ErrUserDisable
	}

	return item, nil
}

// GenerateToken - Generate token
func (a *Login) GenerateToken(ctx context.Context, userUUID string) (*schema.LoginTokenInfo, error) {
	tokenInfo, err := a.Auth.GenerateToken(ctx, userUUID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	item := &schema.LoginTokenInfo{
		AccessToken: tokenInfo.GetAccessToken(),
		TokenType:   tokenInfo.GetTokenType(),
		ExpiresAt:   tokenInfo.GetExpiresAt(),
	}
	return item, nil
}

// DestroyToken - Destroy token
func (a *Login) DestroyToken(ctx context.Context, tokenString string) error {
	err := a.Auth.DestroyToken(ctx, tokenString)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *Login) getAndCheckUser(ctx context.Context, userUUID string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	user, err := a.UserModel.Get(ctx, userUUID, opts...)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.ErrInvalidUser
	} else if user.Status != 1 {
		return nil, errors.ErrUserDisable
	}
	return user, nil
}

// GetLoginInfo - Get current user login information
func (a *Login) GetLoginInfo(ctx context.Context, userUUID string) (*schema.UserLoginInfo, error) {
	if isRoot := common.CheckIsRootUser(ctx, userUUID); isRoot {
		root := common.GetRootUser()
		loginInfo := &schema.UserLoginInfo{
			UserName: root.UserName,
			RealName: root.RealName,
		}
		return loginInfo, nil
	}

	user, err := a.getAndCheckUser(ctx, userUUID, schema.UserQueryOptions{
		IncludeRoles: true,
	})
	if err != nil {
		return nil, err
	}

	loginInfo := &schema.UserLoginInfo{
		UserName: user.UserName,
		RealName: user.RealName,
	}

	if roleIDs := user.Roles.ToRoleIDs(); len(roleIDs) > 0 {
		roles, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
			UUIDs: roleIDs,
		})
		if err != nil {
			return nil, err
		}
		loginInfo.RoleNames = roles.Data.ToNames()
	}
	return loginInfo, nil
}

// QueryUserPermissionTree - Get current user's permission permission tree
func (a *Login) QueryUserPermissionTree(ctx context.Context, userUUID string) ([]*schema.PermissionTree, error) {
	isRoot := common.CheckIsRootUser(ctx, userUUID)
	// If it is a root user, query all displayed permission trees
	if isRoot {
		hidden := 0
		result, err := a.PermissionModel.Query(ctx, schema.PermissionQueryParam{
			Hidden: &hidden,
		}, schema.PermissionQueryOptions{
			IncludeActions: true,
		})
		if err != nil {
			return nil, err
		}
		return result.Data.ToTrees().ToTree(), nil
	}

	roleResult, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
		UserUUID: userUUID,
	}, schema.RoleQueryOptions{
		IncludePermissions: true,
	})
	if err != nil {
		return nil, err
	} else if len(roleResult.Data) == 0 {
		return nil, errors.ErrNoPerm
	}

	// Get role permission permission list
	permissionResult, err := a.PermissionModel.Query(ctx, schema.PermissionQueryParam{
		UUIDs: roleResult.Data.ToPermissionIDs(),
	})
	if err != nil {
		return nil, err
	} else if len(permissionResult.Data) == 0 {
		return nil, errors.ErrNoPerm
	}

	// Split and query the permission tree
	permissionResult, err = a.PermissionModel.Query(ctx, schema.PermissionQueryParam{
		UUIDs: permissionResult.Data.SplitAndGetAllUUIDs(),
	}, schema.PermissionQueryOptions{
		IncludeActions: true,
	})
	if err != nil {
		return nil, err
	} else if len(permissionResult.Data) == 0 {
		return nil, errors.ErrNoPerm
	}

	permissionActions := roleResult.Data.ToPermissionIDActionsMap()
	return permissionResult.Data.ToTrees().ForEach(func(item *schema.PermissionTree, _ int) {
		// Traverse permission action permissions
		var actions []*schema.PermissionAction
		for _, code := range permissionActions[item.UUID] {
			for _, aitem := range item.Actions {
				if aitem.Code == code {
					actions = append(actions, aitem)
					break
				}
			}
		}
		item.Actions = actions
	}).ToTree(), nil
}

// UpdatePassword Update current user login password
func (a *Login) UpdatePassword(ctx context.Context, userUUID string, params schema.UpdatePasswordParam) error {
	if common.CheckIsRootUser(ctx, userUUID) {
		return errors.New400Response("Root user not allowed to update the password")
	}

	user, err := a.getAndCheckUser(ctx, userUUID)
	if err != nil {
		return err
	} else if util.SHA1HashString(params.OldPassword) != user.Password {
		return errors.New400Response("Old password is incorrect")
	}

	params.NewPassword = util.SHA1HashString(params.NewPassword)
	return a.UserModel.UpdatePassword(ctx, userUUID, params.NewPassword)
}

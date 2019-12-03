package implement

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/MayCMF/core/src/account/model"
	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common"
	"github.com/MayCMF/core/src/common/config"
	"github.com/MayCMF/core/src/common/errors"
	comschema "github.com/MayCMF/core/src/common/schema"
	"github.com/MayCMF/core/src/common/util"
)

// NewUser - Create a new user
func NewUser(
	e *casbin.SyncedEnforcer,
	mUser model.IUser,
	mRole model.IRole,
) *User {
	return &User{
		Enforcer:  e,
		UserModel: mUser,
		RoleModel: mRole,
		DeleteHook: func(ctx context.Context, bUser *User, UUID string) error {
			if config.Global().Casbin.Enable {
				_, _ = bUser.Enforcer.DeleteUser(UUID)
			}
			return nil
		},
		SaveHook: func(ctx context.Context, bUser *User, item *schema.User) error {
			if config.Global().Casbin.Enable {
				if item.Status == 1 {
					err := bUser.LoadPolicy(ctx, item)
					if err != nil {
						return err
					}
				} else {
					_, _ = bUser.Enforcer.DeleteUser(item.UUID)
				}
			}
			return nil
		},
	}
}

// User - Manage User
type User struct {
	Enforcer   *casbin.SyncedEnforcer
	UserModel  model.IUser
	RoleModel  model.IRole
	DeleteHook func(context.Context, *User, string) error
	SaveHook   func(context.Context, *User, *schema.User) error
}

// Query - Query data
func (a *User) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	return a.UserModel.Query(ctx, params, opts...)
}

// QueryShow - Query display item data
func (a *User) QueryShow(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserShowQueryResult, error) {
	userResult, err := a.UserModel.Query(ctx, params, opts...)
	if err != nil {
		return nil, err
	} else if userResult == nil {
		return nil, nil
	}

	result := &schema.UserShowQueryResult{
		PageResult: userResult.PageResult,
	}
	if len(userResult.Data) == 0 {
		return result, nil
	}

	roleResult, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
		UUIDs: userResult.Data.ToRoleIDs(),
	})
	if err != nil {
		return nil, err
	}
	result.Data = userResult.Data.ToUserShows(roleResult.Data.ToMap())
	return result, nil
}

// Get - Get specified data
func (a *User) Get(ctx context.Context, UUID string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	item, err := a.UserModel.Get(ctx, UUID, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}
	return item, nil
}

func (a *User) checkUserName(ctx context.Context, userName string) error {
	if userName == common.GetRootUser().UserName {
		return errors.New400Response("Username is illegal")
	}

	result, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		UserName: userName,
	}, schema.UserQueryOptions{
		PageParam: &comschema.PaginationParam{PageSize: -1},
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("Uername already exists")
	}
	return nil
}

func (a *User) getUpdate(ctx context.Context, UUID string) (*schema.User, error) {
	nitem, err := a.Get(ctx, UUID, schema.UserQueryOptions{
		IncludeRoles: true,
	})
	if err != nil {
		return nil, err
	}

	if hook := a.SaveHook; hook != nil {
		if err := hook(ctx, a, nitem); err != nil {
			return nil, err
		}
	}

	return nitem, nil
}

// Create - Create user
func (a *User) Create(ctx context.Context, item schema.User) (*schema.User, error) {
	if item.Password == "" {
		return nil, errors.New400Response("Password is not allowed to be empty")
	}

	err := a.checkUserName(ctx, item.UserName)
	if err != nil {
		return nil, err
	}

	item.Password = util.SHA1HashString(item.Password)
	item.UUID = util.MustUUID()
	err = a.UserModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return a.getUpdate(ctx, item.UUID)
}

// Update - Update User
func (a *User) Update(ctx context.Context, UUID string, item schema.User) (*schema.User, error) {
	oldItem, err := a.UserModel.Get(ctx, UUID)
	if err != nil {
		return nil, err
	} else if oldItem == nil {
		return nil, errors.ErrNotFound
	} else if oldItem.UserName != item.UserName {
		err := a.checkUserName(ctx, item.UserName)
		if err != nil {
			return nil, err
		}
	}

	if item.Password != "" {
		item.Password = util.SHA1HashString(item.Password)
	}

	err = a.UserModel.Update(ctx, UUID, item)
	if err != nil {
		return nil, err
	}

	return a.getUpdate(ctx, UUID)
}

// Delete - Delete User
func (a *User) Delete(ctx context.Context, UUID string) error {
	oldItem, err := a.UserModel.Get(ctx, UUID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	err = a.UserModel.Delete(ctx, UUID)
	if err != nil {
		return err
	}

	if hook := a.DeleteHook; hook != nil {
		if err := hook(ctx, a, UUID); err != nil {
			return err
		}
	}

	return nil
}

// UpdateStatus - Update status
func (a *User) UpdateStatus(ctx context.Context, UUID string, status int) error {
	oldItem, err := a.UserModel.Get(ctx, UUID, schema.UserQueryOptions{
		IncludeRoles: true,
	})
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}
	oldItem.Status = status

	err = a.UserModel.UpdateStatus(ctx, UUID, status)
	if err != nil {
		return err
	}

	if hook := a.SaveHook; hook != nil {
		if err := hook(ctx, a, oldItem); err != nil {
			return err
		}
	}

	return nil
}

// LoadPolicy - Load user permission policy
func (a *User) LoadPolicy(ctx context.Context, item *schema.User) error {
	_, _ = a.Enforcer.DeleteRolesForUser(item.UUID)
	for _, roleID := range item.Roles.ToRoleIDs() {
		_, _ = a.Enforcer.AddRoleForUser(item.UUID, roleID)
	}
	return nil
}

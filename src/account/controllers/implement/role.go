package implement

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/MayCMF/core/src/account/model"
	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common/config"
	"github.com/MayCMF/core/src/common/errors"
	comschema "github.com/MayCMF/core/src/common/schema"
	"github.com/MayCMF/core/src/common/util"
)

// NewRole - Create a role management instance
func NewRole(
	e *casbin.SyncedEnforcer,
	mRole model.IRole,
	mPermission model.IPermission,
	mUser model.IUser,
) *Role {
	return &Role{
		Enforcer:        e,
		RoleModel:       mRole,
		PermissionModel: mPermission,
		UserModel:       mUser,
		DeleteHook: func(ctx context.Context, bRole *Role, UUID string) error {
			if config.Global().Casbin.Enable {
				_, _ = bRole.Enforcer.DeletePermissionsForUser(UUID)
			}
			return nil
		},
		SaveHook: func(ctx context.Context, bRole *Role, item *schema.Role) error {
			if config.Global().Casbin.Enable {
				err := bRole.LoadPolicy(ctx, item)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
}

// Role Manage Role
type Role struct {
	Enforcer        *casbin.SyncedEnforcer
	RoleModel       model.IRole
	PermissionModel model.IPermission
	UserModel       model.IUser
	DeleteHook      func(context.Context, *Role, string) error
	SaveHook        func(context.Context, *Role, *schema.Role) error
}

// Query - Query data
func (a *Role) Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error) {
	return a.RoleModel.Query(ctx, params, opts...)
}

// Get - Get specified data
func (a *Role) Get(ctx context.Context, UUID string, opts ...schema.RoleQueryOptions) (*schema.Role, error) {
	item, err := a.RoleModel.Get(ctx, UUID, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

func (a *Role) checkName(ctx context.Context, name string) error {
	result, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
		Name: name,
	}, schema.RoleQueryOptions{
		PageParam: &comschema.PaginationParam{PageSize: -1},
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("The role name already exists")
	}
	return nil
}

func (a *Role) getUpdate(ctx context.Context, UUID string) (*schema.Role, error) {
	nitem, err := a.Get(ctx, UUID, schema.RoleQueryOptions{
		IncludePermissions: true,
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

// Create - Create Role
func (a *Role) Create(ctx context.Context, item schema.Role) (*schema.Role, error) {
	err := a.checkName(ctx, item.Name)
	if err != nil {
		return nil, err
	}

	item.UUID = util.MustUUID()
	err = a.RoleModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return a.getUpdate(ctx, item.UUID)
}

// Update - update role
func (a *Role) Update(ctx context.Context, UUID string, item schema.Role) (*schema.Role, error) {
	oldItem, err := a.RoleModel.Get(ctx, UUID)
	if err != nil {
		return nil, err
	} else if oldItem == nil {
		return nil, errors.ErrNotFound
	} else if oldItem.Name != item.Name {
		err := a.checkName(ctx, item.Name)
		if err != nil {
			return nil, err
		}
	}

	err = a.RoleModel.Update(ctx, UUID, item)
	if err != nil {
		return nil, err
	}

	return a.getUpdate(ctx, UUID)
}

// Delete - delete role
func (a *Role) Delete(ctx context.Context, UUID string) error {
	oldItem, err := a.RoleModel.Get(ctx, UUID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	// If the user has been given the role, it is not allowed to delete
	userResult, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		RoleIDs: []string{UUID},
	}, schema.UserQueryOptions{
		PageParam: &comschema.PaginationParam{PageIndex: -1},
	})
	if err != nil {
		return err
	} else if userResult.PageResult.Total > 0 {
		return errors.New400Response("This role has been assigned to the user and is not allowed to delete")
	}

	err = a.RoleModel.Delete(ctx, UUID)
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

// GetPermissionResources - Get resource permissions
func (a *Role) GetPermissionResources(ctx context.Context, item *schema.Role) (schema.PermissionResources, error) {
	result, err := a.PermissionModel.Query(ctx, schema.PermissionQueryParam{
		UUIDs: item.Permissions.ToPermissionIDs(),
	}, schema.PermissionQueryOptions{
		IncludeResources: true,
	})
	if err != nil {
		return nil, err
	}

	var data schema.PermissionResources
	permissionMap := result.Data.ToMap()
	for _, item := range item.Permissions {
		mitem, ok := permissionMap[item.PermissionID]
		if !ok {
			continue
		}
		resMap := mitem.Resources.ToMap()
		for _, res := range item.Resources {
			ritem, ok := resMap[res]
			if !ok || ritem.Path == "" || ritem.Method == "" {
				continue
			}
			data = append(data, ritem)
		}
	}
	return data, nil
}

// LoadPolicy - Load role permission policy
func (a *Role) LoadPolicy(ctx context.Context, item *schema.Role) error {
	resources, err := a.GetPermissionResources(ctx, item)
	if err != nil {
		return err
	}

	roleID := item.UUID
	_, _ = a.Enforcer.DeletePermissionsForUser(roleID)
	for _, item := range resources {
		_, _ = a.Enforcer.AddPermissionForUser(roleID, item.Path, item.Method)
	}

	return nil
}

package model

import (
	"context"

	"github.com/MayCMF/core/src/account/model/impl/gorm/entity"
	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common/errors"
	"github.com/MayCMF/core/src/common/model"
	"github.com/jinzhu/gorm"
)

// NewRole - Create a role store instance
func NewRole(db *gorm.DB) *Role {
	return &Role{db}
}

// Role - Role storage
type Role struct {
	db *gorm.DB
}

func (a *Role) getQueryOption(opts ...schema.RoleQueryOptions) schema.RoleQueryOptions {
	var opt schema.RoleQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query - Query data
func (a *Role) Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error) {
	db := entity.GetRoleDB(ctx, a.db)
	if v := params.UUIDs; len(v) > 0 {
		db = db.Where("record_id IN(?)", v)
	}
	if v := params.Name; v != "" {
		db = db.Where("name=?", v)
	}
	if v := params.LikeName; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.UserUUID; v != "" {
		subQuery := entity.GetUserRoleDB(ctx, a.db).Where("user_uuid=?", v).Select("role_id").SubQuery()
		db = db.Where("record_id IN(?)", subQuery)
	}
	db = db.Order("sequence DESC,id DESC")

	opt := a.getQueryOption(opts...)
	var list entity.Roles
	pr, err := model.WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.RoleQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaRoles(),
	}

	err = a.fillSchameRoles(ctx, qr.Data, opts...)
	if err != nil {
		return nil, err
	}

	return qr, nil
}

// Populate the role object
func (a *Role) fillSchameRoles(ctx context.Context, items []*schema.Role, opts ...schema.RoleQueryOptions) error {
	opt := a.getQueryOption(opts...)

	if opt.IncludePermissions {

		roleIDs := make([]string, len(items))
		for i, item := range items {
			roleIDs[i] = item.UUID
		}

		var permissionList entity.RolePermissions
		if opt.IncludePermissions {
			items, err := a.queryPermissions(ctx, roleIDs...)
			if err != nil {
				return err
			}
			permissionList = items
		}

		for i, item := range items {
			if len(permissionList) > 0 {
				items[i].Permissions = permissionList.GetByRoleID(item.UUID)
			}
		}
	}
	return nil
}

// Get - Query specified data
func (a *Role) Get(ctx context.Context, UUID string, opts ...schema.RoleQueryOptions) (*schema.Role, error) {
	var role entity.Role
	ok, err := model.FindOne(ctx, entity.GetRoleDB(ctx, a.db).Where("record_id=?", UUID), &role)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	sitem := role.ToSchemaRole()
	err = a.fillSchameRoles(ctx, []*schema.Role{sitem}, opts...)
	if err != nil {
		return nil, err
	}

	return sitem, nil
}

// Create - Create data
func (a *Role) Create(ctx context.Context, item schema.Role) error {
	return model.ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaRole(item)
		result := entity.GetRoleDB(ctx, a.db).Create(sitem.ToRole())
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		for _, item := range sitem.ToRolePermissions() {
			result := entity.GetRolePermissionDB(ctx, a.db).Create(item)
			if err := result.Error; err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	})
}

// Compare and get permission items that need to be added, modified, and deleted
func (a *Role) compareUpdatePermission(oldList, newList entity.RolePermissions) (clist, dlist, ulist entity.RolePermissions) {
	oldMap, newMap := oldList.ToMap(), newList.ToMap()

	for _, nitem := range newList {
		if _, ok := oldMap[nitem.PermissionID]; ok {
			ulist = append(ulist, nitem)
			continue
		}
		clist = append(clist, nitem)
	}

	for _, oitem := range oldList {
		if _, ok := newMap[oitem.PermissionID]; !ok {
			dlist = append(dlist, oitem)
		}
	}
	return
}

// Update permission data
func (a *Role) updatePermissions(ctx context.Context, roleID string, items entity.RolePermissions) error {
	list, err := a.queryPermissions(ctx, roleID)
	if err != nil {
		return err
	}

	clist, dlist, ulist := a.compareUpdatePermission(list, items)
	for _, item := range clist {
		result := entity.GetRolePermissionDB(ctx, a.db).Create(item)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
	}

	for _, item := range dlist {
		result := entity.GetRolePermissionDB(ctx, a.db).Where("role_id=? AND permission_id=?", roleID, item.PermissionID).Delete(entity.RolePermission{})
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
	}

	for _, item := range ulist {
		result := entity.GetRolePermissionDB(ctx, a.db).Where("role_id=? AND permission_id=?", roleID, item.PermissionID).Omit("role_id", "permission_id").Updates(item)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Update - Update data
func (a *Role) Update(ctx context.Context, UUID string, item schema.Role) error {
	return model.ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaRole(item)
		result := entity.GetRoleDB(ctx, a.db).Where("record_id=?", UUID).Omit("record_id", "creator").Updates(sitem.ToRole())
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		err := a.updatePermissions(ctx, UUID, sitem.ToRolePermissions())
		if err != nil {
			return err
		}

		return nil
	})
}

// Delete - delete data
func (a *Role) Delete(ctx context.Context, UUID string) error {
	return model.ExecTrans(ctx, a.db, func(ctx context.Context) error {
		result := entity.GetRoleDB(ctx, a.db).Where("record_id=?", UUID).Delete(entity.Role{})
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		result = entity.GetRolePermissionDB(ctx, a.db).Where("role_id=?", UUID).Delete(entity.RolePermission{})
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
}

func (a *Role) queryPermissions(ctx context.Context, roleIDs ...string) (entity.RolePermissions, error) {
	var list entity.RolePermissions
	result := entity.GetRolePermissionDB(ctx, a.db).Where("role_id IN(?)", roleIDs).Find(&list)
	if err := result.Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return list, nil
}

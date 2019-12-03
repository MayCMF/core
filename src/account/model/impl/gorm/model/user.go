package model

import (
	"context"

	"github.com/MayCMF/core/src/account/model/impl/gorm/entity"
	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common/errors"
	"github.com/MayCMF/core/src/common/model"
	"github.com/jinzhu/gorm"
)

// NewUser - Create a user store instance
func NewUser(db *gorm.DB) *User {
	return &User{db}
}

// User - User storage
type User struct {
	db *gorm.DB
}

func (a *User) getQueryOption(opts ...schema.UserQueryOptions) schema.UserQueryOptions {
	var opt schema.UserQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query - Query data
func (a *User) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	db := entity.GetUserDB(ctx, a.db)
	if v := params.UserName; v != "" {
		db = db.Where("user_name=?", v)
	}
	if v := params.LikeUserName; v != "" {
		db = db.Where("user_name LIKE ?", "%"+v+"%")
	}
	if v := params.LikeRealName; v != "" {
		db = db.Where("real_name LIKE ?", "%"+v+"%")
	}
	if v := params.Status; v > 0 {
		db = db.Where("status=?", v)
	}
	if v := params.RoleIDs; len(v) > 0 {
		subQuery := entity.GetUserRoleDB(ctx, a.db).Select("user_uuid").Where("role_id IN(?)", v).SubQuery()
		db = db.Where("record_id IN(?)", subQuery)
	}
	db = db.Order("id DESC")

	opt := a.getQueryOption(opts...)
	var list entity.Users
	pr, err := model.WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.UserQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaUsers(),
	}

	err = a.fillSchemaUsers(ctx, qr.Data, opts...)
	if err != nil {
		return nil, err
	}

	return qr, nil
}

func (a *User) fillSchemaUsers(ctx context.Context, items []*schema.User, opts ...schema.UserQueryOptions) error {
	opt := a.getQueryOption(opts...)

	if opt.IncludeRoles {
		userUUIDs := make([]string, len(items))
		for i, item := range items {
			userUUIDs[i] = item.UUID
		}

		var roleList entity.UserRoles
		if opt.IncludeRoles {
			items, err := a.queryRoles(ctx, userUUIDs...)
			if err != nil {
				return err
			}
			roleList = items
		}

		for i, item := range items {
			if len(roleList) > 0 {
				items[i].Roles = roleList.GetByUserUUID(item.UUID)
			}
		}
	}

	return nil
}

// Get - Query specified data
func (a *User) Get(ctx context.Context, UUID string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	var item entity.User
	ok, err := model.FindOne(ctx, entity.GetUserDB(ctx, a.db).Where("record_id=?", UUID), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	sitem := item.ToSchemaUser()
	err = a.fillSchemaUsers(ctx, []*schema.User{sitem}, opts...)
	if err != nil {
		return nil, err
	}

	return sitem, nil
}

// Create - Create data
func (a *User) Create(ctx context.Context, item schema.User) error {
	return model.ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaUser(item)
		result := entity.GetUserDB(ctx, a.db).Create(sitem.ToUser())
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		for _, eitem := range sitem.ToUserRoles() {
			result := entity.GetUserRoleDB(ctx, a.db).Create(eitem)
			if err := result.Error; err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})
}

// Compare and get the role data that needs to be added, modified, and deleted.
func (a *User) compareUpdateRole(oldList, newList []*entity.UserRole) (clist, dlist, ulist []*entity.UserRole) {
	for _, nitem := range newList {
		exists := false
		for _, oitem := range oldList {
			if oitem.RoleID == nitem.RoleID {
				exists = true
				ulist = append(ulist, nitem)
				break
			}
		}
		if !exists {
			clist = append(clist, nitem)
		}
	}

	for _, oitem := range oldList {
		exists := false
		for _, nitem := range newList {
			if nitem.RoleID == oitem.RoleID {
				exists = true
				break
			}
		}
		if !exists {
			dlist = append(dlist, oitem)
		}
	}

	return
}

// Update - Update data
func (a *User) Update(ctx context.Context, UUID string, item schema.User) error {
	return model.ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaUser(item)
		omits := []string{"record_id", "creator"}
		if sitem.Password == "" {
			omits = append(omits, "password")
		}

		result := entity.GetUserDB(ctx, a.db).Where("record_id=?", UUID).Omit(omits...).Updates(sitem.ToUser())
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		roles, err := a.queryRoles(ctx, UUID)
		if err != nil {
			return err
		}

		clist, dlist, ulist := a.compareUpdateRole(roles, sitem.ToUserRoles())
		for _, item := range clist {
			result := entity.GetUserRoleDB(ctx, a.db).Create(item)
			if err := result.Error; err != nil {
				return errors.WithStack(err)
			}
		}

		for _, item := range dlist {
			result := entity.GetUserRoleDB(ctx, a.db).Where("user_uuid=? AND role_id=?", UUID, item.RoleID).Delete(entity.UserRole{})
			if err := result.Error; err != nil {
				return errors.WithStack(err)
			}
		}

		for _, item := range ulist {
			result := entity.GetUserRoleDB(ctx, a.db).Where("user_uuid=? AND role_id=?", UUID, item.RoleID).Omit("user_uuid", "role_id").Updates(item)
			if err := result.Error; err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})
}

// Delete - Delete data
func (a *User) Delete(ctx context.Context, UUID string) error {
	return model.ExecTrans(ctx, a.db, func(ctx context.Context) error {
		result := entity.GetUserDB(ctx, a.db).Where("record_id=?", UUID).Delete(entity.User{})
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		result = entity.GetUserRoleDB(ctx, a.db).Where("user_uuid=?", UUID).Delete(entity.UserRole{})
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
}

// UpdateStatus - Update status
func (a *User) UpdateStatus(ctx context.Context, UUID string, status int) error {
	result := entity.GetUserDB(ctx, a.db).Where("record_id=?", UUID).Update("status", status)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdatePassword - Update password
func (a *User) UpdatePassword(ctx context.Context, UUID, password string) error {
	result := entity.GetUserDB(ctx, a.db).Where("record_id=?", UUID).Update("password", password)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *User) queryRoles(ctx context.Context, userUUIDs ...string) (entity.UserRoles, error) {
	var list entity.UserRoles
	result := entity.GetUserRoleDB(ctx, a.db).Where("user_uuid IN(?)", userUUIDs).Find(&list)
	if err := result.Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return list, nil
}

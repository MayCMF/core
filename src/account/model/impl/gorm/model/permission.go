package model

import (
	"context"

	"github.com/MayCMF/core/src/account/model/impl/gorm/entity"
	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common/errors"
	"github.com/MayCMF/core/src/common/model"
	"github.com/jinzhu/gorm"
)

// NewPermission - Create a permission storage instance
func NewPermission(db *gorm.DB) *Permission {
	return &Permission{db}
}

// Permission - Permission storage
type Permission struct {
	db *gorm.DB
}

func (a *Permission) getQueryOption(opts ...schema.PermissionQueryOptions) schema.PermissionQueryOptions {
	var opt schema.PermissionQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query - Query data
func (a *Permission) Query(ctx context.Context, params schema.PermissionQueryParam, opts ...schema.PermissionQueryOptions) (*schema.PermissionQueryResult, error) {
	db := entity.GetPermissionDB(ctx, a.db)
	if v := params.UUIDs; len(v) > 0 {
		db = db.Where("record_id IN(?)", v)
	}
	if v := params.Name; v != "" {
		db = db.Where("name=?", v)
	}
	if v := params.LikeName; v != "" {
		db = db.Where("name LIKE ?", "%"+v+"%")
	}
	if v := params.ParentID; v != nil {
		db = db.Where("parent_id=?", *v)
	}
	if v := params.PrefixParentPath; v != "" {
		db = db.Where("parent_path LIKE ?", v+"%")
	}
	if v := params.Hidden; v != nil {
		db = db.Where("hidden=?", *v)
	}
	db = db.Order("sequence DESC,id DESC")

	opt := a.getQueryOption(opts...)
	var list entity.Permissions
	pr, err := model.WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.PermissionQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaPermissions(),
	}

	err = a.fillSchemaPermissions(ctx, qr.Data, opts...)
	if err != nil {
		return nil, err
	}

	return qr, nil
}

// Populate permission object data
func (a *Permission) fillSchemaPermissions(ctx context.Context, items []*schema.Permission, opts ...schema.PermissionQueryOptions) error {
	opt := a.getQueryOption(opts...)

	if opt.IncludeActions || opt.IncludeResources {

		permissionIDs := make([]string, len(items))
		for i, item := range items {
			permissionIDs[i] = item.UUID
		}

		var actionList entity.PermissionActions
		var resourceList entity.PermissionResources
		if opt.IncludeActions {
			items, err := a.queryActions(ctx, permissionIDs...)
			if err != nil {
				return err
			}
			actionList = items
		}

		if opt.IncludeResources {
			items, err := a.queryResources(ctx, permissionIDs...)
			if err != nil {
				return err
			}
			resourceList = items
		}

		for i, item := range items {
			if len(actionList) > 0 {
				items[i].Actions = actionList.GetByPermissionID(item.UUID)
			}
			if len(resourceList) > 0 {
				items[i].Resources = resourceList.GetByPermissionID(item.UUID)
			}
		}
	}

	return nil
}

// Get - Query specified data
func (a *Permission) Get(ctx context.Context, UUID string, opts ...schema.PermissionQueryOptions) (*schema.Permission, error) {
	var item entity.Permission
	ok, err := model.FindOne(ctx, entity.GetPermissionDB(ctx, a.db).Where("record_id=?", UUID), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	sitem := item.ToSchemaPermission()
	err = a.fillSchemaPermissions(ctx, []*schema.Permission{sitem}, opts...)
	if err != nil {
		return nil, err
	}

	return sitem, nil
}

// Create - Create data
func (a *Permission) Create(ctx context.Context, item schema.Permission) error {
	return model.ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaPermission(item)
		result := entity.GetPermissionDB(ctx, a.db).Create(sitem.ToPermission())
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		for _, item := range sitem.ToPermissionActions() {
			result := entity.GetPermissionActionDB(ctx, a.db).Create(item)
			if err := result.Error; err != nil {
				return errors.WithStack(err)
			}
		}

		for _, item := range sitem.ToPermissionResources() {
			result := entity.GetPermissionResourceDB(ctx, a.db).Create(item)
			if err := result.Error; err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	})
}

// Compare and get action items that need to be added, modified, or deleted
func (a *Permission) compareUpdateAction(oldList, newList entity.PermissionActions) (clist, dlist, ulist []*entity.PermissionAction) {
	oldMap, newMap := oldList.ToMap(), newList.ToMap()

	for _, nitem := range newList {
		if _, ok := oldMap[nitem.Code]; ok {
			ulist = append(ulist, nitem)
			continue
		}
		clist = append(clist, nitem)
	}

	for _, oitem := range oldList {
		if _, ok := newMap[oitem.Code]; !ok {
			dlist = append(dlist, oitem)
		}
	}
	return
}

// Update action data
func (a *Permission) updateActions(ctx context.Context, permissionID string, items entity.PermissionActions) error {
	list, err := a.queryActions(ctx, permissionID)
	if err != nil {
		return err
	}

	clist, dlist, ulist := a.compareUpdateAction(list, items)
	for _, item := range clist {
		result := entity.GetPermissionActionDB(ctx, a.db).Create(item)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
	}

	for _, item := range dlist {
		result := entity.GetPermissionActionDB(ctx, a.db).Where("permission_id=? AND code=?", permissionID, item.Code).Delete(entity.PermissionAction{})
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
	}

	for _, item := range ulist {
		result := entity.GetPermissionActionDB(ctx, a.db).Where("permission_id=? AND code=?", permissionID, item.Code).Omit("permission_id", "code").Updates(item)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Compare and get resource items that need to be added, modified, and deleted
func (a *Permission) compareUpdateResource(oldList, newList entity.PermissionResources) (clist, dlist, ulist []*entity.PermissionResource) {
	oldMap, newMap := oldList.ToMap(), newList.ToMap()

	for _, nitem := range newList {
		if _, ok := oldMap[nitem.Code]; ok {
			ulist = append(ulist, nitem)
			continue
		}
		clist = append(clist, nitem)
	}

	for _, oitem := range oldList {
		if _, ok := newMap[oitem.Code]; !ok {
			dlist = append(dlist, oitem)
		}
	}
	return
}

// Update resource data
func (a *Permission) updateResources(ctx context.Context, permissionID string, items entity.PermissionResources) error {
	list, err := a.queryResources(ctx, permissionID)
	if err != nil {
		return err
	}

	clist, dlist, ulist := a.compareUpdateResource(list, items)
	for _, item := range clist {
		result := entity.GetPermissionResourceDB(ctx, a.db).Create(item)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
	}

	for _, item := range dlist {
		result := entity.GetPermissionResourceDB(ctx, a.db).Where("permission_id=? AND code=?", permissionID, item.Code).Delete(entity.PermissionResource{})
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
	}

	for _, item := range ulist {
		result := entity.GetPermissionResourceDB(ctx, a.db).Where("permission_id=? AND code=?", permissionID, item.Code).Omit("permission_id", "code").Updates(item)
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Update - update data
func (a *Permission) Update(ctx context.Context, UUID string, item schema.Permission) error {
	return model.ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaPermission(item)
		result := entity.GetPermissionDB(ctx, a.db).Where("record_id=?", UUID).Omit("record_id", "creator").Updates(sitem.ToPermission())
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		err := a.updateActions(ctx, UUID, sitem.ToPermissionActions())
		if err != nil {
			return err
		}

		err = a.updateResources(ctx, UUID, sitem.ToPermissionResources())
		if err != nil {
			return err
		}

		return nil
	})
}

// UpdateParentPath - Update parent path
func (a *Permission) UpdateParentPath(ctx context.Context, UUID, parentPath string) error {
	result := entity.GetPermissionDB(ctx, a.db).Where("record_id=?", UUID).Update("parent_path", parentPath)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete - Delete data
func (a *Permission) Delete(ctx context.Context, UUID string) error {
	return model.ExecTrans(ctx, a.db, func(ctx context.Context) error {
		result := entity.GetPermissionDB(ctx, a.db).Where("record_id=?", UUID).Delete(entity.Permission{})
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		result = entity.GetPermissionActionDB(ctx, a.db).Where("permission_id=?", UUID).Delete(entity.PermissionAction{})
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		result = entity.GetPermissionResourceDB(ctx, a.db).Where("permission_id=?", UUID).Delete(entity.PermissionResource{})
		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
}

func (a *Permission) queryActions(ctx context.Context, permissionIDs ...string) (entity.PermissionActions, error) {
	var list entity.PermissionActions
	result := entity.GetPermissionActionDB(ctx, a.db).Where("permission_id IN(?)", permissionIDs).Find(&list)
	if err := result.Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return list, nil
}

func (a *Permission) queryResources(ctx context.Context, permissionIDs ...string) (entity.PermissionResources, error) {
	var list entity.PermissionResources
	result := entity.GetPermissionResourceDB(ctx, a.db).Where("permission_id IN(?)", permissionIDs).Find(&list)
	if err := result.Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return list, nil
}

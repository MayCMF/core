package implement

import (
	"context"

	"github.com/MayCMF/core/src/account/model"
	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common"
	"github.com/MayCMF/core/src/common/errors"
	comschema "github.com/MayCMF/core/src/common/schema"
	"github.com/MayCMF/core/src/common/util"
	transaction "github.com/MayCMF/core/src/transaction/model"
)

// NewPermission - Create a Permission management instance
func NewPermission(
	trans transaction.ITrans,
	mPermission model.IPermission,
) *Permission {
	return &Permission{
		TransModel:      trans,
		PermissionModel: mPermission,
	}
}

// Permission - Manage Permission
type Permission struct {
	TransModel      transaction.ITrans
	PermissionModel model.IPermission
}

// Query - Get Data
func (a *Permission) Query(ctx context.Context, params schema.PermissionQueryParam, opts ...schema.PermissionQueryOptions) (*schema.PermissionQueryResult, error) {
	return a.PermissionModel.Query(ctx, params, opts...)
}

// Get - get specified data
func (a *Permission) Get(ctx context.Context, UUID string, opts ...schema.PermissionQueryOptions) (*schema.Permission, error) {
	item, err := a.PermissionModel.Get(ctx, UUID, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}
	return item, nil
}

func (a *Permission) getSep() string {
	return "/"
}

func (a *Permission) joinParentPath(ppath, code string) string {
	if ppath != "" {
		ppath += a.getSep()
	}
	return ppath + code
}

// Get the parent path
func (a *Permission) getParentPath(ctx context.Context, parentID string) (string, error) {
	if parentID == "" {
		return "", nil
	}

	pitem, err := a.PermissionModel.Get(ctx, parentID)
	if err != nil {
		return "", err
	} else if pitem == nil {
		return "", errors.ErrInvalidParent
	}

	return a.joinParentPath(pitem.ParentPath, pitem.UUID), nil
}

func (a *Permission) getUpdate(ctx context.Context, UUID string) (*schema.Permission, error) {
	return a.Get(ctx, UUID, schema.PermissionQueryOptions{
		IncludeActions:   true,
		IncludeResources: true,
	})
}

func (a *Permission) checkName(ctx context.Context, item schema.Permission) error {
	result, err := a.PermissionModel.Query(ctx, schema.PermissionQueryParam{
		ParentID: &item.ParentID,
		Name:     item.Name,
	}, schema.PermissionQueryOptions{
		PageParam: &comschema.PaginationParam{PageSize: -1},
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("Permission name is already exists")
	}
	return nil
}

// Create - Create Permission
func (a *Permission) Create(ctx context.Context, item schema.Permission) (*schema.Permission, error) {
	if err := a.checkName(ctx, item); err != nil {
		return nil, err
	}

	parentPath, err := a.getParentPath(ctx, item.ParentID)
	if err != nil {
		return nil, err
	}

	item.ParentPath = parentPath
	item.UUID = util.MustUUID()
	err = a.PermissionModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return a.getUpdate(ctx, item.UUID)
}

// Update - update Permission
func (a *Permission) Update(ctx context.Context, UUID string, item schema.Permission) (*schema.Permission, error) {
	if UUID == item.ParentID {
		return nil, errors.ErrInvalidParent
	}

	oldItem, err := a.PermissionModel.Get(ctx, UUID)
	if err != nil {
		return nil, err
	} else if oldItem == nil {
		return nil, errors.ErrNotFound
	} else if oldItem.Name != item.Name {
		if err := a.checkName(ctx, item); err != nil {
			return nil, err
		}
	}
	item.ParentPath = oldItem.ParentPath

	err = common.ExecTrans(ctx, a.TransModel, func(ctx context.Context) error {
		// If the parent is updated, you need to update the current node and the parent path below the node.
		if item.ParentID != oldItem.ParentID {
			parentPath, err := a.getParentPath(ctx, item.ParentID)
			if err != nil {
				return err
			}
			item.ParentPath = parentPath

			opath := a.joinParentPath(oldItem.ParentPath, oldItem.UUID)
			result, err := a.PermissionModel.Query(ctx, schema.PermissionQueryParam{
				PrefixParentPath: opath,
			})
			if err != nil {
				return err
			}

			npath := a.joinParentPath(item.ParentPath, item.UUID)
			for _, permission := range result.Data {
				npath2 := npath + permission.ParentPath[len(opath):]
				err = a.PermissionModel.UpdateParentPath(ctx, permission.UUID, npath2)
				if err != nil {
					return err
				}
			}
		}

		return a.PermissionModel.Update(ctx, UUID, item)
	})
	if err != nil {
		return nil, err
	}
	return a.getUpdate(ctx, UUID)
}

// Delete - Delete permission
func (a *Permission) Delete(ctx context.Context, UUID string) error {
	oldItem, err := a.PermissionModel.Get(ctx, UUID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	result, err := a.PermissionModel.Query(ctx, schema.PermissionQueryParam{
		ParentID: &UUID,
	}, schema.PermissionQueryOptions{PageParam: &comschema.PaginationParam{PageSize: -1}})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.ErrNotAllowDeleteWithChild
	}

	return a.PermissionModel.Delete(ctx, UUID)
}

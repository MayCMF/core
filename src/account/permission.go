package account

import (
	"context"
	"os"

	"github.com/MayCMF/core/src/account/controllers"
	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common/config"
	comschema "github.com/MayCMF/core/src/common/schema"
	"github.com/MayCMF/core/src/common/util"
	transaction "github.com/MayCMF/core/src/transaction/controllers"
	"go.uber.org/dig"
)

// InitPermission - Initialize Permission data
func InitPermission(ctx context.Context, container *dig.Container) error {
	if c := config.Global().Permission; c.Enable && c.Data != "" {
		return initPermissionData(ctx, container)
	}

	return nil
}

// initPermissionData - Initialize permission data
func initPermissionData(ctx context.Context, container *dig.Container) error {
	return container.Invoke(func(trans transaction.ITrans, permission controllers.IPermission) error {
		// Check if there is permission data, initialize if it does not exist
		permissionResult, err := permission.Query(ctx, schema.PermissionQueryParam{}, schema.PermissionQueryOptions{
			PageParam: &comschema.PaginationParam{PageIndex: -1},
		})
		if err != nil {
			return err
		} else if permissionResult.PageResult.Total > 0 {
			return nil
		}

		data, err := readPermissionData()
		if err != nil {
			return err
		}

		return createPermissions(ctx, trans, permission, "", data)
	})
}

func readPermissionData() (schema.PermissionTrees, error) {
	file, err := os.Open(config.Global().Permission.Data)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data schema.PermissionTrees
	err = util.JSONNewDecoder(file).Decode(&data)
	return data, err
}

func createPermissions(ctx context.Context, trans transaction.ITrans, permission controllers.IPermission, parentID string, list schema.PermissionTrees) error {
	return trans.Exec(ctx, func(ctx context.Context) error {
		for _, item := range list {
			sitem := schema.Permission{
				Name:      item.Name,
				Sequence:  item.Sequence,
				Icon:      item.Icon,
				Router:    item.Router,
				Hidden:    item.Hidden,
				ParentID:  parentID,
				Actions:   item.Actions,
				Resources: item.Resources,
			}
			nsitem, err := permission.Create(ctx, sitem)
			if err != nil {
				return err
			}

			if item.Children != nil && len(*item.Children) > 0 {
				err := createPermissions(ctx, trans, permission, nsitem.UUID, *item.Children)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

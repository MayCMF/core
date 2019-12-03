package account

import (
	"context"
	"fmt"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/MayCMF/core/src/account/controllers"
	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common/config"
	"github.com/MayCMF/core/src/common/logger"
	"go.uber.org/dig"
)

// NewCasbinEnforcer - Create casbin validator
func NewCasbinEnforcer() *casbin.SyncedEnforcer {
	cfg := config.Global().Casbin
	if !cfg.Enable {
		return nil
	}

	e, err := casbin.NewSyncedEnforcer(cfg.Model)
	handleError(err)

	e.EnableAutoSave(false)
	e.EnableAutoBuildRoleLinks(true)

	if cfg.Debug {
		e.EnableLog(true)
	}
	return e
}

// InitCasbinEnforcer - Initialize the casbin checker
func InitCasbinEnforcer(container *dig.Container) error {
	cfg := config.Global().Casbin
	if !cfg.Enable {
		return nil
	}

	return container.Invoke(func(e *casbin.SyncedEnforcer, bRole controllers.IRole, bUser controllers.IUser) error {
		adapter := NewCasbinAdapter(bRole, bUser)

		if cfg.AutoLoad {
			e.InitWithModelAndAdapter(e.GetModel(), adapter)
			e.StartAutoLoadPolicy(time.Duration(cfg.AutoLoadInternal) * time.Second)
		} else {
			err := adapter.LoadPolicy(e.GetModel())
			if err != nil {
				return err
			}
		}

		err := e.BuildRoleLinks()
		if err != nil {
			return err
		}

		return nil
	})
}

// ReleaseCasbinEnforcer - Release casbin resources
func ReleaseCasbinEnforcer(container *dig.Container) {
	cfg := config.Global().Casbin
	if !cfg.Enable || !cfg.AutoLoad {
		return
	}

	container.Invoke(func(e *casbin.SyncedEnforcer) {
		e.StopAutoLoadPolicy()
	})
}

// NewCasbinAdapter - Create a casbin adapter
func NewCasbinAdapter(bRole controllers.IRole, bUser controllers.IUser) *CasbinAdapter {
	return &CasbinAdapter{
		RoleBll: bRole,
		UserBll: bUser,
	}
}

// CasbinAdapter - Casbin adapter
type CasbinAdapter struct {
	RoleBll controllers.IRole
	UserBll controllers.IUser
}

// LoadPolicy - Load all policy rules from the storage.
func (a *CasbinAdapter) LoadPolicy(model model.Model) error {
	ctx := context.Background()
	err := a.loadRolePolicy(ctx, model)
	if err != nil {
		logger.Errorf(ctx, "Load casbin role policy error: %s", err.Error())
		return err
	}

	err = a.loadUserPolicy(ctx, model)
	if err != nil {
		logger.Errorf(ctx, "Load casbin user policy error: %s", err.Error())
		return err
	}
	return nil
}

func (a *CasbinAdapter) loadRolePolicy(ctx context.Context, model model.Model) error {
	// Load role strategy
	roleResult, err := a.RoleBll.Query(ctx, schema.RoleQueryParam{}, schema.RoleQueryOptions{
		IncludePermissions: true,
	})
	if err != nil {
		return err
	}

	for _, item := range roleResult.Data {
		resources, err := a.RoleBll.GetPermissionResources(ctx, item)
		if err != nil {
			return err
		}

		for _, ritem := range resources {
			if ritem.Path == "" || ritem.Method == "" {
				continue
			}

			line := fmt.Sprintf("p,%s,%s,%s", item.UUID, ritem.Path, ritem.Method)
			persist.LoadPolicyLine(line, model)
		}
	}

	return nil
}

func (a *CasbinAdapter) loadUserPolicy(ctx context.Context, model model.Model) error {
	result, err := a.UserBll.Query(ctx, schema.UserQueryParam{
		Status: 1,
	}, schema.UserQueryOptions{IncludeRoles: true})
	if err != nil {
		return err
	}

	for _, item := range result.Data {
		for _, roleID := range item.Roles.ToRoleIDs() {
			line := fmt.Sprintf("g,%s,%s", item.UUID, roleID)
			persist.LoadPolicyLine(line, model)
		}
	}

	return nil
}

// SavePolicy saves all policy rules to the storage.
func (a *CasbinAdapter) SavePolicy(model model.Model) error {
	return nil
}

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}

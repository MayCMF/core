package entity

import (
	"context"

	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common/entity"
	"github.com/jinzhu/gorm"
)

// GetPermissionDB - Get permission storage
func GetPermissionDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(ctx, defDB, Permission{})
}

// GetPermissionActionDB - Get permission action storage
func GetPermissionActionDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(ctx, defDB, PermissionAction{})
}

// GetPermissionResourceDB - Get permission resource storage
func GetPermissionResourceDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(ctx, defDB, PermissionResource{})
}

// SchemaPermission - Permission object
type SchemaPermission schema.Permission

// ToPermission - Convert to permission entity
func (a SchemaPermission) ToPermission() *Permission {
	item := &Permission{
		UUID:       a.UUID,
		Name:       &a.Name,
		Sequence:   &a.Sequence,
		Icon:       &a.Icon,
		Router:     &a.Router,
		Hidden:     &a.Hidden,
		ParentID:   &a.ParentID,
		ParentPath: &a.ParentPath,
		Creator:    &a.Creator,
	}
	return item
}

// ToPermissionActions - Convert to permission action list
func (a SchemaPermission) ToPermissionActions() []*PermissionAction {
	list := make([]*PermissionAction, len(a.Actions))
	for i, item := range a.Actions {
		list[i] = SchemaPermissionAction(*item).ToPermissionAction(a.UUID)
	}
	return list
}

// ToPermissionResources - Convert to permission resource list
func (a SchemaPermission) ToPermissionResources() []*PermissionResource {
	list := make([]*PermissionResource, len(a.Resources))
	for i, item := range a.Resources {
		list[i] = SchemaPermissionResource(*item).ToPermissionResource(a.UUID)
	}
	return list
}

// Permission - Permission entity
type Permission struct {
	entity.Model
	UUID       string  `gorm:"column:record_id;size:36;index;"`    // Record internal code
	Name       *string `gorm:"column:name;size:50;index;"`         // Permission name
	Sequence   *int    `gorm:"column:sequence;index;"`             // Sort value
	Icon       *string `gorm:"column:icon;size:255;"`              // Permission icon
	Router     *string `gorm:"column:router;size:255;"`            // Access routing
	Hidden     *int    `gorm:"column:hidden;index;"`               // Hide permission (0: don't hide 1: hide)
	ParentID   *string `gorm:"column:parent_id;size:36;index;"`    // Parent inner code
	ParentPath *string `gorm:"column:parent_path;size:518;index;"` // Parent path
	Creator    *string `gorm:"column:creator;size:36;"`            // Creator
}

func (a Permission) String() string {
	return entity.ToString(a)
}

// TableName - Table Name
func (a Permission) TableName() string {
	return a.Model.TableName("permission")
}

// ToSchemaPermission - Convert to permission object
func (a Permission) ToSchemaPermission() *schema.Permission {
	item := &schema.Permission{
		UUID:       a.UUID,
		Name:       *a.Name,
		Sequence:   *a.Sequence,
		Icon:       *a.Icon,
		Router:     *a.Router,
		ParentID:   *a.ParentID,
		ParentPath: *a.ParentPath,
		Creator:    *a.Creator,
	}
	if a.Hidden != nil {
		item.Hidden = *a.Hidden
	}
	return item
}

// Permissions - Permission entity list
type Permissions []*Permission

// ToSchemaPermissions - Convert to permission object list
func (a Permissions) ToSchemaPermissions() []*schema.Permission {
	list := make([]*schema.Permission, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaPermission()
	}
	return list
}

// SchemaPermissionAction Permission action object
type SchemaPermissionAction schema.PermissionAction

// ToPermissionAction - Convert to permission action entity
func (a SchemaPermissionAction) ToPermissionAction(permissionID string) *PermissionAction {
	return &PermissionAction{
		PermissionID: permissionID,
		Code:         a.Code,
		Name:         a.Name,
	}
}

// PermissionAction - Permission action associated entity
type PermissionAction struct {
	entity.Model
	PermissionID string `gorm:"column:permission_id;size:36;index;"` // Permission ID
	Code         string `gorm:"column:code;size:50;index;"`          // Action number
	Name         string `gorm:"column:name;size:50;"`                // Action name
}

// TableName - Table Name
func (a PermissionAction) TableName() string {
	return a.Model.TableName("permission_action")
}

// ToSchemaPermissionAction - Convert to permission action object
func (a PermissionAction) ToSchemaPermissionAction() *schema.PermissionAction {
	return &schema.PermissionAction{
		Code: a.Code,
		Name: a.Name,
	}
}

// PermissionActions - Permission action associated entity list
type PermissionActions []*PermissionAction

// GetByPermissionID - Get permission action list based on permission ID
func (a PermissionActions) GetByPermissionID(permissionID string) []*schema.PermissionAction {
	var list []*schema.PermissionAction
	for _, item := range a {
		if item.PermissionID == permissionID {
			list = append(list, item.ToSchemaPermissionAction())
		}
	}
	return list
}

// ToSchemaPermissionActions - Convert to permission action list
func (a PermissionActions) ToSchemaPermissionActions() []*schema.PermissionAction {
	list := make([]*schema.PermissionAction, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaPermissionAction()
	}
	return list
}

// ToMap - Convert to key-value mapping
func (a PermissionActions) ToMap() map[string]*PermissionAction {
	m := make(map[string]*PermissionAction)
	for _, item := range a {
		m[item.Code] = item
	}
	return m
}

// SchemaPermissionResource - Permission resource object
type SchemaPermissionResource schema.PermissionResource

// ToPermissionResource - Convert to permission resource entity
func (a SchemaPermissionResource) ToPermissionResource(permissionID string) *PermissionResource {
	return &PermissionResource{
		PermissionID: permissionID,
		Code:         a.Code,
		Name:         a.Name,
		Method:       a.Method,
		Path:         a.Path,
	}
}

// PermissionResource - Permission resource associated entity
type PermissionResource struct {
	entity.Model
	PermissionID string `gorm:"column:permission_id;size:36;index;"` // Permission ID
	Code         string `gorm:"column:code;size:50;index;"`          // Resource number
	Name         string `gorm:"column:name;size:50;"`                // Resource Name
	Method       string `gorm:"column:method;size:50;"`              // Request method
	Path         string `gorm:"column:path;size:255;"`               // Request path
}

// TableName - Table Name
func (a PermissionResource) TableName() string {
	return a.Model.TableName("permission_resource")
}

// ToSchemaPermissionResource - Convert to permission resource object
func (a PermissionResource) ToSchemaPermissionResource() *schema.PermissionResource {
	return &schema.PermissionResource{
		Code:   a.Code,
		Name:   a.Name,
		Method: a.Method,
		Path:   a.Path,
	}
}

// PermissionResources - Permission resource associated entity list
type PermissionResources []*PermissionResource

// GetByPermissionID - Get permission resource list according to permission ID
func (a PermissionResources) GetByPermissionID(permissionID string) []*schema.PermissionResource {
	var list []*schema.PermissionResource
	for _, item := range a {
		if item.PermissionID == permissionID {
			list = append(list, item.ToSchemaPermissionResource())
		}
	}
	return list
}

// ToSchemaPermissionResources - Convert to permission resource list
func (a PermissionResources) ToSchemaPermissionResources() []*schema.PermissionResource {
	list := make([]*schema.PermissionResource, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaPermissionResource()
	}
	return list
}

// ToMap - Convert to key-value mapping
func (a PermissionResources) ToMap() map[string]*PermissionResource {
	m := make(map[string]*PermissionResource)
	for _, item := range a {
		m[item.Code] = item
	}
	return m
}

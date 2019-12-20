package entity

import (
	"context"
	"strings"

	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common/entity"
	"github.com/jinzhu/gorm"
)

// GetRoleDB - Get the role store
func GetRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(ctx, defDB, Role{})
}

// GetRolePermissionDB - Get the role permission associative storage
func GetRolePermissionDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(ctx, defDB, RolePermission{})
}

// SchemaRole - Role object
type SchemaRole schema.Role

// ToRole - Convert to a role entity
func (a SchemaRole) ToRole() *Role {
	item := &Role{
		UUID:     a.UUID,
		Name:     &a.Name,
		Sequence: &a.Sequence,
		Memo:     &a.Memo,
		Creator:  &a.Creator,
	}
	return item
}

// ToRolePermissions - Convert to role permission entity list
func (a SchemaRole) ToRolePermissions() []*RolePermission {
	list := make([]*RolePermission, len(a.Permissions))
	for i, item := range a.Permissions {
		list[i] = SchemaRolePermission(*item).ToRolePermission(a.UUID)
	}
	return list
}

// Role - Role entity
type Role struct {
	entity.Model
	UUID     string  `gorm:"column:record_id;size:36;index;"` // Record internal code
	Name     *string `gorm:"column:name;size:100;index;"`     // Role Name
	Sequence *int    `gorm:"column:sequence;index;"`          // Sort value
	Memo     *string `gorm:"column:memo;size:200;"`           // Remarks
	Creator  *string `gorm:"column:creator;size:36;"`         // Creator
}

func (a Role) String() string {
	return entity.ToString(a)
}

// TableName - Table Name
func (a Role) TableName() string {
	return a.Model.TableName("role")
}

// ToSchemaRole - Convert to a role object
func (a Role) ToSchemaRole() *schema.Role {
	item := &schema.Role{
		UUID:      a.UUID,
		Name:      *a.Name,
		Sequence:  *a.Sequence,
		Memo:      *a.Memo,
		Creator:   *a.Creator,
		CreatedAt: a.CreatedAt,
	}
	return item
}

// Roles - List of role entities
type Roles []*Role

// ToSchemaRoles - Convert to a list of role objects
func (a Roles) ToSchemaRoles() []*schema.Role {
	list := make([]*schema.Role, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRole()
	}
	return list
}

// SchemaRolePermission - Role permission object
type SchemaRolePermission schema.RolePermission

// ToRolePermission - Convert to role permission entity
func (a SchemaRolePermission) ToRolePermission(roleID string) *RolePermission {
	item := &RolePermission{
		RoleID:       roleID,
		PermissionID: a.PermissionID,
	}

	var action string
	if v := a.Actions; len(v) > 0 {
		action = strings.Join(v, ",")
	}
	item.Action = &action

	var resource string
	if v := a.Resources; len(v) > 0 {
		resource = strings.Join(v, ",")
	}
	item.Resource = &resource

	return item
}

// RolePermission - Role permission associated entity
type RolePermission struct {
	entity.Model
	RoleID       string  `gorm:"column:role_id;size:36;index;"`       // Role inner code
	PermissionID string  `gorm:"column:permission_id;size:36;index;"` // Permission internal code
	Action       *string `gorm:"column:action;size:2048;"`            // Action permissions (multiple separated by commas)
	Resource     *string `gorm:"column:resource;size:2048;"`          // Resource permissions (multiple separated by commas)
}

// TableName - Table Name
func (a RolePermission) TableName() string {
	return a.Model.TableName("role_permission")
}

// ToSchemaRolePermission - Convert to a role permission object
func (a RolePermission) ToSchemaRolePermission() *schema.RolePermission {
	item := &schema.RolePermission{
		PermissionID: a.PermissionID,
	}

	if v := a.Action; v != nil && *v != "" {
		item.Actions = strings.Split(*v, ",")
	}
	if v := a.Resource; v != nil && *v != "" {
		item.Resources = strings.Split(*v, ",")
	}

	return item
}

// RolePermissions - Role permission associated entity list
type RolePermissions []*RolePermission

// GetByRoleID - Get a list of role permission objects based on the role ID
func (a RolePermissions) GetByRoleID(roleID string) []*schema.RolePermission {
	var list []*schema.RolePermission
	for _, item := range a {
		if item.RoleID == roleID {
			list = append(list, item.ToSchemaRolePermission())
		}
	}
	return list
}

// ToSchemaRolePermissions - Convert to a list of role permission objects
func (a RolePermissions) ToSchemaRolePermissions() []*schema.RolePermission {
	list := make([]*schema.RolePermission, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRolePermission()
	}
	return list
}

// ToMap - Convert to key-value mapping
func (a RolePermissions) ToMap() map[string]*RolePermission {
	m := make(map[string]*RolePermission)
	for _, item := range a {
		m[item.PermissionID] = item
	}
	return m
}

package schema

import (
	"github.com/MayCMF/core/src/common/schema"
	"time"
)

// Role - Role object
type Role struct {
	UUID    string          `json:"record_id"`                           // Record ID
	Name        string          `json:"name" binding:"required"`             // Role Name
	Sequence    int             `json:"sequence"`                            // Sort value
	Memo        string          `json:"memo"`                                // Remarks
	Creator     string          `json:"creator"`                             // Creator
	CreatedAt   time.Time       `json:"created_at"`                          // Creation time
	Permissions RolePermissions `json:"permissions" binding:"required,gt=0"` // Permission permission
}

// RolePermission - Role permission object
type RolePermission struct {
	PermissionID string   `json:"permission_id"` // Permission ID
	Actions      []string `json:"actions"`       // Action permission list
	Resources    []string `json:"resources"`     // Resource permission list
}

// RoleQueryParam - Query conditions
type RoleQueryParam struct {
	UUIDs []string // Record ID list
	Name      string   // Role Name
	LikeName  string   // Role name (fuzzy query)
	UserUUID  string   // User UUID
}

// RoleQueryOptions Query optional parameter items
type RoleQueryOptions struct {
	PageParam          *schema.PaginationParam // Paging parameter
	IncludePermissions bool                    // Contains permission permissions
}

// RoleQueryResult - Search result
type RoleQueryResult struct {
	Data       Roles
	PageResult *schema.PaginationResult
}

// Roles - Role object list
type Roles []*Role

// ToPermissionIDs - Get all the permission IDs (do not go heavy)
func (a Roles) ToPermissionIDs() []string {
	var idList []string
	for _, item := range a {
		idList = append(idList, item.Permissions.ToPermissionIDs()...)
	}
	return idList
}

func (a Roles) mergeStrings(olds, news []string) []string {
	for _, n := range news {
		exists := false
		for _, o := range olds {
			if o == n {
				exists = true
				break
			}
		}
		if !exists {
			olds = append(olds, n)
		}
	}
	return olds
}

// ToPermissionIDActionsMap - Action permission list mapping converted to permission ID
func (a Roles) ToPermissionIDActionsMap() map[string][]string {
	m := make(map[string][]string)
	for _, item := range a {
		for _, permission := range item.Permissions {
			v, ok := m[permission.PermissionID]
			if ok {
				m[permission.PermissionID] = a.mergeStrings(v, permission.Actions)
				continue
			}
			m[permission.PermissionID] = permission.Actions
		}
	}
	return m
}

// ToNames - Get a list of role names
func (a Roles) ToNames() []string {
	names := make([]string, len(a))
	for i, item := range a {
		names[i] = item.Name
	}
	return names
}

// ToMap - Convert to key-value store
func (a Roles) ToMap() map[string]*Role {
	m := make(map[string]*Role)
	for _, item := range a {
		m[item.UUID] = item
	}
	return m
}

// RolePermissions - Role permission list
type RolePermissions []*RolePermission

// ToPermissionIDs - Convert to permission ID list
func (a RolePermissions) ToPermissionIDs() []string {
	list := make([]string, len(a))
	for i, item := range a {
		list[i] = item.PermissionID
	}
	return list
}

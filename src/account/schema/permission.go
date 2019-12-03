package schema

import (
	"github.com/MayCMF/core/src/common/schema"
	"strings"
	"time"
)

// Permission - Permission object
type Permission struct {
	UUID       string              `json:"record_id"`               // Record ID
	Name       string              `json:"name" binding:"required"` // Permission name
	Sequence   int                 `json:"sequence"`                // Sort value
	Icon       string              `json:"icon"`                    // Permission icon
	Router     string              `json:"router"`                  // Access routing
	Hidden     int                 `json:"hidden"`                  // Hide Permission (0: don't hide 1: hide)
	ParentID   string              `json:"parent_id"`               // Parent ID
	ParentPath string              `json:"parent_path"`             // Parent path
	Creator    string              `json:"creator"`                 // Creator
	CreatedAt  time.Time           `json:"created_at"`              // Creation time
	Actions    PermissionActions   `json:"actions"`                 // Action list
	Resources  PermissionResources `json:"resources"`               // Resource list
}

// PermissionAction - Permission action object
type PermissionAction struct {
	Code string `json:"code"` // Action number
	Name string `json:"name"` // Action name
}

// PermissionResource - Permission resource object
type PermissionResource struct {
	Code   string `json:"code"`   // Resource number
	Name   string `json:"name"`   // Resource Name
	Method string `json:"method"` // Request method
	Path   string `json:"path"`   // Request path
}

// PermissionQueryParam - Query conditions
type PermissionQueryParam struct {
	UUIDs            []string // Record ID list
	LikeName         string   // Permission name (fuzzy query)
	Name             string   // Permission name
	ParentID         *string  // Parent ID
	PrefixParentPath string   // Parent path (prefix fuzzy query)
	Hidden           *int     // Hidden Permission
}

// PermissionQueryOptions - Query optional parameter items
type PermissionQueryOptions struct {
	PageParam        *schema.PaginationParam // Paging parameter
	IncludeActions   bool                    // Contains action list
	IncludeResources bool                    // Include resource list
}

// PermissionQueryResult - Search result
type PermissionQueryResult struct {
	Data       Permissions
	PageResult *schema.PaginationResult
}

// Permissions - Permission list
type Permissions []*Permission

// ToMap - Convert to key-value mapping
func (a Permissions) ToMap() map[string]*Permission {
	m := make(map[string]*Permission)
	for _, item := range a {
		m[item.UUID] = item
	}
	return m
}

// SplitAndGetAllUUIDs - Split parent path and get all record IDs
func (a Permissions) SplitAndGetAllUUIDs() []string {
	var UUIDs []string
	for _, item := range a {
		UUIDs = append(UUIDs, item.UUID)
		if item.ParentPath == "" {
			continue
		}

		pps := strings.Split(item.ParentPath, "/")
		for _, pp := range pps {
			var exists bool
			for _, UUID := range UUIDs {
				if pp == UUID {
					exists = true
					break
				}
			}
			if !exists {
				UUIDs = append(UUIDs, pp)
			}
		}
	}
	return UUIDs
}

// ToTrees - Convert to Permission list
func (a Permissions) ToTrees() PermissionTrees {
	list := make(PermissionTrees, len(a))
	for i, item := range a {
		list[i] = &PermissionTree{
			UUID:       item.UUID,
			Name:       item.Name,
			Sequence:   item.Sequence,
			Icon:       item.Icon,
			Router:     item.Router,
			Hidden:     item.Hidden,
			ParentID:   item.ParentID,
			ParentPath: item.ParentPath,
			Actions:    item.Actions,
			Resources:  item.Resources,
		}
	}
	return list
}

func (a Permissions) fillLeafNodeID(tree *[]*PermissionTree, leafNodeIDs *[]string) {
	for _, node := range *tree {
		if node.Children == nil || len(*node.Children) == 0 {
			*leafNodeIDs = append(*leafNodeIDs, node.UUID)
			continue
		}
		a.fillLeafNodeID(node.Children, leafNodeIDs)
	}
}

// ToLeafUUIDs - Convert to leaf node record ID list
func (a Permissions) ToLeafUUIDs() []string {
	var leafNodeIDs []string
	tree := a.ToTrees().ToTree()
	a.fillLeafNodeID(&tree, &leafNodeIDs)
	return leafNodeIDs
}

// PermissionResources - Permission resource list
type PermissionResources []*PermissionResource

// ForEach - Traversing resource data
func (a PermissionResources) ForEach(fn func(*PermissionResource, int)) PermissionResources {
	for i, item := range a {
		fn(item, i)
	}
	return a
}

// ToMap - Convert to key-value mapping
func (a PermissionResources) ToMap() map[string]*PermissionResource {
	m := make(map[string]*PermissionResource)
	for _, item := range a {
		m[item.Code] = item
	}
	return m
}

// PermissionActions - Permission action list
type PermissionActions []*PermissionAction

// PermissionTree - Permission tree
type PermissionTree struct {
	UUID       string              `json:"record_id"`               // Record ID
	Name       string              `json:"name" binding:"required"` // Permission name
	Sequence   int                 `json:"sequence"`                // Sort value
	Icon       string              `json:"icon"`                    // Permission icon
	Router     string              `json:"router"`                  // Access routing
	Hidden     int                 `json:"hidden"`                  // Hide Permission (0: don't hide 1: hide)
	ParentID   string              `json:"parent_id"`               // Parent ID
	ParentPath string              `json:"parent_path"`             // Parent path
	Resources  PermissionResources `json:"resources"`               // Resource list
	Actions    PermissionActions   `json:"actions"`                 // Action list
	Children   *[]*PermissionTree  `json:"children,omitempty"`      // Child tree
}

// PermissionTrees - Permission tree list
type PermissionTrees []*PermissionTree

// ForEach - Through the data entry
func (a PermissionTrees) ForEach(fn func(*PermissionTree, int)) PermissionTrees {
	for i, item := range a {
		fn(item, i)
	}
	return a
}

// ToTree - Convert to tree structure
func (a PermissionTrees) ToTree() []*PermissionTree {
	mi := make(map[string]*PermissionTree)
	for _, item := range a {
		mi[item.UUID] = item
	}

	var list []*PermissionTree
	for _, item := range a {
		if item.ParentID == "" {
			list = append(list, item)
			continue
		}
		if pitem, ok := mi[item.ParentID]; ok {
			if pitem.Children == nil {
				var children []*PermissionTree
				children = append(children, item)
				pitem.Children = &children
				continue
			}
			*pitem.Children = append(*pitem.Children, item)
		}
	}
	return list
}

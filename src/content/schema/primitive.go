package schema

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/MayCMF/core/src/common/schema"
)

// Primitive - Primitive object
type Primitive struct {
	ID         uint            `json:"id"`                      // Primitive ID
	UUID       string          `json:"uuid"`                    // UUID
	UID        int             `json:"uid" binding:"required"`  // User ID
	Slug       string          `json:"slug" binding:"required"` // Slug short machine name
	Parent     string          `json:"parent"`                  // Parent
	ParentPath string          `json:"parent_path"`             // Parent path
	Options    json.RawMessage `json:"options"`                 // Options in Jeson Format
	CreatedAt  time.Time       `json:"created_at"`              // Creation time
	UpdatedAt  time.Time       `json:"updated_at"`              // Updated time
	Variations Variations      `json:"variations"`              // Primitive Body with Languages
}

// Primitive Body - Primitive Body object
type PrimitiveBody struct {
	Slug      string    `json:"slug"`                        // Primitive Slug
	UID       int       `json:"uid"`                         // User ID
	Lang      string    `json:"language" binding:"required"` // Language Code Identifier
	Title     string    `json:"title" binding:"required"`    // Primitive Title
	Body      string    `json:"body"`                        // Primitive Body
	CreatedAt time.Time `json:"created_at"`                  // Creation time
	UpdatedAt time.Time `json:"updated_at"`                  // Updated time
}

// PrimitiveQueryParam - Query conditions
type PrimitiveQueryParam struct {
	UUIDs            []string // UUID list
	UID              int      // User ID, Creator
	Slug             string   // Slug list
	Slugs            []string // Slug list
	Lang             string   // Language of body
	Title            string   // Title
	Parent           *string  // Parent ID
	PrefixParentPath string   // Parent path (prefix fuzzy query)
	LikeSlug         string   // Slug (fuzzy query)
}

// PrimitiveQueryOptions - Primitive object query optional parameter item
type PrimitiveQueryOptions struct {
	PageParam         *schema.PaginationParam // Paging parameter
	IncludeVariations bool                    // Contains action list
}

// PrimitiveQueryResult - Primitive object query result
type PrimitiveQueryResult struct {
	Data       Primitives
	PageResult *schema.PaginationResult
}

// Primitives - Primitive list
type Primitives []*Primitive

// ToMap - Convert to key-value mapping
func (a Primitives) ToMap() map[string]*Primitive {
	m := make(map[string]*Primitive)
	for _, item := range a {
		m[item.UUID] = item
	}
	return m
}

// SplitAndGetAllSlugs - Split parent path and get all Slugs
func (a Primitives) SplitAndGetAllSlugs() []string {
	var Slugs []string
	for _, item := range a {
		Slugs = append(Slugs, item.Slug)
		if item.ParentPath == "" {
			continue
		}

		pps := strings.Split(item.ParentPath, "/")
		for _, pp := range pps {
			var exists bool
			for _, Slug := range Slugs {
				if pp == Slug {
					exists = true
					break
				}
			}
			if !exists {
				Slugs = append(Slugs, pp)
			}
		}
	}
	return Slugs
}

// ToTrees - Convert to Permission list
func (a Primitives) ToTrees() PrimitiveTrees {
	list := make(PrimitiveTrees, len(a))
	for i, item := range a {
		list[i] = &PrimitiveTree{
			ID:         item.ID,
			UUID:       item.UUID,
			UID:        item.UID,
			Slug:       item.Slug,
			Parent:     item.Parent,
			ParentPath: item.ParentPath,
			Options:    item.Options,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
			Variations: item.Variations,
		}
	}
	return list
}

func (a Primitives) fillLeafNodeID(tree *[]*PrimitiveTree, leafNodeIDs *[]string) {
	for _, node := range *tree {
		if node.Children == nil || len(*node.Children) == 0 {
			*leafNodeIDs = append(*leafNodeIDs, node.Slug)
			continue
		}
		a.fillLeafNodeID(node.Children, leafNodeIDs)
	}
}

// ToLeafSlugs - Convert to leaf node record ID list
func (a Primitives) ToLeafSlugs() []string {
	var leafNodeIDs []string
	tree := a.ToTrees().ToTree()
	a.fillLeafNodeID(&tree, &leafNodeIDs)
	return leafNodeIDs
}

// PermissionActions - Permission action list
type Variations []*PrimitiveBody

// PermissionTree - Permission tree
type PrimitiveTree struct {
	ID         uint              `json:"id"`
	UUID       string            `json:"uuid"`                   // Record UUID
	UID        int               `json:"uid" binding:"required"` // User ID
	Slug       string            `json:"slug"`                   // Sort value
	Parent     string            `json:"parent"`                 // Permission icon
	ParentPath string            `json:"parent_path"`            // Access routing
	Options    json.RawMessage   `json:"options"`                // Hide Permission (0: don't hide 1: hide)
	CreatedAt  time.Time         `json:"created_at"`             // Parent ID
	UpdatedAt  time.Time         `json:"updated_at"`             // Parent path
	Variations Variations        `json:"variations"`             // Resource list           // Action list
	Children   *[]*PrimitiveTree `json:"children,omitempty"`     // Child tree
}

// PermissionTrees - Primitive Tree list
type PrimitiveTrees []*PrimitiveTree

// ForEach - Through the data entry
func (a PrimitiveTrees) ForEach(fn func(*PrimitiveTree, int)) PrimitiveTrees {
	for i, item := range a {
		fn(item, i)
	}
	return a
}

// ToTree - Convert to tree structure
func (a PrimitiveTrees) ToTree() []*PrimitiveTree {
	mi := make(map[string]*PrimitiveTree)
	for _, item := range a {
		mi[item.Slug] = item
	}

	var list []*PrimitiveTree
	for _, item := range a {
		if item.Parent == "" {
			list = append(list, item)
			continue
		}
		if pitem, ok := mi[item.Parent]; ok {
			if pitem.Children == nil {
				var children []*PrimitiveTree
				children = append(children, item)
				pitem.Children = &children
				continue
			}
			*pitem.Children = append(*pitem.Children, item)
		}
	}
	return list
}

package schema

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/MayCMF/core/src/common/schema"
)

// Node - Node object
type Node struct {
	ID         uint            `json:"id"`                     // Node ID
	UUID       string          `json:"uuid"`                   // UUID
	Primitive  string          `json:"primitive"`              // Primitive Slug
	UID        int             `json:"uid" binding:"required"` // User ID
	Slug       string          `json:"slug"`                   // Slug short machine name
	Status     int             `json:"status"`                 // Node status (Published: 1, Draft: 0)
	Parent     string          `json:"parent"`                 // Parent
	ParentPath string          `json:"parent_path"`            // Parent path
	References json.RawMessage `json:"references"`             // References in JSON Format with reference fields and Primitives
	CreatedAt  time.Time       `json:"created_at"`             // Creation time
	UpdatedAt  time.Time       `json:"updated_at"`             // Updated time
	NodeBodies NodeBodies      `json:"variations"`             // Node Body with Languages
}

// Node Body - Node Body object
type NodeBody struct {
	NID       string    `json:"nid"`                         // Node Slug
	UID       int       `json:"uid"`                         // User ID
	Lang      string    `json:"language" binding:"required"` // Language Code Identifier
	Title     string    `json:"title" binding:"required"`    // Node Title
	Body      string    `json:"body"`                        // Node Body
	CreatedAt time.Time `json:"created_at"`                  // Creation time
	UpdatedAt time.Time `json:"updated_at"`                  // Updated time
}

// NodeQueryParam - Query conditions
type NodeQueryParam struct {
	UUIDs            []string // UUID list
	UID              int      // User ID, Creator
	Slug             string   // Short machine name
	Primitive        string   // Primitive Slug
	Slugs            []string // Slug list
	Lang             string   // Language of body
	Title            string   // Title
	Parent           *string  // Parent ID
	PrefixParentPath string   // Parent path (prefix fuzzy query)
	LikeSlug         string   // Slug (fuzzy query)
}

// NodeQueryOptions - Node object query optional parameter item
type NodeQueryOptions struct {
	PageParam         *schema.PaginationParam // Paging parameter
	IncludeNodeBodies bool                    // Contains Node Bodies List
}

// NodeQueryResult - Node object query result
type NodeQueryResult struct {
	Data       Nodes
	PageResult *schema.PaginationResult
}

// Nodes - Node list
type Nodes []*Node

// ToMap - Convert to key-value mapping
func (a Nodes) ToMap() map[string]*Node {
	m := make(map[string]*Node)
	for _, item := range a {
		m[item.UUID] = item
	}
	return m
}

// SplitAndGetAllSlugs - Split parent path and get all Slugs
func (a Nodes) SplitAndGetAllSlugs() []string {
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
func (a Nodes) ToTrees() NodeTrees {
	list := make(NodeTrees, len(a))
	for i, item := range a {
		list[i] = &NodeTree{
			ID:         item.ID,
			UUID:       item.UUID,
			UID:        item.UID,
			Primitive:  item.Primitive,
			Slug:       item.Slug,
			Status:     item.Status,
			Parent:     item.Parent,
			ParentPath: item.ParentPath,
			References: item.References,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
			NodeBodies: item.NodeBodies,
		}
	}
	return list
}

func (a Nodes) fillLeafNodeID(tree *[]*NodeTree, leafNodeIDs *[]string) {
	for _, node := range *tree {
		if node.Children == nil || len(*node.Children) == 0 {
			*leafNodeIDs = append(*leafNodeIDs, node.Slug)
			continue
		}
		a.fillLeafNodeID(node.Children, leafNodeIDs)
	}
}

// ToLeafSlugs - Convert to leaf node record ID list
func (a Nodes) ToLeafSlugs() []string {
	var leafNodeIDs []string
	tree := a.ToTrees().ToTree()
	a.fillLeafNodeID(&tree, &leafNodeIDs)
	return leafNodeIDs
}

// PermissionActions - Permission action list
type NodeBodies []*NodeBody

// PermissionTree - Permission tree
type NodeTree struct {
	ID         uint            `json:"id"`
	UUID       string          `json:"uuid"`                         // Record UUID
	UID        int             `json:"uid" binding:"required"`       // User ID
	Primitive  string          `json:"primitive" binding:"required"` // Primitive Slug
	Slug       string          `json:"slug"`                         // Sort value
	Parent     string          `json:"parent"`                       // Permission icon
	ParentPath string          `json:"parent_path"`                  // Access routing
	Status     int             `json:"status"`                       // Status (0: not published 1: published)
	NodeBodies NodeBodies      `json:"variations"`                   // Node Language Bodies
	References json.RawMessage `json:"references"`                   // References to fields or other Nodes
	CreatedAt  time.Time       `json:"created_at"`                   // Created Time
	UpdatedAt  time.Time       `json:"updated_at"`                   // Updated Time
	Children   *[]*NodeTree    `json:"children,omitempty"`           // Child tree
}

// PermissionTrees - Node Tree list
type NodeTrees []*NodeTree

// ForEach - Through the data entry
func (a NodeTrees) ForEach(fn func(*NodeTree, int)) NodeTrees {
	for i, item := range a {
		fn(item, i)
	}
	return a
}

// ToTree - Convert to tree structure
func (a NodeTrees) ToTree() []*NodeTree {
	mi := make(map[string]*NodeTree)
	for _, item := range a {
		mi[item.Slug] = item
	}

	var list []*NodeTree
	for _, item := range a {
		if item.Parent == "" {
			list = append(list, item)
			continue
		}
		if pitem, ok := mi[item.Parent]; ok {
			if pitem.Children == nil {
				var children []*NodeTree
				children = append(children, item)
				pitem.Children = &children
				continue
			}
			*pitem.Children = append(*pitem.Children, item)
		}
	}
	return list
}

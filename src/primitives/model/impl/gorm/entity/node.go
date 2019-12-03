package entity

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	account "github.com/MayCMF/core/src/account/model/impl/gorm/entity"
	"github.com/MayCMF/core/src/common/entity"
	i18n "github.com/MayCMF/core/src/i18n/model/impl/gorm/entity"
	"github.com/MayCMF/core/src/primitives/schema"
	"github.com/jinzhu/gorm"
)

// GetNodeDB - Get the Node store
func GetNodeDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(ctx, defDB, Node{})
}

// GetNodeDB - Get the Node store
func GetNodeBodyDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(ctx, defDB, NodeBody{})
}

// SchemaNode - Node object
type SchemaNode schema.Node

// ToNode - Convert to Node entity
func (a SchemaNode) ToNode() *Node {
	item := &Node{
		UUID:       a.UUID,
		UID:        a.UID,
		Primitive:  a.Primitive,
		Slug:       a.Slug,
		Parent:     a.Parent,
		ParentPath: a.ParentPath,
		Status:     a.Status,
		References: a.References,
	}
	return item
}

// ToPermissionActions - Convert to permission action list
func (a SchemaNode) ToNodeBodies() []*NodeBody {
	list := make([]*NodeBody, len(a.NodeBodies))
	for i, item := range a.NodeBodies {
		list[i] = SchemaNodeBody(*item).ToNodeBody(a.Slug)
	}
	return list
}

// Node - Node entity
type Node struct {
	entity.Model
	UUID       string          `gorm:"column:uuid;size:36;index;"`               // UUID
	User       account.User    `gorm:"foreignkey:UID;association_foreignkey:ID"` // Creator User ID
	UID        int             `gorm:"column:uid;"`                              // Creator User ID
	Primitive  string          `gorm:"column:primitive;size:100;"`               // Primitive Slug
	Slug       string          `gorm:"column:slug;size:100;unique_index;"`       // Slug short machine name
	Parent     string          `gorm:"column:parent;size:100;unique_index;"`     // Slug short machine name
	ParentPath string          `gorm:"column:parent_path"`                       // Parent path
	Status     int             `gorm:"column:status"`                            // Staus (1: published, 0: unpublished)
	References json.RawMessage `gorm:"column:references;type:jsonb;"`            // References in JSON Format
	NodeBodies NodeBodies
}

func (a Node) String() string {
	return entity.ToString(a)
}

// TableName - Table Name
func (a Node) TableName() string {
	return a.Model.TableName("node")
}

// ToSchemaNode - Convert to Node object
func (a Node) ToSchemaNode() *schema.Node {
	item := &schema.Node{
		ID:         a.ID,
		UUID:       a.UUID,
		UID:        a.UID,
		Primitive:  a.Primitive,
		Slug:       a.Slug,
		Parent:     a.Parent,
		ParentPath: a.ParentPath,
		References: a.References,
		CreatedAt:  a.CreatedAt,
	}
	return item
}

// Nodes - Node list
type Nodes []*Node

// ToSchemaNodes - Convert to Node object list
func (a Nodes) ToSchemaNodes() []*schema.Node {
	list := make([]*schema.Node, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaNode()
	}
	return list
}

// SchemaNodeBody NodeBody action object
type SchemaNodeBody schema.NodeBody

// ToNodeBody - Convert to Node Body entity
func (a SchemaNodeBody) ToNodeBody(NID string) *NodeBody {
	return &NodeBody{
		NID:       a.NID,
		UID:       a.UID,
		Lang:      a.Lang,
		Title:     &a.Title,
		Body:      &a.Body,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

// Node Body - Node Body object
type NodeBody struct {
	Node      Node          `gorm:"foreignkey:Slug;association_foreignkey:Slug"` // Node ID
	NID       string        `gorm:"column:nid"`                                  // Node ID
	User      account.User  `gorm:"foreignkey:UID;association_foreignkey:ID"`    // Creator User ID
	UID       int           `gorm:"column:uid;"`                                 // Creator User ID
	Language  i18n.Language `gorm:"foreignkey:Lang;association_foreignkey:Code"` // Language Code Identifieru se Code as foreign key
	Lang      string        `gorm:"column:language"`                             // Language Code Identifieru se Code as foreign key
	Title     *string       `gorm:"column:title" binding:"required"`             // Node Title
	Body      *string       `gorm:"column:body"`                                 // Node Body
	CreatedAt time.Time     `gorm:"column:created_at"`                           // Creation time
	UpdatedAt time.Time     `gorm:"column:updated_at"`                           // Updated time
}

// TableName - Table Name
func (a NodeBody) TableName() string {
	return fmt.Sprintf("%s%s", entity.GetTablePrefix(), "node_body")
}

// ToSchemaNodeBody - Convert to Node Body object
func (a NodeBody) ToSchemaNodeBody() *schema.NodeBody {
	item := &schema.NodeBody{
		NID:       a.NID,
		UID:       a.UID,
		Lang:      a.Lang,
		Title:     *a.Title,
		Body:      *a.Body,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
	return item
}

// NodeBodies - NodeBody lassociated entity ist
type NodeBodies []*NodeBody

// GetByNodeID - Get Node Body list based on Node ID
func (a NodeBodies) GetByNodeID(UUID string) []*schema.NodeBody {
	var list []*schema.NodeBody
	for _, item := range a {
		if item.NID == UUID {
			list = append(list, item.ToSchemaNodeBody())
		}
	}
	return list
}

// ToSchemaNodeBodies - Convert to Node Body variations action list
func (a NodeBodies) ToSchemaNodeBodies() []*schema.NodeBody {
	list := make([]*schema.NodeBody, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaNodeBody()
	}
	return list
}

// ToMap - Convert to key-value mapping
func (a NodeBodies) ToMap() map[string]*NodeBody {
	m := make(map[string]*NodeBody)
	for _, item := range a {
		m[item.Lang] = item
	}
	return m
}

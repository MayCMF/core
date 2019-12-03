package model

import (
	"context"

	"github.com/MayCMF/core/src/common/errors"
	"github.com/MayCMF/core/src/common/model"
	"github.com/MayCMF/core/src/primitives/model/impl/gorm/entity"
	"github.com/MayCMF/core/src/primitives/schema"
	"github.com/jinzhu/gorm"
)

// NewNode - Create a Node storage instance
func NewNode(db *gorm.DB) *Node {
	return &Node{db}
}

// Node - Node storage
type Node struct {
	db *gorm.DB
}

func (a *Node) getQueryOption(opts ...schema.NodeQueryOptions) schema.NodeQueryOptions {
	var opt schema.NodeQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query - Query data
func (a *Node) Query(ctx context.Context, params schema.NodeQueryParam, opts ...schema.NodeQueryOptions) (*schema.NodeQueryResult, error) {
	db := entity.GetNodeDB(ctx, a.db)
	if v := params.UUIDs; len(v) > 0 {
		db = db.Where("uuid=?", v)
	}
	if v := params.LikeSlug; v != "" {
		db = db.Where("slug LIKE ?", "%"+v+"%")
	}
	db = db.Order("id DESC")

	opt := a.getQueryOption(opts...)
	var list entity.Nodes
	pr, err := model.WrapPageQuery(ctx, db, opt.PageParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.NodeQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaNodes(),
	}

	err = a.fillSchemaNodes(ctx, qr.Data, opts...)
	if err != nil {
		return nil, err
	}

	return qr, nil
}

// Populate Node object data
func (a *Node) fillSchemaNodes(ctx context.Context, items []*schema.Node, opts ...schema.NodeQueryOptions) error {
	opt := a.getQueryOption(opts...)

	if opt.IncludeNodeBodies {

		nodeIDs := make([]string, len(items))
		for i, item := range items {
			nodeIDs[i] = item.UUID
		}

		var bodyList entity.NodeBodies
		if opt.IncludeNodeBodies {
			items, err := a.queryNodeBodies(ctx, nodeIDs...)
			if err != nil {
				return err
			}
			bodyList = items
		}

		for i, item := range items {
			if len(bodyList) > 0 {
				items[i].NodeBodies = bodyList.GetByNodeID(item.UUID)
			}
		}
	}

	return nil
}

// Get - Query specified data
func (a *Node) Get(ctx context.Context, UUID string, opts ...schema.NodeQueryOptions) (*schema.Node, error) {
	var item entity.Node
	ok, err := model.FindOne(ctx, entity.GetNodeDB(ctx, a.db).Where("uuid=?", UUID), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	sitem := item.ToSchemaNode()
	err = a.fillSchemaNodes(ctx, []*schema.Node{sitem}, opts...)
	if err != nil {
		return nil, err
	}

	return sitem, nil
}

// Create - Create data
func (a *Node) Create(ctx context.Context, item schema.Node) error {
	return model.ExecTrans(ctx, a.db, func(ctx context.Context) error {
		sitem := entity.SchemaNode(item)
		result := entity.GetNodeDB(ctx, a.db).Create(sitem.ToNode())

		if err := result.Error; err != nil {
			return errors.WithStack(err)
		}

		for _, item := range sitem.ToNodeBodies() {
			item.NID = sitem.UUID
			result := entity.GetNodeBodyDB(ctx, a.db).Create(item)
			if err := result.Error; err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	})
}

// Update - Update data
func (a *Node) Update(ctx context.Context, UUID string, item schema.Node) error {
	node := entity.SchemaNode(item).ToNode()
	result := entity.GetNodeDB(ctx, a.db).Where("uuid=?", UUID).Omit("uuid", "creator").Updates(node)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete - delete data
func (a *Node) Delete(ctx context.Context, UUID string) error {
	result := entity.GetNodeDB(ctx, a.db).Where("uuid=?", UUID).Delete(entity.Node{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *Node) queryNodeBodies(ctx context.Context, nodeIDs ...string) (entity.NodeBodies, error) {
	var list entity.NodeBodies
	result := entity.GetNodeBodyDB(ctx, a.db).Where("nid IN(?)", nodeIDs).Find(&list)
	if err := result.Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return list, nil
}

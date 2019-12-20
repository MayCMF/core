package implement

import (
	"context"

	"github.com/MayCMF/core/src/common/errors"
	commonschema "github.com/MayCMF/core/src/common/schema"
	"github.com/MayCMF/core/src/common/util"
	"github.com/MayCMF/core/src/primitives/model"
	"github.com/MayCMF/core/src/primitives/schema"
)

// NewNode - Create a Node
func NewNode(mNode model.INode) *Node {
	return &Node{
		NodeModel: mNode,
	}
}

// Node - Sample program
type Node struct {
	NodeModel model.INode
}

// Query - Query data
func (a *Node) Query(ctx context.Context, params schema.NodeQueryParam, opts ...schema.NodeQueryOptions) (*schema.NodeQueryResult, error) {
	return a.NodeModel.Query(ctx, params, opts...)
}

// Get - Get specified data
func (a *Node) Get(ctx context.Context, UUID string, opts ...schema.NodeQueryOptions) (*schema.Node, error) {
	item, err := a.NodeModel.Get(ctx, UUID, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

func (a *Node) checkSlug(ctx context.Context, slug string) error {
	result, err := a.NodeModel.Query(ctx, schema.NodeQueryParam{
		Slug: slug,
	}, schema.NodeQueryOptions{
		PageParam: &commonschema.PaginationParam{PageSize: -1},
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("Number already exists")
	}
	return nil
}

func (a *Node) getUpdate(ctx context.Context, UUID string) (*schema.Node, error) {
	return a.Get(ctx, UUID)
}

// Create - Create Node data
func (a *Node) Create(ctx context.Context, item schema.Node) (*schema.Node, error) {
	err := a.checkSlug(ctx, item.Slug)
	if err != nil {
		return nil, err
	}

	item.UUID = util.MustUUID()
	err = a.NodeModel.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	return a.getUpdate(ctx, item.UUID)
}

// Update - Update Node data
func (a *Node) Update(ctx context.Context, UUID string, item schema.Node) (*schema.Node, error) {
	oldItem, err := a.NodeModel.Get(ctx, UUID)
	if err != nil {
		return nil, err
	} else if oldItem == nil {
		return nil, errors.ErrNotFound
	} else if oldItem.Slug != item.Slug {
		err := a.checkSlug(ctx, item.Slug)
		if err != nil {
			return nil, err
		}
	}

	err = a.NodeModel.Update(ctx, UUID, item)
	if err != nil {
		return nil, err
	}
	return a.getUpdate(ctx, UUID)
}

// Delete - Delete data
func (a *Node) Delete(ctx context.Context, UUID string) error {
	oldItem, err := a.NodeModel.Get(ctx, UUID)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.NodeModel.Delete(ctx, UUID)
}

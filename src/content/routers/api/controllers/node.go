package controllers

import (
	"github.com/MayCMF/core/src/common/ginplus"
	"github.com/MayCMF/core/src/primitives/controllers"
	"github.com/MayCMF/core/src/primitives/schema"
	"github.com/gin-gonic/gin"
)

// NewNode - Create a Node controller
func NewNode(bNode controllers.INode) *Node {
	return &Node{
		NodeBll: bNode,
	}
}

// Node - Sample program
type Node struct {
	NodeBll controllers.INode
}

// Query - Query data
// @Tags Node
// @Summary Query data
// @Param Authorization header string false "Bearer User Token"
// @Param current query int true "Page Index" default(1)
// @Param pageSize query int true "Paging Size" default(10)
// @Param code query string false "Numbering"
// @Param name query string false "Name"
// @Param status query int false "Status (1: Enable 2: Disable)"
// @Success 200 {array} schema.Node "Search result: {list:List data,pagination:{current:Page index, pageSize: Page size, total: The total number}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/node [get]
func (a *Node) Query(c *gin.Context) {
	// var params schema.NodeQueryParam
	// params.LikeSlug = c.Query("slug")

	result, err := a.NodeBll.Query(ginplus.NewContext(c), schema.NodeQueryParam{}, schema.NodeQueryOptions{
		PageParam:         ginplus.GetPaginationParam(c),
		IncludeNodeBodies: true,
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResList(c, result.Data.ToTrees().ToTree())
}

// Get - Query specified data
// @Tags Node
// @Summary Query specified data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Success 200 {object} schema.Node
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 404 {object} schema.HTTPError "{error:{code:0,message: Resource does not exist.}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/node/{id} [get]
func (a *Node) Get(c *gin.Context) {
	item, err := a.NodeBll.Get(ginplus.NewContext(c), c.Param("id"), schema.NodeQueryOptions{
		IncludeNodeBodies: true,
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item)
}

// Create - Create data
// @Tags Node
// @Summary Create data
// @Param Authorization header string false "Bearer User Token"
// @Param body body schema.Node true "Create data"
// @Success 200 {object} schema.Node
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/node [post]
func (a *Node) Create(c *gin.Context) {
	var item schema.Node
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	item.UID = ginplus.GetUserID(c)

	nitem, err := a.NodeBll.Create(ginplus.NewContext(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Update - Update data
// @Tags Node
// @Summary Update data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Param body body schema.Node true "Update data"
// @Success 200 {object} schema.Node
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/node/{id} [put]
func (a *Node) Update(c *gin.Context) {
	var item schema.Node
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	nitem, err := a.NodeBll.Update(ginplus.NewContext(c), c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Delete - Delete data
// @Tags Node
// @Summary Delete data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/node/{id} [delete]
func (a *Node) Delete(c *gin.Context) {
	err := a.NodeBll.Delete(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// // Publish - Enable data
// // @Tags Node
// // @Summary Enable data
// // @Param Authorization header string false "Bearer User Token"
// // @Param id path string true "Record ID"
// // @Success 200 {object} schema.HTTPStatus "{status:OK}"
// // @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// // @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// // @Router /api/v1/node/{id}/publish [patch]
// func (a *Node) Enable(c *gin.Context) {
// 	err := a.NodeBll.UpdateStatus(ginplus.NewContext(c), c.Param("id"), 1)
// 	if err != nil {
// 		ginplus.ResError(c, err)
// 		return
// 	}
// 	ginplus.ResOK(c)
// }

// // Unpublish - Unpublish node
// // @Tags Node
// // @Summary Unpublish node
// // @Param Authorization header string false "Bearer User Token"
// // @Param id path string true "Record ID"
// // @Success 200 {object} schema.HTTPStatus "{status: OK}"
// // @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// // @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// // @Router /api/v1/node/{id}/unpublish [patch]
// func (a *Node) Disable(c *gin.Context) {
// 	err := a.NodeBll.UpdateStatus(ginplus.NewContext(c), c.Param("id"), 0)
// 	if err != nil {
// 		ginplus.ResError(c, err)
// 		return
// 	}
// 	ginplus.ResOK(c)
// }

package controllers

import (
	"github.com/MayCMF/core/src/common/ginplus"
	"github.com/MayCMF/core/src/primitives/controllers"
	"github.com/MayCMF/core/src/primitives/schema"
	"github.com/gin-gonic/gin"
)

// NewPrimitive - Create a primitive controller
func NewPrimitive(bPrimitive controllers.IPrimitive) *Primitive {
	return &Primitive{
		PrimitiveBll: bPrimitive,
	}
}

// Primitive - Sample program
type Primitive struct {
	PrimitiveBll controllers.IPrimitive
}

// Query - Query data
// @Tags Primitive
// @Summary Query data
// @Param Authorization header string false "Bearer User Token"
// @Param current query int true "Page Index" default(1)
// @Param pageSize query int true "Paging Size" default(10)
// @Param code query string false "Numbering"
// @Param name query string false "Name"
// @Param status query int false "Status (1: Enable 2: Disable)"
// @Success 200 {array} schema.Primitive "Search result: {list:List data,pagination:{current:Page index, pageSize: Page size, total: The total number}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/primitives [get]
func (a *Primitive) Query(c *gin.Context) {
	// var params schema.PrimitiveQueryParam
	// params.LikeSlug = c.Query("slug")

	result, err := a.PrimitiveBll.Query(ginplus.NewContext(c), schema.PrimitiveQueryParam{}, schema.PrimitiveQueryOptions{
		PageParam:         ginplus.GetPaginationParam(c),
		IncludeVariations: true,
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResList(c, result.Data.ToTrees().ToTree())
}

// Get - Query specified data
// @Tags Primitive
// @Summary Query specified data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Success 200 {object} schema.Primitive
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 404 {object} schema.HTTPError "{error:{code:0,message: Resource does not exist.}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/primitives/{id} [get]
func (a *Primitive) Get(c *gin.Context) {
	item, err := a.PrimitiveBll.Get(ginplus.NewContext(c), c.Param("id"), schema.PrimitiveQueryOptions{
		IncludeVariations: true,
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item)
}

// Create - Create data
// @Tags Primitive
// @Summary Create data
// @Param Authorization header string false "Bearer User Token"
// @Param body body schema.Primitive true "Create data"
// @Success 200 {object} schema.Primitive
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/primitives [post]
func (a *Primitive) Create(c *gin.Context) {
	var item schema.Primitive
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	item.UID = ginplus.GetUserID(c)

	nitem, err := a.PrimitiveBll.Create(ginplus.NewContext(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Update - Update data
// @Tags Primitive
// @Summary Update data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Param body body schema.Primitive true "Update data"
// @Success 200 {object} schema.Primitive
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/primitives/{id} [put]
func (a *Primitive) Update(c *gin.Context) {
	var item schema.Primitive
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	nitem, err := a.PrimitiveBll.Update(ginplus.NewContext(c), c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Delete - Delete data
// @Tags Primitive
// @Summary Delete data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/primitives/{id} [delete]
func (a *Primitive) Delete(c *gin.Context) {
	err := a.PrimitiveBll.Delete(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

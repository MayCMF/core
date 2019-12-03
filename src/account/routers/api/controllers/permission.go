package controllers

import (
	"github.com/MayCMF/core/src/account/controllers"
	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common/ginplus"
	"github.com/MayCMF/core/src/common/util"
	"github.com/gin-gonic/gin"
)

// NewPermission - Create a Permission management controller
func NewPermission(bPermission controllers.IPermission) *Permission {
	return &Permission{
		PermissionBll: bPermission,
	}
}

// Permission - Manage Permission
type Permission struct {
	PermissionBll controllers.IPermission
}

// Query - Query data
// @Tags Manage Permission
// @Summary Query data
// @Param Authorization header string false "Bearer User Token"
// @Param current query int true "Paging index" default(1)
// @Param pageSize query int true "Paging Size" default(10)
// @Param name query string false "Name"
// @Param hidden query int false "Hide permission (0: don't hide 1: hide)"
// @Param parentID query string false "Parent ID"
// @Success 200 {array} schema.Permission "Search result: {list:List data,pagination:{current:Page index,pageSize:Page size,total:total number}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/permissions [get]
func (a *Permission) Query(c *gin.Context) {
	params := schema.PermissionQueryParam{
		LikeName: c.Query("name"),
	}

	if v := c.Query("parentID"); v != "" {
		params.ParentID = &v
	}

	if v := c.Query("hidden"); v != "" {
		if hidden := util.S(v).DefaultInt(0); hidden > -1 {
			params.Hidden = &hidden
		}
	}

	result, err := a.PermissionBll.Query(ginplus.NewContext(c), params, schema.PermissionQueryOptions{
		PageParam: ginplus.GetPaginationParam(c),
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResPage(c, result.Data, result.PageResult)
}

// QueryTree - Query permission tree
// @Tags Manage Permission
// @Summary Query permission tree
// @Param Authorization header string false "Bearer User Token"
// @Param includeActions query int false "Whether to include action data (1 is)"
// @Param includeResources query int false "Whether to include resource data (1 is)"
// @Success 200 {array} schema.PermissionTree "Search result: {list: List data}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/permissions.tree [get]
func (a *Permission) QueryTree(c *gin.Context) {
	result, err := a.PermissionBll.Query(ginplus.NewContext(c), schema.PermissionQueryParam{}, schema.PermissionQueryOptions{
		IncludeActions:   c.Query("includeActions") == "1",
		IncludeResources: c.Query("includeResources") == "1",
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResList(c, result.Data.ToTrees().ToTree())
}

// Get - Query specified data
// @Tags Manage Permission
// @Summary Query specified data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Success 200 {object} schema.Permission
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 404 {object} schema.HTTPError "{error:{code:0,message:Resource does not exist.}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/permissions/{id} [get]
func (a *Permission) Get(c *gin.Context) {
	item, err := a.PermissionBll.Get(ginplus.NewContext(c), c.Param("id"), schema.PermissionQueryOptions{
		IncludeActions:   true,
		IncludeResources: true,
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item)
}

// Create - Create data
// @Tags Manage Permission
// @Summary Create data
// @Param Authorization header string false "Bearer User Token"
// @Param body body schema.Permission true "Create data"
// @Success 200 {object} schema.Permission
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/permissions [post]
func (a *Permission) Create(c *gin.Context) {
	var item schema.Permission
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	item.Creator = ginplus.GetUserUUID(c)
	nitem, err := a.PermissionBll.Create(ginplus.NewContext(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Update - Update data
// @Tags Manage Permission
// @Summary Update data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Param body body schema.Permission true "Update data"
// @Success 200 {object} schema.Permission
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/permissions/{id} [put]
func (a *Permission) Update(c *gin.Context) {
	var item schema.Permission
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	nitem, err := a.PermissionBll.Update(ginplus.NewContext(c), c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Delete - Delete data
// @Tags Manage Permission
// @Summary Delete data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/permissions/{id} [delete]
func (a *Permission) Delete(c *gin.Context) {
	err := a.PermissionBll.Delete(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

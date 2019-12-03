package controllers

import (
	"github.com/MayCMF/core/src/account/controllers"
	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common/ginplus"
	"github.com/gin-gonic/gin"
)

// NewRole - Create Role management controller
func NewRole(bRole controllers.IRole) *Role {
	return &Role{
		RoleBll: bRole,
	}
}

// Role - Manage Role
type Role struct {
	RoleBll controllers.IRole
}

// Query - Query data
// @Tags Manage Role
// @Summary Query data
// @Param Authorization header string false "Bearer User Token"
// @Param current query int true "Paging index" default(1)
// @Param pageSize query int true "Paging Size" default(10)
// @Param name query string false "Role name (fuzzy query)"
// @Success 200 {array} schema.Role "Search result: {list:List data,pagination:{current:Page index,pageSize:Page size,total:Total number}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/roles [get]
func (a *Role) Query(c *gin.Context) {
	var params schema.RoleQueryParam
	params.LikeName = c.Query("name")

	result, err := a.RoleBll.Query(ginplus.NewContext(c), params, schema.RoleQueryOptions{
		PageParam: ginplus.GetPaginationParam(c),
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResPage(c, result.Data, result.PageResult)
}

// QuerySelect - Query selection data
// @Tags Manage Role
// @Summary Query selection data
// @Param Authorization header string false "Bearer User Token"
// @Success 200 {array} schema.Role "Search result: {list: Role list}"
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Unknown query type}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/roles.select [get]
func (a *Role) QuerySelect(c *gin.Context) {
	result, err := a.RoleBll.Query(ginplus.NewContext(c), schema.RoleQueryParam{})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResList(c, result.Data)
}

// Get - Query specified data
// @Tags Manage Role
// @Summary Query specified data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Success 200 {object} schema.Role
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 404 {object} schema.HTTPError "{error:{code:0,message:Resource does not exist.}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/roles/{id} [get]
func (a *Role) Get(c *gin.Context) {
	item, err := a.RoleBll.Get(ginplus.NewContext(c), c.Param("id"), schema.RoleQueryOptions{
		IncludePermissions: true,
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item)
}

// Create - Create data
// @Tags Manage Role
// @Summary Create data
// @Param Authorization header string false "Bearer User Token"
// @Param body body schema.Role true "Create data"
// @Success 200 {object} schema.Role
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/roles [post]
func (a *Role) Create(c *gin.Context) {
	var item schema.Role
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	item.Creator = ginplus.GetUserUUID(c)
	nitem, err := a.RoleBll.Create(ginplus.NewContext(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, nitem)
}

// Update - Update data
// @Tags Manage Role
// @Summary Update data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Param body body schema.Role true "Update data"
// @Success 200 {object} schema.Role
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/roles/{id} [put]
func (a *Role) Update(c *gin.Context) {
	var item schema.Role
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	nitem, err := a.RoleBll.Update(ginplus.NewContext(c), c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Delete - Delete data
// @Tags Manage Role
// @Summary Delete data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/roles/{id} [delete]
func (a *Role) Delete(c *gin.Context) {
	err := a.RoleBll.Delete(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

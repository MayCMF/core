package controllers

import (
	"strings"

	"github.com/MayCMF/core/src/account/controllers"
	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common/ginplus"
	"github.com/MayCMF/core/src/common/util"
	"github.com/gin-gonic/gin"
)

// NewUser - Create a User Management Controller
func NewUser(bUser controllers.IUser) *User {
	return &User{
		UserBll: bUser,
	}
}

// User - Manage Users
type User struct {
	UserBll controllers.IUser
}

// Query - Query data
// @Tags Manage Users
// @Summary Query data
// @Param Authorization header string false "Bearer User Token"
// @Param current query int true "Paging index" default(1)
// @Param pageSize query int true "Paging Size" default(10)
// @Param userName query string false "Username (fuzzy query)"
// @Param realName query string false "Real name (fuzzy query)"
// @Param roleIDs query string false "Role ID (multiple separated by commas)"
// @Param status query int false "Status (1: Enable 2: Disable)"
// @Success 200 {array} schema.UserShow "Search result: {list:List data,pagination:{current:Page index,pageSize:Page size,total:Total number}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/users [get]
func (a *User) Query(c *gin.Context) {
	var params schema.UserQueryParam
	params.LikeUserName = c.Query("userName")
	params.LikeRealName = c.Query("realName")
	if v := util.S(c.Query("status")).DefaultInt(0); v > 0 {
		params.Status = v
	}

	if v := c.Query("roleIDs"); v != "" {
		params.RoleIDs = strings.Split(v, ",")
	}

	result, err := a.UserBll.QueryShow(ginplus.NewContext(c), params, schema.UserQueryOptions{
		IncludeRoles: true,
		PageParam:    ginplus.GetPaginationParam(c),
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResPage(c, result.Data, result.PageResult)
}

// Get - Query specified data
// @Tags Manage Users
// @Summary Query specified data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Success 200 {object} schema.User
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 404 {object} schema.HTTPError "{error:{code:0,message: Resource does not exist.}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/users/{id} [get]
func (a *User) Get(c *gin.Context) {
	item, err := a.UserBll.Get(ginplus.NewContext(c), c.Param("id"), schema.UserQueryOptions{
		IncludeRoles: true,
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item.CleanSecure())
}

// Create - Create data
// @Tags Manage Users
// @Summary Create data
// @Param Authorization header string false "Bearer User Token"
// @Param body body schema.User true "Create data"
// @Success 200 {object} schema.User
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/users [post]
func (a *User) Create(c *gin.Context) {
	var item schema.User
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	item.Creator = ginplus.GetUserUUID(c)
	nitem, err := a.UserBll.Create(ginplus.NewContext(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem.CleanSecure())
}

// Update - Update data
// @Tags Manage Users
// @Summary Update data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Param body body schema.User true "Update data"
// @Success 200 {object} schema.User
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/users/{id} [put]
func (a *User) Update(c *gin.Context) {
	var item schema.User
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	nitem, err := a.UserBll.Update(ginplus.NewContext(c), c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem.CleanSecure())
}

// Delete - Delete data
// @Tags Manage Users
// @Summary Delete data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/users/{id} [delete]
func (a *User) Delete(c *gin.Context) {
	err := a.UserBll.Delete(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Enable - Enable data
// @Tags Manage Users
// @Summary Enable data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/users/{id}/enable [patch]
func (a *User) Enable(c *gin.Context) {
	err := a.UserBll.UpdateStatus(ginplus.NewContext(c), c.Param("id"), 1)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Disable - Disable data
// @Tags Manage Users
// @Summary Disable data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/users/{id}/disable [patch]
func (a *User) Disable(c *gin.Context) {
	err := a.UserBll.UpdateStatus(ginplus.NewContext(c), c.Param("id"), 2)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

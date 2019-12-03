package controllers

import (
	"github.com/MayCMF/core/src/common/ginplus"
	"github.com/MayCMF/core/src/common/util"
	"github.com/MayCMF/src/filemanager/controllers"
	"github.com/MayCMF/src/filemanager/schema"
	"github.com/gin-gonic/gin"
)

// NewFile - Create a File controller
func NewFile(bFile controllers.IFile) *File {
	return &File{
		FileBll: bFile,
	}
}

// File - Sample File entity
type File struct {
	FileBll controllers.IFile
}

// Query - Query data
// @Tags File
// @Summary Query data
// @Param Authorization header string false "Bearer User Token"
// @Param current query int true "Page Index" default(1)
// @Param pageSize query int true "Paging Size" default(10)
// @Param code query string false "Numbering"
// @Param name query string false "Name"
// @Param status query int false "Status (1: Enable 2: Disable)"
// @Success 200 {array} schema.File "Search result: {list:List data,pagination:{current:Page index, pageSize: Page size, total: The total number}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/file [get]
func (a *File) Query(c *gin.Context) {
	var params schema.FileQueryParam
	params.LikeCode = c.Query("code")
	params.LikeName = c.Query("name")
	params.Status = util.S(c.Query("status")).DefaultInt(0)

	result, err := a.FileBll.Query(ginplus.NewContext(c), params, schema.FileQueryOptions{
		PageParam: ginplus.GetPaginationParam(c),
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResPage(c, result.Data, result.PageResult)
}

// Get - Query specified data
// @Tags File
// @Summary Query specified data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Success 200 {object} schema.File
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 404 {object} schema.HTTPError "{error:{code:0,message: Resource does not exist.}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/file/{id} [get]
func (a *File) Get(c *gin.Context) {
	item, err := a.FileBll.Get(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item)
}

// Create - Create data
// @Tags File
// @Summary Create data
// @Param Authorization header string false "Bearer User Token"
// @Param body body schema.File true "Create data"
// @Success 200 {object} schema.File
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/file [post]
func (a *File) Create(c *gin.Context) {
	var item schema.File
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	item.Creator = ginplus.GetUserUUID(c)
	nitem, err := a.FileBll.Create(ginplus.NewContext(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Update - Update data
// @Tags File
// @Summary Update data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Param body body schema.File true "Update data"
// @Success 200 {object} schema.File
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/file/{id} [put]
func (a *File) Update(c *gin.Context) {
	var item schema.File
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	nitem, err := a.FileBll.Update(ginplus.NewContext(c), c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Delete - Delete data
// @Tags File
// @Summary Delete data
// @Param Authorization header string false "Bearer User Token"
// @Param id path string true "Record ID"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/file/{id} [delete]
func (a *File) Delete(c *gin.Context) {
	err := a.FileBll.Delete(ginplus.NewContext(c), c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

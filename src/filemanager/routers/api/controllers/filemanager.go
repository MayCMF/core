package controllers

import (
	"path"
	"time"

	"github.com/MayCMF/core/src/common/ginplus"
	"github.com/MayCMF/core/src/filemanager/controllers"
	"github.com/MayCMF/core/src/filemanager/schema"
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
// @Param filename query string false "Numbering"
// @Param name query string false "Name"
// @Param status query int false "Status (1: Enable 2: Disable)"
// @Success 200 {array} schema.File "Search result: {list:List data,pagination:{current:Page index, pageSize: Page size, total: The total number}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/file [get]
func (a *File) Query(c *gin.Context) {
	var params schema.FileQueryParam
	params.Filename = c.Query("filename")
	params.LikeFilename = c.Query("filename")
	params.Uri = c.Query("uri")

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

	item.UUID = ginplus.GetUserUUID(c)
	item.UID = uint(ginplus.GetUserID(c))
	nitem, err := a.FileBll.Create(ginplus.NewContext(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Upload - Upload File
// @Tags File
// @Summary Upload File
// @Param Authorization header string false "Bearer User Token"
// @Param body body schema.File true "Create data"
// @Success 200 {object} schema.File
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/file/upload [post]
func (a *File) Upload(c *gin.Context) {
	var item schema.File

	// Multipart form
	// uid, err := strconv.ParseUint(c.PostForm("UserID"), 10, 32)
	// item.UID = uint(uid)
	// item.UserUUID = ginplus.GetUserUUID(c)
	uri := c.PostForm("URL")
	form, err := c.MultipartForm()
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	files := form.File["MayFile"]
	dst := "static" + item.Uri
	for _, file := range files {
		// Upload the file to specific dst.
		controllers.CheckDir(dst + uri)
		item.Uri = dst + uri + "/" + time.Now().Format("20060102-1504") + "_" + file.Filename
		item.Filename = file.Filename
		item.Filesize = file.Size
		item.Filemime = file.Header.Get("Content-Type")
		item.FileExt = path.Ext(file.Filename)
		item.UID = uint(ginplus.GetUserID(c))
		nitem, err := a.FileBll.Upload(ginplus.NewContext(c), item)
		if err != nil {
			ginplus.ResError(c, err)
			return
		}
		c.SaveUploadedFile(file, item.Uri)
		ginplus.ResSuccess(c, nitem)
	}
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

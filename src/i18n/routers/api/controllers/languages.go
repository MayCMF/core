package controllers

import (
	"github.com/MayCMF/core/src/common/ginplus"
	"github.com/MayCMF/core/src/common/util"
	"github.com/MayCMF/core/src/i18n/controllers"
	"github.com/MayCMF/core/src/i18n/schema"
	"github.com/gin-gonic/gin"
)

// NewLanguage - Create a language controller
func NewLanguage(bLanguage controllers.ILanguage) *Language {
	return &Language{
		LanguageBll: bLanguage,
	}
}

// Language - Sample program
type Language struct {
	LanguageBll controllers.ILanguage
}

// Query - Query data
// @Tags Language
// @Summary Query data
// @Param Authorization header string false "Bearer User Token"
// @Param current query int true "Page Index" default(1)
// @Param pageSize query int true "Paging Size" default(10)
// @Param code query string false "Numbering"
// @Param name query string false "Name"
// @Param status query int false "Status (1: Enable 2: Disable)"
// @Success 200 {array} schema.Language "Search result: {list:List data,pagination:{current:Page index, pageSize: Page size, total: The total number}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/language [get]
func (a *Language) Query(c *gin.Context) {
	var params schema.LanguageQueryParam
	params.LikeCode = c.Query("code")
	params.LikeName = c.Query("name")
	params.Status = util.S(c.Query("active")).DefaultInt(0)

	result, err := a.LanguageBll.Query(ginplus.NewContext(c), params, schema.LanguageQueryOptions{
		PageParam: ginplus.GetPaginationParam(c),
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResPage(c, result.Data, result.PageResult)
}

// Get - Query specified data
// @Tags Language
// @Summary Query specified data
// @Param Authorization header string false "Bearer User Token"
// @Param code path string true "code"
// @Success 200 {object} schema.Language
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 404 {object} schema.HTTPError "{error:{code:0,message: Resource does not exist.}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/language/{code} [get]
func (a *Language) Get(c *gin.Context) {
	item, err := a.LanguageBll.Get(ginplus.NewContext(c), c.Param("code"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item)
}

// Create - Create data
// @Tags Language
// @Summary Create data
// @Param Authorization header string false "Bearer User Token"
// @Param body body schema.Language true "Create data"
// @Success 200 {object} schema.Language
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/language [post]
func (a *Language) Create(c *gin.Context) {
	var item schema.Language
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	nitem, err := a.LanguageBll.Create(ginplus.NewContext(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Update - Update data
// @Tags Language
// @Summary Update data
// @Param Authorization header string false "Bearer User Token"
// @Param code path string true "code"
// @Param body body schema.Language true "Update data"
// @Success 200 {object} schema.Language
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/language/{code} [put]
func (a *Language) Update(c *gin.Context) {
	var item schema.Language
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	nitem, err := a.LanguageBll.Update(ginplus.NewContext(c), c.Param("code"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Delete - Delete data
// @Tags Language
// @Summary Delete data
// @Param Authorization header string false "Bearer User Token"
// @Param code path string true "code"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/language/{code} [delete]
func (a *Language) Delete(c *gin.Context) {
	err := a.LanguageBll.Delete(ginplus.NewContext(c), c.Param("code"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Enable - Enable data
// @Tags Language
// @Summary Enable data
// @Param Authorization header string false "Bearer User Token"
// @Param code path string true "code"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/language/{code}/enable [patch]
func (a *Language) Enable(c *gin.Context) {
	err := a.LanguageBll.UpdateStatus(ginplus.NewContext(c), c.Param("code"), 1)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Disable - Disable data
// @Tags Language
// @Summary Disable data
// @Param Authorization header string false "Bearer User Token"
// @Param code path string true "code"
// @Success 200 {object} schema.HTTPStatus "{status: OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/language/{code}/disable [patch]
func (a *Language) Disable(c *gin.Context) {
	err := a.LanguageBll.UpdateStatus(ginplus.NewContext(c), c.Param("code"), 0)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

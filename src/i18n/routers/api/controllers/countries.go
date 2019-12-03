package controllers

import (
	"github.com/MayCMF/core/src/common/ginplus"
	"github.com/MayCMF/core/src/i18n/controllers"
	"github.com/MayCMF/core/src/i18n/schema"
	"github.com/gin-gonic/gin"
)

// NewCountry - Create a country controller
func NewCountry(bCountry controllers.ICountry) *Country {
	return &Country{
		CountryBll: bCountry,
	}
}

// Country - Sample program
type Country struct {
	CountryBll controllers.ICountry
}

// Query - Query data
// @Tags Country
// @Summary Query data
// @Param Authorization header string false "Bearer User Token"
// @Param current query int true "Page Index" default(1)
// @Param pageSize query int true "Paging Size" default(10)
// @Param code query string false "Numbering"
// @Param name query string false "Name"
// @Param status query int false "Status (1: Enable 2: Disable)"
// @Success 200 {array} schema.Country "Search result: {list:List data,pagination:{current:Page index, pageSize: Page size, total: The total number}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/countries [get]
func (a *Country) Query(c *gin.Context) {
	var params schema.CountryQueryParam
	params.LikeCode = c.Query("code")
	params.LikeName = c.Query("name")

	result, err := a.CountryBll.Query(ginplus.NewContext(c), params, schema.CountryQueryOptions{
		PageParam: ginplus.GetPaginationParam(c),
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResPage(c, result.Data, result.PageResult)
}

// Get - Query specified data
// @Tags Country
// @Summary Query specified data
// @Param Authorization header string false "Bearer User Token"
// @Param code path string true "code"
// @Success 200 {object} schema.Country
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 404 {object} schema.HTTPError "{error:{code:0,message: Resource does not exist.}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/countries/{code} [get]
func (a *Country) Get(c *gin.Context) {
	item, err := a.CountryBll.Get(ginplus.NewContext(c), c.Param("code"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item)
}

// Create - Create data
// @Tags Country
// @Summary Create data
// @Param Authorization header string false "Bearer User Token"
// @Param body body schema.Country true "Create data"
// @Success 200 {object} schema.Country
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/countries [post]
func (a *Country) Create(c *gin.Context) {
	var item schema.Country
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	nitem, err := a.CountryBll.Create(ginplus.NewContext(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Update - Update data
// @Tags Country
// @Summary Update data
// @Param Authorization header string false "Bearer User Token"
// @Param code path string true "code"
// @Param body body schema.Country true "Update data"
// @Success 200 {object} schema.Country
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/countries/{code} [put]
func (a *Country) Update(c *gin.Context) {
	var item schema.Country
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	nitem, err := a.CountryBll.Update(ginplus.NewContext(c), c.Param("code"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, nitem)
}

// Delete - Delete data
// @Tags Country
// @Summary Delete data
// @Param Authorization header string false "Bearer User Token"
// @Param code path string true "code"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/countries/{code} [delete]
func (a *Country) Delete(c *gin.Context) {
	err := a.CountryBll.Delete(ginplus.NewContext(c), c.Param("code"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

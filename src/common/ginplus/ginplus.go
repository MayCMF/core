package ginplus

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	icontext "github.com/MayCMF/core/src/common/context"
	"github.com/MayCMF/core/src/common/errors"
	"github.com/MayCMF/core/src/common/logger"
	"github.com/MayCMF/core/src/common/schema"
	"github.com/MayCMF/core/src/common/util"
	"github.com/gin-gonic/gin"
)

// Define the keys in the context
const (
	prefix = "maycms"
	// UserIDKey - Key in the storage context (user ID)
	UserIDKey = prefix + "/user-id"
	// UserUUIDKey - Key in the storage context (user UUID)
	UserUUIDKey = prefix + "/user-uuid"
	// TraceIDKey - Key in storage context (tracking ID)
	TraceIDKey = prefix + "/trace-id"
	// ResBodyKey - The key in the storage context (response to the Body data)
	ResBodyKey = prefix + "/res-body"
)

// NewContext - Package context entry
func NewContext(c *gin.Context) context.Context {
	parent := context.Background()

	if v := GetTraceID(c); v != "" {
		parent = icontext.NewTraceID(parent, v)
		parent = logger.NewTraceIDContext(parent, GetTraceID(c))
	}

	if v := GetUserUUID(c); v != "" {
		parent = icontext.NewUserUUID(parent, v)
		parent = logger.NewUserUUIDContext(parent, v)
	}

	return parent
}

// GetToken - Get user token
func GetToken(c *gin.Context) string {
	var token string
	auth := c.GetHeader("Authorization")
	prefix := "Bearer "
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token
}

// GetPageIndex - Get paged page index
func GetPageIndex(c *gin.Context) int {
	defaultVal := 1
	if v := c.Query("current"); v != "" {
		if iv := util.S(v).DefaultInt(defaultVal); iv > 0 {
			return iv
		}
	}
	return defaultVal
}

// GetPageSize - Get the page size of the page (up to 50)
func GetPageSize(c *gin.Context) int {
	defaultVal := 50
	if v := c.Query("pageSize"); v != "" {
		if iv := util.S(v).DefaultInt(defaultVal); iv > 0 {
			if iv > 50 {
				iv = 50
			}
			return iv
		}
	}
	return defaultVal
}

// GetPaginationParam - Get paging query parameters
func GetPaginationParam(c *gin.Context) *schema.PaginationParam {
	return &schema.PaginationParam{
		PageIndex: GetPageIndex(c),
		PageSize:  GetPageSize(c),
	}
}

// GetTraceID - Get tracking ID
func GetTraceID(c *gin.Context) string {
	return c.GetString(TraceIDKey)
}

// GetUserID - Get user ID
func GetUserID(c *gin.Context) int {
	return c.GetInt(UserIDKey)
}

// SetUserID - Set user ID
func SetUserID(c *gin.Context, ID int) {
	c.Set(UserIDKey, ID)
}

// GetUserUUID - Get user UUID
func GetUserUUID(c *gin.Context) string {
	return c.GetString(UserUUIDKey)
}

// SetUserUUID - Set user UUID
func SetUserUUID(c *gin.Context, userUUID string) {
	c.Set(UserUUIDKey, userUUID)
}

// ParseJSON - Parse request JSON
func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.Wrap400Response(err, "Parse request parameter error")
	}
	return nil
}

// ResPage - Response paging data
func ResPage(c *gin.Context, v interface{}, pr *schema.PaginationResult) {
	list := schema.HTTPList{
		List: v,
		Pagination: &schema.HTTPPagination{
			Current:  GetPageIndex(c),
			PageSize: GetPageSize(c),
		},
	}
	if pr != nil {
		list.Pagination.Total = pr.Total
	}

	ResSuccess(c, list)
}

// ResList - Response list data
func ResList(c *gin.Context, v interface{}) {
	ResSuccess(c, schema.HTTPList{List: v})
}

// ResOK - Respond OK
func ResOK(c *gin.Context) {
	ResSuccess(c, schema.HTTPStatus{Status: schema.OKStatusText.String()})
}

// ResSuccess - Successful response
func ResSuccess(c *gin.Context, v interface{}) {
	ResJSON(c, http.StatusOK, v)
}

// ResJSON - Respond to JSON data
func ResJSON(c *gin.Context, status int, v interface{}) {
	buf, err := util.JSONMarshal(v)
	if err != nil {
		panic(err)
	}
	c.Set(ResBodyKey, buf)
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}

// ResError - Response error
func ResError(c *gin.Context, err error, status ...int) {
	var res *errors.ResponseError
	if err != nil {
		if e, ok := err.(*errors.ResponseError); ok {
			res = e
		} else {
			res = errors.UnWrapResponse(errors.Wrap500Response(err))
		}
	} else {
		res = errors.UnWrapResponse(errors.ErrInternalServer)
	}

	if len(status) > 0 {
		res.StatusCode = status[0]
	}

	if err := res.ERR; err != nil {
		if status := res.StatusCode; status >= 400 && status < 500 {
			logger.StartSpan(NewContext(c)).Warnf(err.Error())
		} else if status >= 500 {
			span := logger.StartSpan(NewContext(c))
			span = span.WithField("stack", fmt.Sprintf("%+v", err))
			span.Errorf(err.Error())
		}
	}

	eitem := schema.HTTPErrorItem{
		Code:    res.Code,
		Message: res.Message,
	}
	ResJSON(c, res.StatusCode, schema.HTTPError{Error: eitem})
}

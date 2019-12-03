package middleware

import (
	"fmt"
	"strings"

	"github.com/MayCMF/core/src/common/errors"
	"github.com/MayCMF/core/src/common/ginplus"
	"github.com/gin-gonic/gin"
)

// NoMethodHandler - The handler for the request method was not found
func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ginplus.ResError(c, errors.ErrMethodNotAllow)
	}
}

// NoRouteHandler - The handler for request routing was not found
func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ginplus.ResError(c, errors.ErrNotFound)
	}
}

// SkipperFunc - Define middleware skip function
type SkipperFunc func(*gin.Context) bool

// AllowPathPrefixSkipper - Check if the request path contains the specified prefix, skip if it is included
func AllowPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

// AllowPathPrefixNoSkipper - Check if the request path contains the specified prefix, if not, skip it
func AllowPathPrefixNoSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return false
			}
		}
		return true
	}
}

// AllowMethodAndPathPrefixSkipper Check if the request method and path contain the specified prefix, skip if not included
func AllowMethodAndPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := JoinRouter(c.Request.Method, c.Request.URL.Path)
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

// JoinRouter - Splicing route
func JoinRouter(method, path string) string {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(method), path)
}

// SkipHandler - Unified processing of skip functions
func SkipHandler(c *gin.Context, skippers ...SkipperFunc) bool {
	for _, skipper := range skippers {
		if skipper(c) {
			return true
		}
	}
	return false
}

// EmptyMiddleware - Middleware that does not perform business processing
func EmptyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

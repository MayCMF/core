package middleware

import (
	"github.com/MayCMF/core/src/common/ginplus"
	"github.com/MayCMF/core/src/common/util"
	"github.com/gin-gonic/gin"
)

// TraceMiddleware - Tracking ID middleware
func TraceMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		// Get the request ID first from the request header, if not, use the UUID
		traceID := c.GetHeader("X-Request-Id")
		if traceID == "" {
			traceID = util.NewTraceID()
		}
		c.Set(ginplus.TraceIDKey, traceID)
		c.Next()
	}
}

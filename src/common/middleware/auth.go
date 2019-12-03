package middleware

import (
	"github.com/MayCMF/core/src/common/auth"
	"github.com/MayCMF/core/src/common/config"
	"github.com/MayCMF/core/src/common/errors"
	"github.com/MayCMF/core/src/common/ginplus"
	"github.com/gin-gonic/gin"
)

// UserAuthMiddleware - User authorization middleware
func UserAuthMiddleware(a auth.Auther, skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if t := ginplus.GetToken(c); t != "" {
			id, err := a.ParseUserUUID(ginplus.NewContext(c), t)
			if err != nil {
				if err == auth.ErrInvalidToken {
					ginplus.ResError(c, errors.ErrInvalidToken)
					return
				}

				e := errors.UnWrapResponse(errors.ErrInvalidToken)
				ginplus.ResError(c, errors.WrapResponse(err, e.Code, e.Message, e.StatusCode))
				return
			} else if id != "" {
				c.Set(ginplus.UserUUIDKey, id)
				c.Next()
				return
			}
		}

		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		cfg := config.Global()
		if cfg.IsDebugMode() {
			c.Set(ginplus.UserUUIDKey, cfg.Root.UserName)
			c.Next()
			return
		}
		ginplus.ResError(c, errors.ErrInvalidToken)
	}
}

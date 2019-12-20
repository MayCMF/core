package api

import (
	"github.com/MayCMF/core/src/account/routers/api/controllers"
	"github.com/MayCMF/core/src/common/auth"
	"github.com/MayCMF/core/src/common/middleware"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

// RegisterRouter - Registration /api routing
func RegisterRouter(app *gin.Engine, container *dig.Container) error {
	err := controllers.Inject(container)
	if err != nil {
		return err
	}

	return container.Invoke(func(
		a auth.Auther,
		e *casbin.SyncedEnforcer,
		cLogin *controllers.Login,
		cPermission *controllers.Permission,
		cRole *controllers.Role,
		cUser *controllers.User,
	) error {

		g := app.Group("/api")

		// User identity authorization
		g.Use(middleware.UserAuthMiddleware(a,
			middleware.AllowPathPrefixSkipper("/api/v1/pub/login"),
		))

		// Casbin permission check middleware
		g.Use(middleware.CasbinMiddleware(e,
			middleware.AllowPathPrefixSkipper("/api/v1/pub"),
		))

		// Request frequency limit middleware
		g.Use(middleware.RateLimiterMiddleware())

		v1 := g.Group("/v1")
		{
			pub := v1.Group("/pub")
			{
				// [PUBLIC]/api/v1/pub/login
				gLogin := pub.Group("login")
				{
					gLogin.GET("captchaid", cLogin.GetCaptcha)
					gLogin.GET("captcha", cLogin.ResCaptcha)
					gLogin.POST("", cLogin.Login)
					gLogin.POST("exit", cLogin.Logout)
				}

				// [PUBLIC]/api/v1/pub/refresh-token
				pub.POST("/refresh-token", cLogin.RefreshToken)

				// [PUBLIC]/api/v1/pub/current
				gCurrent := pub.Group("current")
				{
					gCurrent.PUT("password", cLogin.UpdatePassword)
					gCurrent.GET("user", cLogin.GetUserInfo)
					gCurrent.GET("permission.tree", cLogin.QueryUserPermissionTree)
				}

			}

			// [REGISTERED]/api/v1/permissions
			gPermission := v1.Group("permissions")
			{
				gPermission.GET("", cPermission.Query)
				gPermission.GET(":id", cPermission.Get)
				gPermission.POST("", cPermission.Create)
				gPermission.PUT(":id", cPermission.Update)
				gPermission.DELETE(":id", cPermission.Delete)
			}
			v1.GET("/permissions.tree", cPermission.QueryTree)

			// [REGISTERED]/api/v1/roles
			gRole := v1.Group("roles")
			{
				gRole.GET("", cRole.Query)
				gRole.GET(":id", cRole.Get)
				gRole.POST("", cRole.Create)
				gRole.PUT(":id", cRole.Update)
				gRole.DELETE(":id", cRole.Delete)
			}
			v1.GET("/roles.select", cRole.QuerySelect)

			// [REGISTERED]/api/v1/users
			gUser := v1.Group("users")
			{
				gUser.GET("", cUser.Query)
				gUser.GET(":id", cUser.Get)
				gUser.POST("", cUser.Create)
				gUser.PUT(":id", cUser.Update)
				gUser.DELETE(":id", cUser.Delete)
				gUser.PATCH(":id/enable", cUser.Enable)
				gUser.PATCH(":id/disable", cUser.Disable)
			}
		}

		return nil
	})
}

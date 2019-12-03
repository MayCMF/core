package api

import (
	"github.com/MayCMF/core/src/common/middleware"
	"github.com/MayCMF/src/filemanager/routers/api/controllers"
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
		cDemo *controllers.Demo,
	) error {

		g := app.Group("/api")

		// Request frequency limit middleware
		g.Use(middleware.RateLimiterMiddleware())

		v1 := g.Group("/v1")
		{

			// [REGISTERED]/api/v1/example
			gDemo := v1.Group("example")
			{
				gDemo.GET("", cDemo.Query)
				gDemo.GET(":id", cDemo.Get)
				gDemo.POST("", cDemo.Create)
				gDemo.PUT(":id", cDemo.Update)
				gDemo.DELETE(":id", cDemo.Delete)
				gDemo.PATCH(":id/enable", cDemo.Enable)
				gDemo.PATCH(":id/disable", cDemo.Disable)
			}
		}

		return nil
	})
}

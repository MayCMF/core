package api

import (
	"github.com/MayCMF/core/src/common/middleware"
	"github.com/MayCMF/core/src/filemanager/routers/api/controllers"
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
		cFile *controllers.File,
	) error {

		g := app.Group("/api")

		// Request frequency limit middleware
		g.Use(middleware.RateLimiterMiddleware())

		v1 := g.Group("/v1")
		{

			// [REGISTERED]/api/v1/file
			gFile := v1.Group("file")
			{
				gFile.GET("", cFile.Query)
				gFile.GET(":id", cFile.Get)
				gFile.POST("", cFile.Create)
				gFile.POST("/upload", cFile.Upload)
				gFile.PUT(":id", cFile.Update)
				gFile.DELETE(":id", cFile.Delete)
			}
		}

		return nil
	})
}

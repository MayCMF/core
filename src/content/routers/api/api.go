package api

import (
	"github.com/MayCMF/core/src/common/middleware"
	"github.com/MayCMF/core/src/primitives/routers/api/controllers"
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
		cPrimitive *controllers.Primitive,
		cNode *controllers.Node,
	) error {

		g := app.Group("/api")

		// Request frequency limit middleware
		g.Use(middleware.RateLimiterMiddleware())

		v1 := g.Group("/v1")
		{

			// [REGISTERED]/api/v1/primitive
			gPrimitive := v1.Group("primitive")
			{
				gPrimitive.GET("", cPrimitive.Query)
				gPrimitive.GET(":id", cPrimitive.Get)
				gPrimitive.POST("", cPrimitive.Create)
				gPrimitive.PUT(":id", cPrimitive.Update)
				gPrimitive.DELETE(":id", cPrimitive.Delete)
			}

			// [REGISTERED]/api/v1/node
			gNode := v1.Group("node")
			{
				gNode.GET("", cNode.Query)
				gNode.GET(":id", cNode.Get)
				gNode.POST("", cNode.Create)
				gNode.PUT(":id", cNode.Update)
				gNode.DELETE(":id", cNode.Delete)
				// gNode.PATCH(":id/publish", cNode.Publish)
				// gNode.PATCH(":id/unpublish", cNode.Unpublish)
			}
		}

		return nil
	})
}

package api

import (
	"github.com/MayCMF/core/src/common/middleware"
	"github.com/MayCMF/core/src/i18n/routers/api/controllers"
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
		cLanguage *controllers.Language,
		cCountry *controllers.Country,
	) error {

		g := app.Group("/api")

		// Request frequency limit middleware
		g.Use(middleware.RateLimiterMiddleware())

		v1 := g.Group("/v1")
		{

			// [REGISTERED]/api/v1/language
			gLanguage := v1.Group("language")
			{
				gLanguage.GET("", cLanguage.Query)
				gLanguage.GET(":code", cLanguage.Get)
				gLanguage.POST("", cLanguage.Create)
				gLanguage.PUT(":code", cLanguage.Update)
				gLanguage.DELETE(":code", cLanguage.Delete)
				gLanguage.PATCH(":code/enable", cLanguage.Enable)
				gLanguage.PATCH(":code/disable", cLanguage.Disable)
			}

			// [REGISTERED]/api/v1/countries
			gCountry := v1.Group("countries")
			{
				gCountry.GET("", cCountry.Query)
				gCountry.GET(":code", cCountry.Get)
				gCountry.POST("", cCountry.Create)
				gCountry.PUT(":code", cCountry.Update)
				gCountry.DELETE(":code", cCountry.Delete)
			}
		}

		return nil
	})
}

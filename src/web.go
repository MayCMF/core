package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	accountApi "github.com/MayCMF/core/src/account/routers/api"
	i18nApi "github.com/MayCMF/core/src/i18n/routers/api"
	primitivesApi "github.com/MayCMF/core/src/primitives/routers/api"

	"github.com/MayCMF/core/src/common/config"
	"github.com/MayCMF/core/src/common/logger"
	"github.com/MayCMF/core/src/common/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

// InitWeb - Initialize the web engine
func InitWeb(container *dig.Container) *gin.Engine {
	cfg := config.Global()
	gin.SetMode(cfg.RunMode)

	app := gin.New()
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	apiPrefixes := []string{"/api/"}

	// Tracking ID
	app.Use(middleware.TraceMiddleware(middleware.AllowPathPrefixNoSkipper(apiPrefixes...)))

	// Access log
	app.Use(middleware.LoggerMiddleware(middleware.AllowPathPrefixNoSkipper(apiPrefixes...)))

	// Crash recovery
	app.Use(middleware.RecoveryMiddleware())

	// Cross-domain request
	if cfg.CORS.Enable {
		app.Use(middleware.CORSMiddleware())
	}

	// Registration Account /api routing
	err := accountApi.RegisterRouter(app, container)
	handleError(err)

	// Registration i18n (Languages&Countries) /api routing
	i18nApi.RegisterRouter(app, container)
	// Registration Primitives /api routing
	primitivesApi.RegisterRouter(app, container)

	// Swagger document
	if dir := cfg.Swagger; dir != "" {
		app.Static("/swagger", dir)
	}

	// Static site
	if dir := cfg.WWW; dir != "" {
		app.Use(middleware.WWWMiddleware(dir))
	}

	return app
}

// InitHTTPServer - Initialize HTTP service
func InitHTTPServer(ctx context.Context, container *dig.Container) func() {
	cfg := config.Global().HTTP
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      InitWeb(container),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		logger.Printf(ctx, "HTTP service starts and Listen address is：[%s]", addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Errorf(ctx, err.Error())
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.ShutdownTimeout))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Errorf(ctx, err.Error())
		}
	}
}

package app

import (
	"context"
	"fmt"
	"os"

	"github.com/MayCMF/core/src/account"
	"github.com/MayCMF/core/src/filemanager"
	"github.com/MayCMF/core/src/i18n"
	"github.com/MayCMF/core/src/primitives"

	"github.com/MayCMF/core/src/common/auth"
	"github.com/MayCMF/core/src/common/boot"
	"github.com/MayCMF/core/src/common/config"
	"github.com/MayCMF/core/src/common/logger"

	"github.com/MayCMF/core/src/transaction"
	"go.uber.org/dig"
)

type options struct {
	ConfigFile     string
	ModelFile      string
	WWWDir         string
	SwaggerDir     string
	PermissionFile string
	Version        string
}

// Option - Defining configuration items
type Option func(*options)

// SetConfigFile - Setting profile
func SetConfigFile(s string) Option {
	return func(o *options) {
		o.ConfigFile = s
	}
}

// SetModelFile - Set the casbin model configuration file
func SetModelFile(s string) Option {
	return func(o *options) {
		o.ModelFile = s
	}
}

// SetWWWDir - Set static site directory
func SetWWWDir(s string) Option {
	return func(o *options) {
		o.WWWDir = s
	}
}

// SetSwaggerDir - Set the swagger directory
func SetSwaggerDir(s string) Option {
	return func(o *options) {
		o.SwaggerDir = s
	}
}

// SetPermissionFile - Setting Permission data file
func SetPermissionFile(s string) Option {
	return func(o *options) {
		o.PermissionFile = s
	}
}

// SetVersion - Set the version number
func SetVersion(s string) Option {
	return func(o *options) {
		o.Version = s
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// Init - Application initialization
func Init(ctx context.Context, opts ...Option) func() {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	err := config.LoadGlobal(o.ConfigFile)
	handleError(err)

	cfg := config.Global()

	logger.Printf(ctx, "Service startup, running mode: %s，version number: %s，Process number: %d", cfg.RunMode, o.Version, os.Getpid())

	if v := o.ModelFile; v != "" {
		cfg.Casbin.Model = v
	}
	if v := o.WWWDir; v != "" {
		cfg.WWW = v
	}
	if v := o.SwaggerDir; v != "" {
		cfg.Swagger = v
	}
	if v := o.PermissionFile; v != "" {
		cfg.Permission.Data = v
	}

	loggerCall, err := boot.InitLogger()
	handleError(err)

	err = boot.InitMonitor()
	if err != nil {
		logger.Errorf(ctx, err.Error())
	}

	// Initialize the graphics verification code
	account.InitCaptcha()

	// Create a dependency injection container
	container, containerCall := BuildContainer()

	// Initialization Languages
	err = i18n.InitLanguages(ctx, container)
	handleError(err)
	fmt.Printf("IMPORT LANGUAGE DATA")

	// Initialization Permission
	err = account.InitPermission(ctx, container)
	handleError(err)

	// Initialize the HTTP service
	httpCall := InitHTTPServer(ctx, container)
	handleError(err)

	return func() {
		if httpCall != nil {
			httpCall()
		}
		if containerCall != nil {
			containerCall()
		}
		if loggerCall != nil {
			loggerCall()
		}
	}
}

// BuildContainer Create a dependency injection container
func BuildContainer() (*dig.Container, func()) {
	// Create a dependency injection container
	container := dig.New()

	// Injection authentication module
	auther, err := account.InitAuth()
	handleError(err)

	container.Provide(func() auth.Auther {
		return auther
	})

	// Inject casbin
	container.Provide(account.NewCasbinEnforcer)

	// ---------------------------------------------------
	// Injection memory module
	storeCall, err := InitStore(container)
	handleError(err)

	err = transaction.InjectControllers(container)
	handleError(err)

	err = account.InjectControllers(container)
	handleError(err)

	// Initialize casbin
	err = account.InitCasbinEnforcer(container)
	handleError(err)

	err = i18n.InjectControllers(container)
	handleError(err)

	err = primitives.InjectControllers(container)
	handleError(err)

	err = filemanager.InjectControllers(container)
	handleError(err)

	// ---------------------------------------------------
	return container, func() {
		if auther != nil {
			_ = auther.Release()
		}

		// Release resources
		account.ReleaseCasbinEnforcer(container)

		if storeCall != nil {
			storeCall()
		}
	}
}

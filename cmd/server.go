package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	app "github.com/MayCMF/core/src"
	"github.com/MayCMF/core/src/common/logger"
	"github.com/MayCMF/core/src/common/util"
)

// VERSION - version number,
// The version number can be specified by compiling: go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "0.1.0"

var (
	configFile     string
	modelFile      string
	wwwDir         string
	swaggerDir     string
	permissionFile string
)

func init() {
	flag.StringVar(&configFile, "c", "./configs/config.toml", "Configuration file(.json,.yaml,.toml)")
	flag.StringVar(&modelFile, "m", "./configs/model.conf", "Casbin's access control model(.conf)")
	flag.StringVar(&wwwDir, "www", "www", "Static site directory")
	flag.StringVar(&swaggerDir, "swagger", "docs/swagger", "Swagger directory")
	flag.StringVar(&permissionFile, "permission", "./configs/menu.json", "Permission data file(.json)")
}

func main() {
	flag.Parse()

	if configFile == "" {
		panic("Please use -c to specify the configuration file")
	}

	var state int32 = 1
	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Initialize log parameters
	logger.SetVersion(VERSION)
	logger.SetTraceIDFunc(util.NewTraceID)
	ctx := logger.NewTraceIDContext(context.Background(), util.NewTraceID())
	span := logger.StartSpanWithCall(ctx)

	call := app.Init(ctx,
		app.SetConfigFile(configFile),
		app.SetModelFile(modelFile),
		app.SetWWWDir(wwwDir),
		app.SetSwaggerDir(swaggerDir),
		app.SetPermissionFile(permissionFile),
		app.SetVersion(VERSION))

EXIT:
	for {
		sig := <-sc
		span().Printf("Get the signal[%s]", sig.String())

		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			atomic.StoreInt32(&state, 0)
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	if call != nil {
		call()
	}

	span().Printf("Service exit")
	time.Sleep(time.Second)
	os.Exit(int(atomic.LoadInt32(&state)))
}

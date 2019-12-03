package boot

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MayCMF/core/src/common/config"
	"github.com/MayCMF/core/src/common/logger"
	loggerhook "github.com/MayCMF/core/src/common/logger/hook"
	loggergormhook "github.com/MayCMF/core/src/common/logger/hook/gorm"
)

// InitLogger - Initialization log
func InitLogger() (func(), error) {
	c := config.Global().Log
	logger.SetLevel(c.Level)
	logger.SetFormatter(c.Format)

	// Set log output
	var file *os.File
	if c.Output != "" {
		switch c.Output {
		case "stdout":
			logger.SetOutput(os.Stdout)
		case "stderr":
			logger.SetOutput(os.Stderr)
		case "file":
			if name := c.OutputFile; name != "" {
				os.MkdirAll(filepath.Dir(name), 0777)

				f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					return nil, err
				}
				logger.SetOutput(f)
				file = f
			}
		}
	}

	var hook *loggerhook.Hook
	if c.EnableHook {
		switch c.Hook {
		case "gorm":
			hc := config.Global().LogGormHook

			var dsn string
			switch hc.DBType {
			case "mysql":
				dsn = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", config.Global().MySQL.User, config.Global().MySQL.Password, config.Global().MySQL.Host, config.Global().MySQL.Port, config.Global().MySQL.DBName)
			case "sqlite3":
				dsn = config.Global().Sqlite3.Path
			case "postgres":
				dsn = fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", config.Global().MySQL.User, config.Global().MySQL.Password, config.Global().MySQL.Host, config.Global().MySQL.DBName)
			default:
				return nil, errors.New("LOGER: Not supported database")
			}

			h := loggerhook.New(loggergormhook.New(&loggergormhook.Config{
				DBType:       hc.DBType,
				DSN:          dsn,
				MaxLifetime:  hc.MaxLifetime,
				MaxOpenConns: hc.MaxOpenConns,
				MaxIdleConns: hc.MaxIdleConns,
				TableName:    hc.Table,
			}),
				loggerhook.SetMaxWorkers(c.HookMaxThread),
				loggerhook.SetMaxQueues(c.HookMaxBuffer),
			)
			logger.AddHook(h)
			hook = h
		}
	}

	return func() {
		if file != nil {
			file.Close()
		}

		if hook != nil {
			hook.Flush()
		}
	}, nil
}

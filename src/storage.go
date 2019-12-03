package app

import (
	"errors"
	"fmt"
	"time"

	accountIject "github.com/MayCMF/core/src/account"
	filemanagerIject "github.com/MayCMF/core/src/filemanager"
	i18nIject "github.com/MayCMF/core/src/i18n"
	primitivesIject "github.com/MayCMF/core/src/primitives"
	"github.com/MayCMF/core/src/transaction"

	"github.com/MayCMF/core/src/common/config"
	"github.com/MayCMF/core/src/common/entity"
	"github.com/jinzhu/gorm"
	"go.uber.org/dig"

	// Gorm storage injection
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Config - Configuration parameter
type Config struct {
	Debug        bool
	DBType       string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
}

// SetTablePrefix - Set the table name prefix
func SetTablePrefix(prefix string) {
	entity.SetTablePrefix(prefix)
}

// InitStore - Initialize storage
func InitStore(container *dig.Container) (func(), error) {
	var storeCall func()
	cfg := config.Global()

	switch cfg.Store {
	case "gorm":
		db, err := initGorm()
		if err != nil {
			return nil, err
		}

		storeCall = func() {
			db.Close()
		}

		SetTablePrefix(cfg.Gorm.TablePrefix)

		if cfg.Gorm.EnableAutoMigrate {
			err = AutoMigrate(db)
			if err != nil {
				return nil, err
			}
		}

		// Inject DB
		container.Provide(func() *gorm.DB {
			return db
		})

		transaction.InjectStarage(container)
		accountIject.InjectStarage(container)
		i18nIject.InjectStarage(container)
		primitivesIject.InjectStarage(container)
		filemanagerIject.InjectStarage(container)

	default:
		return nil, errors.New("Unknown storage")
	}

	return storeCall, nil
}

// initGorm - Instantiate gorm storage
func initGorm() (db *gorm.DB, err error) {
	storConf := config.Global()

	fmt.Printf("config: %#v", storConf.Gorm)

	switch storConf.Gorm.DBType {
	case "mysql":
		db, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", storConf.MySQL.User, storConf.MySQL.Password, storConf.MySQL.Host, storConf.MySQL.Port, storConf.MySQL.DBName))
		db = db.Set("gorm:table_options", "CHARSET=utf8, ENGINE=InnoDB")
	case "sqlite3":
		db, err = gorm.Open("sqlite3", fmt.Sprintf("%v/%v", storConf.Sqlite3.Dir, storConf.Sqlite3.Name))
	case "postgres":
		db, err = gorm.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=%v", storConf.Postgres.User, storConf.Postgres.Password, storConf.Postgres.Host, storConf.Postgres.DBName, storConf.Postgres.SSLMode))
	default:
		return nil, errors.New("STORAGE: Not supported database")
	}

	if err != nil {
		return nil, err
	}

	if storConf.Gorm.Debug {
		db = db.Debug()
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(storConf.Gorm.MaxIdleConns)
	db.DB().SetMaxOpenConns(storConf.Gorm.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(storConf.Gorm.MaxLifetime) * time.Second)
	return db, nil
}

package entity

import (
	"context"
	"fmt"
	"time"

	"github.com/MayCMF/core/src/common/config"
	icontext "github.com/MayCMF/core/src/common/context"
	"github.com/MayCMF/core/src/common/util"
	"github.com/jinzhu/gorm"
)

// Table name prefix
var tablePrefix string

// SetTablePrefix - Set the table name prefix
func SetTablePrefix(prefix string) {
	tablePrefix = prefix
}

// GetTablePrefix - Get the table name prefix
func GetTablePrefix() string {
	return tablePrefix
}

// Model base model
type Model struct {
	ID        uint       `gorm:"column:id;primary_key;auto_increment;"`
	CreatedAt time.Time  `gorm:"column:created_at;"`
	UpdatedAt time.Time  `gorm:"column:updated_at;"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index;"`
}

// TableName table name
func (Model) TableName(name string) string {
	return fmt.Sprintf("%s%s", GetTablePrefix(), name)
}

func ToString(v interface{}) string {
	return util.JSONMarshalToString(v)
}

func getDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	trans, ok := icontext.FromTrans(ctx)
	if ok {
		db, ok := trans.(*gorm.DB)
		if ok {
			if icontext.FromTransLock(ctx) {
				if dbType := config.Global().Gorm.DBType; dbType == "mysql" ||
					dbType == "postgres" {
					db = db.Set("gorm:query_option", "FOR UPDATE")
				}
			}
			return db
		}
	}
	return defDB
}

func GetDBWithModel(ctx context.Context, defDB *gorm.DB, m interface{}) *gorm.DB {
	return getDB(ctx, defDB).Model(m)
}

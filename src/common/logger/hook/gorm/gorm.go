package gorm

import (
	"time"

	"github.com/MayCMF/core/src/common/logger"
	"github.com/MayCMF/core/src/common/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var tableName string

// Config - Configuration parameter
type Config struct {
	DBType       string
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	TableName    string
}

// New Create a gorm-based hook instance (requires a table name)
func New(c *Config) *Hook {
	tableName = c.TableName

	db, err := gorm.Open(c.DBType, c.DSN)
	if err != nil {
		panic(err)
	}

	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)

	db.AutoMigrate(new(LogItem))
	return &Hook{
		db: db,
	}
}

// Hook - Gorm log hook
type Hook struct {
	db *gorm.DB
}

// Exec - Execution log write
func (h *Hook) Exec(entry *logrus.Entry) error {
	item := &LogItem{
		Level:     entry.Level.String(),
		Message:   entry.Message,
		CreatedAt: entry.Time,
	}

	data := entry.Data
	if v, ok := data[logger.TraceIDKey]; ok {
		item.TraceID, _ = v.(string)
		delete(data, logger.TraceIDKey)
	}
	if v, ok := data[logger.UserUUIDKey]; ok {
		item.UserUUID, _ = v.(string)
		delete(data, logger.UserUUIDKey)
	}
	if v, ok := data[logger.SpanTitleKey]; ok {
		item.SpanTitle, _ = v.(string)
		delete(data, logger.SpanTitleKey)
	}
	if v, ok := data[logger.SpanFunctionKey]; ok {
		item.SpanFunction, _ = v.(string)
		delete(data, logger.SpanFunctionKey)
	}
	if v, ok := data[logger.VersionKey]; ok {
		item.Version, _ = v.(string)
		delete(data, logger.VersionKey)
	}

	if len(data) > 0 {
		item.Data = util.JSONMarshalToString(data)
	}

	result := h.db.Create(item)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

// Close - Close hook
func (h *Hook) Close() error {
	return h.db.Close()
}

// LogItem - Store log entries
type LogItem struct {
	ID           uint      `gorm:"column:id;primary_key;auto_increment;"` // id
	Level        string    `gorm:"column:level;size:20;index;"`           // Log level
	Message      string    `gorm:"column:message;size:1024;"`             // Message
	TraceID      string    `gorm:"column:trace_id;size:128;index;"`       // Tracking ID
	UserUUID     string    `gorm:"column:user_uuid;size:36;index;"`       // User ID
	SpanTitle    string    `gorm:"column:span_title;size:256;"`           // Tracking unit title
	SpanFunction string    `gorm:"column:span_function;size:256;"`        // Tracking unit function name
	Data         string    `gorm:"column:data;type:text;"`                // Log data(json)
	Version      string    `gorm:"column:version;index;size:32;"`         // Service version number
	CreatedAt    time.Time `gorm:"column:created_at;index"`               // Creation time
}

// TableName - Table Name
func (LogItem) TableName() string {
	return tableName
}

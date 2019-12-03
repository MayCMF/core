package entity

import (
	"context"
	"fmt"

	"github.com/MayCMF/core/src/common/entity"
	"github.com/MayCMF/core/src/i18n/schema"
	"github.com/jinzhu/gorm"
)

// GetLanguageDB - Get the language store
func GetLanguageDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(ctx, defDB, Language{})
}

// SchemaLanguage - Language object
type SchemaLanguage schema.Language

// ToLanguage - Convert to language entity
func (a SchemaLanguage) ToLanguage() *Language {
	item := &Language{
		Code:    &a.Code,
		Name:    &a.Name,
		Native:  &a.Native,
		Rtl:     a.Rtl,
		Default: a.Default,
		Active:  a.Active,
	}
	return item
}

// Language - Language entity
type Language struct {
	Code    *string `gorm:"column:code;size:50;unique_index;"` // Number
	Name    *string `gorm:"column:name;size:100;index;"`       // Name
	Native  *string `gorm:"column:native;size:100;index;"`     // Native
	Rtl     bool    `gorm:"column:rtl"`                        // RTL
	Default bool    `gorm:"column:default"`                    // Default language for project
	Active  bool    `gorm:"column:active"`                     // Status language can be added as content
}

func (a Language) String() string {
	return entity.ToString(a)
}

// TableName - Table Name
func (a Language) TableName() string {
	return fmt.Sprintf("%s%s", entity.GetTablePrefix(), "languages")
}

// ToSchemaLanguage - Convert to language object
func (a Language) ToSchemaLanguage() *schema.Language {
	item := &schema.Language{
		Code:    *a.Code,
		Name:    *a.Name,
		Native:  *a.Native,
		Rtl:     a.Rtl,
		Default: a.Default,
		Active:  a.Active,
	}
	return item
}

// Languages - Language list
type Languages []*Language

// ToSchemaLanguages - Convert to language object list
func (a Languages) ToSchemaLanguages() []*schema.Language {
	list := make([]*schema.Language, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaLanguage()
	}
	return list
}

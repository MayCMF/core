package entity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MayCMF/core/src/common/entity"
	"github.com/MayCMF/core/src/i18n/schema"
	"github.com/jinzhu/gorm"
)

// GetCountryDB - Get the country store
func GetCountryDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(ctx, db, Country{})
}

// SchemaCountry - Country object
type SchemaCountry schema.Country

// ToCountry - Convert to country entity
func (a SchemaCountry) ToCountry() *Country {
	item := &Country{
		Code:      &a.Code,
		Name:      &a.Name,
		Native:    &a.Native,
		Phone:     &a.Phone,
		Continent: &a.Continent,
		Capital:   &a.Capital,
		Currency:  a.Currency,
		Languages: a.Languages,
		Timezones: a.Timezones,
		LatLng:    a.LatLng,
		Emoji:     &a.Emoji,
		EmojiU:    &a.EmojiU,
	}
	return item
}

// Country - Country entity
type Country struct {
	Code      *string         `gorm:"column:code;size:50;index;"`      // Number
	Name      *string         `gorm:"column:name;size:100;index;"`     // Name
	Native    *string         `gorm:"column:native;size:100;index;"`   // Native language
	Phone     *string         `gorm:"column:phone;size:20;index;"`     // Phone number
	Continent *string         `gorm:"column:continent;size:10;index;"` // Continent where country located
	Capital   *string         `gorm:"column:capital;size:100;index;"`  // Country's Capital
	Currency  json.RawMessage `gorm:"type:jsonb;"`                     // Country's Currency
	Languages json.RawMessage `gorm:"column:languages;type:json"`
	Timezones json.RawMessage `gorm:"column:timezones;type:jsonb"`
	LatLng    json.RawMessage `gorm:"column:latlng;type:jsonb"`
	Emoji     *string         `gorm:"column:emoji"`
	EmojiU    *string         `gorm:"column:emojiU"`
}

func (a Country) String() string {
	return entity.ToString(a)
}

// TableName - Table Name
func (a Country) TableName() string {
	return fmt.Sprintf("%s%s", entity.GetTablePrefix(), "countries")
}

// ToSchemaCountry - Convert to country object
func (a Country) ToSchemaCountry() *schema.Country {
	item := &schema.Country{
		Code:      *a.Code,
		Name:      *a.Name,
		Native:    *a.Native,
		Phone:     *a.Phone,
		Continent: *a.Continent,
		Capital:   *a.Capital,
		Currency:  a.Currency,
		Languages: a.Languages,
		Timezones: a.Timezones,
		LatLng:    a.LatLng,
		Emoji:     *a.Emoji,
		EmojiU:    *a.EmojiU,
	}
	return item
}

// Countries - Countries list
type Countries []*Country

// ToSchemaCountries - Convert to country object list
func (a Countries) ToSchemaCountries() []*schema.Country {
	list := make([]*schema.Country, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaCountry()
	}
	return list
}

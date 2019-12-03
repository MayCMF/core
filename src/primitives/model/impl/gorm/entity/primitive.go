package entity

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	account "github.com/MayCMF/core/src/account/model/impl/gorm/entity"
	"github.com/MayCMF/core/src/common/entity"
	i18n "github.com/MayCMF/core/src/i18n/model/impl/gorm/entity"
	"github.com/MayCMF/core/src/primitives/schema"
	"github.com/jinzhu/gorm"
)

// GetPrimitiveDB - Get the Primitive store
func GetPrimitiveDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(ctx, defDB, Primitive{})
}

// GetPrimitiveDB - Get the Primitive store
func GetPrimitiveBodyDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(ctx, defDB, PrimitiveBody{})
}

// SchemaPrimitive - Primitive object
type SchemaPrimitive schema.Primitive

// ToPrimitive - Convert to Primitive entity
func (a SchemaPrimitive) ToPrimitive() *Primitive {
	item := &Primitive{
		UUID:       a.UUID,
		UID:        a.UID,
		Slug:       a.Slug,
		Parent:     a.Parent,
		ParentPath: a.ParentPath,
		Options:    a.Options,
	}
	return item
}

// ToPermissionActions - Convert to permission action list
func (a SchemaPrimitive) ToVariations() []*PrimitiveBody {
	list := make([]*PrimitiveBody, len(a.Variations))
	for i, item := range a.Variations {
		list[i] = SchemaPrimitiveBody(*item).ToPrimitiveBody(a.Slug)
	}
	return list
}

// Primitive - Primitive entity
type Primitive struct {
	entity.Model
	UUID       string          `gorm:"column:uuid;size:36;index;"`               // UUID
	User       account.User    `gorm:"foreignkey:UID;association_foreignkey:ID"` // Creator User ID
	UID        int             `gorm:"column:uid;"`                              // Creator User ID
	Slug       string          `gorm:"column:slug;size:100;unique_index;"`       // Slug short machine name
	Parent     string          `gorm:"column:parent;size:100;unique_index;"`     // Slug short machine name
	ParentPath string          `gorm:"column:parent_path"`                       // Parent path
	Options    json.RawMessage `gorm:"column:options;type:jsonb;"`               // Options in Jeson Format
	Variations Variations
}

func (a Primitive) String() string {
	return entity.ToString(a)
}

// TableName - Table Name
func (a Primitive) TableName() string {
	return a.Model.TableName("primitives")
}

// ToSchemaPrimitive - Convert to Primitive object
func (a Primitive) ToSchemaPrimitive() *schema.Primitive {
	item := &schema.Primitive{
		ID:         a.ID,
		UUID:       a.UUID,
		UID:        a.UID,
		Slug:       a.Slug,
		Parent:     a.Parent,
		ParentPath: a.ParentPath,
		Options:    a.Options,
		CreatedAt:  a.CreatedAt,
	}
	return item
}

// Primitives - Primitive list
type Primitives []*Primitive

// ToSchemaPrimitives - Convert to Primitive object list
func (a Primitives) ToSchemaPrimitives() []*schema.Primitive {
	list := make([]*schema.Primitive, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaPrimitive()
	}
	return list
}

// SchemaPrimitiveBody PrimitiveBody action object
type SchemaPrimitiveBody schema.PrimitiveBody

// ToPrimitiveBody - Convert to Primitive Body entity
func (a SchemaPrimitiveBody) ToPrimitiveBody(Slug string) *PrimitiveBody {
	return &PrimitiveBody{
		Slug:      a.Slug,
		UID:       a.UID,
		Lang:      a.Lang,
		Title:     &a.Title,
		Body:      &a.Body,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

// Primitive Body - Primitive Body object
type PrimitiveBody struct {
	Primitive Primitive     `gorm:"foreignkey:Slug;association_foreignkey:Slug"` // Primitive ID Slug
	Slug      string        `gorm:"column:slug"`                                 // Primitive ID Slug
	User      account.User  `gorm:"foreignkey:UID;association_foreignkey:ID"`    // Creator User ID
	UID       int           `gorm:"column:uid;"`                                 // Creator User ID
	Language  i18n.Language `gorm:"foreignkey:Lang;association_foreignkey:Code"` // Language Code Identifieru se Code as foreign key
	Lang      string        `gorm:"column:language"`                             // Language Code Identifieru se Code as foreign key
	Title     *string       `gorm:"column:title" binding:"required"`             // Primitive Title
	Body      *string       `gorm:"column:body"`                                 // Primitive Body
	CreatedAt time.Time     `gorm:"column:created_at"`                           // Creation time
	UpdatedAt time.Time     `gorm:"column:updated_at"`                           // Updated time
}

// TableName - Table Name
func (a PrimitiveBody) TableName() string {
	return fmt.Sprintf("%s%s", entity.GetTablePrefix(), "primitives_body")
}

// ToSchemaPrimitiveBody - Convert to Primitive Body object
func (a PrimitiveBody) ToSchemaPrimitiveBody() *schema.PrimitiveBody {
	item := &schema.PrimitiveBody{
		Slug:      a.Slug,
		UID:       a.UID,
		Lang:      a.Lang,
		Title:     *a.Title,
		Body:      *a.Body,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
	return item
}

// Variations - PrimitiveBody lassociated entity ist
type Variations []*PrimitiveBody

// GetByPrimitiveID - Get Primitive Body list based on Primitive ID
func (a Variations) GetByPrimitiveID(Slug string) []*schema.PrimitiveBody {
	var list []*schema.PrimitiveBody
	for _, item := range a {
		if item.Slug == Slug {
			list = append(list, item.ToSchemaPrimitiveBody())
		}
	}
	return list
}

// ToSchemaVariations - Convert to Primitive Body variations action list
func (a Variations) ToSchemaVariations() []*schema.PrimitiveBody {
	list := make([]*schema.PrimitiveBody, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaPrimitiveBody()
	}
	return list
}

// ToMap - Convert to key-value mapping
func (a Variations) ToMap() map[string]*PrimitiveBody {
	m := make(map[string]*PrimitiveBody)
	for _, item := range a {
		m[item.Lang] = item
	}
	return m
}

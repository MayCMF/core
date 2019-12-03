package entity

import (
	"context"

	"github.com/MayCMF/core/src/common/entity"
	"github.com/MayCMF/src/filemanager/schema"
	"github.com/jinzhu/gorm"
)

// GetFileDB - Get the File store
func GetFileDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(ctx, defDB, File{})
}

// SchemaFile - File object
type SchemaFile schema.File

// ToFile - Convert to File entity
func (a SchemaFile) ToFile() *File {
	item := &File{
		UUID:      a.UUID,
		UID:       &a.UID,
		Filename:  &a.Filename,
		Uri:       &a.Uri,
		Filemime:  &a.Filemime,
		Filesize:  &a.Filesize,
	}
	return item
}

// File - File entity
type File struct {
	entity.Model
	UUID      string  `gorm:"column:uuid;size:36;index;"`       // UUID code
	UID       *uint `gorm:"column:uid;size:50;index;"`          // User ID
	Filename  *string `gorm:"column:filename;size:100;index;"`  // File Name
	Uri       *string `gorm:"column:uri;size:200;"`             // File URI
	Filemime  *int    `gorm:"column:filemime;index;"`           // Filemime (image/jpeg, application/msword etc)
	Filesize  *string `gorm:"column:filesize;size:100;"`        // Filesize in bytes
}

func (a File) String() string {
	return entity.ToString(a)
}

// TableName - Table Name
func (a File) TableName() string {
	return a.Model.TableName("filemanager")
}

// ToSchemaFile - Convert to File object
func (a File) ToSchemaFile() *schema.File {
	item := &schema.File{
		UUID:      a.UUID,
		UID:       *a.UID,
		Filename:  *a.Filename,
		Uri:       *a.Uri,
		Filemime:  *a.Filemime,
		Filesize:  *a.Filesize,
		CreatedAt: a.CreatedAt,
	}
	return item
}

// Files - File list
type Files []*File

// ToSchemaFiles - Convert to File object list
func (a Files) ToSchemaFiles() []*schema.File {
	list := make([]*schema.File, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaFile()
	}
	return list
}

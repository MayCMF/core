package app

import (
	account "github.com/MayCMF/core/src/account/model/impl/gorm/entity"
	filemanager "github.com/MayCMF/core/src/filemanager/model/impl/gorm/entity"
	i18n "github.com/MayCMF/core/src/i18n/model/impl/gorm/entity"
	primitives "github.com/MayCMF/core/src/primitives/model/impl/gorm/entity"
	"github.com/jinzhu/gorm"
)

// AutoMigrate - Automatic mapping data table
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		new(account.User),
		new(account.UserRole),
		new(account.Role),
		new(account.RolePermission),
		new(account.Permission),
		new(account.PermissionAction),
		new(account.PermissionResource),
		new(i18n.Language),
		new(i18n.Country),
		new(primitives.Primitive),
		new(primitives.PrimitiveBody),
		new(primitives.Node),
		new(primitives.NodeBody),
		new(filemanager.File),
	).Error
}

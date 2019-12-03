package app

import (
	aentity "github.com/MayCMF/core/src/account/model/impl/gorm/entity"
	i18n "github.com/MayCMF/core/src/i18n/model/impl/gorm/entity"
	primitives "github.com/MayCMF/core/src/primitives/model/impl/gorm/entity"
	"github.com/jinzhu/gorm"
)

// AutoMigrate - Automatic mapping data table
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		new(aentity.User),
		new(aentity.UserRole),
		new(aentity.Role),
		new(aentity.RolePermission),
		new(aentity.Permission),
		new(aentity.PermissionAction),
		new(aentity.PermissionResource),
		new(i18n.Language),
		new(i18n.Country),
		new(primitives.Primitive),
		new(primitives.PrimitiveBody),
		new(primitives.Node),
		new(primitives.NodeBody),
	).Error
}

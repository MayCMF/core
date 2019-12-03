package entity

import (
	"context"

	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common/entity"
	"github.com/jinzhu/gorm"
)

// GetUserDB - Get user storage
func GetUserDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(ctx, defDB, User{})
}

// GetUserRoleDB - Get user role association storage
func GetUserRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return entity.GetDBWithModel(ctx, defDB, UserRole{})
}

// SchemaUser - User object
type SchemaUser schema.User

// ToUser - Convert to user entity
func (a SchemaUser) ToUser() *User {
	item := &User{
		UUID: a.UUID,
		UserName: &a.UserName,
		RealName: &a.RealName,
		Password: &a.Password,
		Status:   &a.Status,
		Creator:  &a.Creator,
		Email:    &a.Email,
		Phone:    &a.Phone,
	}
	return item
}

// ToUserRoles - Convert to user role association list
func (a SchemaUser) ToUserRoles() []*UserRole {
	list := make([]*UserRole, len(a.Roles))
	for i, item := range a.Roles {
		list[i] = &UserRole{
			UserUUID: a.UUID,
			RoleID:   item.RoleID,
		}
	}
	return list
}

// User - User entity
type User struct {
	entity.Model
	UUID string  `gorm:"column:record_id;size:36;index;"` // Record internal code
	UserName *string `gorm:"column:user_name;size:64;index;"` // UserName
	RealName *string `gorm:"column:real_name;size:64;index;"` // RealName
	Password *string `gorm:"column:password;size:40;"`        // Password (sha1 (md5 (plain text)) encryption)
	Email    *string `gorm:"column:email;not null;unique"`    // Email
	Phone    *string `gorm:"column:phone;size:20;index;"`     // Phone
	Status   *int    `gorm:"column:status;index;"`            // Status (1: Enable 2: Disable)
	Creator  *string `gorm:"column:creator;size:36;"`         // Creator
}

func (a User) String() string {
	return entity.ToString(a)
}

// TableName - Table Name
func (a User) TableName() string {
	return a.Model.TableName("user")
}

// ToSchemaUser - Convert to user object
func (a User) ToSchemaUser() *schema.User {
	item := &schema.User{
		UUID:  a.UUID,
		UserName:  *a.UserName,
		RealName:  *a.RealName,
		Password:  *a.Password,
		Status:    *a.Status,
		Creator:   *a.Creator,
		Email:     *a.Email,
		Phone:     *a.Phone,
		CreatedAt: a.CreatedAt,
	}
	return item
}

// Users - User entity list
type Users []*User

// ToSchemaUsers - Convert to user object list
func (a Users) ToSchemaUsers() []*schema.User {
	list := make([]*schema.User, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUser()
	}
	return list
}

// UserRole - User role association entity
type UserRole struct {
	entity.Model
	UserUUID string `gorm:"column:user_uuid;size:36;index;"` // User ID
	RoleID   string `gorm:"column:role_id;size:36;index;"`   // Role ID
}

// TableName - Table Name
func (a UserRole) TableName() string {
	return a.Model.TableName("user_role")
}

// ToSchemaUserRole - Convert to user role object
func (a UserRole) ToSchemaUserRole() *schema.UserRole {
	return &schema.UserRole{
		RoleID: a.RoleID,
	}
}

// UserRoles - User role association list
type UserRoles []*UserRole

// GetByUserUUID - Get the list of user role objects based on the user ID
func (a UserRoles) GetByUserUUID(userUUID string) []*schema.UserRole {
	var list []*schema.UserRole
	for _, item := range a {
		if item.UserUUID == userUUID {
			list = append(list, item.ToSchemaUserRole())
		}
	}
	return list
}

// ToSchemaUserRoles - Convert to a list of user role objects
func (a UserRoles) ToSchemaUserRoles() []*schema.UserRole {
	list := make([]*schema.UserRole, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUserRole()
	}
	return list
}

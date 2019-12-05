package schema

import (
	"github.com/MayCMF/core/src/common/schema"
	"time"
)

// User - User object
type User struct {
	ID        uint      `json:"id"`                                    // Record ID
	UUID      string    `json:"record_id"`                             // Record ID
	UserName  string    `json:"user_name" binding:"required"`          // UserName
	RealName  string    `json:"real_name" binding:"required"`          // RealName
	Password  string    `json:"password"`                              // Password
	Phone     string    `json:"phone"`                                 // Phone number
	Email     string    `json:"email"`                                 // Email
	Status    int       `json:"status" binding:"required,max=2,min=1"` // User Status (1: Enable 2: Disable)
	Creator   string    `json:"creator"`                               // Creator
	CreatedAt time.Time `json:"created_at"`                            // Creation time
	Roles     UserRoles `json:"roles" binding:"required,gt=0"`         // Role authorization
}

// CleanSecure - Clean up safety data
func (a *User) CleanSecure() *User {
	a.Password = ""
	return a
}

// UserRole - User role
type UserRole struct {
	RoleID string `json:"role_id" swaggo:"true, Role ID"`
}

// UserQueryParam - Query conditions
type UserQueryParam struct {
	UserName     string   // UserName
	LikeUserName string   // Username (fuzzy query)
	LikeRealName string   // Real name (fuzzy query)
	Status       int      // User Status (1: Enable 2: Disable)
	RoleIDs      []string // Role ID list
}

// UserQueryOptions - Query optional parameter items
type UserQueryOptions struct {
	PageParam    *schema.PaginationParam // Paging parameter
	IncludeRoles bool                    // Include role permissions
}

// UserQueryResult - User Query result
type UserQueryResult struct {
	Data       Users
	PageResult *schema.PaginationResult
}

// Users - User object list
type Users []*User

// ToRoleIDs - Get a list of role IDs
func (a Users) ToRoleIDs() []string {
	var roleIDs []string
	for _, item := range a {
		roleIDs = append(roleIDs, item.Roles.ToRoleIDs()...)
	}
	return roleIDs
}

// ToUserShows - Convert to user display list
func (a Users) ToUserShows(mroles map[string]*Role) UserShows {
	list := make(UserShows, len(a))

	for i, item := range a {
		showItem := &UserShow{
			UUID:      item.UUID,
			RealName:  item.RealName,
			UserName:  item.UserName,
			Email:     item.Email,
			Phone:     item.Phone,
			Status:    item.Status,
			CreatedAt: item.CreatedAt,
		}

		var roles Roles
		for _, roleID := range item.Roles.ToRoleIDs() {
			if v, ok := mroles[roleID]; ok {
				roles = append(roles, v)
			}
		}
		showItem.Roles = roles
		list[i] = showItem
	}

	return list
}

// UserRoles - User role list
type UserRoles []*UserRole

// ToRoleIDs - Convert to a list of role IDs
func (a UserRoles) ToRoleIDs() []string {
	list := make([]string, len(a))
	for i, item := range a {
		list[i] = item.RoleID
	}
	return list
}

// UserShow - User display item
type UserShow struct {
	UUID      string    `json:"record_id"`  // Record ID
	UserName  string    `json:"user_name"`  // UserName
	RealName  string    `json:"real_name"`  // RealName
	Phone     string    `json:"phone"`      // Phone
	Email     string    `json:"email"`      // Email
	Status    int       `json:"status"`     // User Status (1: Enable 2: Disable)
	CreatedAt time.Time `json:"created_at"` // Creation time
	Roles     []*Role   `json:"roles"`      // Roles List
}

// UserShows - User display item list
type UserShows []*UserShow

// UserShowQueryResult - User display item query result
type UserShowQueryResult struct {
	Data       UserShows
	PageResult *schema.PaginationResult
}

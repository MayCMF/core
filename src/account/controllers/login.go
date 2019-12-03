package controllers

import (
	"context"
	"net/http"

	"github.com/MayCMF/core/src/account/schema"
)

// ILogin - Login business logic interface
type ILogin interface {
	// Get graphic verification code information
	GetCaptcha(ctx context.Context, length int) (*schema.LoginCaptcha, error)
	// Generate and respond to a captcha
	ResCaptcha(ctx context.Context, w http.ResponseWriter, captchaID string, width, height int) error
	// Login authentication
	Verify(ctx context.Context, userName, password string) (*schema.User, error)
	// Generate token
	GenerateToken(ctx context.Context, userUUID string) (*schema.LoginTokenInfo, error)
	// Destroy token
	DestroyToken(ctx context.Context, tokenString string) error
	// Get user login information
	GetLoginInfo(ctx context.Context, userUUID string) (*schema.UserLoginInfo, error)
	// Query the user's permission Permission tree
	QueryUserPermissionTree(ctx context.Context, userUUID string) ([]*schema.PermissionTree, error)
	// Update user login password
	UpdatePassword(ctx context.Context, userUUID string, params schema.UpdatePasswordParam) error
}

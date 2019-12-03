package schema

// LoginParam - Login parameter
type LoginParam struct {
	UserName    string `json:"user_name" binding:"required"`    // UserName
	Password    string `json:"password" binding:"required"`     // Password (md5 encryption)
	CaptchaID   string `json:"captcha_id" binding:"required"`   // Verification code ID
	CaptchaCode string `json:"captcha_code" binding:"required"` // Verification code
}

// UserLoginInfo - User login information
type UserLoginInfo struct {
	UserName  string   `json:"user_name"`  // UserName
	RealName  string   `json:"real_name"`  // RealName
	RoleNames []string `json:"role_names"` // List of role names
}

// UpdatePasswordParam - Update password request parameters
type UpdatePasswordParam struct {
	OldPassword string `json:"old_password" binding:"required"` // Old password (md5 encryption)
	NewPassword string `json:"new_password" binding:"required"` // Old password (md5 encryption)
}

// LoginCaptcha - Login verification code
type LoginCaptcha struct {
	CaptchaID string `json:"captcha_id"` // Verification code ID
}

// LoginTokenInfo - Login token information
type LoginTokenInfo struct {
	AccessToken string `json:"access_token"` // Access token
	TokenType   string `json:"token_type"`   // Token type
	ExpiresAt   int64  `json:"expires_at"`   // Token expiration timestamp
}

package controllers

import (
	"github.com/LyricTian/captcha"
	"github.com/MayCMF/core/src/account/controllers"
	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common/config"
	"github.com/MayCMF/core/src/common/errors"
	"github.com/MayCMF/core/src/common/ginplus"
	"github.com/MayCMF/core/src/common/logger"
	"github.com/gin-gonic/gin"
)

// NewLogin - Create login manage controller
func NewLogin(bLogin controllers.ILogin) *Login {
	return &Login{
		LoginBll: bLogin,
	}
}

// Login - Manage Login
type Login struct {
	LoginBll controllers.ILogin
}

// GetCaptcha - Get verification code information
// @Tags Manage Login
// @Summary Get verification code information
// @Success 200 {object} schema.LoginCaptcha
// @Router /api/v1/pub/login/captchaid [get]
func (a *Login) GetCaptcha(c *gin.Context) {
	item, err := a.LoginBll.GetCaptcha(ginplus.NewContext(c), config.Global().Captcha.Length)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item)
}

// ResCaptcha - Response graphic verification code
// @Tags Manage Login
// @Summary Response graphic verification code
// @Param id query string true "Verification code ID"
// @Param reload query string false "Reload"
// @Produce image/png
// @Success 200 "Captcha"
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/pub/login/captcha [get]
func (a *Login) ResCaptcha(c *gin.Context) {
	captchaID := c.Query("id")
	if captchaID == "" {
		ginplus.ResError(c, errors.New400Response("Please provide a verification code ID"))
		return
	}

	if c.Query("reload") != "" {
		if !captcha.Reload(captchaID) {
			ginplus.ResError(c, errors.New400Response("No verification code ID found"))
			return
		}
	}

	cfg := config.Global().Captcha
	err := a.LoginBll.ResCaptcha(ginplus.NewContext(c), c.Writer, captchaID, cfg.Width, cfg.Height)
	if err != nil {
		ginplus.ResError(c, err)
	}
}

// Login - User login
// @Tags Manage Login
// @Summary User login
// @Param body body schema.LoginParam true "Request parameter"
// @Success 200 {object} schema.LoginTokenInfo
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/pub/login [post]
func (a *Login) Login(c *gin.Context) {
	var item schema.LoginParam
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	if !captcha.VerifyString(item.CaptchaID, item.CaptchaCode) {
		ginplus.ResError(c, errors.New400Response("Invalid verification code"))
		return
	}

	user, err := a.LoginBll.Verify(ginplus.NewContext(c), item.UserName, item.Password)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	userUUID := user.UUID
	// Put user ID into context
	ginplus.SetUserUUID(c, userUUID)

	tokenInfo, err := a.LoginBll.GenerateToken(ginplus.NewContext(c), userUUID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	logger.StartSpan(ginplus.NewContext(c), logger.SetSpanTitle("User login"), logger.SetSpanFuncName("Login")).Infof("Login system")
	ginplus.ResSuccess(c, tokenInfo)
}

// Logout - User logout
// @Tags Manage Login
// @Summary User logout
// @Success 200 {object} schema.HTTPStatus "{status: OK}"
// @Router /api/v1/pub/login/exit [post]
func (a *Login) Logout(c *gin.Context) {
	// Check if the user is logged in, and if so, destroy
	userUUID := ginplus.GetUserUUID(c)
	if userUUID != "" {
		ctx := ginplus.NewContext(c)
		err := a.LoginBll.DestroyToken(ctx, ginplus.GetToken(c))
		if err != nil {
			logger.Errorf(ctx, err.Error())
		}
		logger.StartSpan(ginplus.NewContext(c), logger.SetSpanTitle("Logout"), logger.SetSpanFuncName("Logout")).Infof("Logout system")
	}
	ginplus.ResOK(c)
}

// RefreshToken - Refresh token
// @Tags Manage Login
// @Summary Refresh token
// @Param Authorization header string false "Bearer User Token"
// @Success 200 {object} schema.LoginTokenInfo
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/pub/refresh-token [post]
func (a *Login) RefreshToken(c *gin.Context) {
	tokenInfo, err := a.LoginBll.GenerateToken(ginplus.NewContext(c), ginplus.GetUserUUID(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, tokenInfo)
}

// GetUserInfo - Get current user information
// @Tags Manage Login
// @Summary Get current user information
// @Param Authorization header string false "Bearer User Token"
// @Success 200 {object} schema.UserLoginInfo
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/pub/current/user [get]
func (a *Login) GetUserInfo(c *gin.Context) {
	info, err := a.LoginBll.GetLoginInfo(ginplus.NewContext(c), ginplus.GetUserUUID(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, info)
}

// QueryUserPermissionTree - Query current user permission tree
// @Tags Manage Login
// @Summary Query current user permission tree
// @Param Authorization header string false "Bearer User Token"
// @Success 200 {object} schema.Permission "Search result: {list:Permission tree}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/pub/current/permissiontree [get]
func (a *Login) QueryUserPermissionTree(c *gin.Context) {
	permissions, err := a.LoginBll.QueryUserPermissionTree(ginplus.NewContext(c), ginplus.GetUserUUID(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResList(c, permissions)
}

// UpdatePassword - Update personal password
// @Tags Manage Login
// @Summary Update personal password
// @Param Authorization header string false "Bearer User Token"
// @Param body body schema.UpdatePasswordParam true "Request parameter"
// @Success 200 {object} schema.HTTPStatus "{status:OK}"
// @Failure 400 {object} schema.HTTPError "{error:{code:0,message: Invalid request parameter}}"
// @Failure 401 {object} schema.HTTPError "{error:{code:0,message: Unauthorized}}"
// @Failure 500 {object} schema.HTTPError "{error:{code:0,message: Server Error}}"
// @Router /api/v1/pub/current/password [put]
func (a *Login) UpdatePassword(c *gin.Context) {
	var item schema.UpdatePasswordParam
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err := a.LoginBll.UpdatePassword(ginplus.NewContext(c), ginplus.GetUserUUID(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

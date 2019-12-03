package account

import (
	"github.com/MayCMF/core/src/common/config"
	"github.com/MayCMF/core/src/common/logger"
	"github.com/LyricTian/captcha"
	"github.com/LyricTian/captcha/store"
)

// InitCaptcha - Initialize the graphics verification code
func InitCaptcha() {
	cfg := config.Global().Captcha
	if cfg.Store == "redis" {
		rc := config.Global().Redis
		captcha.SetCustomStore(store.NewRedisStore(&store.RedisOptions{
			Addr:     rc.Addr,
			Password: rc.Password,
			DB:       cfg.RedisDB,
		}, captcha.Expiration, logger.StandardLogger(), cfg.RedisPrefix))
	}
}

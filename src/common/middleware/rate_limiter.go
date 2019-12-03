package middleware

import (
	"strconv"
	"time"

	"github.com/MayCMF/core/src/common/config"
	"github.com/MayCMF/core/src/common/errors"
	"github.com/MayCMF/core/src/common/ginplus"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
	"golang.org/x/time/rate"
)

// RateLimiterMiddleware - Request frequency limit middleware
func RateLimiterMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	cfg := config.Global().RateLimiter
	if !cfg.Enable {
		return EmptyMiddleware()
	}

	rc := config.Global().Redis
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"localhost": rc.Addr,
		},
		Password: rc.Password,
		DB:       cfg.RedisDB,
	})

	limiter := redis_rate.NewLimiter(ring)
	limiter.Fallback = rate.NewLimiter(rate.Inf, 0)

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		userUUID := ginplus.GetUserUUID(c)
		if userUUID == "" {
			c.Next()
			return
		}

		limit := cfg.Count
		rate, delay, allowed := limiter.AllowMinute(userUUID, limit)
		if !allowed {
			h := c.Writer.Header()
			h.Set("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
			h.Set("X-RateLimit-Remaining", strconv.FormatInt(limit-rate, 10))
			delaySec := int64(delay / time.Second)
			h.Set("X-RateLimit-Delay", strconv.FormatInt(delaySec, 10))
			ginplus.ResError(c, errors.ErrTooManyRequests)
			return
		}

		c.Next()
	}
}

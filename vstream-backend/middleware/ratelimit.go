package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func newLimiter(period time.Duration, limit int64) gin.HandlerFunc {
	rate := limiter.Rate{Period: period, Limit: limit}
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	return mgin.NewMiddleware(instance)
}

// LoginLimiter — 5 attempt / 10 menit per IP
func LoginLimiter() gin.HandlerFunc {
	return newLimiter(10*time.Minute, 5)
}

// PlaybackLogLimiter — 30 request / menit per IP
func PlaybackLogLimiter() gin.HandlerFunc {
	return newLimiter(time.Minute, 30)
}

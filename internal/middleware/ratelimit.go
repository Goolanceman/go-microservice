package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type KeyFunc func(c *gin.Context) string

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
	keyFunc  KeyFunc
}

func NewRateLimiter(r rate.Limit, b int, keyFunc KeyFunc) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
		keyFunc:  keyFunc,
	}
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = limiter
	}
	return limiter
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := rl.keyFunc(c)
		limiter := rl.getLimiter(key)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too Many Requests",
			})
			return
		}
		c.Next()
	}
}

/*
	// IP-based limiter: 10 requests per minute
	ipLimiter := middleware.NewRateLimiter(
		rate.Every(time.Minute/10),
		10,
		func(c *gin.Context) string { return c.ClientIP() },
	)

	// IP + Path-based limiter: 5 requests per minute per route
	ipPathLimiter := middleware.NewRateLimiter(
		rate.Every(time.Minute/5),
		5,
		func(c *gin.Context) string { return c.ClientIP() + "_" + c.FullPath() },
	)

	// Apply middlewares
	r.Use(ipLimiter.Middleware())
	r.Use(ipPathLimiter.Middleware())
*/

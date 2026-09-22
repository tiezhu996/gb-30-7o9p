package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gbadopt/gbadopt/internal/constants"
	"github.com/gbadopt/gbadopt/internal/dto"
)

type bucket struct {
	count   int
	resetAt time.Time
}

// RateLimiter is a per-IP token bucket.
type RateLimiter struct {
	mu     sync.Mutex
	limits map[string]*bucket
	reqs   int
	window time.Duration
}

// NewRateLimiter creates a limiter.
func NewRateLimiter(reqs int, window time.Duration) *RateLimiter {
	return &RateLimiter{limits: make(map[string]*bucket), reqs: reqs, window: window}
}

// Limit returns a middleware enforcing the rate limit per client IP.
func (r *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		r.mu.Lock()
		b, ok := r.limits[ip]
		if !ok || now.After(b.resetAt) {
			b = &bucket{count: 0, resetAt: now.Add(r.window)}
			r.limits[ip] = b
		}
		b.count++
		if b.count > r.reqs {
			r.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests,
				dto.Fail(constants.CodeRateLimited, constants.MsgRateLimited))
			return
		}
		r.mu.Unlock()
		c.Next()
	}
}

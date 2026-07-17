package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jordisetiawan/insurance-auth-service/internal/utils"
)

type visitor struct {
	count     int
	windowEnd time.Time
}

// RateLimiter returns a middleware that allows at most `limit` requests per
// `window` for each client IP. It's a simple in-memory fixed-window limiter
// meant to slow down brute-force attempts against auth endpoints
// (login/register). It is per-instance state and resets on restart, so it
// is not a substitute for a shared limiter behind a load balancer.
func RateLimiter(limit int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	visitors := make(map[string]*visitor)

	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		for now := range ticker.C {
			mu.Lock()
			for ip, v := range visitors {
				if now.After(v.windowEnd) {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists || now.After(v.windowEnd) {
			v = &visitor{windowEnd: now.Add(window)}
			visitors[ip] = v
		}
		v.count++
		count := v.count
		mu.Unlock()

		if count > limit {
			utils.ErrorResponse(c, http.StatusTooManyRequests, "Too many requests, please try again later", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

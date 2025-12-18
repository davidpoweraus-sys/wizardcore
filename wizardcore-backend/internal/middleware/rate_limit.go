package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"github.com/gin-gonic/gin"
)

// RateLimiter holds the rate limiter per IP.
type RateLimiter struct {
	ips map[string]*rate.Limiter
	mu  sync.RWMutex
	r   rate.Limit
	b   int
}

// NewRateLimiter creates a new rate limiter with the given requests per second and burst size.
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	return &RateLimiter{
		ips: make(map[string]*rate.Limiter),
		r:   rate.Limit(rps),
		b:   burst,
	}
}

// GetLimiter returns the rate limiter for the given IP, creating one if it doesn't exist.
func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.ips[ip] = limiter
		// Clean up old entries after a while (optional)
		go func(ip string) {
			time.Sleep(5 * time.Minute)
			rl.mu.Lock()
			delete(rl.ips, ip)
			rl.mu.Unlock()
		}(ip)
	}

	return limiter
}

// RateLimitMiddleware returns a Gin middleware that rate‑limits requests per IP.
func RateLimitMiddleware(rps float64, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(rps, burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			ip = "unknown"
		}

		if !limiter.GetLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
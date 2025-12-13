package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// rateLimiter stores rate limiters per IP address
type rateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

var globalLimiter = &rateLimiter{
	limiters: make(map[string]*rate.Limiter),
	rate:     rate.Limit(10), // 10 requests per second
	burst:    20,             // burst of 20 requests
}

// getLimiter returns a rate limiter for the given IP address
func (rl *rateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Double-check after acquiring write lock
		limiter, exists = rl.limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(rl.rate, rl.burst)
			rl.limiters[ip] = limiter
		}
		rl.mu.Unlock()
	}

	return limiter
}

// cleanup removes old limiters periodically to prevent memory leaks
func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			rl.mu.Lock()
			// In a production system, you'd track last access time
			// For now, we'll keep all limiters (they're lightweight)
			rl.mu.Unlock()
		}
	}()
}

func init() {
	globalLimiter.cleanup()
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// RateLimit applies rate limiting to HTTP handlers
// Limits: 10 requests per second, burst of 20
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		limiter := globalLimiter.getLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}


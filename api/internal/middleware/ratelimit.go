package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/teyk0o/loupi/api/internal/models"
)

const maxVisitors = 10000

// visitor tracks request counts for rate limiting.
type visitor struct {
	count    int
	lastSeen time.Time
}

// RateLimiter provides in-memory rate limiting per IP address.
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
	cancel   context.CancelFunc
}

// NewRateLimiter creates a new rate limiter with the given limit per window.
// Pass a cancellable context to stop the background cleanup goroutine.
func NewRateLimiter(ctx context.Context, limit int, window time.Duration) *RateLimiter {
	ctx, cancel := context.WithCancel(ctx)
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
		cancel:   cancel,
	}
	go rl.cleanup(ctx)
	return rl
}

// Stop cancels the background cleanup goroutine.
func (rl *RateLimiter) Stop() {
	rl.cancel()
}

// Middleware returns a Gin middleware that enforces the rate limit.
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		if !exists || time.Since(v.lastSeen) > rl.window {
			// Evict oldest entries if map is too large
			if len(rl.visitors) >= maxVisitors {
				rl.evictOldest()
			}
			rl.visitors[ip] = &visitor{count: 1, lastSeen: time.Now()}
			rl.mu.Unlock()
			c.Next()
			return
		}

		v.count++
		v.lastSeen = time.Now()

		if v.count > rl.limit {
			rl.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, models.ErrorResponse{
				Error:   "rate_limit_exceeded",
				Message: "Too many requests, please try again later",
			})
			return
		}

		rl.mu.Unlock()
		c.Next()
	}
}

// evictOldest removes the oldest visitor entry. Must be called with mu held.
func (rl *RateLimiter) evictOldest() {
	var oldestIP string
	var oldestTime time.Time
	first := true
	for ip, v := range rl.visitors {
		if first || v.lastSeen.Before(oldestTime) {
			oldestIP = ip
			oldestTime = v.lastSeen
			first = false
		}
	}
	if oldestIP != "" {
		delete(rl.visitors, oldestIP)
	}
}

// cleanup periodically removes stale visitor entries.
func (rl *RateLimiter) cleanup(ctx context.Context) {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rl.mu.Lock()
			for ip, v := range rl.visitors {
				if time.Since(v.lastSeen) > rl.window {
					delete(rl.visitors, ip)
				}
			}
			rl.mu.Unlock()
		}
	}
}

// LoginRateLimiter tracks failed login attempts per email and enforces account-level lockout.
type LoginRateLimiter struct {
	mu       sync.Mutex
	attempts map[string]*loginAttempt
	maxFails int
	lockout  time.Duration
	cancel   context.CancelFunc
}

type loginAttempt struct {
	failures int
	lockedAt time.Time
	lastFail time.Time
}

// NewLoginRateLimiter creates a rate limiter for login attempts.
func NewLoginRateLimiter(ctx context.Context, maxFails int, lockout time.Duration) *LoginRateLimiter {
	ctx, cancel := context.WithCancel(ctx)
	lr := &LoginRateLimiter{
		attempts: make(map[string]*loginAttempt),
		maxFails: maxFails,
		lockout:  lockout,
		cancel:   cancel,
	}
	go lr.cleanup(ctx)
	return lr
}

// Stop cancels the background cleanup goroutine.
func (lr *LoginRateLimiter) Stop() {
	lr.cancel()
}

// IsLocked returns true if the email is currently locked out.
func (lr *LoginRateLimiter) IsLocked(email string) bool {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	a, exists := lr.attempts[email]
	if !exists {
		return false
	}
	if a.failures >= lr.maxFails {
		if time.Since(a.lockedAt) < lr.lockout {
			return true
		}
		// Lockout expired, reset
		delete(lr.attempts, email)
	}
	return false
}

// RecordFailure records a failed login attempt for the given email.
func (lr *LoginRateLimiter) RecordFailure(email string) {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	a, exists := lr.attempts[email]
	if !exists {
		lr.attempts[email] = &loginAttempt{failures: 1, lastFail: time.Now()}
		return
	}
	a.failures++
	a.lastFail = time.Now()
	if a.failures >= lr.maxFails {
		a.lockedAt = time.Now()
	}
}

// RecordSuccess clears failed attempts for the given email.
func (lr *LoginRateLimiter) RecordSuccess(email string) {
	lr.mu.Lock()
	defer lr.mu.Unlock()
	delete(lr.attempts, email)
}

// cleanup periodically removes expired lockout entries.
func (lr *LoginRateLimiter) cleanup(ctx context.Context) {
	ticker := time.NewTicker(lr.lockout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			lr.mu.Lock()
			for email, a := range lr.attempts {
				if time.Since(a.lastFail) > lr.lockout {
					delete(lr.attempts, email)
				}
			}
			lr.mu.Unlock()
		}
	}
}

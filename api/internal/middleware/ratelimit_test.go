package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiter_AllowsUnderLimit(t *testing.T) {
	r := gin.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rl := NewRateLimiter(ctx, 5, 1*time.Minute)
	r.Use(rl.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}
}

func TestRateLimiter_BlocksOverLimit(t *testing.T) {
	r := gin.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rl := NewRateLimiter(ctx, 3, 1*time.Minute)
	r.Use(rl.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Exhaust the limit
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.2:1234"
		r.ServeHTTP(w, req)
	}

	// Next request should be blocked
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.2:1234"
	r.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w.Code)
	}
}

func TestRateLimiter_DifferentIPsIndependent(t *testing.T) {
	r := gin.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rl := NewRateLimiter(ctx, 2, 1*time.Minute)
	r.Use(rl.Middleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Exhaust limit for IP 1
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "10.0.0.1:1234"
		r.ServeHTTP(w, req)
	}

	// IP 2 should still work
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "10.0.0.2:1234"
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for different IP, got %d", w.Code)
	}
}

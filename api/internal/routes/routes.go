// Package routes configures all API route groups and endpoints.
package routes

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/teyk0o/loupi/api/internal/config"
	"github.com/teyk0o/loupi/api/internal/handlers"
	"github.com/teyk0o/loupi/api/internal/middleware"
	"github.com/teyk0o/loupi/api/internal/services"
)

// Setup configures all routes on the given Gin engine.
func Setup(r *gin.Engine, cfg *config.Config, authService *services.AuthService) {
	// Global middleware
	r.Use(middleware.CORS(cfg.AllowedOrigins))
	r.Use(middleware.SecurityHeaders())

	// Rate limiter for auth endpoints (10 requests per minute)
	authRateLimiter := middleware.NewRateLimiter(10, 1*time.Minute)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Health check
	r.GET("/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "loupi-api"})
	})

	// Auth routes (public)
	auth := r.Group("/v1/auth")
	auth.Use(authRateLimiter.Middleware())
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}

	// Auth routes (protected)
	authProtected := r.Group("/v1/auth")
	authProtected.Use(middleware.Auth(authService))
	{
		authProtected.GET("/me", authHandler.Me)
		authProtected.DELETE("/account", authHandler.DeleteAccount)
	}

	// Protected API routes (placeholder groups for future handlers)
	// api := r.Group("/v1")
	// api.Use(middleware.Auth(authService))
	// {
	//     // Meals, symptoms, wellness routes will be added here
	// }
}

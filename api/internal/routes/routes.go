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
func Setup(
	r *gin.Engine,
	cfg *config.Config,
	authService *services.AuthService,
	mealService *services.MealService,
	symptomService *services.SymptomService,
	wellnessService *services.WellnessService,
) {
	// Global middleware
	r.Use(middleware.CORS(cfg.AllowedOrigins))
	r.Use(middleware.SecurityHeaders())

	// Rate limiter for auth endpoints (10 requests per minute)
	authRateLimiter := middleware.NewRateLimiter(10, 1*time.Minute)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	mealHandler := handlers.NewMealHandler(mealService)
	symptomHandler := handlers.NewSymptomHandler(symptomService)
	wellnessHandler := handlers.NewWellnessHandler(wellnessService)

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

	// Protected API routes
	api := r.Group("/v1")
	api.Use(middleware.Auth(authService))
	{
		// Meals
		api.GET("/meals", mealHandler.ListByDate)
		api.POST("/meals", mealHandler.Create)
		api.GET("/meals/:id", mealHandler.GetByID)
		api.PUT("/meals/:id", mealHandler.Update)
		api.DELETE("/meals/:id", mealHandler.Delete)
		api.GET("/meals/:id/check-ins", mealHandler.GetCheckins)
		api.POST("/meals/:id/check-ins", mealHandler.CreateCheckin)

		// Symptom check-ins (update/delete by check-in ID)
		api.PUT("/check-ins/:id", symptomHandler.UpdateCheckin)
		api.DELETE("/check-ins/:id", symptomHandler.DeleteCheckin)

		// Standalone symptoms
		api.GET("/symptoms", symptomHandler.ListByDate)
		api.POST("/symptoms", symptomHandler.Create)
		api.PUT("/symptoms/:id", symptomHandler.Update)
		api.DELETE("/symptoms/:id", symptomHandler.Delete)

		// Wellness
		api.GET("/wellness", wellnessHandler.GetByDate)
		api.POST("/wellness", wellnessHandler.Upsert)
		api.PUT("/wellness/:id", wellnessHandler.Update)
	}
}

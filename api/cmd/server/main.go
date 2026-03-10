// Package main is the entry point for the Loupi API server.
package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/teyk0o/loupi/api/internal/config"
	"github.com/teyk0o/loupi/api/internal/database"
	"github.com/teyk0o/loupi/api/internal/routes"
	"github.com/teyk0o/loupi/api/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	if !cfg.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize services
	authService := services.NewAuthService(db, cfg)

	// Setup router
	r := gin.Default()
	routes.Setup(r, cfg, authService)

	// Start server
	addr := ":" + cfg.Port
	log.Printf("Loupi API starting on %s (%s)", addr, cfg.Env)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

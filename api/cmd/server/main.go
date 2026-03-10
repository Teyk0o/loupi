// Package main is the entry point for the Loupi API server.
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/teyk0o/loupi/api/internal/config"
	"github.com/teyk0o/loupi/api/internal/database"
	"github.com/teyk0o/loupi/api/internal/models"
	"github.com/teyk0o/loupi/api/internal/routes"
	"github.com/teyk0o/loupi/api/internal/services"
	"github.com/teyk0o/loupi/api/internal/utils"
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

	// Register custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := models.RegisterPasswordComplexity(v); err != nil {
			log.Fatalf("Failed to register custom validator: %v", err)
		}
	}

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Connect to Redis
	rdb, err := database.ConnectRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()

	// Run migrations
	if err := database.RunMigrations(db, "migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize encryption
	encryptor, err := utils.NewEncryptor(cfg.EncryptionKey)
	if err != nil {
		log.Fatalf("Failed to initialize encryption: %v", err)
	}

	// Initialize audit logger
	audit := utils.NewAuditLogger(db)

	// Initialize services
	authService := services.NewAuthService(db, rdb, cfg)
	mealService := services.NewMealService(db, encryptor)
	symptomService := services.NewSymptomService(db, encryptor)
	wellnessService := services.NewWellnessService(db, encryptor)
	customOptionService := services.NewCustomOptionService(db)

	// Setup router
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	routes.Setup(r, cfg, authService, mealService, symptomService, wellnessService, customOptionService, audit)

	// Start server
	addr := ":" + cfg.Port
	log.Printf("Loupi API starting on %s (%s)", addr, cfg.Env)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

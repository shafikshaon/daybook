package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"daybook-backend/config"
	"daybook-backend/database"
	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	// Load configuration
	logger.Infof(ctx, "Loading application configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Errorf(ctx, "Failed to load configuration: %v", err)
		os.Exit(1)
	}
	logger.Infof(ctx, "Configuration loaded successfully")

	// Set logger level based on server mode
	if cfg.Server.Mode == "release" {
		logger.SetLevel(logger.InfoLevel)
	} else {
		logger.SetLevel(logger.DebugLevel)
	}

	// Initialize database
	logger.Infof(ctx, "Initializing database connection...")
	if err := database.InitDatabase(cfg); err != nil {
		logger.Errorf(ctx, "Failed to initialize database: %v", err)
		os.Exit(1)
	}

	// Initialize Redis (optional)
	logger.Infof(ctx, "Initializing Redis connection...")
	if err := database.InitRedis(cfg); err != nil {
		logger.Warnf(ctx, "Redis initialization failed: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)
	logger.Infof(ctx, "Gin mode set to: %s", cfg.Server.Mode)

	// Create router
	router := gin.New()

	// Add recovery middleware
	router.Use(gin.Recovery())

	// Add tracing middleware (must be first to capture all requests)
	router.Use(middleware.TracingMiddleware())

	// Setup CORS
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		ExposeHeaders:    cfg.CORS.ExposeHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           time.Duration(cfg.CORS.MaxAge) * time.Hour,
	}
	router.Use(cors.New(corsConfig))
	logger.Infof(ctx, "CORS middleware configured")

	// Setup routes
	logger.Infof(ctx, "Setting up application routes...")
	routes.SetupRoutes(router)
	logger.Infof(ctx, "Routes configured successfully")

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		shutdownCtx := context.Background()
		logger.Infof(shutdownCtx, "Shutdown signal received, gracefully shutting down server...")

		// Close database connections
		if err := database.CloseDatabase(); err != nil {
			logger.Errorf(shutdownCtx, "Error closing database: %v", err)
		}

		if err := database.CloseRedis(); err != nil {
			logger.Errorf(shutdownCtx, "Error closing Redis: %v", err)
		}

		logger.Infof(shutdownCtx, "Server stopped gracefully")
		os.Exit(0)
	}()

	// Start server
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	logger.Infof(ctx, "===================================")
	logger.Infof(ctx, "Starting Daybook Backend Server")
	logger.Infof(ctx, "===================================")
	logger.Infof(ctx, "Server Port: %s", port)
	logger.Infof(ctx, "Server Mode: %s", cfg.Server.Mode)
	logger.Infof(ctx, "API Documentation: http://localhost:%s/health", port)
	logger.Infof(ctx, "===================================")

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		logger.Errorf(ctx, "Failed to start server: %v", err)
		os.Exit(1)
	}
}

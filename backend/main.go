package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"daybook-backend/config"
	"daybook-backend/container"
	"daybook-backend/database"
	"daybook-backend/logger"
	"daybook-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	// Create context with trace ID for startup logs
	ctx := logger.CreateContext("")

	// Load configuration
	logger.Infof(ctx, "Loading application configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Errorf(ctx, "Failed to load configuration: %v", err)
		os.Exit(1)
	}
	logger.Infof(ctx, "Configuration loaded successfully")

	// Initialize Datadog tracer if enabled
	if cfg.Datadog.Enabled {
		logger.Infof(ctx, "Initializing Datadog APM tracer...")
		tracer.Start(
			tracer.WithService(cfg.Datadog.ServiceName),
			tracer.WithEnv(cfg.Datadog.Environment),
			tracer.WithAgentAddr(fmt.Sprintf("%s:%s", cfg.Datadog.AgentHost, cfg.Datadog.AgentPort)),
		)
		logger.Infof(ctx, "Datadog APM tracer initialized for service: %s (env: %s)", cfg.Datadog.ServiceName, cfg.Datadog.Environment)
	}

	if err := logger.InitLogger(false, "daybook.log", cfg.Server.Mode); err != nil {
		panic(err)
	}
	defer func() {
		err := logger.Close()
		if err != nil {
			panic(err)
		}
	}()

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

	// Initialize dependency injection container
	logger.Infof(ctx, "Initializing dependency injection container...")
	appContainer := container.NewContainer(database.DB, cfg)
	logger.Infof(ctx, "Dependency injection container initialized")

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)
	logger.Infof(ctx, "Gin mode set to: %s", cfg.Server.Mode)

	// Create router
	router := gin.New()

	// Add recovery middleware
	router.Use(gin.Recovery())

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

	// Add Datadog tracing middleware if enabled
	if cfg.Datadog.Enabled {
		router.Use(gintrace.Middleware(cfg.Datadog.ServiceName))
		logger.Infof(ctx, "Datadog APM middleware configured")
	}

	// Setup routes
	logger.Infof(ctx, "Setting up application routes...")
	routes.SetupRoutes(router, appContainer)
	logger.Infof(ctx, "Routes configured successfully")

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		shutdownCtx := logger.CreateContext("")
		logger.Infof(shutdownCtx, "Shutdown signal received, gracefully shutting down server...")

		// Close database connections
		if err := database.CloseDatabase(); err != nil {
			logger.Errorf(shutdownCtx, "Error closing database: %v", err)
		}

		if err := database.CloseRedis(); err != nil {
			logger.Errorf(shutdownCtx, "Error closing Redis: %v", err)
		}

		// Stop Datadog tracer if enabled
		if cfg.Datadog.Enabled {
			logger.Infof(shutdownCtx, "Stopping Datadog APM tracer...")
			tracer.Stop()
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

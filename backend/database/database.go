package database

import (
	"fmt"

	"daybook-backend/config"
	customLogger "daybook-backend/logger"
	"daybook-backend/models"

	"github.com/go-redis/redis/v8"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis.v8"
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
)

func InitDatabase(cfg *config.Config) error {
	var err error

	// Initialize PostgreSQL
	dsn := cfg.Database.GetDSN()

	// Create custom GORM logger with 200ms slow query threshold
	// Use logger context to ensure trace IDs are populated
	ctx := customLogger.CreateContext("")
	customLogger.Infof(ctx, "Connecting to database: %s", dsn)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: customLogger.NewDefaultCustomLogger(),
	})
	if err != nil {
		customLogger.Errorf(ctx, "Failed to connect to database: %v", err)
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	customLogger.Infof(ctx, "Database connection established successfully")

	// Add Datadog tracing plugin if enabled
	if cfg.Datadog.Enabled {
		if err := DB.Use(gormtrace.NewTracePlugin(gormtrace.WithServiceName(cfg.Datadog.ServiceName))); err != nil {
			customLogger.Warnf(ctx, "Failed to enable Datadog GORM tracing: %v", err)
		} else {
			customLogger.Infof(ctx, "Datadog GORM tracing enabled")
		}
	}

	// Enable UUID extension
	DB.WithContext(ctx).Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	customLogger.Infof(ctx, "UUID extension enabled")

	// Auto-migrate all models
	customLogger.Infof(ctx, "Starting database migration...")
	err = DB.WithContext(ctx).AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.AccountType{},
		&models.Category{},
		&models.Transaction{},
		&models.RecurringTransaction{},
		&models.Tag{},
		&models.CreditCard{},
		&models.CreditCardTransaction{},
		&models.CreditCardPayment{},
		&models.Statement{},
		&models.Reward{},
		&models.Budget{},
		&models.Reconciliation{},
		&models.ReconciliationTransaction{},
		&models.Goal{},
		&models.GoalHolding{},
		&models.GoalContribution{},
		&models.Settings{},
		&models.DebtRecord{},
		&models.DebtPayment{},
		&models.LendRecord{},
		&models.LendPayment{},
		&models.Asset{},
		&models.ServiceRecord{},
		&models.AssetAttachment{},
		&models.ActivityLog{},
	)
	if err != nil {
		customLogger.Errorf(ctx, "Database migration failed: %v", err)
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	customLogger.Infof(ctx, "Database migrated successfully")

	return nil
}

func InitRedis(cfg *config.Config) error {
	ctx := customLogger.CreateContext("")

	// Create Redis client with optional Datadog tracing
	if cfg.Datadog.Enabled {
		// Use Datadog's wrapped client for tracing
		wrappedClient := redistrace.NewClient(&redis.Options{
			Addr:     cfg.Redis.GetAddr(),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}, redistrace.WithServiceName(cfg.Datadog.ServiceName+"-redis"))

		// Type assert to *redis.Client
		if client, ok := wrappedClient.(*redis.Client); ok {
			RedisClient = client
			customLogger.Infof(ctx, "Datadog Redis tracing enabled")
		} else {
			// Fallback to standard client if type assertion fails
			RedisClient = redis.NewClient(&redis.Options{
				Addr:     cfg.Redis.GetAddr(),
				Password: cfg.Redis.Password,
				DB:       cfg.Redis.DB,
			})
			customLogger.Warnf(ctx, "Could not enable Datadog Redis tracing, using standard client")
		}
	} else {
		// Standard Redis client without tracing
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.GetAddr(),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
	}

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		customLogger.Warnf(ctx, "Failed to connect to Redis: %v", err)
		customLogger.Warnf(ctx, "Application will continue without Redis caching")
		RedisClient = nil
		return nil // Don't fail if Redis is not available
	}

	customLogger.Infof(ctx, "Redis connected successfully at %s", cfg.Redis.GetAddr())
	return nil
}

func CloseDatabase() error {
	ctx := customLogger.CreateContext("")

	if DB != nil {
		customLogger.Infof(ctx, "Closing database connection...")
		sqlDB, err := DB.DB()
		if err != nil {
			customLogger.Errorf(ctx, "Error getting database instance: %v", err)
			return err
		}
		err = sqlDB.Close()
		if err != nil {
			customLogger.Errorf(ctx, "Error closing database: %v", err)
			return err
		}
		customLogger.Infof(ctx, "Database connection closed successfully")
	}
	return nil
}

func CloseRedis() error {
	ctx := customLogger.CreateContext("")

	if RedisClient != nil {
		customLogger.Infof(ctx, "Closing Redis connection...")
		err := RedisClient.Close()
		if err != nil {
			customLogger.Errorf(ctx, "Error closing Redis: %v", err)
			return err
		}
		customLogger.Infof(ctx, "Redis connection closed successfully")
	}
	return nil
}

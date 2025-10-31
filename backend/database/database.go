package database

import (
	"context"
	"fmt"
	"log"

	"daybook-backend/config"
	"daybook-backend/models"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
	ctx         = context.Background()
)

func InitDatabase(cfg *config.Config) error {
	var err error

	// Initialize PostgreSQL
	dsn := cfg.Database.GetDSN()
	log.Printf("Connecting to database: %s\n", dsn)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable UUID extension
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	// Auto-migrate all models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.AccountType{},
		&models.Transaction{},
		&models.RecurringTransaction{},
		&models.Tag{},
		&models.CreditCard{},
		&models.CreditCardTransaction{},
		&models.CreditCardPayment{},
		&models.Statement{},
		&models.Reward{},
		&models.Bill{},
		&models.BillPayment{},
		&models.Budget{},
		&models.Reconciliation{},
		&models.ReconciliationTransaction{},
		&models.Goal{},
		&models.GoalHolding{},
		&models.GoalContribution{},
		&models.Settings{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrated successfully")

	return nil
}

func InitRedis(cfg *config.Config) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.GetAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v\n", err)
		log.Println("Application will continue without Redis caching")
		RedisClient = nil
		return nil // Don't fail if Redis is not available
	}

	log.Println("Redis connected successfully")
	return nil
}

func CloseDatabase() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

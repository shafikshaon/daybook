package repository

import (
	"context"

	"gorm.io/gorm"
)

// TransactionManager handles database transactions
// It provides a way to execute multiple repository operations atomically
type TransactionManager interface {
	// WithTransaction executes the provided function within a database transaction
	// If the function returns an error, the transaction is rolled back
	// Otherwise, the transaction is committed
	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) error
}

type gormTransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &gormTransactionManager{db: db}
}

// WithTransaction executes a function within a database transaction
func (tm *gormTransactionManager) WithTransaction(ctx context.Context, fn func(context.Context, *gorm.DB) error) error {
	// Begin transaction
	tx := tm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Execute function
	if err := fn(ctx, tx); err != nil {
		// Rollback on error
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

package repository

import (
	"context"
	"time"

	"daybook-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RecurringTransactionRepository handles recurring transaction data access
type RecurringTransactionRepository interface {
	BaseRepository[models.RecurringTransaction]
	FindEnabled(ctx context.Context, userID uuid.UUID) ([]models.RecurringTransaction, error)
	UpdateLastProcessed(ctx context.Context, id uuid.UUID, processedTime time.Time) error
}

type recurringTransactionRepository struct {
	*GormBaseRepository[models.RecurringTransaction]
}

// NewRecurringTransactionRepository creates a new recurring transaction repository
func NewRecurringTransactionRepository(db *gorm.DB) RecurringTransactionRepository {
	return &recurringTransactionRepository{
		GormBaseRepository: NewGormBaseRepository[models.RecurringTransaction](db),
	}
}

// FindEnabled retrieves all enabled recurring transactions for a user
func (r *recurringTransactionRepository) FindEnabled(ctx context.Context, userID uuid.UUID) ([]models.RecurringTransaction, error) {
	var recurringTransactions []models.RecurringTransaction
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND enabled = ?", userID, true).
		Find(&recurringTransactions).Error
	return recurringTransactions, err
}

// UpdateLastProcessed updates the last processed timestamp for a recurring transaction
func (r *recurringTransactionRepository) UpdateLastProcessed(ctx context.Context, id uuid.UUID, processedTime time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.RecurringTransaction{}).
		Where("id = ?", id).
		Update("last_processed", processedTime).Error
}

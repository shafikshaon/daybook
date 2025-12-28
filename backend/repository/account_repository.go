package repository

import (
	"context"
	"daybook-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccountRepository handles account database operations
type AccountRepository interface {
	BaseRepository[models.Account]
}

type accountRepository struct {
	*GormBaseRepository[models.Account]
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{
		GormBaseRepository: NewGormBaseRepository[models.Account](db),
	}
}

// Override FindAll to order by created_at DESC
func (r *accountRepository) FindAll(ctx context.Context, userID uuid.UUID) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&accounts).Error
	return accounts, err
}

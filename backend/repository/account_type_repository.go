package repository

import (
	"context"

	"daybook-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccountTypeRepository handles account type database operations
type AccountTypeRepository interface {
	BaseRepository[models.AccountType]

	// CountAccountsUsingType counts how many accounts use this account type
	CountAccountsUsingType(ctx context.Context, userID uuid.UUID, typeName string) (int64, error)
}

type accountTypeRepository struct {
	*GormBaseRepository[models.AccountType]
}

// NewAccountTypeRepository creates a new account type repository
func NewAccountTypeRepository(db *gorm.DB) AccountTypeRepository {
	return &accountTypeRepository{
		GormBaseRepository: NewGormBaseRepository[models.AccountType](db),
	}
}

// CountAccountsUsingType counts how many accounts use this account type
func (r *accountTypeRepository) CountAccountsUsingType(ctx context.Context, userID uuid.UUID, typeName string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Account{}).
		Where("user_id = ? AND type = ?", userID, typeName).
		Count(&count).Error
	return count, err
}

// Override FindAll to order by sort_order
func (r *accountTypeRepository) FindAll(ctx context.Context, userID uuid.UUID) ([]models.AccountType, error) {
	var accountTypes []models.AccountType
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("sort_order ASC, name ASC").
		Find(&accountTypes).Error
	return accountTypes, err
}

package repository

import (
	"context"

	"daybook-backend/models"

	"gorm.io/gorm"
)

// LendFilters represents query filters for lends
type LendFilters struct {
	Status *string
}

// LendRepository handles lend data access
type LendRepository interface {
	BaseRepository[models.LendRecord]

	// FindWithFilters retrieves lends with optional filters
	FindWithFilters(ctx context.Context, userID uint, filters LendFilters) ([]models.LendRecord, error)

	// FindPaymentsByLend retrieves all payments for a specific lend
	FindPaymentsByLend(ctx context.Context, lendID, userID uint) ([]models.LendPayment, error)

	// CreatePayment creates a new lend payment record
	CreatePayment(ctx context.Context, payment *models.LendPayment) error

	// GetAccountName retrieves the account name by ID
	GetAccountName(ctx context.Context, accountID uint) (string, error)
}

type lendRepository struct {
	*GormBaseRepository[models.LendRecord]
}

// NewLendRepository creates a new lend repository
func NewLendRepository(db *gorm.DB) LendRepository {
	return &lendRepository{
		GormBaseRepository: NewGormBaseRepository[models.LendRecord](db),
	}
}

// FindWithFilters retrieves lends with optional filters
func (r *lendRepository) FindWithFilters(ctx context.Context, userID uint, filters LendFilters) ([]models.LendRecord, error) {
	var lends []models.LendRecord

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	err := query.Order("lent_date DESC, created_at DESC").Find(&lends).Error
	return lends, err
}

// FindPaymentsByLend retrieves all payments for a specific lend
func (r *lendRepository) FindPaymentsByLend(ctx context.Context, lendID, userID uint) ([]models.LendPayment, error) {
	var payments []models.LendPayment
	err := r.db.WithContext(ctx).
		Where("lend_id = ? AND user_id = ?", lendID, userID).
		Order("payment_date DESC, created_at DESC").
		Find(&payments).Error
	return payments, err
}

// CreatePayment creates a new lend payment record
func (r *lendRepository) CreatePayment(ctx context.Context, payment *models.LendPayment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// GetAccountName retrieves the account name by ID
func (r *lendRepository) GetAccountName(ctx context.Context, accountID uint) (string, error) {
	var account models.Account
	err := r.db.WithContext(ctx).Select("name").Where("id = ?", accountID).First(&account).Error
	if err != nil {
		return "", err
	}
	return account.Name, nil
}

// WithTx returns a new repository instance with the transaction
func (r *lendRepository) WithTx(tx *gorm.DB) BaseRepository[models.LendRecord] {
	return &lendRepository{
		GormBaseRepository: NewGormBaseRepository[models.LendRecord](tx),
	}
}

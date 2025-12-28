package repository

import (
	"context"

	"daybook-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DebtFilters represents query filters for debts
type DebtFilters struct {
	Status *string
}

// DebtRepository handles debt data access
type DebtRepository interface {
	BaseRepository[models.DebtRecord]

	// FindWithFilters retrieves debts with optional filters
	FindWithFilters(ctx context.Context, userID uuid.UUID, filters DebtFilters) ([]models.DebtRecord, error)

	// FindPaymentsByDebt retrieves all payments for a specific debt
	FindPaymentsByDebt(ctx context.Context, debtID, userID uuid.UUID) ([]models.DebtPayment, error)

	// CreatePayment creates a new debt payment record
	CreatePayment(ctx context.Context, payment *models.DebtPayment) error

	// GetAccountName retrieves the account name by ID
	GetAccountName(ctx context.Context, accountID uuid.UUID) (string, error)
}

type debtRepository struct {
	*GormBaseRepository[models.DebtRecord]
}

// NewDebtRepository creates a new debt repository
func NewDebtRepository(db *gorm.DB) DebtRepository {
	return &debtRepository{
		GormBaseRepository: NewGormBaseRepository[models.DebtRecord](db),
	}
}

// FindWithFilters retrieves debts with optional filters
func (r *debtRepository) FindWithFilters(ctx context.Context, userID uuid.UUID, filters DebtFilters) ([]models.DebtRecord, error) {
	var debts []models.DebtRecord

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	err := query.Order("borrowed_date DESC, created_at DESC").Find(&debts).Error
	return debts, err
}

// FindPaymentsByDebt retrieves all payments for a specific debt
func (r *debtRepository) FindPaymentsByDebt(ctx context.Context, debtID, userID uuid.UUID) ([]models.DebtPayment, error) {
	var payments []models.DebtPayment
	err := r.db.WithContext(ctx).
		Where("debt_id = ? AND user_id = ?", debtID, userID).
		Order("payment_date DESC, created_at DESC").
		Find(&payments).Error
	return payments, err
}

// CreatePayment creates a new debt payment record
func (r *debtRepository) CreatePayment(ctx context.Context, payment *models.DebtPayment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// GetAccountName retrieves the account name by ID
func (r *debtRepository) GetAccountName(ctx context.Context, accountID uuid.UUID) (string, error) {
	var account models.Account
	err := r.db.WithContext(ctx).Select("name").Where("id = ?", accountID).First(&account).Error
	if err != nil {
		return "", err
	}
	return account.Name, nil
}

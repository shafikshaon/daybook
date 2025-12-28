package repository

import (
	"context"

	"daybook-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreditCardRepository handles credit card data access
type CreditCardRepository interface {
	BaseRepository[models.CreditCard]

	// Transaction operations
	CreateTransaction(ctx context.Context, transaction *models.CreditCardTransaction) error
	FindTransactionsByCard(ctx context.Context, cardID, userID uuid.UUID) ([]models.CreditCardTransaction, error)
	FindTransactionByID(ctx context.Context, transactionID, userID uuid.UUID) (*models.CreditCardTransaction, error)
	DeleteTransaction(ctx context.Context, transactionID, userID uuid.UUID) error

	// Payment operations
	CreatePayment(ctx context.Context, payment *models.CreditCardPayment) error
	FindPaymentsByCard(ctx context.Context, cardID, userID uuid.UUID) ([]models.CreditCardPayment, error)

	// Statement operations
	CreateStatement(ctx context.Context, statement *models.Statement) error
	FindStatementsByCard(ctx context.Context, cardID, userID uuid.UUID) ([]models.Statement, error)

	// Reward operations
	CreateReward(ctx context.Context, reward *models.Reward) error
	FindRewardsByCard(ctx context.Context, cardID, userID uuid.UUID) ([]models.Reward, error)
	FindRewardsByUser(ctx context.Context, userID uuid.UUID) ([]models.Reward, error)
}

type creditCardRepository struct {
	*GormBaseRepository[models.CreditCard]
}

// NewCreditCardRepository creates a new credit card repository
func NewCreditCardRepository(db *gorm.DB) CreditCardRepository {
	return &creditCardRepository{
		GormBaseRepository: NewGormBaseRepository[models.CreditCard](db),
	}
}

// CreateTransaction creates a new credit card transaction
func (r *creditCardRepository) CreateTransaction(ctx context.Context, transaction *models.CreditCardTransaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

// FindTransactionsByCard retrieves all transactions for a credit card
func (r *creditCardRepository) FindTransactionsByCard(ctx context.Context, cardID, userID uuid.UUID) ([]models.CreditCardTransaction, error) {
	var transactions []models.CreditCardTransaction
	err := r.db.WithContext(ctx).
		Where("card_id = ? AND user_id = ?", cardID, userID).
		Order("date DESC, created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

// FindTransactionByID retrieves a specific credit card transaction
func (r *creditCardRepository) FindTransactionByID(ctx context.Context, transactionID, userID uuid.UUID) (*models.CreditCardTransaction, error) {
	var transaction models.CreditCardTransaction
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", transactionID, userID).
		First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// DeleteTransaction deletes a credit card transaction
func (r *creditCardRepository) DeleteTransaction(ctx context.Context, transactionID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", transactionID, userID).
		Delete(&models.CreditCardTransaction{}).Error
}

// CreatePayment creates a new credit card payment
func (r *creditCardRepository) CreatePayment(ctx context.Context, payment *models.CreditCardPayment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

// FindPaymentsByCard retrieves all payments for a credit card
func (r *creditCardRepository) FindPaymentsByCard(ctx context.Context, cardID, userID uuid.UUID) ([]models.CreditCardPayment, error) {
	var payments []models.CreditCardPayment
	err := r.db.WithContext(ctx).
		Where("card_id = ? AND user_id = ?", cardID, userID).
		Order("payment_date DESC, created_at DESC").
		Find(&payments).Error
	return payments, err
}

// CreateStatement creates a new statement
func (r *creditCardRepository) CreateStatement(ctx context.Context, statement *models.Statement) error {
	return r.db.WithContext(ctx).Create(statement).Error
}

// FindStatementsByCard retrieves all statements for a credit card
func (r *creditCardRepository) FindStatementsByCard(ctx context.Context, cardID, userID uuid.UUID) ([]models.Statement, error) {
	var statements []models.Statement
	err := r.db.WithContext(ctx).
		Where("card_id = ? AND user_id = ?", cardID, userID).
		Order("statement_date DESC").
		Find(&statements).Error
	return statements, err
}

// CreateReward creates a new reward
func (r *creditCardRepository) CreateReward(ctx context.Context, reward *models.Reward) error {
	return r.db.WithContext(ctx).Create(reward).Error
}

// FindRewardsByCard retrieves all rewards for a credit card
func (r *creditCardRepository) FindRewardsByCard(ctx context.Context, cardID, userID uuid.UUID) ([]models.Reward, error) {
	var rewards []models.Reward
	err := r.db.WithContext(ctx).
		Where("card_id = ? AND user_id = ?", cardID, userID).
		Order("earned_date DESC").
		Find(&rewards).Error
	return rewards, err
}

// FindRewardsByUser retrieves all rewards for a user
func (r *creditCardRepository) FindRewardsByUser(ctx context.Context, userID uuid.UUID) ([]models.Reward, error) {
	var rewards []models.Reward
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("earned_date DESC").
		Find(&rewards).Error
	return rewards, err
}

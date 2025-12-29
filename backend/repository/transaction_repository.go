package repository

import (
	"context"
	"time"

	"daybook-backend/models"

	"gorm.io/gorm"
)

// TransactionFilters holds filter parameters for transaction queries
type TransactionFilters struct {
	Type            *string
	CategoryID      *string
	AccountID       *string
	StartDate       *time.Time
	EndDate         *time.Time
	IncludeTracking bool
}

// PaginationParams holds pagination parameters
type PaginationParams struct {
	Page  int
	Limit int
}

// TransactionListResult holds paginated transaction results
type TransactionListResult struct {
	Transactions []models.Transaction
	TotalCount   int64
}

// TransactionStats holds transaction statistics
type TransactionStats struct {
	TotalIncome   float64
	TotalExpense  float64
	TotalTransfer float64
	NetAmount     float64
}

// TransactionRepository handles transaction data access
type TransactionRepository interface {
	BaseRepository[models.Transaction]
	FindWithFilters(ctx context.Context, userID uint, filters TransactionFilters, pagination PaginationParams) (*TransactionListResult, error)
	CountWithFilters(ctx context.Context, userID uint, filters TransactionFilters) (int64, error)
	CalculateStats(ctx context.Context, userID uint, filters TransactionFilters) (*TransactionStats, error)
	BulkCreate(ctx context.Context, transactions []models.Transaction) error
}

type transactionRepository struct {
	*GormBaseRepository[models.Transaction]
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		GormBaseRepository: NewGormBaseRepository[models.Transaction](db),
	}
}

// FindWithFilters retrieves transactions with filtering and pagination
func (r *transactionRepository) FindWithFilters(ctx context.Context, userID uint, filters TransactionFilters, pagination PaginationParams) (*TransactionListResult, error) {
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	// Exclude tracking transactions unless explicitly requested
	if !filters.IncludeTracking {
		query = query.Where("type != ?", "tracking")
	}

	// Apply filters
	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}
	if filters.AccountID != nil {
		query = query.Where("account_id = ?", *filters.AccountID)
	}
	if filters.StartDate != nil {
		query = query.Where("date >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("date <= ?", *filters.EndDate)
	}

	// Get total count
	var totalCount int64
	if err := query.Model(&models.Transaction{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.Limit
	var transactions []models.Transaction
	if err := query.Order("date DESC, created_at DESC").
		Limit(pagination.Limit).
		Offset(offset).
		Find(&transactions).Error; err != nil {
		return nil, err
	}

	return &TransactionListResult{
		Transactions: transactions,
		TotalCount:   totalCount,
	}, nil
}

// CountWithFilters counts transactions matching the filters
func (r *transactionRepository) CountWithFilters(ctx context.Context, userID uint, filters TransactionFilters) (int64, error) {
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if !filters.IncludeTracking {
		query = query.Where("type != ?", "tracking")
	}
	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}
	if filters.AccountID != nil {
		query = query.Where("account_id = ?", *filters.AccountID)
	}
	if filters.StartDate != nil {
		query = query.Where("date >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("date <= ?", *filters.EndDate)
	}

	var count int64
	err := query.Model(&models.Transaction{}).Count(&count).Error
	return count, err
}

// CalculateStats calculates transaction statistics based on filters
func (r *transactionRepository) CalculateStats(ctx context.Context, userID uint, filters TransactionFilters) (*TransactionStats, error) {
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if !filters.IncludeTracking {
		query = query.Where("type != ?", "tracking")
	}
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}
	if filters.AccountID != nil {
		query = query.Where("account_id = ?", *filters.AccountID)
	}
	if filters.StartDate != nil {
		query = query.Where("date >= ?", *filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("date <= ?", *filters.EndDate)
	}

	stats := &TransactionStats{}

	// Calculate total income
	var incomeSum float64
	query.Model(&models.Transaction{}).
		Where("type = ?", "income").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&incomeSum)
	stats.TotalIncome = incomeSum

	// Calculate total expense
	var expenseSum float64
	query.Model(&models.Transaction{}).
		Where("type = ?", "expense").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&expenseSum)
	stats.TotalExpense = expenseSum

	// Calculate total transfers
	var transferSum float64
	query.Model(&models.Transaction{}).
		Where("type = ?", "transfer").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&transferSum)
	stats.TotalTransfer = transferSum

	stats.NetAmount = incomeSum - expenseSum

	return stats, nil
}

// BulkCreate creates multiple transactions in a single operation
func (r *transactionRepository) BulkCreate(ctx context.Context, transactions []models.Transaction) error {
	return r.db.WithContext(ctx).Create(&transactions).Error
}

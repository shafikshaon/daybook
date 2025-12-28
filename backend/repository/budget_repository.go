package repository

import (
	"context"
	"time"

	"daybook-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BudgetFilters represents filter options for budget queries
type BudgetFilters struct {
	Enabled    *bool
	Period     string
	CategoryID string
}

// BudgetRepository handles budget database operations
type BudgetRepository interface {
	BaseRepository[models.Budget]

	// FindWithFilters retrieves budgets with optional filters
	FindWithFilters(ctx context.Context, userID uuid.UUID, filters BudgetFilters) ([]models.Budget, error)

	// CalculateTotalSpent calculates total spending for a category in a date range
	CalculateTotalSpent(ctx context.Context, userID uuid.UUID, categoryID string, startDate, endDate time.Time) (float64, error)
}

type budgetRepository struct {
	*GormBaseRepository[models.Budget]
}

// NewBudgetRepository creates a new budget repository
func NewBudgetRepository(db *gorm.DB) BudgetRepository {
	return &budgetRepository{
		GormBaseRepository: NewGormBaseRepository[models.Budget](db),
	}
}

// FindWithFilters retrieves budgets with optional filters
func (r *budgetRepository) FindWithFilters(ctx context.Context, userID uuid.UUID, filters BudgetFilters) ([]models.Budget, error) {
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if filters.Enabled != nil {
		query = query.Where("enabled = ?", *filters.Enabled)
	}

	if filters.Period != "" {
		query = query.Where("period = ?", filters.Period)
	}

	if filters.CategoryID != "" {
		query = query.Where("category_id = ?", filters.CategoryID)
	}

	var budgets []models.Budget
	err := query.Order("created_at DESC").Find(&budgets).Error
	return budgets, err
}

// CalculateTotalSpent calculates total spending for a category in a date range
func (r *budgetRepository) CalculateTotalSpent(ctx context.Context, userID uuid.UUID, categoryID string, startDate, endDate time.Time) (float64, error) {
	var totalSpent float64
	err := r.db.WithContext(ctx).
		Model(&models.Transaction{}).
		Where("user_id = ? AND category_id = ? AND type = ? AND date >= ? AND date < ?",
			userID, categoryID, "expense", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Row().Scan(&totalSpent)
	return totalSpent, err
}

// Override FindAll to order by created_at DESC
func (r *budgetRepository) FindAll(ctx context.Context, userID uuid.UUID) ([]models.Budget, error) {
	var budgets []models.Budget
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&budgets).Error
	return budgets, err
}

package repository

import (
	"context"
	"daybook-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CategoryRepository handles category database operations
type CategoryRepository interface {
	BaseRepository[models.Category]

	// FindByType retrieves categories filtered by type
	FindByType(ctx context.Context, userID uuid.UUID, categoryType string) ([]models.Category, error)

	// FindByNameAndType finds a category by name and type for a user
	FindByNameAndType(ctx context.Context, userID uuid.UUID, name, categoryType string) (*models.Category, error)

	// GetMaxOrder gets the maximum order value for a given type
	GetMaxOrder(ctx context.Context, userID uuid.UUID, categoryType string) (int, error)

	// BulkUpdateOrder updates the order of multiple categories
	BulkUpdateOrder(ctx context.Context, categoryOrders []CategoryOrder) error
}

// CategoryOrder represents a category ID and its new order
type CategoryOrder struct {
	ID    uuid.UUID
	Order int
}

type categoryRepository struct {
	*GormBaseRepository[models.Category]
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		GormBaseRepository: NewGormBaseRepository[models.Category](db),
	}
}

// FindByType retrieves categories of a specific type
func (r *categoryRepository) FindByType(ctx context.Context, userID uuid.UUID, categoryType string) ([]models.Category, error) {
	var categories []models.Category
	err := r.Query(ctx, userID).
		Where("type = ?", categoryType).
		Order("\"order\" ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// FindByNameAndType finds a category by name and type
func (r *categoryRepository) FindByNameAndType(ctx context.Context, userID uuid.UUID, name, categoryType string) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND name = ? AND type = ?", userID, name, categoryType).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetMaxOrder gets the maximum order value for a given type
func (r *categoryRepository) GetMaxOrder(ctx context.Context, userID uuid.UUID, categoryType string) (int, error) {
	var maxOrder int
	err := r.db.WithContext(ctx).
		Model(&models.Category{}).
		Where("user_id = ? AND type = ?", userID, categoryType).
		Select("COALESCE(MAX(\"order\"), 0)").
		Scan(&maxOrder).Error
	return maxOrder, err
}

// BulkUpdateOrder updates the order of multiple categories
func (r *categoryRepository) BulkUpdateOrder(ctx context.Context, categoryOrders []CategoryOrder) error {
	for _, catOrder := range categoryOrders {
		if err := r.db.WithContext(ctx).
			Model(&models.Category{}).
			Where("id = ?", catOrder.ID).
			Update("order", catOrder.Order).Error; err != nil {
			return err
		}
	}
	return nil
}

// Override FindAll to order by order field
func (r *categoryRepository) FindAll(ctx context.Context, userID uuid.UUID) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("type ASC, \"order\" ASC").
		Find(&categories).Error
	return categories, err
}

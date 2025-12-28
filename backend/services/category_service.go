package services

import (
	"context"
	"errors"
	"fmt"

	"daybook-backend/models"
	"daybook-backend/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CategoryService handles category business logic
type CategoryService interface {
	// ListCategories retrieves all categories for a user, optionally filtered by type
	ListCategories(ctx context.Context, userID uuid.UUID, categoryType string) ([]models.Category, error)

	// GetCategory retrieves a specific category by ID
	GetCategory(ctx context.Context, categoryID, userID uuid.UUID) (*models.Category, error)

	// CreateCategory creates a new category
	CreateCategory(ctx context.Context, category *models.Category) (*models.Category, error)

	// UpdateCategory updates an existing category
	UpdateCategory(ctx context.Context, categoryID, userID uuid.UUID, updateData *models.Category) (*models.Category, error)

	// DeleteCategory deletes a category
	DeleteCategory(ctx context.Context, categoryID, userID uuid.UUID) error

	// ReorderCategories updates the order of multiple categories
	ReorderCategories(ctx context.Context, userID uuid.UUID, categoryOrders []repository.CategoryOrder) error

	// GetAvailableIcons returns the list of available icons
	GetAvailableIcons() map[string][]string
}

type categoryService struct {
	repo           repository.CategoryRepository
	activityLogger ActivityLogService
}

// NewCategoryService creates a new category service
func NewCategoryService(
	repo repository.CategoryRepository,
	activityLogger ActivityLogService,
) CategoryService {
	return &categoryService{
		repo:           repo,
		activityLogger: activityLogger,
	}
}

// ListCategories retrieves categories, optionally filtered by type
func (s *categoryService) ListCategories(ctx context.Context, userID uuid.UUID, categoryType string) ([]models.Category, error) {
	if categoryType != "" {
		// Validate category type
		if categoryType != "income" && categoryType != "expense" && categoryType != "transfer" {
			return nil, errors.New("invalid category type")
		}
		return s.repo.FindByType(ctx, userID, categoryType)
	}
	return s.repo.FindAll(ctx, userID)
}

// GetCategory retrieves a specific category
func (s *categoryService) GetCategory(ctx context.Context, categoryID, userID uuid.UUID) (*models.Category, error) {
	return s.repo.FindByID(ctx, categoryID, userID)
}

// CreateCategory creates a new category
func (s *categoryService) CreateCategory(ctx context.Context, category *models.Category) (*models.Category, error) {
	// Validate category type
	if category.Type != "income" && category.Type != "expense" && category.Type != "transfer" {
		return nil, errors.New("category type must be 'income', 'expense', or 'transfer'")
	}

	// Check for duplicate name within the same type
	existing, err := s.repo.FindByNameAndType(ctx, category.UserID, category.Name, category.Type)
	if err == nil && existing != nil {
		return nil, errors.New("category with this name already exists for this type")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// User-created categories cannot be default
	category.IsDefault = false

	// Get the next order value for this type
	maxOrder, err := s.repo.GetMaxOrder(ctx, category.UserID, category.Type)
	if err != nil {
		return nil, err
	}
	category.Order = maxOrder + 1

	// Create the category
	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		category.UserID,
		models.ActionCreate,
		models.ModuleCategory,
		"Category",
		category.ID,
		"Created category: "+category.Name,
		nil,
	)

	return category, nil
}

// UpdateCategory updates an existing category
func (s *categoryService) UpdateCategory(ctx context.Context, categoryID, userID uuid.UUID, updateData *models.Category) (*models.Category, error) {
	// Fetch existing category
	existing, err := s.repo.FindByID(ctx, categoryID, userID)
	if err != nil {
		return nil, err
	}

	// Validate category type if being updated
	if updateData.Type != "" && updateData.Type != "income" && updateData.Type != "expense" && updateData.Type != "transfer" {
		return nil, errors.New("category type must be 'income', 'expense', or 'transfer'")
	}

	// Check for duplicate name (excluding current category)
	if updateData.Name != "" && updateData.Type != "" {
		duplicate, err := s.repo.FindByNameAndType(ctx, userID, updateData.Name, updateData.Type)
		if err == nil && duplicate != nil && duplicate.ID != categoryID {
			return nil, errors.New("category with this name already exists for this type")
		}
	}

	// Update allowed fields
	if updateData.Name != "" {
		existing.Name = updateData.Name
	}
	if updateData.Type != "" {
		existing.Type = updateData.Type
	}
	if updateData.Icon != "" {
		existing.Icon = updateData.Icon
	}
	if updateData.Color != "" {
		existing.Color = updateData.Color
	}
	existing.Description = updateData.Description

	// Cannot modify IsDefault field
	// existing.IsDefault remains unchanged

	// Save updates
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionUpdate,
		models.ModuleCategory,
		"Category",
		existing.ID,
		"Updated category: "+existing.Name,
		nil,
	)

	return existing, nil
}

// DeleteCategory deletes a category
func (s *categoryService) DeleteCategory(ctx context.Context, categoryID, userID uuid.UUID) error {
	// Fetch the category to get its name for logging
	category, err := s.repo.FindByID(ctx, categoryID, userID)
	if err != nil {
		return err
	}

	// TODO: Add validation to check if category is used in transactions, budgets, etc.
	// This should be done when we migrate those entities

	// Delete the category
	if err := s.repo.Delete(ctx, categoryID, userID); err != nil {
		return err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionDelete,
		models.ModuleCategory,
		"Category",
		category.ID,
		"Deleted category: "+category.Name,
		nil,
	)

	return nil
}

// ReorderCategories updates the order of multiple categories
func (s *categoryService) ReorderCategories(ctx context.Context, userID uuid.UUID, categoryOrders []repository.CategoryOrder) error {
	// Validate that all categories belong to the user
	for _, catOrder := range categoryOrders {
		_, err := s.repo.FindByID(ctx, catOrder.ID, userID)
		if err != nil {
			return fmt.Errorf("category not found or unauthorized: %s", catOrder.ID)
		}
	}

	// Bulk update orders
	if err := s.repo.BulkUpdateOrder(ctx, categoryOrders); err != nil {
		return err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionUpdate,
		models.ModuleCategory,
		"Category",
		uuid.Nil,
		"Reordered categories",
		nil,
	)

	return nil
}

// GetAvailableIcons returns available icons
func (s *categoryService) GetAvailableIcons() map[string][]string {
	return models.GetAvailableIcons()
}

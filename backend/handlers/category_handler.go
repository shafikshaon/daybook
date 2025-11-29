package handlers

import (
	"net/http"

	"daybook-backend/database"
	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListCategories returns all categories for the authenticated user
func ListCategories(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListCategories - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Fetching categories for user: %s", userID)
	query := database.DB.WithContext(ctx).Where("user_id = ?", userID)

	// Filter by type if specified
	if categoryType := c.Query("type"); categoryType != "" {
		if categoryType == "income" || categoryType == "expense" || categoryType == "transfer" {
			query = query.Where("type = ?", categoryType)
		}
	}

	var categories []models.Category
	if err := query.Order("type ASC, name ASC").Find(&categories).Error; err != nil {
		logger.Errorf(ctx, "Failed to fetch categories: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}

	logger.Infof(ctx, "Categories retrieved successfully for user: %s, count: %d", userID, len(categories))
	utilities.SuccessResponse(c, categories, "Categories retrieved successfully")
}

// GetCategory returns a specific category by ID
func GetCategory(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetCategory - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid category ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	logger.Debugf(ctx, "Fetching category %s for user: %s", categoryID, userID)
	var category models.Category
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		logger.Warnf(ctx, "Category not found: %s, error: %v", categoryID, err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Category not found")
		return
	}

	logger.Infof(ctx, "Category retrieved successfully: %s for user: %s", categoryID, userID)
	utilities.SuccessResponse(c, category, "Category retrieved successfully")
}

// CreateCategory creates a new category
func CreateCategory(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateCategory - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Parsing category request for user: %s", userID)
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		logger.Warnf(ctx, "Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate category type
	if category.Type != "income" && category.Type != "expense" && category.Type != "transfer" {
		logger.Warnf(ctx, "Invalid category type: %s", category.Type)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Category type must be 'income', 'expense', or 'transfer'")
		return
	}

	category.UserID = userID
	category.IsDefault = false // User-created categories are never default

	// Check for duplicate category name for this user and type
	var existingCount int64
	database.DB.WithContext(ctx).Model(&models.Category{}).
		Where("user_id = ? AND name = ? AND type = ?", userID, category.Name, category.Type).
		Count(&existingCount)

	if existingCount > 0 {
		logger.Warnf(ctx, "Category with name '%s' already exists for type '%s'", category.Name, category.Type)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Category with this name already exists for this type")
		return
	}

	logger.Debugf(ctx, "Creating category in database")
	if err := database.DB.WithContext(ctx).Create(&category).Error; err != nil {
		logger.Errorf(ctx, "Failed to create category: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create category")
		return
	}

	// Log category creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleCategory,
		"Category", category.ID, "Created category: "+category.Name, nil)

	logger.Infof(ctx, "Category created successfully: %s for user: %s", category.ID, userID)
	utilities.CreatedResponse(c, category, "Category created successfully")
}

// UpdateCategory updates an existing category
func UpdateCategory(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateCategory - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid category ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	logger.Debugf(ctx, "Fetching existing category %s for user: %s", categoryID, userID)
	var existingCategory models.Category
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", categoryID, userID).First(&existingCategory).Error; err != nil {
		logger.Warnf(ctx, "Category not found: %s, error: %v", categoryID, err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Category not found")
		return
	}

	logger.Debugf(ctx, "Parsing update data for category: %s", categoryID)
	var updateData models.Category
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate category type
	if updateData.Type != "income" && updateData.Type != "expense" && updateData.Type != "transfer" {
		logger.Warnf(ctx, "Invalid category type: %s", updateData.Type)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Category type must be 'income', 'expense', or 'transfer'")
		return
	}

	// Check for duplicate category name (excluding current category)
	var existingCount int64
	database.DB.WithContext(ctx).Model(&models.Category{}).
		Where("user_id = ? AND name = ? AND type = ? AND id != ?", userID, updateData.Name, updateData.Type, categoryID).
		Count(&existingCount)

	if existingCount > 0 {
		logger.Warnf(ctx, "Category with name '%s' already exists for type '%s'", updateData.Name, updateData.Type)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Category with this name already exists for this type")
		return
	}

	// Update category fields
	existingCategory.Name = updateData.Name
	existingCategory.Type = updateData.Type
	existingCategory.Icon = updateData.Icon
	existingCategory.Color = updateData.Color
	existingCategory.Description = updateData.Description

	logger.Debugf(ctx, "Updating category: %s", categoryID)
	if err := database.DB.WithContext(ctx).Save(&existingCategory).Error; err != nil {
		logger.Errorf(ctx, "Failed to update category: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update category")
		return
	}

	// Log category update activity
	utilities.LogEntityActivity(c, userID, models.ActionUpdate, models.ModuleCategory,
		"Category", existingCategory.ID, "Updated category: "+existingCategory.Name, nil)

	logger.Infof(ctx, "Category updated successfully: %s for user: %s", categoryID, userID)
	utilities.SuccessResponse(c, existingCategory, "Category updated successfully")
}

// DeleteCategory deletes a category
func DeleteCategory(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteCategory - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid category ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid category ID")
		return
	}

	logger.Debugf(ctx, "Fetching category %s for deletion", categoryID)
	var category models.Category
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		logger.Warnf(ctx, "Category not found: %s, error: %v", categoryID, err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Category not found")
		return
	}

	// Check if category is being used in transactions
	var transactionCount int64
	database.DB.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ? AND category_id = ?", userID, categoryID.String()).
		Count(&transactionCount)

	if transactionCount > 0 {
		logger.Warnf(ctx, "Cannot delete category: %s - used in %d transactions", categoryID, transactionCount)
		utilities.ErrorResponse(c, http.StatusBadRequest,
			"Cannot delete category as it is being used in transactions. Please reassign transactions first.")
		return
	}

	// Check if category is being used in recurring transactions
	var recurringCount int64
	database.DB.WithContext(ctx).Model(&models.RecurringTransaction{}).
		Where("user_id = ? AND template_category_id = ?", userID, categoryID.String()).
		Count(&recurringCount)

	if recurringCount > 0 {
		logger.Warnf(ctx, "Cannot delete category: %s - used in %d recurring transactions", categoryID, recurringCount)
		utilities.ErrorResponse(c, http.StatusBadRequest,
			"Cannot delete category as it is being used in recurring transactions. Please reassign or delete recurring transactions first.")
		return
	}

	// Check if category is being used in budgets
	var budgetCount int64
	database.DB.WithContext(ctx).Model(&models.Budget{}).
		Where("user_id = ? AND category_id = ?", userID, categoryID.String()).
		Count(&budgetCount)

	if budgetCount > 0 {
		logger.Warnf(ctx, "Cannot delete category: %s - used in %d budgets", categoryID, budgetCount)
		utilities.ErrorResponse(c, http.StatusBadRequest,
			"Cannot delete category as it is being used in budgets. Please delete or update budgets first.")
		return
	}

	logger.Debugf(ctx, "Deleting category: %s", categoryID)
	if err := database.DB.WithContext(ctx).Delete(&category).Error; err != nil {
		logger.Errorf(ctx, "Failed to delete category: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete category")
		return
	}

	// Log category deletion activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, models.ModuleCategory,
		"Category", category.ID, "Deleted category: "+category.Name, nil)

	logger.Infof(ctx, "Category deleted successfully: %s for user: %s", categoryID, userID)
	utilities.SuccessResponse(c, nil, "Category deleted successfully")
}

// GetAvailableIcons returns the list of available icons
func GetAvailableIcons(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetAvailableIcons - Entry")

	_, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	icons := models.GetAvailableIcons()
	logger.Infof(ctx, "Available icons retrieved successfully")
	utilities.SuccessResponse(c, icons, "Available icons retrieved successfully")
}

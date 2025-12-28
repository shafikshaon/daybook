package handlers

import (
	"net/http"

	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/repository"
	"daybook-backend/services"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CategoryHandler handles category-related HTTP requests
type CategoryHandler struct {
	service services.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// ListCategories returns all categories for the authenticated user
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListCategories - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get type filter from query parameter
	categoryType := c.Query("type")

	logger.Debugf(ctx, "Fetching categories for user: %s, type: %s", userID, categoryType)
	categories, err := h.service.ListCategories(ctx, userID, categoryType)
	if err != nil {
		logger.Errorf(ctx, "Failed to fetch categories: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}

	logger.Infof(ctx, "Categories retrieved successfully for user: %s, count: %d", userID, len(categories))
	utilities.SuccessResponse(c, categories, "Categories retrieved successfully")
}

// GetCategory returns a specific category by ID
func (h *CategoryHandler) GetCategory(c *gin.Context) {
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
	category, err := h.service.GetCategory(ctx, categoryID, userID)
	if err != nil {
		logger.Warnf(ctx, "Category not found: %s, error: %v", categoryID, err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Category not found")
		return
	}

	logger.Infof(ctx, "Category retrieved successfully: %s for user: %s", categoryID, userID)
	utilities.SuccessResponse(c, category, "Category retrieved successfully")
}

// CreateCategory creates a new category
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
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

	category.UserID = userID

	logger.Debugf(ctx, "Creating category for user: %s", userID)
	created, err := h.service.CreateCategory(ctx, &category)
	if err != nil {
		logger.Errorf(ctx, "Failed to create category: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Category created successfully: %s for user: %s", created.ID, userID)
	utilities.CreatedResponse(c, created, "Category created successfully")
}

// UpdateCategory updates an existing category
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
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

	logger.Debugf(ctx, "Parsing update request for category %s", categoryID)
	var updateData models.Category
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "Updating category %s for user: %s", categoryID, userID)
	updated, err := h.service.UpdateCategory(ctx, categoryID, userID, &updateData)
	if err != nil {
		logger.Errorf(ctx, "Failed to update category: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Category updated successfully: %s for user: %s", categoryID, userID)
	utilities.SuccessResponse(c, updated, "Category updated successfully")
}

// DeleteCategory deletes a category
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
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

	logger.Debugf(ctx, "Deleting category %s for user: %s", categoryID, userID)
	if err := h.service.DeleteCategory(ctx, categoryID, userID); err != nil {
		logger.Errorf(ctx, "Failed to delete category: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Category deleted successfully: %s for user: %s", categoryID, userID)
	utilities.SuccessResponse(c, nil, "Category deleted successfully")
}

// ReorderCategories updates the order of multiple categories in bulk
func (h *CategoryHandler) ReorderCategories(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ReorderCategories - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Request payload structure
	type ReorderRequest struct {
		Categories []struct {
			ID    uuid.UUID `json:"id" binding:"required"`
			Order int       `json:"order" binding:"required"`
		} `json:"categories" binding:"required,min=1"`
	}

	var request ReorderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Warnf(ctx, "Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "Reordering %d categories for user: %s", len(request.Categories), userID)

	// Convert to repository format
	categoryOrders := make([]repository.CategoryOrder, len(request.Categories))
	for i, cat := range request.Categories {
		categoryOrders[i] = repository.CategoryOrder{
			ID:    cat.ID,
			Order: cat.Order,
		}
	}

	if err := h.service.ReorderCategories(ctx, userID, categoryOrders); err != nil {
		logger.Errorf(ctx, "Failed to reorder categories: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Categories reordered successfully for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Categories reordered successfully")
}

// GetAvailableIcons returns the list of available icons
func (h *CategoryHandler) GetAvailableIcons(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetAvailableIcons - Entry")

	_, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	icons := h.service.GetAvailableIcons()
	logger.Infof(ctx, "Available icons retrieved successfully")
	utilities.SuccessResponse(c, icons, "Available icons retrieved successfully")
}

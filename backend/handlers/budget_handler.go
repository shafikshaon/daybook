package handlers

import (
	"net/http"
	"strconv"

	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/repository"
	"daybook-backend/services"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
)

// BudgetHandler handles budget-related HTTP requests
type BudgetHandler struct {
	service services.BudgetService
}

// NewBudgetHandler creates a new budget handler
func NewBudgetHandler(service services.BudgetService) *BudgetHandler {
	return &BudgetHandler{service: service}
}

// ListBudgets returns all budgets for the authenticated user
func (h *BudgetHandler) ListBudgets(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListBudgets - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	// Build filters from query parameters
	filters := repository.BudgetFilters{}

	if enabled := c.Query("enabled"); enabled != "" {
		enabledBool := enabled == "true"
		filters.Enabled = &enabledBool
	}

	if period := c.Query("period"); period != "" {
		filters.Period = period
	}

	if categoryID := c.Query("categoryId"); categoryID != "" {
		if categoryIDUint, err := strconv.ParseUint(categoryID, 10, 32); err == nil {
			filters.CategoryID = uint(categoryIDUint)
		}
	}

	logger.Debugf(ctx, "Fetching budgets with progress from database...")
	budgetsWithProgress, err := h.service.ListBudgetsWithProgress(ctx, userID, filters)
	if err != nil {
		logger.Errorf(ctx, "Database error fetching budgets: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch budgets")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d budgets with progress for user: %s", len(budgetsWithProgress), userID)
	utilities.SuccessResponse(c, budgetsWithProgress, "Budgets retrieved successfully")
}

// GetBudget returns a specific budget by ID
func (h *BudgetHandler) GetBudget(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetBudget - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	budgetIDStr := c.Param("id")
	budgetIDUint, err := strconv.ParseUint(budgetIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "Invalid budget ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid budget ID")
		return
	}
	budgetID := uint(budgetIDUint)

	logger.Debugf(ctx, "Fetching budget: %s", budgetID)
	budget, err := h.service.GetBudget(ctx, budgetID, userID)
	if err != nil {
		logger.Errorf(ctx, "Database error fetching budget: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Budget not found")
		return
	}

	logger.Infof(ctx, "Successfully retrieved budget for user: %s", userID)
	utilities.SuccessResponse(c, budget, "Budget retrieved successfully")
}

// CreateBudget creates a new budget
func (h *BudgetHandler) CreateBudget(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateBudget - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	var budget models.Budget
	if err := c.ShouldBindJSON(&budget); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	budget.UserID = userID

	logger.Debugf(ctx, "Creating budget...")
	created, err := h.service.CreateBudget(ctx, &budget)
	if err != nil {
		logger.Errorf(ctx, "Failed to create budget: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Successfully created budget for user: %s", userID)
	utilities.CreatedResponse(c, created, "Budget created successfully")
}

// UpdateBudget updates an existing budget
func (h *BudgetHandler) UpdateBudget(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateBudget - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	budgetIDStr := c.Param("id")
	budgetIDUint, err := strconv.ParseUint(budgetIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "Invalid budget ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid budget ID")
		return
	}
	budgetID := uint(budgetIDUint)

	var updateData models.Budget
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "Updating budget...")
	updated, err := h.service.UpdateBudget(ctx, budgetID, userID, &updateData)
	if err != nil {
		logger.Errorf(ctx, "Failed to update budget: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Successfully updated budget for user: %s", userID)
	utilities.SuccessResponse(c, updated, "Budget updated successfully")
}

// DeleteBudget deletes a budget
func (h *BudgetHandler) DeleteBudget(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteBudget - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	budgetIDStr := c.Param("id")
	budgetIDUint, err := strconv.ParseUint(budgetIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "Invalid budget ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid budget ID")
		return
	}
	budgetID := uint(budgetIDUint)

	logger.Debugf(ctx, "Deleting budget: %s", budgetID)
	if err := h.service.DeleteBudget(ctx, budgetID, userID); err != nil {
		logger.Errorf(ctx, "Failed to delete budget: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Successfully deleted budget for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Budget deleted successfully")
}

// GetBudgetProgress returns spending progress for a budget
func (h *BudgetHandler) GetBudgetProgress(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetBudgetProgress - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	budgetIDStr := c.Param("id")
	budgetIDUint, err := strconv.ParseUint(budgetIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "Invalid budget ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid budget ID")
		return
	}
	budgetID := uint(budgetIDUint)

	logger.Debugf(ctx, "Calculating budget progress for: %s", budgetID)
	progress, err := h.service.GetBudgetProgress(ctx, budgetID, userID)
	if err != nil {
		logger.Errorf(ctx, "Failed to get budget progress: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Successfully retrieved budget progress for user: %s", userID)
	utilities.SuccessResponse(c, progress, "Budget progress retrieved successfully")
}

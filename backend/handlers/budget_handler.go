package handlers

import (
	"net/http"
	"time"

	"daybook-backend/database"
	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListBudgets returns all budgets for the authenticated user
func ListBudgets(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListBudgets - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	query := database.DB.WithContext(ctx).Where("user_id = ?", userID)

	// Optional filter by enabled status
	if enabled := c.Query("enabled"); enabled != "" {
		query = query.Where("enabled = ?", enabled == "true")
	}

	// Optional filter by period
	if period := c.Query("period"); period != "" {
		query = query.Where("period = ?", period)
	}

	// Optional filter by category
	if categoryID := c.Query("categoryId"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	var budgets []models.Budget
	logger.Debugf(ctx, "Fetching budgets from database...")
	if err := query.Order("created_at DESC").Find(&budgets).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching budgets: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch budgets")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d budgets for user: %s", len(budgets), userID)
	utilities.SuccessResponse(c, budgets, "Budgets retrieved successfully")
}

// GetBudget returns a specific budget by ID
func GetBudget(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetBudget - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid budget ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid budget ID")
		return
	}

	logger.Debugf(ctx, "Fetching budget: %s", budgetID)
	var budget models.Budget
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", budgetID, userID).First(&budget).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching budget: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Budget not found")
		return
	}

	logger.Infof(ctx, "Successfully retrieved budget for user: %s", userID)
	utilities.SuccessResponse(c, budget, "Budget retrieved successfully")
}

// CreateBudget creates a new budget
func CreateBudget(c *gin.Context) {
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

	// Validate custom period dates
	if budget.Period == "custom" {
		if budget.CustomStartDate == nil || budget.CustomEndDate == nil {
			logger.Warnf(ctx, "Validation failed: custom period requires start and end dates")
			utilities.ErrorResponse(c, http.StatusBadRequest, "Custom period requires start and end dates")
			return
		}
	}

	logger.Debugf(ctx, "Creating budget in database...")
	if err := database.DB.WithContext(ctx).Create(&budget).Error; err != nil {
		logger.Errorf(ctx, "Database error creating budget: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create budget")
		return
	}

	// Log budget creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleBudget,
		"Budget", budget.ID, "Created budget for category "+budget.CategoryID, nil)

	logger.Infof(ctx, "Successfully created budget for user: %s", userID)
	utilities.CreatedResponse(c, budget, "Budget created successfully")
}

// UpdateBudget updates an existing budget
func UpdateBudget(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateBudget - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid budget ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid budget ID")
		return
	}

	logger.Debugf(ctx, "Fetching existing budget: %s", budgetID)
	var existingBudget models.Budget
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", budgetID, userID).First(&existingBudget).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching budget: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Budget not found")
		return
	}

	var updateData models.Budget
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update allowed fields
	existingBudget.CategoryID = updateData.CategoryID
	existingBudget.Amount = updateData.Amount
	existingBudget.Period = updateData.Period
	existingBudget.CustomStartDate = updateData.CustomStartDate
	existingBudget.CustomEndDate = updateData.CustomEndDate
	existingBudget.Rollover = updateData.Rollover
	existingBudget.AlertThreshold = updateData.AlertThreshold
	existingBudget.Enabled = updateData.Enabled
	existingBudget.Notes = updateData.Notes

	logger.Debugf(ctx, "Updating budget in database...")
	if err := database.DB.WithContext(ctx).Save(&existingBudget).Error; err != nil {
		logger.Errorf(ctx, "Database error updating budget: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update budget")
		return
	}

	// Log budget update activity
	utilities.LogEntityActivity(c, userID, models.ActionUpdate, models.ModuleBudget,
		"Budget", existingBudget.ID, "Updated budget for category "+existingBudget.CategoryID, nil)

	logger.Infof(ctx, "Successfully updated budget for user: %s", userID)
	utilities.SuccessResponse(c, existingBudget, "Budget updated successfully")
}

// DeleteBudget deletes a budget
func DeleteBudget(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteBudget - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid budget ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid budget ID")
		return
	}

	logger.Debugf(ctx, "Fetching budget to delete: %s", budgetID)
	var budget models.Budget
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", budgetID, userID).First(&budget).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching budget: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Budget not found")
		return
	}

	// Soft delete
	logger.Debugf(ctx, "Deleting budget from database...")
	if err := database.DB.WithContext(ctx).Delete(&budget).Error; err != nil {
		logger.Errorf(ctx, "Database error deleting budget: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete budget")
		return
	}

	// Log budget deletion activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, models.ModuleBudget,
		"Budget", budget.ID, "Deleted budget for category "+budget.CategoryID, nil)

	logger.Infof(ctx, "Successfully deleted budget for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Budget deleted successfully")
}

// GetBudgetProgress returns spending progress for a budget
func GetBudgetProgress(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetBudgetProgress - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid budget ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid budget ID")
		return
	}

	logger.Debugf(ctx, "Fetching budget: %s", budgetID)
	var budget models.Budget
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", budgetID, userID).First(&budget).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching budget: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Budget not found")
		return
	}

	// Calculate date range based on period
	var startDate, endDate time.Time
	now := time.Now()

	switch budget.Period {
	case "weekly":
		// Start of current week (Sunday)
		startDate = now.AddDate(0, 0, -int(now.Weekday()))
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
		endDate = startDate.AddDate(0, 0, 7)

	case "monthly":
		// Start of current month
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, 0)

	case "quarterly":
		// Start of current quarter
		currentMonth := int(now.Month())
		quarterStartMonth := ((currentMonth-1)/3)*3 + 1
		startDate = time.Date(now.Year(), time.Month(quarterStartMonth), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 3, 0)

	case "yearly":
		// Start of current year
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(1, 0, 0)

	case "custom":
		if budget.CustomStartDate != nil && budget.CustomEndDate != nil {
			startDate = *budget.CustomStartDate
			endDate = *budget.CustomEndDate
		} else {
			utilities.ErrorResponse(c, http.StatusBadRequest, "Custom budget dates not set")
			return
		}

	default:
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid budget period")
		return
	}

	// Calculate total spending for the category in the period
	logger.Debugf(ctx, "Calculating budget progress...")
	var totalSpent float64
	database.DB.WithContext(ctx).Model(&models.Transaction{}).
		Where("user_id = ? AND category_id = ? AND type = ? AND date >= ? AND date < ?",
			userID, budget.CategoryID, "expense", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Row().Scan(&totalSpent)

	// Calculate progress
	progress := map[string]interface{}{
		"budget":         budget,
		"totalSpent":     totalSpent,
		"remaining":      budget.Amount - totalSpent,
		"percentageUsed": (totalSpent / budget.Amount) * 100,
		"startDate":      startDate,
		"endDate":        endDate,
		"isOverBudget":   totalSpent > budget.Amount,
		"alertTriggered": (totalSpent / budget.Amount * 100) >= budget.AlertThreshold,
	}

	logger.Infof(ctx, "Successfully retrieved budget progress for user: %s", userID)
	utilities.SuccessResponse(c, progress, "Budget progress retrieved successfully")
}

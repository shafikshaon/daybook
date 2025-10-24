package handlers

import (
	"net/http"
	"time"

	"daybook-backend/database"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListBudgets returns all budgets for the authenticated user
func ListBudgets(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	query := database.DB.Where("user_id = ?", userID)

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
	if err := query.Order("created_at DESC").Find(&budgets).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch budgets")
		return
	}

	utilities.SuccessResponse(c, budgets, "Budgets retrieved successfully")
}

// GetBudget returns a specific budget by ID
func GetBudget(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid budget ID")
		return
	}

	var budget models.Budget
	if err := database.DB.Where("id = ? AND user_id = ?", budgetID, userID).First(&budget).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Budget not found")
		return
	}

	utilities.SuccessResponse(c, budget, "Budget retrieved successfully")
}

// CreateBudget creates a new budget
func CreateBudget(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var budget models.Budget
	if err := c.ShouldBindJSON(&budget); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	budget.UserID = userID

	// Validate custom period dates
	if budget.Period == "custom" {
		if budget.CustomStartDate == nil || budget.CustomEndDate == nil {
			utilities.ErrorResponse(c, http.StatusBadRequest, "Custom period requires start and end dates")
			return
		}
	}

	if err := database.DB.Create(&budget).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create budget")
		return
	}

	utilities.CreatedResponse(c, budget, "Budget created successfully")
}

// UpdateBudget updates an existing budget
func UpdateBudget(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid budget ID")
		return
	}

	var existingBudget models.Budget
	if err := database.DB.Where("id = ? AND user_id = ?", budgetID, userID).First(&existingBudget).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Budget not found")
		return
	}

	var updateData models.Budget
	if err := c.ShouldBindJSON(&updateData); err != nil {
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

	if err := database.DB.Save(&existingBudget).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update budget")
		return
	}

	utilities.SuccessResponse(c, existingBudget, "Budget updated successfully")
}

// DeleteBudget deletes a budget
func DeleteBudget(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid budget ID")
		return
	}

	var budget models.Budget
	if err := database.DB.Where("id = ? AND user_id = ?", budgetID, userID).First(&budget).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Budget not found")
		return
	}

	// Soft delete
	if err := database.DB.Delete(&budget).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete budget")
		return
	}

	utilities.SuccessResponse(c, nil, "Budget deleted successfully")
}

// GetBudgetProgress returns spending progress for a budget
func GetBudgetProgress(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid budget ID")
		return
	}

	var budget models.Budget
	if err := database.DB.Where("id = ? AND user_id = ?", budgetID, userID).First(&budget).Error; err != nil {
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
	var totalSpent float64
	database.DB.Model(&models.Transaction{}).
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

	utilities.SuccessResponse(c, progress, "Budget progress retrieved successfully")
}

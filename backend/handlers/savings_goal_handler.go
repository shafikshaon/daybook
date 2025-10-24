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

// ListSavingsGoals returns all savings goals for the authenticated user
func ListSavingsGoals(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	query := database.DB.Where("user_id = ?", userID)

	// Optional filter by achieved status
	if achieved := c.Query("achieved"); achieved != "" {
		query = query.Where("achieved = ?", achieved == "true")
	}

	// Optional filter by archived status
	if archived := c.Query("archived"); archived != "" {
		query = query.Where("archived = ?", archived == "true")
	}

	// Optional filter by category
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	// Optional filter by priority
	if priority := c.Query("priority"); priority != "" {
		query = query.Where("priority = ?", priority)
	}

	var goals []models.SavingsGoal
	if err := query.Order("priority DESC, created_at DESC").Find(&goals).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch savings goals")
		return
	}

	utilities.SuccessResponse(c, goals, "Savings goals retrieved successfully")
}

// GetSavingsGoal returns a specific savings goal by ID
func GetSavingsGoal(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid savings goal ID")
		return
	}

	var goal models.SavingsGoal
	if err := database.DB.Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Savings goal not found")
		return
	}

	utilities.SuccessResponse(c, goal, "Savings goal retrieved successfully")
}

// CreateSavingsGoal creates a new savings goal
func CreateSavingsGoal(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var goal models.SavingsGoal
	if err := c.ShouldBindJSON(&goal); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	goal.UserID = userID

	if err := database.DB.Create(&goal).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create savings goal")
		return
	}

	utilities.CreatedResponse(c, goal, "Savings goal created successfully")
}

// UpdateSavingsGoal updates an existing savings goal
func UpdateSavingsGoal(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid savings goal ID")
		return
	}

	var existingGoal models.SavingsGoal
	if err := database.DB.Where("id = ? AND user_id = ?", goalID, userID).First(&existingGoal).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Savings goal not found")
		return
	}

	var updateData models.SavingsGoal
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update allowed fields
	existingGoal.Name = updateData.Name
	existingGoal.Description = updateData.Description
	existingGoal.TargetAmount = updateData.TargetAmount
	existingGoal.TargetDate = updateData.TargetDate
	existingGoal.MonthlyContribution = updateData.MonthlyContribution
	existingGoal.Category = updateData.Category
	existingGoal.Priority = updateData.Priority

	if err := database.DB.Save(&existingGoal).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update savings goal")
		return
	}

	utilities.SuccessResponse(c, existingGoal, "Savings goal updated successfully")
}

// DeleteSavingsGoal deletes a savings goal
func DeleteSavingsGoal(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid savings goal ID")
		return
	}

	var goal models.SavingsGoal
	if err := database.DB.Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Savings goal not found")
		return
	}

	// Soft delete
	if err := database.DB.Delete(&goal).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete savings goal")
		return
	}

	utilities.SuccessResponse(c, nil, "Savings goal deleted successfully")
}

// AddContribution adds a contribution to a savings goal
func AddContribution(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid savings goal ID")
		return
	}

	var contributionData struct {
		Amount float64    `json:"amount" binding:"required,gt=0"`
		Date   *time.Time `json:"date"`
		Notes  string     `json:"notes"`
	}

	if err := c.ShouldBindJSON(&contributionData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var goal models.SavingsGoal
	if err := database.DB.Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Savings goal not found")
		return
	}

	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create contribution record
	contribution := models.SavingsContribution{
		UserID: userID,
		GoalID: goalID,
		Amount: contributionData.Amount,
		Date: func() time.Time {
			if contributionData.Date != nil {
				return *contributionData.Date
			}
			return time.Now()
		}(),
		Notes: contributionData.Notes,
	}

	if err := tx.Create(&contribution).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create contribution")
		return
	}

	// Update goal
	goal.CurrentAmount += contributionData.Amount
	goal.LastContribution = contributionData.Amount
	goal.LastContributionDate = &contribution.Date

	// Check if goal is achieved
	if goal.CurrentAmount >= goal.TargetAmount && !goal.Achieved {
		goal.Achieved = true
		now := time.Now()
		goal.AchievedDate = &now
	}

	if err := tx.Save(&goal).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update goal")
		return
	}

	tx.Commit()

	result := map[string]interface{}{
		"goal":         goal,
		"contribution": contribution,
	}

	utilities.SuccessResponse(c, result, "Contribution added successfully")
}

// WithdrawFromGoal withdraws from a savings goal
func WithdrawFromGoal(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid savings goal ID")
		return
	}

	var withdrawalData struct {
		Amount float64    `json:"amount" binding:"required,gt=0"`
		Date   *time.Time `json:"date"`
		Notes  string     `json:"notes"`
	}

	if err := c.ShouldBindJSON(&withdrawalData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var goal models.SavingsGoal
	if err := database.DB.Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Savings goal not found")
		return
	}

	// Verify sufficient funds
	if withdrawalData.Amount > goal.CurrentAmount {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Insufficient funds in savings goal")
		return
	}

	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create withdrawal record (as negative contribution)
	withdrawal := models.SavingsContribution{
		UserID: userID,
		GoalID: goalID,
		Amount: -withdrawalData.Amount, // Negative amount indicates withdrawal
		Date: func() time.Time {
			if withdrawalData.Date != nil {
				return *withdrawalData.Date
			}
			return time.Now()
		}(),
		Notes: withdrawalData.Notes,
	}

	if err := tx.Create(&withdrawal).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create withdrawal record")
		return
	}

	// Update goal
	goal.CurrentAmount -= withdrawalData.Amount

	// If goal was achieved but now isn't, mark as not achieved
	if goal.CurrentAmount < goal.TargetAmount && goal.Achieved {
		goal.Achieved = false
		goal.AchievedDate = nil
	}

	if err := tx.Save(&goal).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update goal")
		return
	}

	tx.Commit()

	result := map[string]interface{}{
		"goal":       goal,
		"withdrawal": withdrawal,
	}

	utilities.SuccessResponse(c, result, "Withdrawal completed successfully")
}

// ListAutomatedRules returns automated rules for the authenticated user
func ListAutomatedRules(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	query := database.DB.Where("user_id = ?", userID)

	// Optional filter by goal
	if goalID := c.Query("goalId"); goalID != "" {
		query = query.Where("goal_id = ?", goalID)
	}

	// Optional filter by enabled status
	if enabled := c.Query("enabled"); enabled != "" {
		query = query.Where("enabled = ?", enabled == "true")
	}

	var rules []models.AutomatedRule
	if err := query.Order("created_at DESC").Find(&rules).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch automated rules")
		return
	}

	utilities.SuccessResponse(c, rules, "Automated rules retrieved successfully")
}

// CreateAutomatedRule creates a new automated rule
func CreateAutomatedRule(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var rule models.AutomatedRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	rule.UserID = userID

	// Verify goal belongs to user
	var goal models.SavingsGoal
	if err := database.DB.Where("id = ? AND user_id = ?", rule.GoalID, userID).First(&goal).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid savings goal ID")
		return
	}

	if err := database.DB.Create(&rule).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create automated rule")
		return
	}

	utilities.CreatedResponse(c, rule, "Automated rule created successfully")
}

package handlers

import (
	"net/http"
	"strconv"
	"time"

	"daybook-backend/database"
	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
)

// ListSavingsGoals returns all savings goals for the authenticated user
func ListSavingsGoals(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListSavingsGoals - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Fetching savings goals for user")
	query := database.DB.WithContext(ctx).Where("user_id = ?", userID)

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
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch savings goals")
		return
	}

	logger.Infof(ctx, "Savings goals retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, goals, "Savings goals retrieved successfully")
}

// GetSavingsGoal returns a specific savings goal by ID
func GetSavingsGoal(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetSavingsGoal - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goalIDStr := c.Param("id")
	goalIDUint, err := strconv.ParseUint(goalIDStr, 10, 32)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid savings goal ID")
		return
	}
	goalID := uint(goalIDUint)

	logger.Debugf(ctx, "Fetching savings goal with ID: %s", goalID)
	var goal models.SavingsGoal
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Savings goal not found")
		return
	}

	logger.Infof(ctx, "Savings goal retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, goal, "Savings goal retrieved successfully")
}

// CreateSavingsGoal creates a new savings goal
func CreateSavingsGoal(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateSavingsGoal - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var goal models.SavingsGoal
	if err := c.ShouldBindJSON(&goal); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	goal.UserID = userID

	logger.Debugf(ctx, "Creating savings goal: %s", goal.Name)
	if err := database.DB.WithContext(ctx).Create(&goal).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create savings goal")
		return
	}

	// Log savings goal creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleGoal,
		"SavingsGoal", goal.ID, "Created savings goal: "+goal.Name, nil)

	logger.Infof(ctx, "Savings goal created successfully for user: %s", userID)
	utilities.CreatedResponse(c, goal, "Savings goal created successfully")
}

// UpdateSavingsGoal updates an existing savings goal
func UpdateSavingsGoal(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateSavingsGoal - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goalIDStr := c.Param("id")
	goalIDUint, err := strconv.ParseUint(goalIDStr, 10, 32)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid savings goal ID")
		return
	}
	goalID := uint(goalIDUint)

	logger.Debugf(ctx, "Fetching existing savings goal with ID: %s", goalID)
	var existingGoal models.SavingsGoal
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", goalID, userID).First(&existingGoal).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
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

	logger.Debugf(ctx, "Updating savings goal: %s", existingGoal.Name)
	if err := database.DB.WithContext(ctx).Save(&existingGoal).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update savings goal")
		return
	}

	// Log savings goal update activity
	utilities.LogEntityActivity(c, userID, models.ActionUpdate, models.ModuleGoal,
		"SavingsGoal", existingGoal.ID, "Updated savings goal: "+existingGoal.Name, nil)

	logger.Infof(ctx, "Savings goal updated successfully for user: %s", userID)
	utilities.SuccessResponse(c, existingGoal, "Savings goal updated successfully")
}

// DeleteSavingsGoal deletes a savings goal
func DeleteSavingsGoal(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteSavingsGoal - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goalIDStr := c.Param("id")
	goalIDUint, err := strconv.ParseUint(goalIDStr, 10, 32)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid savings goal ID")
		return
	}
	goalID := uint(goalIDUint)

	logger.Debugf(ctx, "Fetching savings goal with ID: %s for deletion", goalID)
	var goal models.SavingsGoal
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Savings goal not found")
		return
	}

	// Soft delete
	logger.Debugf(ctx, "Deleting savings goal: %s", goal.Name)
	if err := database.DB.WithContext(ctx).Delete(&goal).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete savings goal")
		return
	}

	// Log savings goal deletion activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, models.ModuleGoal,
		"SavingsGoal", goal.ID, "Deleted savings goal: "+goal.Name, nil)

	logger.Infof(ctx, "Savings goal deleted successfully for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Savings goal deleted successfully")
}

// AddContribution adds a contribution to a savings goal
func AddContribution(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "AddContribution - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goalIDStr := c.Param("id")
	goalIDUint, err := strconv.ParseUint(goalIDStr, 10, 32)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid savings goal ID")
		return
	}
	goalID := uint(goalIDUint)

	var contributionData struct {
		Amount    float64    `json:"amount" binding:"required,gt=0"`
		AccountID uint       `json:"accountId" binding:"required"`
		Date      *time.Time `json:"date"`
		Notes     string     `json:"notes"`
	}

	if err := c.ShouldBindJSON(&contributionData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "Fetching savings goal with ID: %s", goalID)
	var goal models.SavingsGoal
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Savings goal not found")
		return
	}

	// Verify account belongs to user
	logger.Debugf(ctx, "Verifying account ownership for account ID: %s", contributionData.AccountID)
	var account models.Account
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", contributionData.AccountID, userID).First(&account).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	// Start transaction
	logger.Debugf(ctx, "Starting transaction for contribution")
	tx := database.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	contributionDate := func() time.Time {
		if contributionData.Date != nil {
			return *contributionData.Date
		}
		return time.Now()
	}()

	// Create contribution record
	contribution := models.SavingsContribution{
		UserID: userID,
		GoalID: goalID,
		Amount: contributionData.Amount,
		Date:   contributionDate,
		Notes:  contributionData.Notes,
	}

	logger.Debugf(ctx, "Creating contribution record")
	if err := tx.Create(&contribution).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create contribution")
		return
	}

	// Create transaction record to track the expense
	transaction := models.Transaction{
		UserID:        userID,
		AccountID:     contributionData.AccountID,
		Type:          "expense",
		Amount:        contributionData.Amount,
		CategoryID:    0, // Use 0 for system transactions
		Date:          models.Date{Time: contributionDate},
		Description:   "Contribution to " + goal.Name,
		SavingsGoalID: &goalID,
		Tags:          []string{"savings_goal"},
	}

	logger.Debugf(ctx, "Creating transaction record")
	if err := tx.Create(&transaction).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	// Update account balance (debit)
	account.Balance -= contributionData.Amount
	logger.Debugf(ctx, "Updating account balance")
	if err := tx.Save(&account).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account balance")
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

	logger.Debugf(ctx, "Updating savings goal")
	if err := tx.Save(&goal).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update goal")
		return
	}

	tx.Commit()

	// Log contribution addition activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleGoal,
		"SavingsContribution", contribution.ID, "Added contribution to savings goal: "+goal.Name, nil)

	result := map[string]interface{}{
		"goal":         goal,
		"contribution": contribution,
		"transaction":  transaction,
	}

	logger.Infof(ctx, "Contribution added successfully for user: %s", userID)
	utilities.SuccessResponse(c, result, "Contribution added successfully")
}

// WithdrawFromGoal withdraws from a savings goal
func WithdrawFromGoal(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "WithdrawFromGoal - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goalIDStr := c.Param("id")
	goalIDUint, err := strconv.ParseUint(goalIDStr, 10, 32)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid savings goal ID")
		return
	}
	goalID := uint(goalIDUint)

	var withdrawalData struct {
		Amount    float64    `json:"amount" binding:"required,gt=0"`
		AccountID uint       `json:"accountId" binding:"required"`
		Date      *time.Time `json:"date"`
		Notes     string     `json:"notes"`
	}

	if err := c.ShouldBindJSON(&withdrawalData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "Fetching savings goal with ID: %s", goalID)
	var goal models.SavingsGoal
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Savings goal not found")
		return
	}

	// Verify sufficient funds
	if withdrawalData.Amount > goal.CurrentAmount {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Insufficient funds in savings goal")
		return
	}

	// Verify account belongs to user
	logger.Debugf(ctx, "Verifying account ownership for account ID: %s", withdrawalData.AccountID)
	var account models.Account
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", withdrawalData.AccountID, userID).First(&account).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	// Start transaction
	logger.Debugf(ctx, "Starting transaction for withdrawal")
	tx := database.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	withdrawalDate := func() time.Time {
		if withdrawalData.Date != nil {
			return *withdrawalData.Date
		}
		return time.Now()
	}()

	// Create withdrawal record (as negative contribution)
	withdrawal := models.SavingsContribution{
		UserID: userID,
		GoalID: goalID,
		Amount: -withdrawalData.Amount, // Negative amount indicates withdrawal
		Date:   withdrawalDate,
		Notes:  withdrawalData.Notes,
	}

	logger.Debugf(ctx, "Creating withdrawal record")
	if err := tx.Create(&withdrawal).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create withdrawal record")
		return
	}

	// Create transaction record to track the income (money coming back)
	transaction := models.Transaction{
		UserID:        userID,
		AccountID:     withdrawalData.AccountID,
		Type:          "income",
		Amount:        withdrawalData.Amount,
		CategoryID:    0, // Use 0 for system transactions
		Date:          models.Date{Time: withdrawalDate},
		Description:   "Withdrawal from " + goal.Name,
		SavingsGoalID: &goalID,
		Tags:          []string{"savings_goal"},
	}

	logger.Debugf(ctx, "Creating transaction record")
	if err := tx.Create(&transaction).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	// Update account balance (credit)
	account.Balance += withdrawalData.Amount
	logger.Debugf(ctx, "Updating account balance")
	if err := tx.Save(&account).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account balance")
		return
	}

	// Update goal
	goal.CurrentAmount -= withdrawalData.Amount

	// If goal was achieved but now isn't, mark as not achieved
	if goal.CurrentAmount < goal.TargetAmount && goal.Achieved {
		goal.Achieved = false
		goal.AchievedDate = nil
	}

	logger.Debugf(ctx, "Updating savings goal")
	if err := tx.Save(&goal).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update goal")
		return
	}

	tx.Commit()

	// Log withdrawal activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, models.ModuleGoal,
		"SavingsContribution", withdrawal.ID, "Withdrew from savings goal: "+goal.Name, nil)

	result := map[string]interface{}{
		"goal":        goal,
		"withdrawal":  withdrawal,
		"transaction": transaction,
	}

	logger.Infof(ctx, "Withdrawal completed successfully for user: %s", userID)
	utilities.SuccessResponse(c, result, "Withdrawal completed successfully")
}

// ListAutomatedRules returns automated rules for the authenticated user
func ListAutomatedRules(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListAutomatedRules - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Fetching automated rules for user")
	query := database.DB.WithContext(ctx).Where("user_id = ?", userID)

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
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch automated rules")
		return
	}

	logger.Infof(ctx, "Automated rules retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, rules, "Automated rules retrieved successfully")
}

// CreateAutomatedRule creates a new automated rule
func CreateAutomatedRule(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateAutomatedRule - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
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
	logger.Debugf(ctx, "Verifying savings goal ownership")
	var goal models.SavingsGoal
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", rule.GoalID, userID).First(&goal).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid savings goal ID")
		return
	}

	logger.Debugf(ctx, "Creating automated rule for goal: %s", goal.Name)
	if err := database.DB.WithContext(ctx).Create(&rule).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create automated rule")
		return
	}

	// Log automated rule creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleGoal,
		"AutomatedRule", rule.ID, "Created automated rule for savings goal: "+goal.Name, nil)

	logger.Infof(ctx, "Automated rule created successfully for user: %s", userID)
	utilities.CreatedResponse(c, rule, "Automated rule created successfully")
}

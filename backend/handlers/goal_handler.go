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

// ListGoals returns all goals for the authenticated user
func ListGoals(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	query := database.DB.Where("user_id = ?", userID)

	// Optional filters
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	if priority := c.Query("priority"); priority != "" {
		query = query.Where("priority = ?", priority)
	}

	var goals []models.Goal
	if err := query.Preload("Holdings").Preload("Contributions").Order("priority DESC, created_at DESC").Find(&goals).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch goals")
		return
	}

	// Calculate progress for each goal
	for i := range goals {
		goals[i].UpdateCurrentAmount(database.DB)
	}

	utilities.SuccessResponse(c, goals, "Goals retrieved successfully")
}

// GetGoal returns a specific goal by ID with all details
func GetGoal(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid goal ID")
		return
	}

	var goal models.Goal
	if err := database.DB.Where("id = ? AND user_id = ?", goalID, userID).
		Preload("Holdings").
		Preload("Contributions").
		First(&goal).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Goal not found")
		return
	}

	// Update current amount
	goal.UpdateCurrentAmount(database.DB)

	utilities.SuccessResponse(c, goal, "Goal retrieved successfully")
}

// CreateGoal creates a new goal
func CreateGoal(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var goal models.Goal
	if err := c.ShouldBindJSON(&goal); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	goal.UserID = userID
	goal.Status = models.GoalStatusActive
	goal.CurrentAmount = 0

	if err := database.DB.Create(&goal).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create goal")
		return
	}

	utilities.CreatedResponse(c, goal, "Goal created successfully")
}

// UpdateGoal updates an existing goal
func UpdateGoal(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid goal ID")
		return
	}

	var existingGoal models.Goal
	if err := database.DB.Where("id = ? AND user_id = ?", goalID, userID).First(&existingGoal).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Goal not found")
		return
	}

	var updateData models.Goal
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update allowed fields
	existingGoal.Name = updateData.Name
	existingGoal.Description = updateData.Description
	existingGoal.Icon = updateData.Icon
	existingGoal.Color = updateData.Color
	existingGoal.Category = updateData.Category
	existingGoal.Priority = updateData.Priority
	existingGoal.TargetAmount = updateData.TargetAmount
	existingGoal.TargetDate = updateData.TargetDate
	existingGoal.MonthlyContribution = updateData.MonthlyContribution
	existingGoal.Status = updateData.Status

	if err := database.DB.Save(&existingGoal).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update goal")
		return
	}

	utilities.SuccessResponse(c, existingGoal, "Goal updated successfully")
}

// DeleteGoal deletes a goal
func DeleteGoal(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid goal ID")
		return
	}

	var goal models.Goal
	if err := database.DB.Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Goal not found")
		return
	}

	// Soft delete
	if err := database.DB.Delete(&goal).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete goal")
		return
	}

	utilities.SuccessResponse(c, nil, "Goal deleted successfully")
}

// AddHolding adds a new holding to a goal
func AddHolding(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid goal ID")
		return
	}

	var holdingData struct {
		models.GoalHolding
		AccountID uuid.UUID `json:"accountId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&holdingData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Verify goal belongs to user
	var goal models.Goal
	if err := database.DB.Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Goal not found")
		return
	}

	// Verify account belongs to user
	var account models.Account
	if err := database.DB.Where("id = ? AND user_id = ?", holdingData.AccountID, userID).First(&account).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	// Check sufficient balance
	if account.Balance < holdingData.Amount {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Insufficient account balance")
		return
	}

	holdingData.GoalHolding.UserID = userID
	holdingData.GoalHolding.GoalID = goalID
	holdingData.GoalHolding.Status = models.HoldingStatusActive

	// Set current value to initial amount if not set
	if holdingData.CurrentValue == 0 {
		holdingData.CurrentValue = holdingData.Amount
	}

	// Update market value for market instruments
	holdingData.GoalHolding.UpdateMarketValue()

	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create holding
	if err := tx.Create(&holdingData.GoalHolding).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create holding")
		return
	}

	// Create transaction record
	categoryID := "goal_holding_added"

	transaction := models.Transaction{
		UserID:      userID,
		AccountID:   holdingData.AccountID,
		Type:        "expense",
		Amount:      holdingData.Amount,
		CategoryID:  categoryID,
		Date:        holdingData.PurchaseDate,
		Description: "Added to " + goal.Name + ": " + holdingData.Name,
		Tags:        []string{"goal", "holding"},
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	holdingData.TransactionID = transaction.ID

	// Update account balance
	account.Balance -= holdingData.Amount
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account balance")
		return
	}

	// Create contribution record
	contribution := models.GoalContribution{
		UserID:        userID,
		GoalID:        goalID,
		HoldingID:     &holdingData.GoalHolding.ID,
		Type:          models.ContributionTypeContribution,
		Amount:        holdingData.Amount,
		Date:          holdingData.PurchaseDate,
		Notes:         "Added " + holdingData.Name,
		TransactionID: transaction.ID,
	}

	if err := tx.Create(&contribution).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create contribution record")
		return
	}

	// Update goal
	goal.UpdateCurrentAmount(tx)
	goal.LastContribution = holdingData.Amount
	goal.LastContributionDate = &holdingData.PurchaseDate

	// Check if goal is achieved
	if goal.CurrentAmount >= goal.TargetAmount && !goal.Achieved {
		goal.Achieved = true
		now := time.Now()
		goal.AchievedDate = &now
		goal.Status = models.GoalStatusAchieved
	}

	if err := tx.Save(&goal).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update goal")
		return
	}

	tx.Commit()

	result := map[string]interface{}{
		"holding":      holdingData.GoalHolding,
		"contribution": contribution,
		"transaction":  transaction,
		"goal":         goal,
	}

	utilities.CreatedResponse(c, result, "Holding added successfully")
}

// UpdateHolding updates a holding (e.g., update stock price)
func UpdateHolding(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	holdingID, err := uuid.Parse(c.Param("holdingId"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid holding ID")
		return
	}

	var existingHolding models.GoalHolding
	if err := database.DB.Where("id = ? AND user_id = ?", holdingID, userID).First(&existingHolding).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Holding not found")
		return
	}

	var updateData models.GoalHolding
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update fields
	if updateData.CurrentPrice != nil {
		existingHolding.CurrentPrice = updateData.CurrentPrice
	}
	if updateData.CurrentValue > 0 {
		existingHolding.CurrentValue = updateData.CurrentValue
	}
	if updateData.Status != "" {
		existingHolding.Status = updateData.Status
	}

	// Recalculate market value
	existingHolding.UpdateMarketValue()

	if err := database.DB.Save(&existingHolding).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update holding")
		return
	}

	// Update goal's current amount
	var goal models.Goal
	if err := database.DB.First(&goal, existingHolding.GoalID).Error; err == nil {
		goal.UpdateCurrentAmount(database.DB)
	}

	utilities.SuccessResponse(c, existingHolding, "Holding updated successfully")
}

// RemoveHolding removes/liquidates a holding
func RemoveHolding(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	holdingID, err := uuid.Parse(c.Param("holdingId"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid holding ID")
		return
	}

	var removeData struct {
		AccountID    uuid.UUID `json:"accountId" binding:"required"`
		CurrentValue float64   `json:"currentValue" binding:"required,gt=0"`
		Date         time.Time `json:"date"`
		Notes        string    `json:"notes"`
	}

	if err := c.ShouldBindJSON(&removeData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var holding models.GoalHolding
	if err := database.DB.Where("id = ? AND user_id = ?", holdingID, userID).First(&holding).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Holding not found")
		return
	}

	// Verify account
	var account models.Account
	if err := database.DB.Where("id = ? AND user_id = ?", removeData.AccountID, userID).First(&account).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	if removeData.Date.IsZero() {
		removeData.Date = time.Now()
	}

	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Mark holding as sold/closed
	holding.Status = models.HoldingStatusSold
	holding.CurrentValue = removeData.CurrentValue
	if err := tx.Save(&holding).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update holding")
		return
	}

	// Create transaction (income as money returns)
	categoryID := "goal_holding_removed"

	transaction := models.Transaction{
		UserID:      userID,
		AccountID:   removeData.AccountID,
		Type:        "income",
		Amount:      removeData.CurrentValue,
		CategoryID:  categoryID,
		Date:        removeData.Date,
		Description: "Sold/Closed: " + holding.Name,
		Tags:        []string{"goal", "holding", "liquidation"},
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	// Credit account
	account.Balance += removeData.CurrentValue
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account balance")
		return
	}

	// Create contribution record
	contribution := models.GoalContribution{
		UserID:        userID,
		GoalID:        holding.GoalID,
		HoldingID:     &holding.ID,
		Type:          models.ContributionTypeWithdrawal,
		Amount:        removeData.CurrentValue,
		Date:          removeData.Date,
		Notes:         removeData.Notes,
		TransactionID: transaction.ID,
	}

	if err := tx.Create(&contribution).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create contribution record")
		return
	}

	// Update goal
	var goal models.Goal
	if err := tx.First(&goal, holding.GoalID).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch goal")
		return
	}

	goal.UpdateCurrentAmount(tx)

	tx.Commit()

	result := map[string]interface{}{
		"holding":      holding,
		"transaction":  transaction,
		"contribution": contribution,
	}

	utilities.SuccessResponse(c, result, "Holding removed successfully")
}

// GetHoldingTypes returns all available holding types
func GetHoldingTypes(c *gin.Context) {
	holdingTypes := map[string]interface{}{
		"Savings": []map[string]string{
			{"value": "savings", "label": "Savings", "icon": "💰"},
			{"value": "fixed_deposit", "label": "Fixed Deposit", "icon": "🏦"},
			{"value": "dps", "label": "DPS (Deposit Pension Scheme)", "icon": "📅"},
			{"value": "recurring_deposit", "label": "Recurring Deposit", "icon": "🔄"},
			{"value": "savings_bond", "label": "Savings Bond", "icon": "🎫"},
			{"value": "ppf", "label": "PPF (Public Provident Fund)", "icon": "🏛️"},
			{"value": "nsc", "label": "NSC (National Savings Certificate)", "icon": "📄"},
		},
		"Investments": []map[string]string{
			{"value": "stocks", "label": "Stocks", "icon": "📈"},
			{"value": "mutual_fund", "label": "Mutual Fund", "icon": "🏛️"},
			{"value": "etf", "label": "ETF", "icon": "📊"},
			{"value": "index_fund", "label": "Index Fund", "icon": "📉"},
			{"value": "bonds", "label": "Bonds", "icon": "📜"},
			{"value": "cryptocurrency", "label": "Cryptocurrency", "icon": "₿"},
		},
		"Alternatives": []map[string]string{
			{"value": "real_estate", "label": "Real Estate", "icon": "🏢"},
			{"value": "reit", "label": "REIT", "icon": "🏗️"},
			{"value": "gold", "label": "Gold", "icon": "🥇"},
			{"value": "commodities", "label": "Commodities", "icon": "🛢️"},
		},
		"Retirement": []map[string]string{
			{"value": "pension_fund", "label": "Pension Fund", "icon": "👴"},
			{"value": "retirement_401k", "label": "401(k) / Retirement", "icon": "🏦"},
			{"value": "provident_fund", "label": "Provident Fund (EPF)", "icon": "💼"},
		},
		"Insurance": []map[string]string{
			{"value": "life_insurance", "label": "Life Insurance", "icon": "🛡️"},
			{"value": "ulip", "label": "ULIP", "icon": "🔗"},
		},
		"Other": []map[string]string{
			{"value": "custom", "label": "Custom Investment", "icon": "💎"},
		},
	}

	utilities.SuccessResponse(c, holdingTypes, "Holding types retrieved successfully")
}

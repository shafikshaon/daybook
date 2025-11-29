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

// ListGoals returns all goals for the authenticated user
func ListGoals(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "ListGoals - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "ListGoals - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	logger.Debugf(ctx, "ListGoals - Fetching goals for user: %s", userID)

	query := database.DB.WithContext(ctx).Where("user_id = ?", userID)

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
	if err := query.Preload("Holdings").Preload("Contributions").Order("name DESC, created_at DESC").Find(&goals).Error; err != nil {
		logger.Errorf(ctx, "ListGoals - Failed to fetch goals: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch goals")
		return
	}

	logger.Debugf(ctx, "ListGoals - Retrieved %d goals", len(goals))

	// Calculate progress for each goal
	for i := range goals {
		goals[i].UpdateCurrentAmount(database.DB.WithContext(ctx))
	}

	logger.Infof(ctx, "ListGoals - Successfully retrieved goals")
	utilities.SuccessResponse(c, goals, "Goals retrieved successfully")
}

// GetGoal returns a specific goal by ID with all details
func GetGoal(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "GetGoal - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "GetGoal - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "GetGoal - Invalid goal ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid goal ID")
		return
	}

	logger.Debugf(ctx, "GetGoal - Fetching goal: %s", goalID)

	var goal models.Goal
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", goalID, userID).
		Preload("Holdings").
		Preload("Contributions").
		First(&goal).Error; err != nil {
		logger.Errorf(ctx, "GetGoal - Failed to fetch goal: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Goal not found")
		return
	}

	// Update current amount
	goal.UpdateCurrentAmount(database.DB.WithContext(ctx))

	logger.Infof(ctx, "GetGoal - Successfully retrieved goal: %s", goalID)
	utilities.SuccessResponse(c, goal, "Goal retrieved successfully")
}

// CreateGoal creates a new goal
func CreateGoal(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "CreateGoal - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "CreateGoal - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	var goal models.Goal
	if err := c.ShouldBindJSON(&goal); err != nil {
		logger.Warnf(ctx, "CreateGoal - Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "CreateGoal - Creating goal: %s", goal.Name)

	goal.UserID = userID
	goal.Status = models.GoalStatusActive
	goal.CurrentAmount = 0

	if err := database.DB.WithContext(ctx).Create(&goal).Error; err != nil {
		logger.Errorf(ctx, "CreateGoal - Failed to create goal: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create goal")
		return
	}

	// Log goal creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleGoal,
		"Goal", goal.ID, "Created goal: "+goal.Name, nil)

	logger.Infof(ctx, "CreateGoal - Successfully created goal: %s", goal.ID)
	utilities.CreatedResponse(c, goal, "Goal created successfully")
}

// UpdateGoal updates an existing goal
func UpdateGoal(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "UpdateGoal - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "UpdateGoal - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "UpdateGoal - Invalid goal ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid goal ID")
		return
	}

	logger.Debugf(ctx, "UpdateGoal - Updating goal: %s", goalID)

	var existingGoal models.Goal
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", goalID, userID).First(&existingGoal).Error; err != nil {
		logger.Errorf(ctx, "UpdateGoal - Goal not found: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Goal not found")
		return
	}

	var updateData models.Goal
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "UpdateGoal - Invalid request body: %v", err)
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

	if err := database.DB.WithContext(ctx).Save(&existingGoal).Error; err != nil {
		logger.Errorf(ctx, "UpdateGoal - Failed to update goal: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update goal")
		return
	}

	// Log goal update activity
	utilities.LogEntityActivity(c, userID, models.ActionUpdate, models.ModuleGoal,
		"Goal", existingGoal.ID, "Updated goal: "+existingGoal.Name, nil)

	logger.Infof(ctx, "UpdateGoal - Successfully updated goal: %s", goalID)
	utilities.SuccessResponse(c, existingGoal, "Goal updated successfully")
}

// DeleteGoal deletes a goal
func DeleteGoal(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "DeleteGoal - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "DeleteGoal - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "DeleteGoal - Invalid goal ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid goal ID")
		return
	}

	logger.Debugf(ctx, "DeleteGoal - Deleting goal: %s", goalID)

	var goal models.Goal
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		logger.Errorf(ctx, "DeleteGoal - Goal not found: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Goal not found")
		return
	}

	// Soft delete
	if err := database.DB.WithContext(ctx).Delete(&goal).Error; err != nil {
		logger.Errorf(ctx, "DeleteGoal - Failed to delete goal: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete goal")
		return
	}

	// Log goal deletion activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, models.ModuleGoal,
		"Goal", goal.ID, "Deleted goal: "+goal.Name, nil)

	logger.Infof(ctx, "DeleteGoal - Successfully deleted goal: %s", goalID)
	utilities.SuccessResponse(c, nil, "Goal deleted successfully")
}

// AddHolding adds a new holding to a goal
func AddHolding(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "AddHolding - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "AddHolding - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "AddHolding - Invalid goal ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid goal ID")
		return
	}

	logger.Debugf(ctx, "AddHolding - Adding holding to goal: %s", goalID)

	var holdingData struct {
		models.GoalHolding
		AccountID  *uuid.UUID `json:"accountId"`  // Pointer to allow null for existing investments
		IsExisting bool       `json:"isExisting"` // Flag for existing/external investments
	}

	if err := c.ShouldBindJSON(&holdingData); err != nil {
		logger.Warnf(ctx, "AddHolding - Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Verify goal belongs to user
	var goal models.Goal
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		logger.Errorf(ctx, "AddHolding - Goal not found: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Goal not found")
		return
	}

	// For existing investments, we don't need to verify account or check balance
	var account models.Account
	if !holdingData.IsExisting {
		// Verify account belongs to user
		if holdingData.AccountID == nil || *holdingData.AccountID == uuid.Nil {
			logger.Warnf(ctx, "AddHolding - Account ID is required for new investments")
			utilities.ErrorResponse(c, http.StatusBadRequest, "Account ID is required for new investments")
			return
		}

		if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", *holdingData.AccountID, userID).First(&account).Error; err != nil {
			logger.Errorf(ctx, "AddHolding - Invalid account ID: %v", err)
			utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
			return
		}

		// Check sufficient balance
		if account.Balance < holdingData.Amount {
			logger.Warnf(ctx, "AddHolding - Insufficient account balance")
			utilities.ErrorResponse(c, http.StatusBadRequest, "Insufficient account balance")
			return
		}
	}

	holdingData.GoalHolding.UserID = userID
	holdingData.GoalHolding.GoalID = goalID
	holdingData.GoalHolding.Status = models.HoldingStatusActive

	// IMPORTANT: Always set currentValue to amount initially
	// The frontend should send this, but ensure it's set
	if holdingData.CurrentValue == 0 {
		holdingData.CurrentValue = holdingData.Amount
	}

	// For market instruments, calculate value based on quantity and price
	// For others (FD/DPS), currentValue should equal amount
	holdingData.GoalHolding.UpdateMarketValue()

	logger.Debugf(ctx, "AddHolding - Starting transaction to create holding")

	// Start transaction
	tx := database.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create holding
	if err := tx.Create(&holdingData.GoalHolding).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "AddHolding - Failed to create holding: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create holding")
		return
	}

	logger.Debugf(ctx, "AddHolding - Created holding: %s", holdingData.GoalHolding.ID)

	// Create transaction record
	var transaction models.Transaction

	if holdingData.IsExisting {
		// For existing investments, create a special "tracking" type transaction
		// This won't show in regular expense/income lists but tracks the holding for reference
		var trackingAccountID uuid.UUID
		if holdingData.AccountID != nil {
			trackingAccountID = *holdingData.AccountID
		}

		transaction = models.Transaction{
			UserID:      userID,
			AccountID:   trackingAccountID,
			Type:        "tracking", // Special type that won't appear in regular transaction lists
			Amount:      holdingData.Amount,
			CategoryID:  "goal_external_holding",
			Date:        holdingData.PurchaseDate,
			Description: "External " + holdingData.Name + " tracked for " + goal.Name,
			Tags:        []string{"goal", "holding", "external", "tracking", "hidden"},
		}
	} else {
		// For new investments, create normal transaction
		transaction = models.Transaction{
			UserID:      userID,
			AccountID:   *holdingData.AccountID,
			Type:        "expense",
			Amount:      holdingData.Amount,
			CategoryID:  "goal_holding_added",
			Date:        holdingData.PurchaseDate,
			Description: "Added to " + goal.Name + ": " + holdingData.Name,
			Tags:        []string{"goal", "holding"},
		}
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	holdingData.TransactionID = transaction.ID

	// Update account balance only for new investments (not existing ones)
	if !holdingData.IsExisting {
		account.Balance -= holdingData.Amount
		if err := tx.Save(&account).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account balance")
			return
		}
	}

	// Create contribution record
	contributionNotes := "Added " + holdingData.Name
	if holdingData.IsExisting {
		contributionNotes = "External holding: " + holdingData.Name
	}

	contribution := models.GoalContribution{
		UserID:        userID,
		GoalID:        goalID,
		HoldingID:     &holdingData.GoalHolding.ID,
		Type:          models.ContributionTypeContribution,
		Amount:        holdingData.Amount,
		Date:          holdingData.PurchaseDate,
		Notes:         contributionNotes,
		TransactionID: transaction.ID,
	}

	if err := tx.Create(&contribution).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create contribution record")
		return
	}

	// Update goal metadata (but not currentAmount yet - we'll do that after commit)
	goal.LastContribution = holdingData.Amount
	goal.LastContributionDate = &holdingData.PurchaseDate

	if err := tx.Save(&goal).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update goal")
		return
	}

	// Commit the transaction first
	if err := tx.Commit().Error; err != nil {
		logger.Errorf(ctx, "AddHolding - Failed to commit transaction: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	logger.Debugf(ctx, "AddHolding - Transaction committed successfully")

	// NOW update currentAmount after everything is committed
	// This ensures the holding is properly saved in the database

	err = goal.UpdateCurrentAmount(database.DB.WithContext(ctx))
	if err != nil {
		logger.Warnf(ctx, "AddHolding - Error updating current amount: %v", err)
	}
	// Reload goal with all relations to return complete data
	var updatedGoal models.Goal
	if err := database.DB.WithContext(ctx).Where("id = ?", goalID).
		Preload("Holdings").
		Preload("Contributions").
		First(&updatedGoal).Error; err == nil {
		goal = updatedGoal
	}

	// Log holding addition activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleGoal,
		"GoalHolding", holdingData.GoalHolding.ID, "Added holding "+holdingData.Name+" to goal "+goal.Name, nil)

	logger.Infof(ctx, "AddHolding - Successfully added holding to goal: %s", goalID)

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
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "UpdateHolding - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "UpdateHolding - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	holdingID, err := uuid.Parse(c.Param("holdingId"))
	if err != nil {
		logger.Warnf(ctx, "UpdateHolding - Invalid holding ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid holding ID")
		return
	}

	logger.Debugf(ctx, "UpdateHolding - Updating holding: %s", holdingID)

	var existingHolding models.GoalHolding
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", holdingID, userID).First(&existingHolding).Error; err != nil {
		logger.Errorf(ctx, "UpdateHolding - Holding not found: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Holding not found")
		return
	}

	var updateData models.GoalHolding
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "UpdateHolding - Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update core fields
	if updateData.Name != "" {
		existingHolding.Name = updateData.Name
	}
	if updateData.CurrentValue > 0 {
		existingHolding.CurrentValue = updateData.CurrentValue
	}
	if updateData.Status != "" {
		existingHolding.Status = updateData.Status
	}
	if !updateData.PurchaseDate.IsZero() {
		existingHolding.PurchaseDate = updateData.PurchaseDate
	}

	// Update market instrument fields
	if updateData.Symbol != nil {
		existingHolding.Symbol = updateData.Symbol
	}
	if updateData.Quantity != nil {
		existingHolding.Quantity = updateData.Quantity
	}
	if updateData.CostBasis != nil {
		existingHolding.CostBasis = updateData.CostBasis
	}
	if updateData.CurrentPrice != nil {
		existingHolding.CurrentPrice = updateData.CurrentPrice
	}

	// Update bank product fields
	if updateData.Institution != nil {
		existingHolding.Institution = updateData.Institution
	}
	if updateData.InterestRate != nil {
		existingHolding.InterestRate = updateData.InterestRate
	}
	if updateData.TenureMonths != nil {
		existingHolding.TenureMonths = updateData.TenureMonths
	}
	if updateData.MaturityDate != nil {
		existingHolding.MaturityDate = updateData.MaturityDate
	}
	if updateData.MaturityAmount != nil {
		existingHolding.MaturityAmount = updateData.MaturityAmount
	}

	// Update DPS field
	if updateData.MonthlyDeposit != nil {
		existingHolding.MonthlyDeposit = updateData.MonthlyDeposit
	}

	// Recalculate market value
	existingHolding.UpdateMarketValue()

	if err := database.DB.WithContext(ctx).Save(&existingHolding).Error; err != nil {
		logger.Errorf(ctx, "UpdateHolding - Failed to save holding: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update holding")
		return
	}

	// Update goal's current amount
	var goal models.Goal
	if err := database.DB.WithContext(ctx).First(&goal, existingHolding.GoalID).Error; err == nil {
		goal.UpdateCurrentAmount(database.DB.WithContext(ctx))
	}

	// Log holding update activity
	utilities.LogEntityActivity(c, userID, models.ActionUpdate, models.ModuleGoal,
		"GoalHolding", existingHolding.ID, "Updated holding: "+existingHolding.Name, nil)

	logger.Infof(ctx, "UpdateHolding - Successfully updated holding: %s", holdingID)
	utilities.SuccessResponse(c, existingHolding, "Holding updated successfully")
}

// RemoveHolding removes/liquidates a holding
func RemoveHolding(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "RemoveHolding - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "RemoveHolding - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	holdingID, err := uuid.Parse(c.Param("holdingId"))
	if err != nil {
		logger.Warnf(ctx, "RemoveHolding - Invalid holding ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid holding ID")
		return
	}

	logger.Debugf(ctx, "RemoveHolding - Removing holding: %s", holdingID)

	var removeData struct {
		AccountID    uuid.UUID   `json:"accountId" binding:"required"`
		CurrentValue float64     `json:"currentValue" binding:"required,gt=0"`
		Date         models.Date `json:"date"`
		Notes        string      `json:"notes"`
	}

	if err := c.ShouldBindJSON(&removeData); err != nil {
		logger.Warnf(ctx, "RemoveHolding - Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var holding models.GoalHolding
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", holdingID, userID).First(&holding).Error; err != nil {
		logger.Errorf(ctx, "RemoveHolding - Holding not found: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Holding not found")
		return
	}

	// Verify account
	var account models.Account
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", removeData.AccountID, userID).First(&account).Error; err != nil {
		logger.Errorf(ctx, "RemoveHolding - Invalid account ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	if removeData.Date.IsZero() {
		removeData.Date = models.Date{Time: time.Now()}
	}

	logger.Debugf(ctx, "RemoveHolding - Starting transaction to remove holding")

	// Start transaction
	tx := database.DB.WithContext(ctx).Begin()
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

	if err := tx.Commit().Error; err != nil {
		logger.Errorf(ctx, "RemoveHolding - Failed to commit transaction: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	logger.Debugf(ctx, "RemoveHolding - Transaction committed successfully")

	// Log holding removal activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, models.ModuleGoal,
		"GoalHolding", holding.ID, "Removed holding: "+holding.Name, nil)

	logger.Infof(ctx, "RemoveHolding - Successfully removed holding: %s", holdingID)

	result := map[string]interface{}{
		"holding":      holding,
		"transaction":  transaction,
		"contribution": contribution,
	}

	utilities.SuccessResponse(c, result, "Holding removed successfully")
}

// GetHoldingTypes returns all available holding types
func GetHoldingTypes(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "GetHoldingTypes - Entry")

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

	logger.Infof(ctx, "GetHoldingTypes - Successfully retrieved holding types")
	utilities.SuccessResponse(c, holdingTypes, "Holding types retrieved successfully")
}

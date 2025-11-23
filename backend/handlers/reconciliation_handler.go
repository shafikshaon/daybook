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

// ListReconciliations returns all reconciliations for a specific account or user
func ListReconciliations(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListReconciliations - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "ListReconciliations - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountID := c.Query("accountId")
	logger.Debugf(ctx, "ListReconciliations - Fetching reconciliations for user: %s", userID)

	var reconciliations []models.Reconciliation
	query := database.DB.WithContext(ctx).Where("user_id = ?", userID).Preload("Account")

	if accountID != "" {
		accID, err := uuid.Parse(accountID)
		if err != nil {
			logger.Warnf(ctx, "ListReconciliations - Invalid account ID: %v", err)
			utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
			return
		}
		logger.Debugf(ctx, "ListReconciliations - Filtering by account ID: %s", accID)
		query = query.Where("account_id = ?", accID)
	}

	if err := query.Order("reconciliation_date DESC").Find(&reconciliations).Error; err != nil {
		logger.Errorf(ctx, "ListReconciliations - Failed to fetch reconciliations: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch reconciliations")
		return
	}

	logger.Infof(ctx, "ListReconciliations - Successfully retrieved %d reconciliations", len(reconciliations))
	utilities.SuccessResponse(c, reconciliations, "Reconciliations retrieved successfully")
}

// GetReconciliation returns a specific reconciliation by ID
func GetReconciliation(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetReconciliation - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "GetReconciliation - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	reconciliationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "GetReconciliation - Invalid reconciliation ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid reconciliation ID")
		return
	}

	logger.Debugf(ctx, "GetReconciliation - Fetching reconciliation: %s for user: %s", reconciliationID, userID)

	var reconciliation models.Reconciliation
	if err := database.DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", reconciliationID, userID).
		Preload("Account").
		Preload("Transactions.Transaction").
		First(&reconciliation).Error; err != nil {
		logger.Warnf(ctx, "GetReconciliation - Reconciliation not found: %s", reconciliationID)
		utilities.ErrorResponse(c, http.StatusNotFound, "Reconciliation not found")
		return
	}

	logger.Infof(ctx, "GetReconciliation - Successfully retrieved reconciliation: %s", reconciliationID)
	utilities.SuccessResponse(c, reconciliation, "Reconciliation retrieved successfully")
}

// CreateReconciliationRequest represents the request body for creating a reconciliation
type CreateReconciliationRequest struct {
	AccountID          uuid.UUID   `json:"accountId" binding:"required"`
	ReconciliationDate time.Time   `json:"reconciliationDate" binding:"required"`
	StatementBalance   float64     `json:"statementBalance" binding:"required"`
	Notes              string      `json:"notes"`
	TransactionIDs     []uuid.UUID `json:"transactionIds"` // Optional: specific transactions to reconcile
}

// CreateReconciliation creates a new reconciliation record
func CreateReconciliation(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateReconciliation - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "CreateReconciliation - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreateReconciliationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf(ctx, "CreateReconciliation - Invalid request data: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "CreateReconciliation - Creating reconciliation for account: %s, user: %s", req.AccountID, userID)

	// Verify account belongs to user
	var account models.Account
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", req.AccountID, userID).First(&account).Error; err != nil {
		logger.Warnf(ctx, "CreateReconciliation - Account not found: %s", req.AccountID)
		utilities.ErrorResponse(c, http.StatusNotFound, "Account not found")
		return
	}

	// Get current book balance (account balance)
	bookBalance := account.Balance
	logger.Debugf(ctx, "CreateReconciliation - Book balance: %.2f, Statement balance: %.2f", bookBalance, req.StatementBalance)

	// Create reconciliation record
	reconciliation := models.Reconciliation{
		UserID:             userID,
		AccountID:          req.AccountID,
		ReconciliationDate: req.ReconciliationDate,
		StatementBalance:   req.StatementBalance,
		BookBalance:        bookBalance,
		Notes:              req.Notes,
	}

	// Start transaction
	tx := database.DB.WithContext(ctx).Begin()

	if err := tx.Create(&reconciliation).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "CreateReconciliation - Failed to create reconciliation: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create reconciliation")
		return
	}

	// If specific transactions are provided, link them to this reconciliation
	if len(req.TransactionIDs) > 0 {
		logger.Debugf(ctx, "CreateReconciliation - Linking %d transactions", len(req.TransactionIDs))
		for _, transactionID := range req.TransactionIDs {
			// Verify transaction belongs to this account and user
			var transaction models.Transaction
			if err := tx.Where("id = ? AND user_id = ? AND account_id = ?", transactionID, userID, req.AccountID).First(&transaction).Error; err != nil {
				tx.Rollback()
				logger.Warnf(ctx, "CreateReconciliation - Invalid transaction ID: %s", transactionID)
				utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID: "+transactionID.String())
				return
			}

			// Create reconciliation transaction link
			reconciliationTransaction := models.ReconciliationTransaction{
				ReconciliationID: reconciliation.ID,
				TransactionID:    transactionID,
			}

			if err := tx.Create(&reconciliationTransaction).Error; err != nil {
				tx.Rollback()
				logger.Errorf(ctx, "CreateReconciliation - Failed to link transactions: %v", err)
				utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to link transactions")
				return
			}

			// Mark transaction as reconciled
			transaction.Reconciled = true
			transaction.ReconciliationID = &reconciliation.ID
			if err := tx.Save(&transaction).Error; err != nil {
				tx.Rollback()
				logger.Errorf(ctx, "CreateReconciliation - Failed to update transaction: %v", err)
				utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update transaction")
				return
			}
		}
	}

	// Update account's last reconciled date if reconciliation is completed
	if reconciliation.Status == models.ReconciliationCompleted {
		logger.Debugf(ctx, "CreateReconciliation - Marking reconciliation as completed")
		account.LastReconciled = &reconciliation.ReconciliationDate
		account.ReconciliationDifference = 0
		if err := tx.Save(&account).Error; err != nil {
			tx.Rollback()
			logger.Errorf(ctx, "CreateReconciliation - Failed to update account: %v", err)
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account")
			return
		}
	} else {
		account.ReconciliationDifference = reconciliation.Difference
		if err := tx.Save(&account).Error; err != nil {
			tx.Rollback()
			logger.Errorf(ctx, "CreateReconciliation - Failed to update account: %v", err)
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account")
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		logger.Errorf(ctx, "CreateReconciliation - Failed to commit transaction: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to commit reconciliation")
		return
	}

	logger.Infof(ctx, "CreateReconciliation - Successfully created reconciliation: %s", reconciliation.ID)

	// Reload reconciliation with relationships
	database.DB.WithContext(ctx).
		Where("id = ?", reconciliation.ID).
		Preload("Account").
		Preload("Transactions.Transaction").
		First(&reconciliation)

	// Log reconciliation creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleReconcile,
		"Reconciliation", reconciliation.ID, "Created reconciliation for account", nil)

	utilities.CreatedResponse(c, reconciliation, "Reconciliation created successfully")
}

// UpdateReconciliation updates an existing reconciliation
func UpdateReconciliation(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateReconciliation - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "UpdateReconciliation - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	reconciliationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "UpdateReconciliation - Invalid reconciliation ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid reconciliation ID")
		return
	}

	logger.Debugf(ctx, "UpdateReconciliation - Updating reconciliation: %s for user: %s", reconciliationID, userID)

	var existingReconciliation models.Reconciliation
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", reconciliationID, userID).First(&existingReconciliation).Error; err != nil {
		logger.Warnf(ctx, "UpdateReconciliation - Reconciliation not found: %s", reconciliationID)
		utilities.ErrorResponse(c, http.StatusNotFound, "Reconciliation not found")
		return
	}

	var updateData models.Reconciliation
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "UpdateReconciliation - Invalid request data: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update allowed fields
	existingReconciliation.ReconciliationDate = updateData.ReconciliationDate
	existingReconciliation.StatementBalance = updateData.StatementBalance
	existingReconciliation.Notes = updateData.Notes
	if updateData.Status != "" {
		logger.Debugf(ctx, "UpdateReconciliation - Updating status to: %s", updateData.Status)
		existingReconciliation.Status = updateData.Status
	}

	if err := database.DB.WithContext(ctx).Save(&existingReconciliation).Error; err != nil {
		logger.Errorf(ctx, "UpdateReconciliation - Failed to update reconciliation: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update reconciliation")
		return
	}

	// Update account's last reconciled date if reconciliation is completed
	if existingReconciliation.Status == models.ReconciliationCompleted {
		logger.Debugf(ctx, "UpdateReconciliation - Reconciliation completed, updating account")
		var account models.Account
		if err := database.DB.WithContext(ctx).Where("id = ?", existingReconciliation.AccountID).First(&account).Error; err == nil {
			account.LastReconciled = &existingReconciliation.ReconciliationDate
			account.ReconciliationDifference = 0
			database.DB.WithContext(ctx).Save(&account)
		}
	}

	logger.Infof(ctx, "UpdateReconciliation - Successfully updated reconciliation: %s", reconciliationID)

	// Log reconciliation update activity
	utilities.LogEntityActivity(c, userID, models.ActionUpdate, models.ModuleReconcile,
		"Reconciliation", existingReconciliation.ID, "Updated reconciliation", nil)

	utilities.SuccessResponse(c, existingReconciliation, "Reconciliation updated successfully")
}

// DeleteReconciliation deletes a reconciliation record
func DeleteReconciliation(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteReconciliation - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "DeleteReconciliation - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	reconciliationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "DeleteReconciliation - Invalid reconciliation ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid reconciliation ID")
		return
	}

	logger.Debugf(ctx, "DeleteReconciliation - Deleting reconciliation: %s for user: %s", reconciliationID, userID)

	var reconciliation models.Reconciliation
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", reconciliationID, userID).First(&reconciliation).Error; err != nil {
		logger.Warnf(ctx, "DeleteReconciliation - Reconciliation not found: %s", reconciliationID)
		utilities.ErrorResponse(c, http.StatusNotFound, "Reconciliation not found")
		return
	}

	// Start transaction
	tx := database.DB.WithContext(ctx).Begin()

	// Unmark all reconciled transactions
	var reconciliationTransactions []models.ReconciliationTransaction
	if err := tx.Where("reconciliation_id = ?", reconciliationID).Find(&reconciliationTransactions).Error; err == nil {
		logger.Debugf(ctx, "DeleteReconciliation - Unmarking %d reconciled transactions", len(reconciliationTransactions))
		for _, rt := range reconciliationTransactions {
			tx.Model(&models.Transaction{}).Where("id = ?", rt.TransactionID).Updates(map[string]interface{}{
				"reconciled":        false,
				"reconciliation_id": nil,
			})
		}
	}

	// Delete reconciliation transaction links
	if err := tx.Where("reconciliation_id = ?", reconciliationID).Delete(&models.ReconciliationTransaction{}).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "DeleteReconciliation - Failed to delete reconciliation transactions: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete reconciliation transactions")
		return
	}

	// Soft delete reconciliation
	if err := tx.Delete(&reconciliation).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "DeleteReconciliation - Failed to delete reconciliation: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete reconciliation")
		return
	}

	tx.Commit()

	logger.Infof(ctx, "DeleteReconciliation - Successfully deleted reconciliation: %s", reconciliationID)

	// Log reconciliation deletion activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, models.ModuleReconcile,
		"Reconciliation", reconciliation.ID, "Deleted reconciliation", nil)

	utilities.SuccessResponse(c, nil, "Reconciliation deleted successfully")
}

// GetUnreconciledTransactions returns all unreconciled transactions for an account
func GetUnreconciledTransactions(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetUnreconciledTransactions - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "GetUnreconciledTransactions - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "GetUnreconciledTransactions - Invalid account ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	logger.Debugf(ctx, "GetUnreconciledTransactions - Fetching unreconciled transactions for account: %s, user: %s", accountID, userID)

	// Verify account belongs to user
	var account models.Account
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		logger.Warnf(ctx, "GetUnreconciledTransactions - Account not found: %s", accountID)
		utilities.ErrorResponse(c, http.StatusNotFound, "Account not found")
		return
	}

	var transactions []models.Transaction
	query := database.DB.WithContext(ctx).Where("account_id = ? AND user_id = ?", accountID, userID)

	// Check if reconciled column exists, if not just return all transactions
	// This handles the case where migration hasn't been run yet
	var columnExists bool
	database.DB.WithContext(ctx).Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'transactions' AND column_name = 'reconciled')").Scan(&columnExists)

	if columnExists {
		logger.Debugf(ctx, "GetUnreconciledTransactions - Filtering by reconciled status")
		query = query.Where("reconciled = ? OR reconciled IS NULL", false)
	}

	if err := query.Order("date DESC").Find(&transactions).Error; err != nil {
		logger.Errorf(ctx, "GetUnreconciledTransactions - Failed to fetch transactions: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

	logger.Infof(ctx, "GetUnreconciledTransactions - Successfully retrieved %d unreconciled transactions", len(transactions))
	utilities.SuccessResponse(c, transactions, "Unreconciled transactions retrieved successfully")
}

// ReconciliationStatsResponse represents reconciliation statistics
type ReconciliationStatsResponse struct {
	TotalReconciliations   int64   `json:"totalReconciliations"`
	CompletedReconciliations int64 `json:"completedReconciliations"`
	PendingReconciliations int64   `json:"pendingReconciliations"`
	DiscrepancyReconciliations int64 `json:"discrepancyReconciliations"`
	LastReconciliationDate *time.Time `json:"lastReconciliationDate"`
	AverageDifference      float64 `json:"averageDifference"`
}

// GetReconciliationStats returns reconciliation statistics for an account
func GetReconciliationStats(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetReconciliationStats - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "GetReconciliationStats - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "GetReconciliationStats - Invalid account ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	logger.Debugf(ctx, "GetReconciliationStats - Fetching reconciliation stats for account: %s, user: %s", accountID, userID)

	// Verify account belongs to user
	var account models.Account
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		logger.Warnf(ctx, "GetReconciliationStats - Account not found: %s", accountID)
		utilities.ErrorResponse(c, http.StatusNotFound, "Account not found")
		return
	}

	var stats ReconciliationStatsResponse

	// Total reconciliations
	database.DB.WithContext(ctx).Model(&models.Reconciliation{}).
		Where("account_id = ? AND user_id = ?", accountID, userID).
		Count(&stats.TotalReconciliations)

	// Completed reconciliations
	database.DB.WithContext(ctx).Model(&models.Reconciliation{}).
		Where("account_id = ? AND user_id = ? AND status = ?", accountID, userID, models.ReconciliationCompleted).
		Count(&stats.CompletedReconciliations)

	// Pending reconciliations
	database.DB.WithContext(ctx).Model(&models.Reconciliation{}).
		Where("account_id = ? AND user_id = ? AND status = ?", accountID, userID, models.ReconciliationPending).
		Count(&stats.PendingReconciliations)

	// Discrepancy reconciliations
	database.DB.WithContext(ctx).Model(&models.Reconciliation{}).
		Where("account_id = ? AND user_id = ? AND status = ?", accountID, userID, models.ReconciliationDiscrepancy).
		Count(&stats.DiscrepancyReconciliations)

	// Last reconciliation date
	var lastReconciliation models.Reconciliation
	if err := database.DB.WithContext(ctx).
		Where("account_id = ? AND user_id = ?", accountID, userID).
		Order("reconciliation_date DESC").
		First(&lastReconciliation).Error; err == nil {
		stats.LastReconciliationDate = &lastReconciliation.ReconciliationDate
	}

	// Average difference
	var avgDiff struct {
		AvgDifference float64
	}
	database.DB.WithContext(ctx).Model(&models.Reconciliation{}).
		Select("AVG(ABS(difference)) as avg_difference").
		Where("account_id = ? AND user_id = ?", accountID, userID).
		Scan(&avgDiff)
	stats.AverageDifference = avgDiff.AvgDifference

	logger.Infof(ctx, "GetReconciliationStats - Successfully retrieved stats. Total: %d, Completed: %d", stats.TotalReconciliations, stats.CompletedReconciliations)
	utilities.SuccessResponse(c, stats, "Reconciliation stats retrieved successfully")
}

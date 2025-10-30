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

// ListReconciliations returns all reconciliations for a specific account or user
func ListReconciliations(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountID := c.Query("accountId")

	var reconciliations []models.Reconciliation
	query := database.DB.Where("user_id = ?", userID).Preload("Account")

	if accountID != "" {
		accID, err := uuid.Parse(accountID)
		if err != nil {
			utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
			return
		}
		query = query.Where("account_id = ?", accID)
	}

	if err := query.Order("reconciliation_date DESC").Find(&reconciliations).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch reconciliations")
		return
	}

	utilities.SuccessResponse(c, reconciliations, "Reconciliations retrieved successfully")
}

// GetReconciliation returns a specific reconciliation by ID
func GetReconciliation(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	reconciliationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid reconciliation ID")
		return
	}

	var reconciliation models.Reconciliation
	if err := database.DB.
		Where("id = ? AND user_id = ?", reconciliationID, userID).
		Preload("Account").
		Preload("Transactions.Transaction").
		First(&reconciliation).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Reconciliation not found")
		return
	}

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
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreateReconciliationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Verify account belongs to user
	var account models.Account
	if err := database.DB.Where("id = ? AND user_id = ?", req.AccountID, userID).First(&account).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Account not found")
		return
	}

	// Get current book balance (account balance)
	bookBalance := account.Balance

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
	tx := database.DB.Begin()

	if err := tx.Create(&reconciliation).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create reconciliation")
		return
	}

	// If specific transactions are provided, link them to this reconciliation
	if len(req.TransactionIDs) > 0 {
		for _, transactionID := range req.TransactionIDs {
			// Verify transaction belongs to this account and user
			var transaction models.Transaction
			if err := tx.Where("id = ? AND user_id = ? AND account_id = ?", transactionID, userID, req.AccountID).First(&transaction).Error; err != nil {
				tx.Rollback()
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
				utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to link transactions")
				return
			}

			// Mark transaction as reconciled
			transaction.Reconciled = true
			transaction.ReconciliationID = &reconciliation.ID
			if err := tx.Save(&transaction).Error; err != nil {
				tx.Rollback()
				utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update transaction")
				return
			}
		}
	}

	// Update account's last reconciled date if reconciliation is completed
	if reconciliation.Status == models.ReconciliationCompleted {
		account.LastReconciled = &reconciliation.ReconciliationDate
		account.ReconciliationDifference = 0
		if err := tx.Save(&account).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account")
			return
		}
	} else {
		account.ReconciliationDifference = reconciliation.Difference
		if err := tx.Save(&account).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account")
			return
		}
	}

	tx.Commit()

	// Reload reconciliation with relationships
	database.DB.
		Where("id = ?", reconciliation.ID).
		Preload("Account").
		Preload("Transactions.Transaction").
		First(&reconciliation)

	utilities.CreatedResponse(c, reconciliation, "Reconciliation created successfully")
}

// UpdateReconciliation updates an existing reconciliation
func UpdateReconciliation(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	reconciliationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid reconciliation ID")
		return
	}

	var existingReconciliation models.Reconciliation
	if err := database.DB.Where("id = ? AND user_id = ?", reconciliationID, userID).First(&existingReconciliation).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Reconciliation not found")
		return
	}

	var updateData models.Reconciliation
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update allowed fields
	existingReconciliation.ReconciliationDate = updateData.ReconciliationDate
	existingReconciliation.StatementBalance = updateData.StatementBalance
	existingReconciliation.Notes = updateData.Notes
	if updateData.Status != "" {
		existingReconciliation.Status = updateData.Status
	}

	if err := database.DB.Save(&existingReconciliation).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update reconciliation")
		return
	}

	// Update account's last reconciled date if reconciliation is completed
	if existingReconciliation.Status == models.ReconciliationCompleted {
		var account models.Account
		if err := database.DB.Where("id = ?", existingReconciliation.AccountID).First(&account).Error; err == nil {
			account.LastReconciled = &existingReconciliation.ReconciliationDate
			account.ReconciliationDifference = 0
			database.DB.Save(&account)
		}
	}

	utilities.SuccessResponse(c, existingReconciliation, "Reconciliation updated successfully")
}

// DeleteReconciliation deletes a reconciliation record
func DeleteReconciliation(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	reconciliationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid reconciliation ID")
		return
	}

	var reconciliation models.Reconciliation
	if err := database.DB.Where("id = ? AND user_id = ?", reconciliationID, userID).First(&reconciliation).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Reconciliation not found")
		return
	}

	// Start transaction
	tx := database.DB.Begin()

	// Unmark all reconciled transactions
	var reconciliationTransactions []models.ReconciliationTransaction
	if err := tx.Where("reconciliation_id = ?", reconciliationID).Find(&reconciliationTransactions).Error; err == nil {
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
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete reconciliation transactions")
		return
	}

	// Soft delete reconciliation
	if err := tx.Delete(&reconciliation).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete reconciliation")
		return
	}

	tx.Commit()

	utilities.SuccessResponse(c, nil, "Reconciliation deleted successfully")
}

// GetUnreconciledTransactions returns all unreconciled transactions for an account
func GetUnreconciledTransactions(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	// Verify account belongs to user
	var account models.Account
	if err := database.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Account not found")
		return
	}

	var transactions []models.Transaction
	query := database.DB.Where("account_id = ? AND user_id = ?", accountID, userID)

	// Check if reconciled column exists, if not just return all transactions
	// This handles the case where migration hasn't been run yet
	var columnExists bool
	database.DB.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'transactions' AND column_name = 'reconciled')").Scan(&columnExists)

	if columnExists {
		query = query.Where("reconciled = ? OR reconciled IS NULL", false)
	}

	if err := query.Order("date DESC").Find(&transactions).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

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
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	// Verify account belongs to user
	var account models.Account
	if err := database.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Account not found")
		return
	}

	var stats ReconciliationStatsResponse

	// Total reconciliations
	database.DB.Model(&models.Reconciliation{}).
		Where("account_id = ? AND user_id = ?", accountID, userID).
		Count(&stats.TotalReconciliations)

	// Completed reconciliations
	database.DB.Model(&models.Reconciliation{}).
		Where("account_id = ? AND user_id = ? AND status = ?", accountID, userID, models.ReconciliationCompleted).
		Count(&stats.CompletedReconciliations)

	// Pending reconciliations
	database.DB.Model(&models.Reconciliation{}).
		Where("account_id = ? AND user_id = ? AND status = ?", accountID, userID, models.ReconciliationPending).
		Count(&stats.PendingReconciliations)

	// Discrepancy reconciliations
	database.DB.Model(&models.Reconciliation{}).
		Where("account_id = ? AND user_id = ? AND status = ?", accountID, userID, models.ReconciliationDiscrepancy).
		Count(&stats.DiscrepancyReconciliations)

	// Last reconciliation date
	var lastReconciliation models.Reconciliation
	if err := database.DB.
		Where("account_id = ? AND user_id = ?", accountID, userID).
		Order("reconciliation_date DESC").
		First(&lastReconciliation).Error; err == nil {
		stats.LastReconciliationDate = &lastReconciliation.ReconciliationDate
	}

	// Average difference
	var avgDiff struct {
		AvgDifference float64
	}
	database.DB.Model(&models.Reconciliation{}).
		Select("AVG(ABS(difference)) as avg_difference").
		Where("account_id = ? AND user_id = ?", accountID, userID).
		Scan(&avgDiff)
	stats.AverageDifference = avgDiff.AvgDifference

	utilities.SuccessResponse(c, stats, "Reconciliation stats retrieved successfully")
}

package handlers

import (
	"net/http"

	"daybook-backend/database"
	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListAccounts returns all accounts for the authenticated user
func ListAccounts(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListAccounts handler - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access to list accounts: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Fetching accounts for user: %s", userID)

	var accounts []models.Account
	if err := database.DB.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&accounts).Error; err != nil {
		logger.Errorf(ctx, "Failed to fetch accounts from database: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch accounts")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d accounts for user: %s", len(accounts), userID)
	utilities.SuccessResponse(c, accounts, "Accounts retrieved successfully")
}

// GetAccount returns a specific account by ID
func GetAccount(c *gin.Context) {
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

	var account models.Account
	if err := database.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Account not found")
		return
	}

	utilities.SuccessResponse(c, account, "Account retrieved successfully")
}

// CreateAccount creates a new account
func CreateAccount(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateAccount handler - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access to create account: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var account models.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		logger.Warnf(ctx, "Invalid create account request: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	account.UserID = userID
	logger.Infof(ctx, "Creating new account '%s' for user: %s", account.Name, userID)

	// Start transaction to ensure atomicity
	logger.Debugf(ctx, "Starting database transaction for account creation")
	tx := database.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.Errorf(ctx, "Failed to start transaction: %v", tx.Error)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	// Create the account
	logger.Debugf(ctx, "Creating account record in database")
	if err := tx.Create(&account).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to create account: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create account")
		return
	}
	logger.Infof(ctx, "Account created successfully with ID: %s", account.ID)

	// If there's an initial balance, create an opening balance transaction
	if account.InitialBalance > 0 {
		logger.Debugf(ctx, "Creating opening balance transaction: %.2f", account.InitialBalance)
		transaction := models.Transaction{
			UserID:      userID,
			AccountID:   account.ID,
			Type:        "income",
			CategoryID:  "opening_balance",
			Amount:      account.InitialBalance,
			Date:        account.CreatedAt,
			Description: "Opening balance for " + account.Name,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			logger.Errorf(ctx, "Failed to create opening balance transaction: %v", err)
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create opening balance transaction")
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	// Log account creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleAccount,
		"Account", account.ID, "Created account: "+account.Name, nil)

	utilities.CreatedResponse(c, account, "Account created successfully")
}

// UpdateAccount updates an existing account
func UpdateAccount(c *gin.Context) {
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

	var existingAccount models.Account
	if err := database.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&existingAccount).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Account not found")
		return
	}

	var updateData models.Account
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update only allowed fields
	existingAccount.Name = updateData.Name
	existingAccount.Type = updateData.Type
	existingAccount.Currency = updateData.Currency
	existingAccount.Description = updateData.Description
	existingAccount.Institution = updateData.Institution
	existingAccount.AccountNumber = updateData.AccountNumber
	existingAccount.Active = updateData.Active

	if err := database.DB.Save(&existingAccount).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account")
		return
	}

	// Log account update activity
	utilities.LogEntityActivity(c, userID, models.ActionUpdate, models.ModuleAccount,
		"Account", existingAccount.ID, "Updated account: "+existingAccount.Name, nil)

	utilities.SuccessResponse(c, existingAccount, "Account updated successfully")
}

// DeleteAccount deletes an account
func DeleteAccount(c *gin.Context) {
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

	var account models.Account
	if err := database.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Account not found")
		return
	}

	// Soft delete
	if err := database.DB.Delete(&account).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete account")
		return
	}

	// Log account deletion activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, models.ModuleAccount,
		"Account", account.ID, "Deleted account: "+account.Name, nil)

	utilities.SuccessResponse(c, nil, "Account deleted successfully")
}

// NOTE: Direct balance updates are NOT allowed in the dual-balance system.
// Balance is automatically updated by transactions only.
// Initial balance is set once during account creation and never changes.
// See BALANCE_SYSTEM.md for detailed documentation.

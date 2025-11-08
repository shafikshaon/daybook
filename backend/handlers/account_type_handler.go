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

// ListAccountTypes returns all account types for the authenticated user
func ListAccountTypes(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListAccountTypes - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Fetching account types for user: %s", userID)
	var accountTypes []models.AccountType
	// Get user's account types
	if err := database.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("sort_order ASC, name ASC").
		Find(&accountTypes).Error; err != nil {
		logger.Errorf(ctx, "Failed to fetch account types: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch account types")
		return
	}

	logger.Infof(ctx, "Account types retrieved successfully for user: %s, count: %d", userID, len(accountTypes))
	utilities.SuccessResponse(c, accountTypes, "Account types retrieved successfully")
}

// GetAccountType returns a specific account type by ID
func GetAccountType(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetAccountType - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountTypeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid account type ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account type ID")
		return
	}

	logger.Debugf(ctx, "Fetching account type %s for user: %s", accountTypeID, userID)
	var accountType models.AccountType
	if err := database.DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", accountTypeID, userID).
		First(&accountType).Error; err != nil {
		logger.Warnf(ctx, "Account type not found: %s, error: %v", accountTypeID, err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Account type not found")
		return
	}

	logger.Infof(ctx, "Account type retrieved successfully: %s for user: %s", accountTypeID, userID)
	utilities.SuccessResponse(c, accountType, "Account type retrieved successfully")
}

// CreateAccountType creates a new account type for the user
func CreateAccountType(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateAccountType - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Parsing account type request for user: %s", userID)
	var accountType models.AccountType
	if err := c.ShouldBindJSON(&accountType); err != nil {
		logger.Warnf(ctx, "Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Set user ID
	accountType.UserID = userID

	logger.Debugf(ctx, "Creating account type '%s' for user: %s", accountType.Name, userID)
	if err := database.DB.WithContext(ctx).Create(&accountType).Error; err != nil {
		logger.Errorf(ctx, "Failed to create account type: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create account type")
		return
	}

	logger.Infof(ctx, "Account type created successfully: %s for user: %s", accountType.ID, userID)
	utilities.CreatedResponse(c, accountType, "Account type created successfully")
}

// UpdateAccountType updates an existing account type
func UpdateAccountType(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateAccountType - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountTypeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid account type ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account type ID")
		return
	}

	logger.Debugf(ctx, "Fetching existing account type %s for user: %s", accountTypeID, userID)
	var existingAccountType models.AccountType
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", accountTypeID, userID).First(&existingAccountType).Error; err != nil {
		logger.Warnf(ctx, "Account type not found: %s, error: %v", accountTypeID, err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Account type not found")
		return
	}

	logger.Debugf(ctx, "Parsing update data for account type: %s", accountTypeID)
	var updateData models.AccountType
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update only allowed fields
	logger.Debugf(ctx, "Updating account type fields for: %s", accountTypeID)
	existingAccountType.Name = updateData.Name
	existingAccountType.Icon = updateData.Icon
	existingAccountType.Description = updateData.Description
	existingAccountType.Active = updateData.Active
	existingAccountType.SortOrder = updateData.SortOrder

	if err := database.DB.WithContext(ctx).Save(&existingAccountType).Error; err != nil {
		logger.Errorf(ctx, "Failed to update account type: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account type")
		return
	}

	logger.Infof(ctx, "Account type updated successfully: %s for user: %s", accountTypeID, userID)
	utilities.SuccessResponse(c, existingAccountType, "Account type updated successfully")
}

// DeleteAccountType deletes an account type (soft delete)
func DeleteAccountType(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteAccountType - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountTypeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid account type ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account type ID")
		return
	}

	logger.Debugf(ctx, "Fetching account type %s for deletion", accountTypeID)
	var accountType models.AccountType
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", accountTypeID, userID).First(&accountType).Error; err != nil {
		logger.Warnf(ctx, "Account type not found: %s, error: %v", accountTypeID, err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Account type not found")
		return
	}

	// Check if any accounts are using this type
	// The account type is stored in lowercase with underscores (e.g., "digital_wallet")
	// Convert the account type name to match the format
	typeValue := utilities.ToSnakeCase(accountType.Name)

	logger.Debugf(ctx, "Checking if account type is in use: %s", typeValue)
	var accountCount int64
	database.DB.WithContext(ctx).Model(&models.Account{}).Where("user_id = ? AND type = ?", userID, typeValue).Count(&accountCount)
	if accountCount > 0 {
		logger.Warnf(ctx, "Cannot delete account type %s: in use by %d accounts", accountTypeID, accountCount)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Cannot delete account type that is in use by existing accounts")
		return
	}

	// Soft delete
	logger.Debugf(ctx, "Deleting account type: %s", accountTypeID)
	if err := database.DB.WithContext(ctx).Delete(&accountType).Error; err != nil {
		logger.Errorf(ctx, "Failed to delete account type: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete account type")
		return
	}

	logger.Infof(ctx, "Account type deleted successfully: %s for user: %s", accountTypeID, userID)
	utilities.SuccessResponse(c, nil, "Account type deleted successfully")
}

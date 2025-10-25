package handlers

import (
	"net/http"

	"daybook-backend/database"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListAccountTypes returns all account types for the authenticated user
func ListAccountTypes(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var accountTypes []models.AccountType
	// Get user's account types
	if err := database.DB.
		Where("user_id = ?", userID).
		Order("sort_order ASC, name ASC").
		Find(&accountTypes).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch account types")
		return
	}

	utilities.SuccessResponse(c, accountTypes, "Account types retrieved successfully")
}

// GetAccountType returns a specific account type by ID
func GetAccountType(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountTypeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account type ID")
		return
	}

	var accountType models.AccountType
	if err := database.DB.
		Where("id = ? AND user_id = ?", accountTypeID, userID).
		First(&accountType).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Account type not found")
		return
	}

	utilities.SuccessResponse(c, accountType, "Account type retrieved successfully")
}

// CreateAccountType creates a new account type for the user
func CreateAccountType(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var accountType models.AccountType
	if err := c.ShouldBindJSON(&accountType); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Set user ID
	accountType.UserID = userID

	if err := database.DB.Create(&accountType).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create account type")
		return
	}

	utilities.CreatedResponse(c, accountType, "Account type created successfully")
}

// UpdateAccountType updates an existing account type
func UpdateAccountType(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountTypeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account type ID")
		return
	}

	var existingAccountType models.AccountType
	if err := database.DB.Where("id = ? AND user_id = ?", accountTypeID, userID).First(&existingAccountType).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Account type not found")
		return
	}

	var updateData models.AccountType
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update only allowed fields
	existingAccountType.Name = updateData.Name
	existingAccountType.Icon = updateData.Icon
	existingAccountType.Description = updateData.Description
	existingAccountType.Active = updateData.Active
	existingAccountType.SortOrder = updateData.SortOrder

	if err := database.DB.Save(&existingAccountType).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account type")
		return
	}

	utilities.SuccessResponse(c, existingAccountType, "Account type updated successfully")
}

// DeleteAccountType deletes an account type (soft delete)
func DeleteAccountType(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountTypeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account type ID")
		return
	}

	var accountType models.AccountType
	if err := database.DB.Where("id = ? AND user_id = ?", accountTypeID, userID).First(&accountType).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Account type not found")
		return
	}

	// Check if any accounts are using this type
	// The account type is stored in lowercase with underscores (e.g., "digital_wallet")
	// Convert the account type name to match the format
	typeValue := utilities.ToSnakeCase(accountType.Name)

	var accountCount int64
	database.DB.Model(&models.Account{}).Where("user_id = ? AND type = ?", userID, typeValue).Count(&accountCount)
	if accountCount > 0 {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Cannot delete account type that is in use by existing accounts")
		return
	}

	// Soft delete
	if err := database.DB.Delete(&accountType).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete account type")
		return
	}

	utilities.SuccessResponse(c, nil, "Account type deleted successfully")
}

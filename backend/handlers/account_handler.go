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

// ListAccounts returns all accounts for the authenticated user
func ListAccounts(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var accounts []models.Account
	if err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&accounts).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch accounts")
		return
	}

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
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var account models.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	account.UserID = userID

	if err := database.DB.Create(&account).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create account")
		return
	}

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

	utilities.SuccessResponse(c, nil, "Account deleted successfully")
}

// UpdateAccountBalance updates the balance of an account
func UpdateAccountBalance(c *gin.Context) {
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

	var balanceUpdate struct {
		Balance float64 `json:"balance" binding:"required"`
	}

	if err := c.ShouldBindJSON(&balanceUpdate); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var account models.Account
	if err := database.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Account not found")
		return
	}

	account.Balance = balanceUpdate.Balance

	if err := database.DB.Save(&account).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update balance")
		return
	}

	utilities.SuccessResponse(c, account, "Account balance updated successfully")
}

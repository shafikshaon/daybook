package handlers

import (
	"net/http"
	"strconv"

	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/services"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
)

// AccountTypeHandler handles account type-related HTTP requests
type AccountTypeHandler struct {
	service services.AccountTypeService
}

// NewAccountTypeHandler creates a new account type handler
func NewAccountTypeHandler(service services.AccountTypeService) *AccountTypeHandler {
	return &AccountTypeHandler{service: service}
}

// ListAccountTypes returns all account types for the authenticated user
func (h *AccountTypeHandler) ListAccountTypes(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListAccountTypes - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Fetching account types for user: %s", userID)
	accountTypes, err := h.service.ListAccountTypes(ctx, userID)
	if err != nil {
		logger.Errorf(ctx, "Failed to fetch account types: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch account types")
		return
	}

	logger.Infof(ctx, "Account types retrieved successfully for user: %s, count: %d", userID, len(accountTypes))
	utilities.SuccessResponse(c, accountTypes, "Account types retrieved successfully")
}

// GetAccountType returns a specific account type by ID
func (h *AccountTypeHandler) GetAccountType(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetAccountType - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountTypeIDStr := c.Param("id")
	accountTypeIDUint, err := strconv.ParseUint(accountTypeIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "Invalid account type ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account type ID")
		return
	}
	accountTypeID := uint(accountTypeIDUint)

	logger.Debugf(ctx, "Fetching account type %s for user: %s", accountTypeID, userID)
	accountType, err := h.service.GetAccountType(ctx, accountTypeID, userID)
	if err != nil {
		logger.Warnf(ctx, "Account type not found: %s, error: %v", accountTypeID, err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Account type not found")
		return
	}

	logger.Infof(ctx, "Account type retrieved successfully: %s for user: %s", accountTypeID, userID)
	utilities.SuccessResponse(c, accountType, "Account type retrieved successfully")
}

// CreateAccountType creates a new account type for the user
func (h *AccountTypeHandler) CreateAccountType(c *gin.Context) {
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
	created, err := h.service.CreateAccountType(ctx, &accountType)
	if err != nil {
		logger.Errorf(ctx, "Failed to create account type: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Account type created successfully: %s for user: %s", created.ID, userID)
	utilities.CreatedResponse(c, created, "Account type created successfully")
}

// UpdateAccountType updates an existing account type
func (h *AccountTypeHandler) UpdateAccountType(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateAccountType - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountTypeIDStr := c.Param("id")
	accountTypeIDUint, err := strconv.ParseUint(accountTypeIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "Invalid account type ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account type ID")
		return
	}
	accountTypeID := uint(accountTypeIDUint)

	logger.Debugf(ctx, "Parsing update data for account type: %s", accountTypeID)
	var updateData models.AccountType
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "Updating account type %s for user: %s", accountTypeID, userID)
	updated, err := h.service.UpdateAccountType(ctx, accountTypeID, userID, &updateData)
	if err != nil {
		logger.Errorf(ctx, "Failed to update account type: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Account type updated successfully: %s for user: %s", accountTypeID, userID)
	utilities.SuccessResponse(c, updated, "Account type updated successfully")
}

// DeleteAccountType deletes an account type (soft delete)
func (h *AccountTypeHandler) DeleteAccountType(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteAccountType - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountTypeIDStr := c.Param("id")
	accountTypeIDUint, err := strconv.ParseUint(accountTypeIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "Invalid account type ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account type ID")
		return
	}
	accountTypeID := uint(accountTypeIDUint)

	logger.Debugf(ctx, "Deleting account type %s for user: %s", accountTypeID, userID)
	if err := h.service.DeleteAccountType(ctx, accountTypeID, userID); err != nil {
		logger.Errorf(ctx, "Failed to delete account type: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Account type deleted successfully: %s for user: %s", accountTypeID, userID)
	utilities.SuccessResponse(c, nil, "Account type deleted successfully")
}

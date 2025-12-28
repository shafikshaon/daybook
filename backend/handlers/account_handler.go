package handlers

import (
	"net/http"

	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/services"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AccountHandler handles account-related HTTP requests
type AccountHandler struct {
	service services.AccountService
}

// NewAccountHandler creates a new account handler
func NewAccountHandler(service services.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

// ListAccounts returns all accounts for the authenticated user
func (h *AccountHandler) ListAccounts(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListAccounts handler - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access to list accounts: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Fetching accounts for user: %s", userID)
	accounts, err := h.service.ListAccounts(ctx, userID)
	if err != nil {
		logger.Errorf(ctx, "Failed to fetch accounts from database: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch accounts")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d accounts for user: %s", len(accounts), userID)
	utilities.SuccessResponse(c, accounts, "Accounts retrieved successfully")
}

// GetAccount returns a specific account by ID
func (h *AccountHandler) GetAccount(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetAccount - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid account ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	logger.Debugf(ctx, "Fetching account %s for user: %s", accountID, userID)
	account, err := h.service.GetAccount(ctx, accountID, userID)
	if err != nil {
		logger.Warnf(ctx, "Account not found: %s, error: %v", accountID, err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Account not found")
		return
	}

	logger.Infof(ctx, "Account retrieved successfully: %s for user: %s", accountID, userID)
	utilities.SuccessResponse(c, account, "Account retrieved successfully")
}

// CreateAccount creates a new account
func (h *AccountHandler) CreateAccount(c *gin.Context) {
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

	created, err := h.service.CreateAccount(ctx, &account)
	if err != nil {
		logger.Errorf(ctx, "Failed to create account: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Account created successfully: %s for user: %s", created.ID, userID)
	utilities.CreatedResponse(c, created, "Account created successfully")
}

// UpdateAccount updates an existing account
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateAccount - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid account ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	logger.Debugf(ctx, "Parsing update data for account: %s", accountID)
	var updateData models.Account
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "Updating account %s for user: %s", accountID, userID)
	updated, err := h.service.UpdateAccount(ctx, accountID, userID, &updateData)
	if err != nil {
		logger.Errorf(ctx, "Failed to update account: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Account updated successfully: %s for user: %s", accountID, userID)
	utilities.SuccessResponse(c, updated, "Account updated successfully")
}

// DeleteAccount deletes an account
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteAccount - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid account ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	logger.Debugf(ctx, "Deleting account %s for user: %s", accountID, userID)
	if err := h.service.DeleteAccount(ctx, accountID, userID); err != nil {
		logger.Errorf(ctx, "Failed to delete account: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Account deleted successfully: %s for user: %s", accountID, userID)
	utilities.SuccessResponse(c, nil, "Account deleted successfully")
}

// NOTE: Direct balance updates are NOT allowed in the dual-balance system.
// Balance is automatically updated by transactions only.
// Initial balance is set once during account creation and never changes.
// See BALANCE_SYSTEM.md for detailed documentation.

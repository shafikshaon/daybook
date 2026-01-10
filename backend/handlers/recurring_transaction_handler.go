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

// RecurringTransactionHandler handles recurring transaction HTTP requests
type RecurringTransactionHandler struct {
	service services.RecurringTransactionService
}

// NewRecurringTransactionHandler creates a new recurring transaction handler
func NewRecurringTransactionHandler(service services.RecurringTransactionService) *RecurringTransactionHandler {
	return &RecurringTransactionHandler{service: service}
}

// ListRecurringTransactions returns all recurring transactions for the authenticated user
func (h *RecurringTransactionHandler) ListRecurringTransactions(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListRecurringTransactions - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	recurringTransactions, err := h.service.ListRecurringTransactions(ctx, userID)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch recurring transactions")
		return
	}

	logger.Infof(ctx, "Recurring transactions retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, recurringTransactions, "Recurring transactions retrieved successfully")
}

// GetRecurringTransaction returns a specific recurring transaction by ID
func (h *RecurringTransactionHandler) GetRecurringTransaction(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetRecurringTransaction - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	recurringIDStr := c.Param("id")
	recurringIDUint, err := strconv.ParseUint(recurringIDStr, 10, 32)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid recurring transaction ID")
		return
	}
	recurringID := uint(recurringIDUint)

	recurringTransaction, err := h.service.GetRecurringTransaction(ctx, recurringID, userID)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Recurring transaction not found")
		return
	}

	logger.Infof(ctx, "Recurring transaction retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, recurringTransaction, "Recurring transaction retrieved successfully")
}

// CreateRecurringTransaction creates a new recurring transaction
func (h *RecurringTransactionHandler) CreateRecurringTransaction(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateRecurringTransaction - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var recurringTransaction models.RecurringTransaction
	if err := c.ShouldBindJSON(&recurringTransaction); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	recurringTransaction.UserID = userID

	createdRecurring, err := h.service.CreateRecurringTransaction(ctx, &recurringTransaction)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Recurring transaction created successfully for user: %s", userID)
	utilities.CreatedResponse(c, createdRecurring, "Recurring transaction created successfully")
}

// UpdateRecurringTransaction updates an existing recurring transaction
func (h *RecurringTransactionHandler) UpdateRecurringTransaction(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateRecurringTransaction - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	recurringIDStr := c.Param("id")
	recurringIDUint, err := strconv.ParseUint(recurringIDStr, 10, 32)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid recurring transaction ID")
		return
	}
	recurringID := uint(recurringIDUint)

	var updateData models.RecurringTransaction
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedRecurring, err := h.service.UpdateRecurringTransaction(ctx, recurringID, userID, &updateData)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Recurring transaction updated successfully for user: %s", userID)
	utilities.SuccessResponse(c, updatedRecurring, "Recurring transaction updated successfully")
}

// DeleteRecurringTransaction deletes a recurring transaction
func (h *RecurringTransactionHandler) DeleteRecurringTransaction(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteRecurringTransaction - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	recurringIDStr := c.Param("id")
	recurringIDUint, err := strconv.ParseUint(recurringIDStr, 10, 32)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid recurring transaction ID")
		return
	}
	recurringID := uint(recurringIDUint)

	if err := h.service.DeleteRecurringTransaction(ctx, recurringID, userID); err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	logger.Infof(ctx, "Recurring transaction deleted successfully for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Recurring transaction deleted successfully")
}

// ProcessRecurringTransactions generates missing transactions for all enabled recurring transactions
func (h *RecurringTransactionHandler) ProcessRecurringTransactions(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ProcessRecurringTransactions - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	result, err := h.service.ProcessRecurringTransactions(ctx, userID)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to process recurring transactions")
		return
	}

	logger.Infof(ctx, "Recurring transactions processed successfully for user: %s (created: %d, skipped: %d, errors: %d)",
		userID, result.Created, result.Skipped, result.Errors)
	utilities.SuccessResponse(c, result, "Processing completed")
}

// UpdateLastProcessed manually updates the last processed date for a recurring transaction
func (h *RecurringTransactionHandler) UpdateLastProcessed(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateLastProcessed - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warnf(ctx, "Invalid ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid recurring transaction ID")
		return
	}

	var request struct {
		LastProcessed string `json:"lastProcessed" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Warnf(ctx, "Validation error: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	recurringTransaction, err := h.service.UpdateLastProcessed(ctx, uint(id), userID, request.LastProcessed)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		if err.Error() == "recurring transaction not found" {
			utilities.ErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update last processed date")
		}
		return
	}

	logger.Infof(ctx, "Last processed date updated successfully for recurring transaction: %d", id)
	utilities.SuccessResponse(c, recurringTransaction, "Last processed date updated successfully")
}

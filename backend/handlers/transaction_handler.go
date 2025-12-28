package handlers

import (
	"net/http"
	"strconv"
	"time"

	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/repository"
	"daybook-backend/services"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TransactionHandler handles transaction HTTP requests
type TransactionHandler struct {
	service services.TransactionService
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// ListTransactions returns all transactions for the authenticated user with optional filters
func (h *TransactionHandler) ListTransactions(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListTransactions - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Build filters
	filters := repository.TransactionFilters{
		IncludeTracking: c.Query("includeTracking") == "true",
	}

	if transactionType := c.Query("type"); transactionType != "" {
		filters.Type = &transactionType
	}
	if categoryID := c.Query("categoryId"); categoryID != "" {
		filters.CategoryID = &categoryID
	}
	if accountID := c.Query("accountId"); accountID != "" {
		filters.AccountID = &accountID
	}
	if startDate := c.Query("startDate"); startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			filters.StartDate = &parsedDate
		}
	}
	if endDate := c.Query("endDate"); endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			filters.EndDate = &parsedDate
		}
	}

	// Build pagination
	pagination := repository.PaginationParams{
		Page:  1,
		Limit: 20,
	}

	if pageParam := c.Query("page"); pageParam != "" {
		if parsedPage, err := strconv.Atoi(pageParam); err == nil && parsedPage > 0 {
			pagination.Page = parsedPage
		}
	}

	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			// Validate limit is one of the allowed values: 20, 50, 100, 500
			switch parsedLimit {
			case 20, 50, 100, 500:
				pagination.Limit = parsedLimit
			default:
				pagination.Limit = 20
			}
		}
	}

	// Get transactions
	response, err := h.service.ListTransactions(ctx, userID, filters, pagination)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

	logger.Infof(ctx, "Transactions retrieved successfully for user: %s, count: %d, total: %d",
		userID, len(response.Transactions), response.Pagination.TotalCount)
	utilities.SuccessResponse(c, response, "Transactions retrieved successfully")
}

// GetTransaction returns a specific transaction by ID
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetTransaction - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid transaction ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	response, err := h.service.GetTransaction(ctx, transactionID, userID)
	if err != nil {
		logger.Warnf(ctx, "Transaction not found: %s, error: %v", transactionID, err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Transaction not found")
		return
	}

	logger.Infof(ctx, "Transaction retrieved successfully: %s for user: %s", transactionID, userID)
	utilities.SuccessResponse(c, response, "Transaction retrieved successfully")
}

// CreateTransaction creates a new transaction
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateTransaction - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		logger.Warnf(ctx, "Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	transaction.UserID = userID

	createdTransaction, err := h.service.CreateTransaction(ctx, &transaction)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Transaction created successfully: %s for user: %s", createdTransaction.ID, userID)
	utilities.CreatedResponse(c, createdTransaction, "Transaction created successfully")
}

// UpdateTransaction updates an existing transaction
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateTransaction - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid transaction ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	var updateData models.Transaction
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedTransaction, err := h.service.UpdateTransaction(ctx, transactionID, userID, &updateData)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Transaction updated successfully: %s for user: %s", transactionID, userID)
	utilities.SuccessResponse(c, updatedTransaction, "Transaction updated successfully")
}

// DeleteTransaction deletes a transaction
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteTransaction - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid transaction ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	if err := h.service.DeleteTransaction(ctx, transactionID, userID); err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	logger.Infof(ctx, "Transaction deleted successfully: %s for user: %s", transactionID, userID)
	utilities.SuccessResponse(c, nil, "Transaction deleted successfully")
}

// BulkImportTransactions imports multiple transactions
func (h *TransactionHandler) BulkImportTransactions(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "BulkImportTransactions - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var transactions []models.Transaction
	if err := c.ShouldBindJSON(&transactions); err != nil {
		logger.Warnf(ctx, "Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.BulkImportTransactions(ctx, userID, transactions)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to import transactions")
		return
	}

	logger.Infof(ctx, "Bulk import completed for user: %s (imported: %d, failed: %d)",
		userID, result.Imported, result.Failed)
	utilities.SuccessResponse(c, result, "Bulk import completed")
}

// GetTransactionStats returns transaction statistics
func (h *TransactionHandler) GetTransactionStats(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetTransactionStats - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Build filters (similar to ListTransactions)
	filters := repository.TransactionFilters{
		IncludeTracking: c.Query("includeTracking") == "true",
	}

	if categoryID := c.Query("categoryId"); categoryID != "" {
		filters.CategoryID = &categoryID
	}
	if accountID := c.Query("accountId"); accountID != "" {
		filters.AccountID = &accountID
	}
	if startDate := c.Query("startDate"); startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			filters.StartDate = &parsedDate
		}
	}
	if endDate := c.Query("endDate"); endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			filters.EndDate = &parsedDate
		}
	}

	stats, err := h.service.GetTransactionStats(ctx, userID, filters)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to calculate statistics")
		return
	}

	logger.Infof(ctx, "Transaction stats calculated successfully for user: %s", userID)
	utilities.SuccessResponse(c, stats, "Statistics calculated successfully")
}

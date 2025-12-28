package handlers

import (
	"net/http"

	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/repository"
	"daybook-backend/services"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReconciliationHandler handles reconciliation-related HTTP requests
type ReconciliationHandler struct {
	service services.ReconciliationService
}

// NewReconciliationHandler creates a new reconciliation handler
func NewReconciliationHandler(service services.ReconciliationService) *ReconciliationHandler {
	return &ReconciliationHandler{service: service}
}

// ListReconciliations returns all reconciliations for a specific account or user
func (h *ReconciliationHandler) ListReconciliations(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListReconciliations - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	// Build filters from query parameters
	filters := repository.ReconciliationFilters{}

	if accountID := c.Query("accountId"); accountID != "" {
		accID, err := uuid.Parse(accountID)
		if err != nil {
			logger.Warnf(ctx, "Invalid account ID: %v", err)
			utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
			return
		}
		filters.AccountID = &accID
	}

	logger.Debugf(ctx, "Fetching reconciliations from database...")
	reconciliations, err := h.service.ListReconciliations(ctx, userID, filters)
	if err != nil {
		logger.Errorf(ctx, "Database error fetching reconciliations: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch reconciliations")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d reconciliations for user: %s", len(reconciliations), userID)
	utilities.SuccessResponse(c, reconciliations, "Reconciliations retrieved successfully")
}

// GetReconciliation returns a specific reconciliation by ID
func (h *ReconciliationHandler) GetReconciliation(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetReconciliation - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	reconciliationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid reconciliation ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid reconciliation ID")
		return
	}

	logger.Debugf(ctx, "Fetching reconciliation: %s for user: %s", reconciliationID, userID)
	reconciliation, err := h.service.GetReconciliation(ctx, reconciliationID, userID)
	if err != nil {
		logger.Warnf(ctx, "Reconciliation not found: %s", reconciliationID)
		utilities.ErrorResponse(c, http.StatusNotFound, "Reconciliation not found")
		return
	}

	logger.Infof(ctx, "Successfully retrieved reconciliation for user: %s", userID)
	utilities.SuccessResponse(c, reconciliation, "Reconciliation retrieved successfully")
}

// CreateReconciliation creates a new reconciliation record
func (h *ReconciliationHandler) CreateReconciliation(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateReconciliation - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	var req services.CreateReconciliationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "Creating reconciliation for account: %s", req.AccountID)
	reconciliation, err := h.service.CreateReconciliation(ctx, userID, &req)
	if err != nil {
		logger.Errorf(ctx, "Failed to create reconciliation: %v", err)
		if err.Error() == "account not found" {
			utilities.ErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	logger.Infof(ctx, "Successfully created reconciliation for user: %s", userID)
	utilities.CreatedResponse(c, reconciliation, "Reconciliation created successfully")
}

// UpdateReconciliation updates an existing reconciliation
func (h *ReconciliationHandler) UpdateReconciliation(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateReconciliation - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	reconciliationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid reconciliation ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid reconciliation ID")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s, reconciliation: %s", userID, reconciliationID)

	var updateData models.Reconciliation
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "Updating reconciliation...")
	updated, err := h.service.UpdateReconciliation(ctx, reconciliationID, userID, &updateData)
	if err != nil {
		logger.Errorf(ctx, "Failed to update reconciliation: %v", err)
		if err.Error() == "reconciliation not found" {
			utilities.ErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	logger.Infof(ctx, "Successfully updated reconciliation for user: %s", userID)
	utilities.SuccessResponse(c, updated, "Reconciliation updated successfully")
}

// DeleteReconciliation deletes a reconciliation record
func (h *ReconciliationHandler) DeleteReconciliation(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteReconciliation - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	reconciliationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid reconciliation ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid reconciliation ID")
		return
	}

	logger.Debugf(ctx, "Deleting reconciliation: %s for user: %s", reconciliationID, userID)
	if err := h.service.DeleteReconciliation(ctx, reconciliationID, userID); err != nil {
		logger.Errorf(ctx, "Failed to delete reconciliation: %v", err)
		if err.Error() == "reconciliation not found" {
			utilities.ErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	logger.Infof(ctx, "Successfully deleted reconciliation for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Reconciliation deleted successfully")
}

// GetUnreconciledTransactions returns all unreconciled transactions for an account
func (h *ReconciliationHandler) GetUnreconciledTransactions(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetUnreconciledTransactions - Entry")

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

	logger.Debugf(ctx, "Fetching unreconciled transactions for account: %s, user: %s", accountID, userID)
	transactions, err := h.service.GetUnreconciledTransactions(ctx, accountID, userID)
	if err != nil {
		logger.Errorf(ctx, "Failed to fetch unreconciled transactions: %v", err)
		if err.Error() == "account not found" {
			utilities.ErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch transactions")
		}
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d unreconciled transactions for user: %s", len(transactions), userID)
	utilities.SuccessResponse(c, transactions, "Unreconciled transactions retrieved successfully")
}

// GetReconciliationStats returns reconciliation statistics for an account
func (h *ReconciliationHandler) GetReconciliationStats(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetReconciliationStats - Entry")

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

	logger.Debugf(ctx, "Fetching reconciliation stats for account: %s, user: %s", accountID, userID)
	stats, err := h.service.GetStats(ctx, accountID, userID)
	if err != nil {
		logger.Errorf(ctx, "Failed to get reconciliation stats: %v", err)
		if err.Error() == "account not found" {
			utilities.ErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch stats")
		}
		return
	}

	logger.Infof(ctx, "Successfully retrieved reconciliation stats for user: %s. Total: %d, Completed: %d",
		userID, stats.TotalReconciliations, stats.CompletedReconciliations)
	utilities.SuccessResponse(c, stats, "Reconciliation stats retrieved successfully")
}

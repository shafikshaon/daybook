package handlers

import (
	"net/http"
	"strconv"

	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/repository"
	"daybook-backend/services"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
)

// DebtHandler handles debt-related HTTP requests
type DebtHandler struct {
	service services.DebtService
}

// NewDebtHandler creates a new debt handler
func NewDebtHandler(service services.DebtService) *DebtHandler {
	return &DebtHandler{service: service}
}

// ListDebts returns all debt records for the authenticated user
func (h *DebtHandler) ListDebts(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListDebts - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	// Build filters from query parameters
	filters := repository.DebtFilters{}

	if status := c.Query("status"); status != "" {
		filters.Status = &status
	}

	logger.Debugf(ctx, "Fetching debts from database...")
	debts, err := h.service.ListDebts(ctx, userID, filters)
	if err != nil {
		logger.Errorf(ctx, "Database error fetching debts: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch debts")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d debts for user: %s", len(debts), userID)
	utilities.SuccessResponse(c, debts, "Debts retrieved successfully")
}

// GetDebt returns a specific debt record by ID
func (h *DebtHandler) GetDebt(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetDebt - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	debtIDStr := c.Param("id")
	debtIDUint, err := strconv.ParseUint(debtIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "Invalid debt ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid debt ID")
		return
	}
	debtID := uint(debtIDUint)

	logger.Debugf(ctx, "Fetching debt: %s for user: %s", debtID, userID)
	debt, err := h.service.GetDebt(ctx, debtID, userID)
	if err != nil {
		logger.Errorf(ctx, "Database error fetching debt: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Debt not found")
		return
	}

	logger.Infof(ctx, "Successfully retrieved debt for user: %s", userID)
	utilities.SuccessResponse(c, debt, "Debt retrieved successfully")
}

// CreateDebt creates a new debt record
func (h *DebtHandler) CreateDebt(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateDebt - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	var debt models.DebtRecord
	if err := c.ShouldBindJSON(&debt); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	debt.UserID = userID

	logger.Debugf(ctx, "Creating debt...")
	created, err := h.service.CreateDebt(ctx, &debt)
	if err != nil {
		logger.Errorf(ctx, "Failed to create debt: %v", err)
		if err.Error() == "invalid account ID" {
			utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	logger.Infof(ctx, "Successfully created debt for user: %s", userID)
	utilities.SuccessResponse(c, created, "Debt created successfully")
}

// UpdateDebt updates a debt record
func (h *DebtHandler) UpdateDebt(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateDebt - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	debtIDStr := c.Param("id")
	debtIDUint, err := strconv.ParseUint(debtIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "Invalid debt ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid debt ID")
		return
	}
	debtID := uint(debtIDUint)

	logger.Debugf(ctx, "Processing request for user: %s, debt: %s", userID, debtID)

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "Updating debt...")
	updated, err := h.service.UpdateDebt(ctx, debtID, userID, updateData)
	if err != nil {
		logger.Errorf(ctx, "Failed to update debt: %v", err)
		if err.Error() == "debt not found" {
			utilities.ErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	logger.Infof(ctx, "Successfully updated debt for user: %s", userID)
	utilities.SuccessResponse(c, updated, "Debt updated successfully")
}

// DeleteDebt soft deletes a debt record
func (h *DebtHandler) DeleteDebt(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteDebt - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	debtIDStr := c.Param("id")
	debtIDUint, err := strconv.ParseUint(debtIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "Invalid debt ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid debt ID")
		return
	}
	debtID := uint(debtIDUint)

	logger.Debugf(ctx, "Deleting debt: %s for user: %s", debtID, userID)
	if err := h.service.DeleteDebt(ctx, debtID, userID); err != nil {
		logger.Errorf(ctx, "Failed to delete debt: %v", err)
		if err.Error() == "debt not found" {
			utilities.ErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	logger.Infof(ctx, "Successfully deleted debt for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Debt deleted successfully")
}

// RecordDebtPayment records a payment towards a debt
func (h *DebtHandler) RecordDebtPayment(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "RecordDebtPayment - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	debtIDStr := c.Param("id")
	debtIDUint, err := strconv.ParseUint(debtIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "Invalid debt ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid debt ID")
		return
	}
	debtID := uint(debtIDUint)

	logger.Debugf(ctx, "Processing request for user: %s, debt: %s", userID, debtID)

	var payment models.DebtPayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "Recording debt payment...")
	createdPayment, updatedDebt, err := h.service.RecordPayment(ctx, debtID, userID, &payment)
	if err != nil {
		logger.Errorf(ctx, "Failed to record payment: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Successfully recorded debt payment for user: %s", userID)
	utilities.SuccessResponse(c, map[string]interface{}{
		"payment": createdPayment,
		"debt":    updatedDebt,
	}, "Payment recorded successfully")
}

// ListDebtPayments returns all payments for a specific debt
func (h *DebtHandler) ListDebtPayments(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListDebtPayments - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	debtIDStr := c.Param("id")
	debtIDUint, err := strconv.ParseUint(debtIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "Invalid debt ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid debt ID")
		return
	}
	debtID := uint(debtIDUint)

	logger.Debugf(ctx, "Fetching debt payments for debt: %s, user: %s", debtID, userID)
	payments, err := h.service.ListPayments(ctx, debtID, userID)
	if err != nil {
		logger.Errorf(ctx, "Failed to fetch payments: %v", err)
		if err.Error() == "debt not found" {
			utilities.ErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch payments")
		}
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d payments for user: %s", len(payments), userID)
	utilities.SuccessResponse(c, payments, "Payments retrieved successfully")
}

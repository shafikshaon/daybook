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

// LendHandler handles lend-related HTTP requests
type LendHandler struct {
	service services.LendService
}

// NewLendHandler creates a new lend handler
func NewLendHandler(service services.LendService) *LendHandler {
	return &LendHandler{service: service}
}

// ListLends returns all lend records for the authenticated user
func (h *LendHandler) ListLends(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListLends - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	// Build filters from query parameters
	filters := repository.LendFilters{}

	if status := c.Query("status"); status != "" {
		filters.Status = &status
	}

	logger.Debugf(ctx, "Fetching lends from database...")
	lends, err := h.service.ListLends(ctx, userID, filters)
	if err != nil {
		logger.Errorf(ctx, "Database error fetching lends: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch lends")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d lends for user: %s", len(lends), userID)
	utilities.SuccessResponse(c, lends, "Lends retrieved successfully")
}

// GetLend returns a specific lend record by ID
func (h *LendHandler) GetLend(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetLend - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	lendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid lend ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid lend ID")
		return
	}

	logger.Debugf(ctx, "Fetching lend: %s for user: %s", lendID, userID)
	lend, err := h.service.GetLend(ctx, lendID, userID)
	if err != nil {
		logger.Errorf(ctx, "Database error fetching lend: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Lend not found")
		return
	}

	logger.Infof(ctx, "Successfully retrieved lend for user: %s", userID)
	utilities.SuccessResponse(c, lend, "Lend retrieved successfully")
}

// CreateLend creates a new lend record
func (h *LendHandler) CreateLend(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateLend - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	var lend models.LendRecord
	if err := c.ShouldBindJSON(&lend); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	lend.UserID = userID

	logger.Debugf(ctx, "Creating lend...")
	created, err := h.service.CreateLend(ctx, &lend)
	if err != nil {
		logger.Errorf(ctx, "Failed to create lend: %v", err)
		if err.Error() == "invalid account ID" {
			utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	logger.Infof(ctx, "Successfully created lend for user: %s", userID)
	utilities.SuccessResponse(c, created, "Lend created successfully")
}

// UpdateLend updates a lend record
func (h *LendHandler) UpdateLend(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateLend - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	lendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid lend ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid lend ID")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s, lend: %s", userID, lendID)

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "Updating lend...")
	updated, err := h.service.UpdateLend(ctx, lendID, userID, updateData)
	if err != nil {
		logger.Errorf(ctx, "Failed to update lend: %v", err)
		if err.Error() == "lend not found" {
			utilities.ErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	logger.Infof(ctx, "Successfully updated lend for user: %s", userID)
	utilities.SuccessResponse(c, updated, "Lend updated successfully")
}

// DeleteLend soft deletes a lend record
func (h *LendHandler) DeleteLend(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteLend - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	lendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid lend ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid lend ID")
		return
	}

	logger.Debugf(ctx, "Deleting lend: %s for user: %s", lendID, userID)
	if err := h.service.DeleteLend(ctx, lendID, userID); err != nil {
		logger.Errorf(ctx, "Failed to delete lend: %v", err)
		if err.Error() == "lend not found" {
			utilities.ErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	logger.Infof(ctx, "Successfully deleted lend for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Lend deleted successfully")
}

// RecordLendPayment records a payment received for a lend
func (h *LendHandler) RecordLendPayment(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "RecordLendPayment - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	lendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid lend ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid lend ID")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s, lend: %s", userID, lendID)

	var payment models.LendPayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "Recording lend payment...")
	createdPayment, updatedLend, err := h.service.RecordPayment(ctx, lendID, userID, &payment)
	if err != nil {
		logger.Errorf(ctx, "Failed to record payment: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Successfully recorded lend payment for user: %s", userID)
	utilities.SuccessResponse(c, map[string]interface{}{
		"payment": createdPayment,
		"lend":    updatedLend,
	}, "Payment recorded successfully")
}

// ListLendPayments returns all payments for a specific lend
func (h *LendHandler) ListLendPayments(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListLendPayments - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	lendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid lend ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid lend ID")
		return
	}

	logger.Debugf(ctx, "Fetching lend payments for lend: %s, user: %s", lendID, userID)
	payments, err := h.service.ListPayments(ctx, lendID, userID)
	if err != nil {
		logger.Errorf(ctx, "Failed to fetch payments: %v", err)
		if err.Error() == "lend not found" {
			utilities.ErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch payments")
		}
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d payments for user: %s", len(payments), userID)
	utilities.SuccessResponse(c, payments, "Payments retrieved successfully")
}

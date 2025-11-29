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

// ListLends returns all lend records for the authenticated user
func ListLends(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListLends - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "ListLends - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "ListLends - Fetching lends for user: %s", userID)

	query := database.DB.WithContext(ctx).Where("user_id = ?", userID)

	// Apply status filter
	if status := c.Query("status"); status != "" {
		logger.Debugf(ctx, "ListLends - Filtering by status: %s", status)
		query = query.Where("status = ?", status)
	}

	var lends []models.LendRecord
	if err := query.Order("lent_date DESC, created_at DESC").Find(&lends).Error; err != nil {
		logger.Errorf(ctx, "ListLends - Failed to fetch lends: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch lends")
		return
	}

	// Enrich lends with account names
	type LendResponse struct {
		models.LendRecord
		AccountName *string `json:"accountName,omitempty"`
	}

	enrichedLends := make([]LendResponse, len(lends))
	for i, lend := range lends {
		enrichedLends[i] = LendResponse{LendRecord: lend}

		if lend.AccountID != nil {
			var account models.Account
			if err := database.DB.WithContext(ctx).Select("name").Where("id = ?", *lend.AccountID).First(&account).Error; err == nil {
				enrichedLends[i].AccountName = &account.Name
			}
		}
	}

	logger.Infof(ctx, "ListLends - Successfully retrieved %d lends", len(lends))
	utilities.SuccessResponse(c, enrichedLends, "Lends retrieved successfully")
}

// GetLend returns a specific lend record by ID
func GetLend(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetLend - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "GetLend - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	lendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "GetLend - Invalid lend ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid lend ID")
		return
	}

	logger.Debugf(ctx, "GetLend - Fetching lend: %s for user: %s", lendID, userID)

	var lend models.LendRecord
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", lendID, userID).First(&lend).Error; err != nil {
		logger.Warnf(ctx, "GetLend - Lend not found: %s", lendID)
		utilities.ErrorResponse(c, http.StatusNotFound, "Lend not found")
		return
	}

	// Enrich with account name
	type LendResponse struct {
		models.LendRecord
		AccountName *string `json:"accountName,omitempty"`
	}

	response := LendResponse{LendRecord: lend}

	if lend.AccountID != nil {
		var account models.Account
		if err := database.DB.WithContext(ctx).Select("name").Where("id = ?", *lend.AccountID).First(&account).Error; err == nil {
			response.AccountName = &account.Name
		}
	}

	logger.Infof(ctx, "GetLend - Successfully retrieved lend: %s", lendID)
	utilities.SuccessResponse(c, response, "Lend retrieved successfully")
}

// CreateLend creates a new lend record
func CreateLend(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateLend - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "CreateLend - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var lend models.LendRecord
	if err := c.ShouldBindJSON(&lend); err != nil {
		logger.Warnf(ctx, "CreateLend - Invalid request data: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate lent date is provided
	if lend.LentDate.IsZero() {
		logger.Warnf(ctx, "CreateLend - Lent date is required")
		utilities.ErrorResponse(c, http.StatusBadRequest, "Lent date is required")
		return
	}

	logger.Debugf(ctx, "CreateLend - Creating lend for %s, user: %s", lend.DebtorName, userID)

	lend.UserID = userID

	// If account is specified, verify it belongs to user and create transaction
	if lend.AccountID != nil {
		var account models.Account
		if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", *lend.AccountID, userID).First(&account).Error; err != nil {
			logger.Warnf(ctx, "CreateLend - Invalid account ID: %s", *lend.AccountID)
			utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
			return
		}

		// Start transaction
		tx := database.DB.WithContext(ctx).Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Create the lend record
		if err := tx.Create(&lend).Error; err != nil {
			tx.Rollback()
			logger.Errorf(ctx, "CreateLend - Failed to create lend: %v", err)
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create lend")
			return
		}

		// If not initial lend, create transaction and update account balance
		if !lend.IsInitial {
			logger.Debugf(ctx, "CreateLend - Not initial lend, creating transaction and updating balance")
			// Validate account has sufficient balance
			if account.Balance < lend.OriginalAmount {
				tx.Rollback()
				logger.Warnf(ctx, "CreateLend - Insufficient account balance. Required: %.2f, Available: %.2f", lend.OriginalAmount, account.Balance)
				utilities.ErrorResponse(c, http.StatusBadRequest, "Insufficient account balance")
				return
			}

			// Create transaction record for the lent money
			transaction := models.Transaction{
				UserID:      userID,
				AccountID:   *lend.AccountID,
				Type:        "expense",
				Amount:      lend.OriginalAmount,
				CategoryID:  "lend",
				Date:        lend.LentDate,
				Description: "Lent to " + lend.DebtorName,
			}

			if err := tx.Create(&transaction).Error; err != nil {
				tx.Rollback()
				logger.Errorf(ctx, "CreateLend - Failed to create transaction: %v", err)
				utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
				return
			}

			// Update account balance
			if err := tx.Where("id = ?", *lend.AccountID).First(&account).Error; err != nil {
				tx.Rollback()
				logger.Errorf(ctx, "CreateLend - Failed to fetch account: %v", err)
				utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch account")
				return
			}

			account.Balance -= lend.OriginalAmount

			if err := tx.Save(&account).Error; err != nil {
				tx.Rollback()
				logger.Errorf(ctx, "CreateLend - Failed to update account balance: %v", err)
				utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account balance")
				return
			}
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			logger.Errorf(ctx, "CreateLend - Failed to commit transaction: %v", err)
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction")
			return
		}
	} else {
		logger.Debugf(ctx, "CreateLend - No account specified, creating standalone lend record")
		// No account specified, just create the lend record
		if err := database.DB.WithContext(ctx).Create(&lend).Error; err != nil {
			logger.Errorf(ctx, "CreateLend - Failed to create lend: %v", err)
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create lend")
			return
		}
	}

	logger.Infof(ctx, "CreateLend - Successfully created lend: %s", lend.ID)

	// Log lend creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleLend,
		"Lend", lend.ID, "Created lend record for "+lend.DebtorName, nil)

	utilities.SuccessResponse(c, lend, "Lend created successfully")
}

// UpdateLend updates a lend record
func UpdateLend(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateLend - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "UpdateLend - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	lendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "UpdateLend - Invalid lend ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid lend ID")
		return
	}

	logger.Debugf(ctx, "UpdateLend - Updating lend: %s for user: %s", lendID, userID)

	var existingLend models.LendRecord
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", lendID, userID).First(&existingLend).Error; err != nil {
		logger.Warnf(ctx, "UpdateLend - Lend not found: %s", lendID)
		utilities.ErrorResponse(c, http.StatusNotFound, "Lend not found")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "UpdateLend - Invalid request data: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Prevent updating certain fields
	delete(updateData, "id")
	delete(updateData, "userId")
	delete(updateData, "originalAmount")
	delete(updateData, "remainingAmount")
	delete(updateData, "createdAt")

	if err := database.DB.WithContext(ctx).Model(&existingLend).Updates(updateData).Error; err != nil {
		logger.Errorf(ctx, "UpdateLend - Failed to update lend: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update lend")
		return
	}

	logger.Infof(ctx, "UpdateLend - Successfully updated lend: %s", lendID)

	// Log lend update activity
	utilities.LogEntityActivity(c, userID, models.ActionUpdate, models.ModuleLend,
		"Lend", existingLend.ID, "Updated lend record for "+existingLend.DebtorName, nil)

	utilities.SuccessResponse(c, existingLend, "Lend updated successfully")
}

// DeleteLend soft deletes a lend record
func DeleteLend(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteLend - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "DeleteLend - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	lendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "DeleteLend - Invalid lend ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid lend ID")
		return
	}

	logger.Debugf(ctx, "DeleteLend - Deleting lend: %s for user: %s", lendID, userID)

	var lend models.LendRecord
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", lendID, userID).First(&lend).Error; err != nil {
		logger.Warnf(ctx, "DeleteLend - Lend not found: %s", lendID)
		utilities.ErrorResponse(c, http.StatusNotFound, "Lend not found")
		return
	}

	if err := database.DB.WithContext(ctx).Delete(&lend).Error; err != nil {
		logger.Errorf(ctx, "DeleteLend - Failed to delete lend: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete lend")
		return
	}

	logger.Infof(ctx, "DeleteLend - Successfully deleted lend: %s", lendID)

	// Log lend deletion activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, models.ModuleLend,
		"Lend", lend.ID, "Deleted lend record for "+lend.DebtorName, nil)

	utilities.SuccessResponse(c, nil, "Lend deleted successfully")
}

// RecordLendPayment records a payment received for a lend
func RecordLendPayment(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "RecordLendPayment - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "RecordLendPayment - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	lendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "RecordLendPayment - Invalid lend ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid lend ID")
		return
	}

	var payment models.LendPayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		logger.Warnf(ctx, "RecordLendPayment - Invalid request data: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate payment date is provided
	if payment.PaymentDate.IsZero() {
		logger.Warnf(ctx, "RecordLendPayment - Payment date is required")
		utilities.ErrorResponse(c, http.StatusBadRequest, "Payment date is required")
		return
	}

	logger.Debugf(ctx, "RecordLendPayment - Recording payment for lend: %s, user: %s", lendID, userID)

	payment.UserID = userID
	payment.LendID = lendID

	// Verify account belongs to user
	var account models.Account
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", payment.AccountID, userID).First(&account).Error; err != nil {
		logger.Warnf(ctx, "RecordLendPayment - Invalid account ID: %s", payment.AccountID)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	// Start transaction
	tx := database.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get lend record
	var lend models.LendRecord
	if err := tx.Where("id = ? AND user_id = ?", lendID, userID).First(&lend).Error; err != nil {
		tx.Rollback()
		logger.Warnf(ctx, "RecordLendPayment - Lend not found: %s", lendID)
		utilities.ErrorResponse(c, http.StatusNotFound, "Lend not found")
		return
	}

	// Validate payment amount doesn't exceed remaining amount
	if payment.Amount > lend.RemainingAmount {
		tx.Rollback()
		logger.Warnf(ctx, "RecordLendPayment - Payment amount (%.2f) exceeds remaining lend (%.2f)", payment.Amount, lend.RemainingAmount)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Payment amount exceeds remaining lend")
		return
	}

	// Create payment record
	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "RecordLendPayment - Failed to create payment: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create payment")
		return
	}

	// Update lend remaining amount and status
	lend.RemainingAmount -= payment.Amount
	if lend.RemainingAmount == 0 {
		logger.Debugf(ctx, "RecordLendPayment - Lend fully received")
		lend.Status = "fully_received"
	} else if lend.RemainingAmount < lend.OriginalAmount {
		logger.Debugf(ctx, "RecordLendPayment - Lend partially received. Remaining: %.2f", lend.RemainingAmount)
		lend.Status = "partially_received"
	}

	if err := tx.Save(&lend).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "RecordLendPayment - Failed to update lend: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update lend")
		return
	}

	// Create transaction record for the payment
	description := "Payment from " + lend.DebtorName
	if payment.Description != "" {
		description = payment.Description
	}

	transaction := models.Transaction{
		UserID:      userID,
		AccountID:   payment.AccountID,
		Type:        "income",
		Amount:      payment.Amount,
		CategoryID:  "lend_payment",
		Date:        payment.PaymentDate,
		Description: description,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "RecordLendPayment - Failed to create transaction: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	// Update account balance
	account.Balance += payment.Amount

	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "RecordLendPayment - Failed to update account balance: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account balance")
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		logger.Errorf(ctx, "RecordLendPayment - Failed to commit transaction: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	logger.Infof(ctx, "RecordLendPayment - Successfully recorded payment: %s for lend: %s", payment.ID, lendID)

	// Log lend payment activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleLend,
		"LendPayment", payment.ID, "Recorded payment from "+lend.DebtorName, nil)

	utilities.SuccessResponse(c, map[string]interface{}{
		"payment": payment,
		"lend":    lend,
	}, "Payment recorded successfully")
}

// ListLendPayments returns all payments for a specific lend
func ListLendPayments(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListLendPayments - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "ListLendPayments - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	lendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "ListLendPayments - Invalid lend ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid lend ID")
		return
	}

	logger.Debugf(ctx, "ListLendPayments - Fetching payments for lend: %s, user: %s", lendID, userID)

	// Verify lend belongs to user
	var lend models.LendRecord
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", lendID, userID).First(&lend).Error; err != nil {
		logger.Warnf(ctx, "ListLendPayments - Lend not found: %s", lendID)
		utilities.ErrorResponse(c, http.StatusNotFound, "Lend not found")
		return
	}

	var payments []models.LendPayment
	if err := database.DB.WithContext(ctx).Where("lend_id = ? AND user_id = ?", lendID, userID).
		Order("payment_date DESC, created_at DESC").
		Find(&payments).Error; err != nil {
		logger.Errorf(ctx, "ListLendPayments - Failed to fetch payments: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch payments")
		return
	}

	// Enrich payments with account names
	type PaymentResponse struct {
		models.LendPayment
		AccountName string `json:"accountName"`
	}

	enrichedPayments := make([]PaymentResponse, len(payments))
	for i, payment := range payments {
		enrichedPayments[i] = PaymentResponse{LendPayment: payment}

		var account models.Account
		if err := database.DB.WithContext(ctx).Select("name").Where("id = ?", payment.AccountID).First(&account).Error; err == nil {
			enrichedPayments[i].AccountName = account.Name
		}
	}

	logger.Infof(ctx, "ListLendPayments - Successfully retrieved %d payments", len(payments))
	utilities.SuccessResponse(c, enrichedPayments, "Payments retrieved successfully")
}

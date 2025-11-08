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

// ListLends returns all lend records for the authenticated user
func ListLends(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	query := database.DB.Where("user_id = ?", userID)

	// Apply status filter
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var lends []models.LendRecord
	if err := query.Order("lent_date DESC, created_at DESC").Find(&lends).Error; err != nil {
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
			if err := database.DB.Select("name").Where("id = ?", *lend.AccountID).First(&account).Error; err == nil {
				enrichedLends[i].AccountName = &account.Name
			}
		}
	}

	utilities.SuccessResponse(c, enrichedLends, "Lends retrieved successfully")
}

// GetLend returns a specific lend record by ID
func GetLend(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	lendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid lend ID")
		return
	}

	var lend models.LendRecord
	if err := database.DB.Where("id = ? AND user_id = ?", lendID, userID).First(&lend).Error; err != nil {
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
		if err := database.DB.Select("name").Where("id = ?", *lend.AccountID).First(&account).Error; err == nil {
			response.AccountName = &account.Name
		}
	}

	utilities.SuccessResponse(c, response, "Lend retrieved successfully")
}

// CreateLend creates a new lend record
func CreateLend(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var lend models.LendRecord
	if err := c.ShouldBindJSON(&lend); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate lent date is provided
	if lend.LentDate.IsZero() {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Lent date is required")
		return
	}

	lend.UserID = userID

	// If account is specified, verify it belongs to user and create transaction
	if lend.AccountID != nil {
		var account models.Account
		if err := database.DB.Where("id = ? AND user_id = ?", *lend.AccountID, userID).First(&account).Error; err != nil {
			utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
			return
		}

		// Start transaction
		tx := database.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Create the lend record
		if err := tx.Create(&lend).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create lend")
			return
		}

		// If not initial lend, create transaction and update account balance
		if !lend.IsInitial {
			// Validate account has sufficient balance
			if account.Balance < lend.OriginalAmount {
				tx.Rollback()
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
				Date:        lend.LentDate.Time,
				Description: "Lent to " + lend.DebtorName,
			}

			if err := tx.Create(&transaction).Error; err != nil {
				tx.Rollback()
				utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
				return
			}

			// Update account balance
			if err := tx.Where("id = ?", *lend.AccountID).First(&account).Error; err != nil {
				tx.Rollback()
				utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch account")
				return
			}

			account.Balance -= lend.OriginalAmount

			if err := tx.Save(&account).Error; err != nil {
				tx.Rollback()
				utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account balance")
				return
			}
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction")
			return
		}
	} else {
		// No account specified, just create the lend record
		if err := database.DB.Create(&lend).Error; err != nil {
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create lend")
			return
		}
	}

	utilities.SuccessResponse(c, lend, "Lend created successfully")
}

// UpdateLend updates a lend record
func UpdateLend(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	lendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid lend ID")
		return
	}

	var existingLend models.LendRecord
	if err := database.DB.Where("id = ? AND user_id = ?", lendID, userID).First(&existingLend).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Lend not found")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Prevent updating certain fields
	delete(updateData, "id")
	delete(updateData, "userId")
	delete(updateData, "originalAmount")
	delete(updateData, "remainingAmount")
	delete(updateData, "createdAt")

	if err := database.DB.Model(&existingLend).Updates(updateData).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update lend")
		return
	}

	utilities.SuccessResponse(c, existingLend, "Lend updated successfully")
}

// DeleteLend soft deletes a lend record
func DeleteLend(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	lendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid lend ID")
		return
	}

	var lend models.LendRecord
	if err := database.DB.Where("id = ? AND user_id = ?", lendID, userID).First(&lend).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Lend not found")
		return
	}

	if err := database.DB.Delete(&lend).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete lend")
		return
	}

	utilities.SuccessResponse(c, nil, "Lend deleted successfully")
}

// RecordLendPayment records a payment received for a lend
func RecordLendPayment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	lendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid lend ID")
		return
	}

	var payment models.LendPayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate payment date is provided
	if payment.PaymentDate.IsZero() {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Payment date is required")
		return
	}

	payment.UserID = userID
	payment.LendID = lendID

	// Verify account belongs to user
	var account models.Account
	if err := database.DB.Where("id = ? AND user_id = ?", payment.AccountID, userID).First(&account).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get lend record
	var lend models.LendRecord
	if err := tx.Where("id = ? AND user_id = ?", lendID, userID).First(&lend).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusNotFound, "Lend not found")
		return
	}

	// Validate payment amount doesn't exceed remaining amount
	if payment.Amount > lend.RemainingAmount {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusBadRequest, "Payment amount exceeds remaining lend")
		return
	}

	// Create payment record
	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create payment")
		return
	}

	// Update lend remaining amount and status
	lend.RemainingAmount -= payment.Amount
	if lend.RemainingAmount == 0 {
		lend.Status = "fully_received"
	} else if lend.RemainingAmount < lend.OriginalAmount {
		lend.Status = "partially_received"
	}

	if err := tx.Save(&lend).Error; err != nil {
		tx.Rollback()
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
		Date:        payment.PaymentDate.Time,
		Description: description,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	// Update account balance
	account.Balance += payment.Amount

	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account balance")
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	utilities.SuccessResponse(c, map[string]interface{}{
		"payment": payment,
		"lend":    lend,
	}, "Payment recorded successfully")
}

// ListLendPayments returns all payments for a specific lend
func ListLendPayments(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	lendID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid lend ID")
		return
	}

	// Verify lend belongs to user
	var lend models.LendRecord
	if err := database.DB.Where("id = ? AND user_id = ?", lendID, userID).First(&lend).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Lend not found")
		return
	}

	var payments []models.LendPayment
	if err := database.DB.Where("lend_id = ? AND user_id = ?", lendID, userID).
		Order("payment_date DESC, created_at DESC").
		Find(&payments).Error; err != nil {
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
		if err := database.DB.Select("name").Where("id = ?", payment.AccountID).First(&account).Error; err == nil {
			enrichedPayments[i].AccountName = account.Name
		}
	}

	utilities.SuccessResponse(c, enrichedPayments, "Payments retrieved successfully")
}

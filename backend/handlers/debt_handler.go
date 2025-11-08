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

// ListDebts returns all debt records for the authenticated user
func ListDebts(c *gin.Context) {
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

	var debts []models.DebtRecord
	if err := query.Order("borrowed_date DESC, created_at DESC").Find(&debts).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch debts")
		return
	}

	// Enrich debts with account names
	type DebtResponse struct {
		models.DebtRecord
		AccountName *string `json:"accountName,omitempty"`
	}

	enrichedDebts := make([]DebtResponse, len(debts))
	for i, debt := range debts {
		enrichedDebts[i] = DebtResponse{DebtRecord: debt}

		if debt.AccountID != nil {
			var account models.Account
			if err := database.DB.Select("name").Where("id = ?", *debt.AccountID).First(&account).Error; err == nil {
				enrichedDebts[i].AccountName = &account.Name
			}
		}
	}

	utilities.SuccessResponse(c, enrichedDebts, "Debts retrieved successfully")
}

// GetDebt returns a specific debt record by ID
func GetDebt(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	debtID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid debt ID")
		return
	}

	var debt models.DebtRecord
	if err := database.DB.Where("id = ? AND user_id = ?", debtID, userID).First(&debt).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Debt not found")
		return
	}

	// Enrich with account name
	type DebtResponse struct {
		models.DebtRecord
		AccountName *string `json:"accountName,omitempty"`
	}

	response := DebtResponse{DebtRecord: debt}

	if debt.AccountID != nil {
		var account models.Account
		if err := database.DB.Select("name").Where("id = ?", *debt.AccountID).First(&account).Error; err == nil {
			response.AccountName = &account.Name
		}
	}

	utilities.SuccessResponse(c, response, "Debt retrieved successfully")
}

// CreateDebt creates a new debt record
func CreateDebt(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var debt models.DebtRecord
	if err := c.ShouldBindJSON(&debt); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate borrowed date is provided
	if debt.BorrowedDate.IsZero() {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Borrowed date is required")
		return
	}

	debt.UserID = userID

	// If account is specified, verify it belongs to user and create transaction
	if debt.AccountID != nil {
		var account models.Account
		if err := database.DB.Where("id = ? AND user_id = ?", *debt.AccountID, userID).First(&account).Error; err != nil {
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

		// Create the debt record
		if err := tx.Create(&debt).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create debt")
			return
		}

		// If not initial debt, create transaction and update account balance
		if !debt.IsInitial {
			// Create transaction record for the borrowed money
			transaction := models.Transaction{
				UserID:      userID,
				AccountID:   *debt.AccountID,
				Type:        "income",
				Amount:      debt.OriginalAmount,
				CategoryID:  "debt",
				Date:        debt.BorrowedDate.Time,
				Description: "Borrowed from " + debt.CreditorName,
			}

			if err := tx.Create(&transaction).Error; err != nil {
				tx.Rollback()
				utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
				return
			}

			// Update account balance
			if err := tx.Where("id = ?", *debt.AccountID).First(&account).Error; err != nil {
				tx.Rollback()
				utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch account")
				return
			}

			account.Balance += debt.OriginalAmount

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
		// No account specified, just create the debt record
		if err := database.DB.Create(&debt).Error; err != nil {
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create debt")
			return
		}
	}

	// Log debt creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleDebt,
		"Debt", debt.ID, "Created debt: "+debt.CreditorName, nil)

	utilities.SuccessResponse(c, debt, "Debt created successfully")
}

// UpdateDebt updates a debt record
func UpdateDebt(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	debtID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid debt ID")
		return
	}

	var existingDebt models.DebtRecord
	if err := database.DB.Where("id = ? AND user_id = ?", debtID, userID).First(&existingDebt).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Debt not found")
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

	if err := database.DB.Model(&existingDebt).Updates(updateData).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update debt")
		return
	}

	// Log debt update activity
	utilities.LogEntityActivity(c, userID, models.ActionUpdate, models.ModuleDebt,
		"Debt", existingDebt.ID, "Updated debt: "+existingDebt.CreditorName, nil)

	utilities.SuccessResponse(c, existingDebt, "Debt updated successfully")
}

// DeleteDebt soft deletes a debt record
func DeleteDebt(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	debtID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid debt ID")
		return
	}

	var debt models.DebtRecord
	if err := database.DB.Where("id = ? AND user_id = ?", debtID, userID).First(&debt).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Debt not found")
		return
	}

	if err := database.DB.Delete(&debt).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete debt")
		return
	}

	// Log debt deletion activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, models.ModuleDebt,
		"Debt", debt.ID, "Deleted debt: "+debt.CreditorName, nil)

	utilities.SuccessResponse(c, nil, "Debt deleted successfully")
}

// RecordDebtPayment records a payment towards a debt
func RecordDebtPayment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	debtID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid debt ID")
		return
	}

	var payment models.DebtPayment
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
	payment.DebtID = debtID

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

	// Get debt record
	var debt models.DebtRecord
	if err := tx.Where("id = ? AND user_id = ?", debtID, userID).First(&debt).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusNotFound, "Debt not found")
		return
	}

	// Validate payment amount doesn't exceed remaining amount
	if payment.Amount > debt.RemainingAmount {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusBadRequest, "Payment amount exceeds remaining debt")
		return
	}

	// Validate account has sufficient balance
	if account.Balance < payment.Amount {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusBadRequest, "Insufficient account balance")
		return
	}

	// Create payment record
	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create payment")
		return
	}

	// Update debt remaining amount and status
	debt.RemainingAmount -= payment.Amount
	if debt.RemainingAmount == 0 {
		debt.Status = "fully_paid"
	} else if debt.RemainingAmount < debt.OriginalAmount {
		debt.Status = "partially_paid"
	}

	if err := tx.Save(&debt).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update debt")
		return
	}

	// Create transaction record for the payment
	description := "Payment to " + debt.CreditorName
	if payment.Description != "" {
		description = payment.Description
	}

	transaction := models.Transaction{
		UserID:      userID,
		AccountID:   payment.AccountID,
		Type:        "expense",
		Amount:      payment.Amount,
		CategoryID:  "debt_payment",
		Date:        payment.PaymentDate.Time,
		Description: description,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	// Update account balance
	account.Balance -= payment.Amount

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

	// Log debt payment activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleDebt,
		"DebtPayment", payment.ID, "Recorded debt payment: "+debt.CreditorName, nil)

	utilities.SuccessResponse(c, map[string]interface{}{
		"payment": payment,
		"debt":    debt,
	}, "Payment recorded successfully")
}

// ListDebtPayments returns all payments for a specific debt
func ListDebtPayments(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	debtID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid debt ID")
		return
	}

	// Verify debt belongs to user
	var debt models.DebtRecord
	if err := database.DB.Where("id = ? AND user_id = ?", debtID, userID).First(&debt).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Debt not found")
		return
	}

	var payments []models.DebtPayment
	if err := database.DB.Where("debt_id = ? AND user_id = ?", debtID, userID).
		Order("payment_date DESC, created_at DESC").
		Find(&payments).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch payments")
		return
	}

	// Enrich payments with account names
	type PaymentResponse struct {
		models.DebtPayment
		AccountName string `json:"accountName"`
	}

	enrichedPayments := make([]PaymentResponse, len(payments))
	for i, payment := range payments {
		enrichedPayments[i] = PaymentResponse{DebtPayment: payment}

		var account models.Account
		if err := database.DB.Select("name").Where("id = ?", payment.AccountID).First(&account).Error; err == nil {
			enrichedPayments[i].AccountName = account.Name
		}
	}

	utilities.SuccessResponse(c, enrichedPayments, "Payments retrieved successfully")
}

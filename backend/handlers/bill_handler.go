package handlers

import (
	"net/http"
	"time"

	"daybook-backend/database"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListBills returns all bills for the authenticated user
func ListBills(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	query := database.DB.Where("user_id = ?", userID)

	// Optional filter by active status
	if active := c.Query("active"); active != "" {
		query = query.Where("active = ?", active == "true")
	}

	// Optional filter by category
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	var bills []models.Bill
	if err := query.Order("due_day ASC").Find(&bills).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch bills")
		return
	}

	utilities.SuccessResponse(c, bills, "Bills retrieved successfully")
}

// GetBill returns a specific bill by ID
func GetBill(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	billID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid bill ID")
		return
	}

	var bill models.Bill
	if err := database.DB.Where("id = ? AND user_id = ?", billID, userID).First(&bill).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Bill not found")
		return
	}

	utilities.SuccessResponse(c, bill, "Bill retrieved successfully")
}

// CreateBill creates a new bill
func CreateBill(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var bill models.Bill
	if err := c.ShouldBindJSON(&bill); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	bill.UserID = userID

	if err := database.DB.Create(&bill).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create bill")
		return
	}

	utilities.CreatedResponse(c, bill, "Bill created successfully")
}

// UpdateBill updates an existing bill
func UpdateBill(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	billID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid bill ID")
		return
	}

	var existingBill models.Bill
	if err := database.DB.Where("id = ? AND user_id = ?", billID, userID).First(&existingBill).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Bill not found")
		return
	}

	var updateData models.Bill
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update allowed fields
	existingBill.Name = updateData.Name
	existingBill.Category = updateData.Category
	existingBill.Amount = updateData.Amount
	existingBill.Frequency = updateData.Frequency
	existingBill.StartDate = updateData.StartDate
	existingBill.DueDay = updateData.DueDay
	existingBill.AutoPay = updateData.AutoPay
	existingBill.ReminderDays = updateData.ReminderDays
	existingBill.Active = updateData.Active
	existingBill.Notes = updateData.Notes

	if err := database.DB.Save(&existingBill).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update bill")
		return
	}

	utilities.SuccessResponse(c, existingBill, "Bill updated successfully")
}

// DeleteBill deletes a bill
func DeleteBill(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	billID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid bill ID")
		return
	}

	var bill models.Bill
	if err := database.DB.Where("id = ? AND user_id = ?", billID, userID).First(&bill).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Bill not found")
		return
	}

	// Soft delete
	if err := database.DB.Delete(&bill).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete bill")
		return
	}

	utilities.SuccessResponse(c, nil, "Bill deleted successfully")
}

// PayBill marks a bill as paid and records the payment
func PayBill(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	billID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid bill ID")
		return
	}

	var paymentData struct {
		Amount      float64    `json:"amount" binding:"required,gt=0"`
		PaymentDate *time.Time `json:"paymentDate"`
		AccountID   *uuid.UUID `json:"accountId"`
		Notes       string     `json:"notes"`
	}

	if err := c.ShouldBindJSON(&paymentData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var bill models.Bill
	if err := database.DB.Where("id = ? AND user_id = ?", billID, userID).First(&bill).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Bill not found")
		return
	}

	// If account is specified, verify it belongs to user
	if paymentData.AccountID != nil {
		var account models.Account
		if err := database.DB.Where("id = ? AND user_id = ?", *paymentData.AccountID, userID).First(&account).Error; err != nil {
			utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
			return
		}
	}

	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create bill payment record
	billPayment := models.BillPayment{
		UserID: userID,
		BillID: billID,
		Amount: paymentData.Amount,
		PaymentDate: func() time.Time {
			if paymentData.PaymentDate != nil {
				return *paymentData.PaymentDate
			}
			return time.Now()
		}(),
		AccountID: paymentData.AccountID,
		Notes:     paymentData.Notes,
	}

	if err := tx.Create(&billPayment).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to record payment")
		return
	}

	// Update bill's last payment info
	bill.LastPaidDate = &billPayment.PaymentDate
	bill.LastPaidAmount = paymentData.Amount

	if err := tx.Save(&bill).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update bill")
		return
	}

	tx.Commit()

	result := map[string]interface{}{
		"bill":    bill,
		"payment": billPayment,
	}

	utilities.SuccessResponse(c, result, "Bill payment recorded successfully")
}

// GetBillPayments returns payment history for bills
func GetBillPayments(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	query := database.DB.Where("user_id = ?", userID)

	// Optional filter by bill
	if billID := c.Query("billId"); billID != "" {
		query = query.Where("bill_id = ?", billID)
	}

	// Optional filter by date range
	if startDate := c.Query("startDate"); startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("payment_date >= ?", parsedDate)
		}
	}

	if endDate := c.Query("endDate"); endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("payment_date <= ?", parsedDate)
		}
	}

	var payments []models.BillPayment
	if err := query.Order("payment_date DESC").Find(&payments).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch bill payments")
		return
	}

	utilities.SuccessResponse(c, payments, "Bill payments retrieved successfully")
}

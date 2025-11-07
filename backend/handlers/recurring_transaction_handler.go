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

// ListRecurringTransactions returns all recurring transactions for the authenticated user
func ListRecurringTransactions(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var recurringTransactions []models.RecurringTransaction
	if err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&recurringTransactions).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch recurring transactions")
		return
	}

	utilities.SuccessResponse(c, recurringTransactions, "Recurring transactions retrieved successfully")
}

// GetRecurringTransaction returns a specific recurring transaction by ID
func GetRecurringTransaction(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	recurringID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid recurring transaction ID")
		return
	}

	var recurringTransaction models.RecurringTransaction
	if err := database.DB.Where("id = ? AND user_id = ?", recurringID, userID).First(&recurringTransaction).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Recurring transaction not found")
		return
	}

	utilities.SuccessResponse(c, recurringTransaction, "Recurring transaction retrieved successfully")
}

// CreateRecurringTransaction creates a new recurring transaction
func CreateRecurringTransaction(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var recurringTransaction models.RecurringTransaction
	if err := c.ShouldBindJSON(&recurringTransaction); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	recurringTransaction.UserID = userID

	// Validate required UUID fields
	if recurringTransaction.TransactionTemplate.AccountID == uuid.Nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Account ID is required")
		return
	}

	// Validate transfer-specific requirements
	if recurringTransaction.TransactionTemplate.Type == "transfer" {
		if recurringTransaction.TransactionTemplate.ToAccountID == nil || *recurringTransaction.TransactionTemplate.ToAccountID == uuid.Nil {
			utilities.ErrorResponse(c, http.StatusBadRequest, "To Account ID is required for transfers")
			return
		}
	}

	// Verify account belongs to user
	var account models.Account
	if err := database.DB.Where("id = ? AND user_id = ?", recurringTransaction.TransactionTemplate.AccountID, userID).First(&account).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	// Validate frequency
	validFrequencies := []string{"daily", "weekly", "biweekly", "monthly", "quarterly", "yearly"}
	isValidFrequency := false
	for _, freq := range validFrequencies {
		if recurringTransaction.Frequency == freq {
			isValidFrequency = true
			break
		}
	}
	if !isValidFrequency {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid frequency. Must be one of: daily, weekly, biweekly, monthly, quarterly, yearly")
		return
	}

	// Create the recurring transaction
	if err := database.DB.Create(&recurringTransaction).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create recurring transaction")
		return
	}

	utilities.CreatedResponse(c, recurringTransaction, "Recurring transaction created successfully")
}

// UpdateRecurringTransaction updates an existing recurring transaction
func UpdateRecurringTransaction(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	recurringID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid recurring transaction ID")
		return
	}

	var existingRecurring models.RecurringTransaction
	if err := database.DB.Where("id = ? AND user_id = ?", recurringID, userID).First(&existingRecurring).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Recurring transaction not found")
		return
	}

	var updateData models.RecurringTransaction
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate required UUID fields
	if updateData.TransactionTemplate.AccountID == uuid.Nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Account ID is required")
		return
	}

	// Validate transfer-specific requirements
	if updateData.TransactionTemplate.Type == "transfer" {
		if updateData.TransactionTemplate.ToAccountID == nil || *updateData.TransactionTemplate.ToAccountID == uuid.Nil {
			utilities.ErrorResponse(c, http.StatusBadRequest, "To Account ID is required for transfers")
			return
		}
	}

	// Verify account belongs to user if changed
	if updateData.TransactionTemplate.AccountID != existingRecurring.TransactionTemplate.AccountID {
		var account models.Account
		if err := database.DB.Where("id = ? AND user_id = ?", updateData.TransactionTemplate.AccountID, userID).First(&account).Error; err != nil {
			utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
			return
		}
	}

	// Update fields
	existingRecurring.TransactionTemplate = updateData.TransactionTemplate
	existingRecurring.Frequency = updateData.Frequency
	existingRecurring.StartDate = updateData.StartDate
	existingRecurring.EndDate = updateData.EndDate
	existingRecurring.Enabled = updateData.Enabled

	if err := database.DB.Save(&existingRecurring).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update recurring transaction")
		return
	}

	utilities.SuccessResponse(c, existingRecurring, "Recurring transaction updated successfully")
}

// DeleteRecurringTransaction deletes a recurring transaction
func DeleteRecurringTransaction(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	recurringID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid recurring transaction ID")
		return
	}

	var recurringTransaction models.RecurringTransaction
	if err := database.DB.Where("id = ? AND user_id = ?", recurringID, userID).First(&recurringTransaction).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Recurring transaction not found")
		return
	}

	// Soft delete
	if err := database.DB.Delete(&recurringTransaction).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete recurring transaction")
		return
	}

	utilities.SuccessResponse(c, nil, "Recurring transaction deleted successfully")
}

// ProcessRecurringTransactions generates missing transactions for all enabled recurring transactions
func ProcessRecurringTransactions(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get all enabled recurring transactions for the user
	var recurringTransactions []models.RecurringTransaction
	if err := database.DB.Where("user_id = ? AND enabled = ?", userID, true).Find(&recurringTransactions).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch recurring transactions")
		return
	}

	now := time.Now()
	createdCount := 0
	skippedCount := 0
	errorCount := 0

	// Start a database transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, recurring := range recurringTransactions {
		// Skip if start date is in the future
		if recurring.StartDate.After(now) {
			skippedCount++
			continue
		}

		// Skip if end date has passed
		if recurring.EndDate != nil && recurring.EndDate.Before(now) {
			skippedCount++
			continue
		}

		// Determine the start point for generating transactions
		startFrom := recurring.StartDate
		if recurring.LastProcessed != nil {
			startFrom = *recurring.LastProcessed
		}

		// Generate transactions from startFrom to now
		transactionDates := calculateTransactionDates(startFrom, now, recurring.Frequency, recurring.StartDate)

		for _, txnDate := range transactionDates {
			// Skip if transaction date is before start date
			if txnDate.Before(recurring.StartDate) {
				continue
			}

			// Skip if transaction date is after end date
			if recurring.EndDate != nil && txnDate.After(*recurring.EndDate) {
				continue
			}

			// Check if transaction already exists for this date and recurring ID (duplicate prevention)
			var existingCount int64
			tx.Model(&models.Transaction{}).Where(
				"user_id = ? AND recurring_id = ? AND DATE(date) = DATE(?)",
				userID, recurring.ID, txnDate,
			).Count(&existingCount)

			if existingCount > 0 {
				skippedCount++
				continue
			}

			// Create the transaction from template
			transaction := models.Transaction{
				UserID:        userID,
				AccountID:     recurring.TransactionTemplate.AccountID,
				ToAccountID:   recurring.TransactionTemplate.ToAccountID,
				Type:          recurring.TransactionTemplate.Type,
				Amount:        recurring.TransactionTemplate.Amount,
				CategoryID:    recurring.TransactionTemplate.CategoryID,
				Date:          txnDate,
				Description:   recurring.TransactionTemplate.Description,
				Tags:          recurring.TransactionTemplate.Tags,
				CreditCardID:  recurring.TransactionTemplate.CreditCardID,
				Attachments:   recurring.TransactionTemplate.Attachments,
				RecurringID:   &recurring.ID,
			}

			// Create the transaction
			if err := tx.Create(&transaction).Error; err != nil {
				errorCount++
				continue
			}

			// Update account balance
			isCreditCardTransaction := transaction.CreditCardID != nil

			if isCreditCardTransaction {
				// Update credit card balance
				var creditCard models.CreditCard
				if err := tx.Where("id = ?", transaction.CreditCardID).First(&creditCard).Error; err != nil {
					errorCount++
					continue
				}

				if transaction.Type == "income" {
					creditCard.CurrentBalance += transaction.Amount
				} else if transaction.Type == "expense" {
					creditCard.CurrentBalance += transaction.Amount
				}

				if err := tx.Save(&creditCard).Error; err != nil {
					errorCount++
					continue
				}
			} else {
				// Update account balance
				var account models.Account
				if err := tx.Where("id = ?", transaction.AccountID).First(&account).Error; err != nil {
					errorCount++
					continue
				}

				if transaction.Type == "income" {
					account.Balance += transaction.Amount
				} else if transaction.Type == "expense" {
					account.Balance -= transaction.Amount
				} else if transaction.Type == "transfer" && transaction.ToAccountID != nil {
					account.Balance -= transaction.Amount
					tx.Model(&models.Account{}).Where("id = ?", *transaction.ToAccountID).
						UpdateColumn("balance", database.DB.Raw("balance + ?", transaction.Amount))
				}

				if err := tx.Save(&account).Error; err != nil {
					errorCount++
					continue
				}
			}

			createdCount++
		}

		// Update LastProcessed date
		recurring.LastProcessed = &now
		if err := tx.Save(&recurring).Error; err != nil {
			errorCount++
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to process recurring transactions")
		return
	}

	result := map[string]interface{}{
		"created": createdCount,
		"skipped": skippedCount,
		"errors":  errorCount,
		"message": "Recurring transactions processed successfully",
	}

	utilities.SuccessResponse(c, result, "Processing completed")
}

// calculateTransactionDates calculates all transaction dates between start and end based on frequency
func calculateTransactionDates(start, end time.Time, frequency string, originalStartDate time.Time) []time.Time {
	var dates []time.Time
	current := start

	// For the first run, include the start date if it's the original start date
	if start.Equal(originalStartDate) {
		dates = append(dates, start)
	}

	for {
		// Calculate next date based on frequency
		var next time.Time
		switch frequency {
		case "daily":
			next = current.AddDate(0, 0, 1)
		case "weekly":
			next = current.AddDate(0, 0, 7)
		case "biweekly":
			next = current.AddDate(0, 0, 14)
		case "monthly":
			next = current.AddDate(0, 1, 0)
		case "quarterly":
			next = current.AddDate(0, 3, 0)
		case "yearly":
			next = current.AddDate(1, 0, 0)
		default:
			return dates
		}

		if next.After(end) {
			break
		}

		dates = append(dates, next)
		current = next
	}

	return dates
}

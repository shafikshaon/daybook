package handlers

import (
	"net/http"
	"strconv"
	"time"

	"daybook-backend/database"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListTransactions returns all transactions for the authenticated user with optional filters
func ListTransactions(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	query := database.DB.Where("user_id = ?", userID)

	// Exclude tracking transactions (used for external holdings) unless explicitly requested
	if c.Query("includeTracking") != "true" {
		query = query.Where("type != ?", "tracking")
	}

	// Apply filters
	if transactionType := c.Query("type"); transactionType != "" {
		query = query.Where("type = ?", transactionType)
	}

	if categoryID := c.Query("categoryId"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if accountID := c.Query("accountId"); accountID != "" {
		query = query.Where("account_id = ?", accountID)
	}

	if startDate := c.Query("startDate"); startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			// Set to beginning of day
			startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
			query = query.Where("date >= ?", startOfDay)
		}
	}

	if endDate := c.Query("endDate"); endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			// Set to end of day
			endOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 999999999, parsedDate.Location())
			query = query.Where("date <= ?", endOfDay)
		}
	}

	// Pagination parameters
	page := 1
	limit := 20 // default limit

	if pageParam := c.Query("page"); pageParam != "" {
		if parsedPage, err := strconv.Atoi(pageParam); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			// Validate limit is one of the allowed values: 20, 50, 100, 500
			switch parsedLimit {
			case 20, 50, 100, 500:
				limit = parsedLimit
			default:
				limit = 20 // fallback to default if invalid value
			}
		}
	}

	// Get total count before pagination
	var totalCount int64
	if err := query.Model(&models.Transaction{}).Count(&totalCount).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to count transactions")
		return
	}

	// Calculate offset
	offset := (page - 1) * limit

	var transactions []models.Transaction
	if err := query.Order("date DESC, created_at DESC").Limit(limit).Offset(offset).Find(&transactions).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

	// Enrich transactions with account/credit card names
	type TransactionResponse struct {
		models.Transaction
		AccountName    *string `json:"accountName,omitempty"`
		CreditCardName *string `json:"creditCardName,omitempty"`
		ToAccountName  *string `json:"toAccountName,omitempty"`
	}

	enrichedTransactions := make([]TransactionResponse, len(transactions))
	for i, txn := range transactions {
		enrichedTransactions[i] = TransactionResponse{Transaction: txn}

		// If transaction has a credit card, the accountId is actually the credit card ID
		// So we need to fetch the credit card name as the account name
		if txn.CreditCardID != nil {
			var card models.CreditCard
			if err := database.DB.Select("name").Where("id = ?", *txn.CreditCardID).First(&card).Error; err == nil {
				enrichedTransactions[i].CreditCardName = &card.Name
				enrichedTransactions[i].AccountName = &card.Name // Use credit card name as account name
			}
		} else {
			// For regular transactions, fetch the account name
			var account models.Account
			if err := database.DB.Select("name").Where("id = ?", txn.AccountID).First(&account).Error; err == nil {
				enrichedTransactions[i].AccountName = &account.Name
			}
		}

		// For transfers, fetch the destination account name
		if txn.ToAccountID != nil {
			var toAccount models.Account
			if err := database.DB.Select("name").Where("id = ?", *txn.ToAccountID).First(&toAccount).Error; err == nil {
				enrichedTransactions[i].ToAccountName = &toAccount.Name
			}
		}
	}

	// Calculate pagination metadata
	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	response := map[string]interface{}{
		"transactions": enrichedTransactions,
		"pagination": map[string]interface{}{
			"currentPage": page,
			"limit":       limit,
			"totalCount":  totalCount,
			"totalPages":  totalPages,
			"hasNext":     page < totalPages,
			"hasPrev":     page > 1,
		},
	}

	utilities.SuccessResponse(c, response, "Transactions retrieved successfully")
}

// GetTransaction returns a specific transaction by ID
func GetTransaction(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	var transaction models.Transaction
	if err := database.DB.Where("id = ? AND user_id = ?", transactionID, userID).First(&transaction).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Transaction not found")
		return
	}

	// Enrich transaction with account/credit card names
	type TransactionResponse struct {
		models.Transaction
		AccountName    *string `json:"accountName,omitempty"`
		CreditCardName *string `json:"creditCardName,omitempty"`
		ToAccountName  *string `json:"toAccountName,omitempty"`
	}

	response := TransactionResponse{Transaction: transaction}

	// If transaction has a credit card, the accountId is actually the credit card ID
	// So we need to fetch the credit card name as the account name
	if transaction.CreditCardID != nil {
		var card models.CreditCard
		if err := database.DB.Select("name").Where("id = ?", *transaction.CreditCardID).First(&card).Error; err == nil {
			response.CreditCardName = &card.Name
			response.AccountName = &card.Name // Use credit card name as account name
		}
	} else {
		// For regular transactions, fetch the account name
		var account models.Account
		if err := database.DB.Select("name").Where("id = ?", transaction.AccountID).First(&account).Error; err == nil {
			response.AccountName = &account.Name
		}
	}

	// For transfers, fetch the destination account name
	if transaction.ToAccountID != nil {
		var toAccount models.Account
		if err := database.DB.Select("name").Where("id = ?", *transaction.ToAccountID).First(&toAccount).Error; err == nil {
			response.ToAccountName = &toAccount.Name
		}
	}

	utilities.SuccessResponse(c, response, "Transaction retrieved successfully")
}

// CreateTransaction creates a new transaction
func CreateTransaction(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	transaction.UserID = userID

	// Determine if this is a credit card transaction or account transaction
	isCreditCardTransaction := transaction.CreditCardID != nil

	// Verify account or credit card belongs to user
	if isCreditCardTransaction {
		var creditCard models.CreditCard
		if err := database.DB.Where("id = ? AND user_id = ?", transaction.CreditCardID, userID).First(&creditCard).Error; err != nil {
			utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
			return
		}
	} else {
		var account models.Account
		if err := database.DB.Where("id = ? AND user_id = ?", transaction.AccountID, userID).First(&account).Error; err != nil {
			utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
			return
		}

		// For transfers, verify the destination account
		if transaction.Type == "transfer" && transaction.ToAccountID != nil {
			var toAccount models.Account
			if err := database.DB.Where("id = ? AND user_id = ?", *transaction.ToAccountID, userID).First(&toAccount).Error; err != nil {
				utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid destination account ID")
				return
			}
		}
	}

	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the transaction
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	// Update balance
	if isCreditCardTransaction {
		// Update credit card balance
		var creditCard models.CreditCard
		if err := tx.Where("id = ?", transaction.CreditCardID).First(&creditCard).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch credit card")
			return
		}

		if transaction.Type == "income" {
			creditCard.CurrentBalance += transaction.Amount
		} else if transaction.Type == "expense" {
			creditCard.CurrentBalance += transaction.Amount
		}

		if err := tx.Save(&creditCard).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update credit card balance")
			return
		}
	} else {
		// Update account balance
		var account models.Account
		if err := tx.Where("id = ?", transaction.AccountID).First(&account).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch account")
			return
		}

		if transaction.Type == "income" {
			account.Balance += transaction.Amount
		} else if transaction.Type == "expense" {
			account.Balance -= transaction.Amount
		} else if transaction.Type == "transfer" && transaction.ToAccountID != nil {
			// Deduct from source account
			account.Balance -= transaction.Amount
			// Add to destination account
			tx.Model(&models.Account{}).Where("id = ?", *transaction.ToAccountID).
				UpdateColumn("balance", database.DB.Raw("balance + ?", transaction.Amount))
		}

		if err := tx.Save(&account).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account balance")
			return
		}
	}

	tx.Commit()

	utilities.CreatedResponse(c, transaction, "Transaction created successfully")
}

// UpdateTransaction updates an existing transaction
func UpdateTransaction(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	var existingTransaction models.Transaction
	if err := database.DB.Where("id = ? AND user_id = ?", transactionID, userID).First(&existingTransaction).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Transaction not found")
		return
	}

	var updateData models.Transaction
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Determine if this is a credit card transaction or account transaction
	isCreditCardTransaction := updateData.CreditCardID != nil
	wasOldCreditCardTransaction := existingTransaction.CreditCardID != nil

	// Verify account or credit card belongs to user
	if isCreditCardTransaction {
		var creditCard models.CreditCard
		if err := database.DB.Where("id = ? AND user_id = ?", updateData.CreditCardID, userID).First(&creditCard).Error; err != nil {
			utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
			return
		}
	} else {
		var account models.Account
		if err := database.DB.Where("id = ? AND user_id = ?", updateData.AccountID, userID).First(&account).Error; err != nil {
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

	// Revert old balance changes
	if wasOldCreditCardTransaction {
		// Revert credit card balance
		var oldCreditCard models.CreditCard
		if err := tx.Where("id = ?", existingTransaction.CreditCardID).First(&oldCreditCard).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch old credit card")
			return
		}

		if existingTransaction.Type == "income" {
			oldCreditCard.CurrentBalance -= existingTransaction.Amount
		} else if existingTransaction.Type == "expense" {
			oldCreditCard.CurrentBalance -= existingTransaction.Amount
		}

		if err := tx.Save(&oldCreditCard).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to revert old credit card balance")
			return
		}
	} else {
		// Revert account balance
		var oldAccount models.Account
		if err := tx.Where("id = ?", existingTransaction.AccountID).First(&oldAccount).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch old account")
			return
		}

		if existingTransaction.Type == "income" {
			oldAccount.Balance -= existingTransaction.Amount
		} else if existingTransaction.Type == "expense" {
			oldAccount.Balance += existingTransaction.Amount
		} else if existingTransaction.Type == "transfer" && existingTransaction.ToAccountID != nil {
			oldAccount.Balance += existingTransaction.Amount
			tx.Model(&models.Account{}).Where("id = ?", *existingTransaction.ToAccountID).
				UpdateColumn("balance", database.DB.Raw("balance - ?", existingTransaction.Amount))
		}

		if err := tx.Save(&oldAccount).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to revert old balance")
			return
		}
	}

	// Update transaction
	existingTransaction.AccountID = updateData.AccountID
	existingTransaction.CreditCardID = updateData.CreditCardID
	existingTransaction.ToAccountID = updateData.ToAccountID
	existingTransaction.Type = updateData.Type
	existingTransaction.Amount = updateData.Amount
	existingTransaction.CategoryID = updateData.CategoryID
	existingTransaction.Date = updateData.Date
	existingTransaction.Description = updateData.Description
	existingTransaction.Tags = updateData.Tags
	existingTransaction.Attachments = updateData.Attachments

	if err := tx.Save(&existingTransaction).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update transaction")
		return
	}

	// Apply new balance changes
	if isCreditCardTransaction {
		// Update credit card balance
		var newCreditCard models.CreditCard
		if err := tx.Where("id = ?", updateData.CreditCardID).First(&newCreditCard).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch new credit card")
			return
		}

		if updateData.Type == "income" {
			newCreditCard.CurrentBalance += updateData.Amount
		} else if updateData.Type == "expense" {
			newCreditCard.CurrentBalance += updateData.Amount
		}

		if err := tx.Save(&newCreditCard).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update new credit card balance")
			return
		}
	} else {
		// Update account balance
		var newAccount models.Account
		if err := tx.Where("id = ?", updateData.AccountID).First(&newAccount).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch new account")
			return
		}

		if updateData.Type == "income" {
			newAccount.Balance += updateData.Amount
		} else if updateData.Type == "expense" {
			newAccount.Balance -= updateData.Amount
		} else if updateData.Type == "transfer" && updateData.ToAccountID != nil {
			newAccount.Balance -= updateData.Amount
			tx.Model(&models.Account{}).Where("id = ?", *updateData.ToAccountID).
				UpdateColumn("balance", database.DB.Raw("balance + ?", updateData.Amount))
		}

		if err := tx.Save(&newAccount).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update new balance")
			return
		}
	}

	tx.Commit()

	utilities.SuccessResponse(c, existingTransaction, "Transaction updated successfully")
}

// DeleteTransaction deletes a transaction
func DeleteTransaction(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	var transaction models.Transaction
	if err := database.DB.Where("id = ? AND user_id = ?", transactionID, userID).First(&transaction).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Transaction not found")
		return
	}

	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Determine if this is a credit card transaction or account transaction
	isCreditCardTransaction := transaction.CreditCardID != nil

	// Revert balance changes
	if isCreditCardTransaction {
		// Revert credit card balance
		var creditCard models.CreditCard
		if err := tx.Where("id = ?", transaction.CreditCardID).First(&creditCard).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch credit card")
			return
		}

		if transaction.Type == "income" {
			creditCard.CurrentBalance -= transaction.Amount
		} else if transaction.Type == "expense" {
			creditCard.CurrentBalance -= transaction.Amount
		}

		if err := tx.Save(&creditCard).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update credit card balance")
			return
		}
	} else {
		// Revert account balance
		var account models.Account
		if err := tx.Where("id = ?", transaction.AccountID).First(&account).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch account")
			return
		}

		if transaction.Type == "income" {
			account.Balance -= transaction.Amount
		} else if transaction.Type == "expense" {
			account.Balance += transaction.Amount
		} else if transaction.Type == "transfer" && transaction.ToAccountID != nil {
			account.Balance += transaction.Amount
			tx.Model(&models.Account{}).Where("id = ?", *transaction.ToAccountID).
				UpdateColumn("balance", database.DB.Raw("balance - ?", transaction.Amount))
		}

		if err := tx.Save(&account).Error; err != nil {
			tx.Rollback()
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account balance")
			return
		}
	}

	// Delete transaction (soft delete)
	if err := tx.Delete(&transaction).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete transaction")
		return
	}

	tx.Commit()

	utilities.SuccessResponse(c, nil, "Transaction deleted successfully")
}

// BulkImportTransactions imports multiple transactions at once
func BulkImportTransactions(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var transactions []models.Transaction
	if err := c.ShouldBindJSON(&transactions); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	successCount := 0
	failedCount := 0

	for i := range transactions {
		transactions[i].UserID = userID

		// Verify account belongs to user
		var account models.Account
		if err := tx.Where("id = ? AND user_id = ?", transactions[i].AccountID, userID).First(&account).Error; err != nil {
			failedCount++
			continue
		}

		// Create transaction
		if err := tx.Create(&transactions[i]).Error; err != nil {
			failedCount++
			continue
		}

		// Update account balance
		if transactions[i].Type == "income" {
			account.Balance += transactions[i].Amount
		} else if transactions[i].Type == "expense" {
			account.Balance -= transactions[i].Amount
		} else if transactions[i].Type == "transfer" && transactions[i].ToAccountID != nil {
			account.Balance -= transactions[i].Amount
			tx.Model(&models.Account{}).Where("id = ?", *transactions[i].ToAccountID).
				UpdateColumn("balance", database.DB.Raw("balance + ?", transactions[i].Amount))
		}

		if err := tx.Save(&account).Error; err != nil {
			failedCount++
			continue
		}

		successCount++
	}

	tx.Commit()

	result := map[string]interface{}{
		"successCount": successCount,
		"failedCount":  failedCount,
		"totalCount":   len(transactions),
	}

	utilities.SuccessResponse(c, result, "Bulk import completed")
}

// GetTransactionStats returns transaction statistics
func GetTransactionStats(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get date range from query params
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	query := database.DB.Model(&models.Transaction{}).Where("user_id = ?", userID)

	// Exclude tracking transactions from statistics
	query = query.Where("type != ?", "tracking")

	if startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("date >= ?", parsedDate)
		}
	}

	if endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("date <= ?", parsedDate)
		}
	}

	// Calculate totals by type
	var stats struct {
		TotalIncome      float64
		TotalExpense     float64
		TotalTransfer    float64
		NetIncome        float64
		TransactionCount int64
	}

	// Total income
	query.Where("type = ?", "income").Select("COALESCE(SUM(amount), 0)").Row().Scan(&stats.TotalIncome)

	// Total expense
	query.Where("type = ?", "expense").Select("COALESCE(SUM(amount), 0)").Row().Scan(&stats.TotalExpense)

	// Total transfers
	query.Where("type = ?", "transfer").Select("COALESCE(SUM(amount), 0)").Row().Scan(&stats.TotalTransfer)

	// Net income
	stats.NetIncome = stats.TotalIncome - stats.TotalExpense

	// Transaction count
	query.Count(&stats.TransactionCount)

	utilities.SuccessResponse(c, stats, "Statistics retrieved successfully")
}

package handlers

import (
	"net/http"
	"time"

	"daybook-backend/database"
	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListCreditCards returns all credit cards for the authenticated user
func ListCreditCards(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListCreditCards - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	var cards []models.CreditCard
	logger.Debugf(ctx, "Fetching credit cards from database...")
	if err := database.DB.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&cards).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching credit cards: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch credit cards")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d credit cards for user: %s", len(cards), userID)
	utilities.SuccessResponse(c, cards, "Credit cards retrieved successfully")
}

// GetCreditCard returns a specific credit card by ID
func GetCreditCard(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetCreditCard - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid credit card ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	logger.Debugf(ctx, "Fetching credit card: %s", cardID)
	var card models.CreditCard
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching credit card: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	logger.Infof(ctx, "Successfully retrieved credit card for user: %s", userID)
	utilities.SuccessResponse(c, card, "Credit card retrieved successfully")
}

// CreateCreditCard creates a new credit card
func CreateCreditCard(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateCreditCard - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	var card models.CreditCard
	if err := c.ShouldBindJSON(&card); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	card.UserID = userID

	logger.Debugf(ctx, "Creating credit card in database...")
	if err := database.DB.WithContext(ctx).Create(&card).Error; err != nil {
		logger.Errorf(ctx, "Database error creating credit card: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create credit card")
		return
	}

	// Log credit card creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleCreditCard,
		"CreditCard", card.ID, "Created credit card: "+card.Name, nil)

	logger.Infof(ctx, "Successfully created credit card for user: %s", userID)
	utilities.CreatedResponse(c, card, "Credit card created successfully")
}

// UpdateCreditCard updates an existing credit card
func UpdateCreditCard(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateCreditCard - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid credit card ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	logger.Debugf(ctx, "Fetching existing credit card: %s", cardID)
	var existingCard models.CreditCard
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", cardID, userID).First(&existingCard).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching credit card: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	var updateData models.CreditCard
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update allowed fields
	existingCard.Name = updateData.Name
	existingCard.LastFourDigits = updateData.LastFourDigits
	existingCard.CardNetwork = updateData.CardNetwork
	existingCard.CreditLimit = updateData.CreditLimit
	existingCard.CurrentBalance = updateData.CurrentBalance
	existingCard.APR = updateData.APR
	existingCard.DueDate = updateData.DueDate
	existingCard.StatementDate = updateData.StatementDate
	existingCard.MinimumPayment = updateData.MinimumPayment
	existingCard.RewardsProgram = updateData.RewardsProgram
	existingCard.Active = updateData.Active
	existingCard.Notes = updateData.Notes

	logger.Debugf(ctx, "Updating credit card in database...")
	if err := database.DB.WithContext(ctx).Save(&existingCard).Error; err != nil {
		logger.Errorf(ctx, "Database error updating credit card: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update credit card")
		return
	}

	// Log credit card update activity
	utilities.LogEntityActivity(c, userID, models.ActionUpdate, models.ModuleCreditCard,
		"CreditCard", existingCard.ID, "Updated credit card: "+existingCard.Name, nil)

	logger.Infof(ctx, "Successfully updated credit card for user: %s", userID)
	utilities.SuccessResponse(c, existingCard, "Credit card updated successfully")
}

// DeleteCreditCard deletes a credit card
func DeleteCreditCard(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteCreditCard - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid credit card ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	logger.Debugf(ctx, "Fetching credit card to delete: %s", cardID)
	var card models.CreditCard
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching credit card: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	// Soft delete
	logger.Debugf(ctx, "Deleting credit card from database...")
	if err := database.DB.WithContext(ctx).Delete(&card).Error; err != nil {
		logger.Errorf(ctx, "Database error deleting credit card: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete credit card")
		return
	}

	// Log credit card deletion activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, models.ModuleCreditCard,
		"CreditCard", card.ID, "Deleted credit card: "+card.Name, nil)

	logger.Infof(ctx, "Successfully deleted credit card for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Credit card deleted successfully")
}

// RecordCreditCardTransaction records a new credit card transaction (purchase)
func RecordCreditCardTransaction(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "RecordCreditCardTransaction - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid credit card ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	var ccTransaction models.CreditCardTransaction
	if err := c.ShouldBindJSON(&ccTransaction); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Verify card belongs to user
	logger.Debugf(ctx, "Verifying credit card ownership...")
	var card models.CreditCard
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching credit card: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	ccTransaction.UserID = userID
	ccTransaction.CardID = cardID

	// Start database transaction
	logger.Debugf(ctx, "Starting database transaction...")
	tx := database.DB.Begin()

	// Create entry in main transactions table so it appears in transaction list
	// For credit card transactions, we use the credit card ID as the account_id
	mainTransaction := models.Transaction{
		UserID:       userID,
		AccountID:    cardID,    // Set account_id to credit card ID so transaction list can display card name
		Type:         "expense", // Credit card purchases are expenses
		Amount:       ccTransaction.Amount,
		Date:         ccTransaction.Date,
		Description:  ccTransaction.Description,
		CategoryID:   ccTransaction.CategoryID,
		CreditCardID: &cardID,
		Tags:         ccTransaction.Tags,
		Attachments:  ccTransaction.Attachments,
	}

	// For refunds, it's income
	if ccTransaction.Type == "refund" {
		mainTransaction.Type = "income"
	}

	if err := tx.Create(&mainTransaction).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Database error creating transaction: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction record")
		return
	}

	// Link the credit card transaction to the main transaction
	ccTransaction.TransactionID = mainTransaction.ID

	// Create credit card transaction record
	if err := tx.Create(&ccTransaction).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Database error creating credit card transaction: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to record credit card transaction")
		return
	}

	// Update card balance based on transaction type
	if ccTransaction.Type == "purchase" || ccTransaction.Type == "fee" || ccTransaction.Type == "interest" {
		card.CurrentBalance += ccTransaction.Amount
	} else if ccTransaction.Type == "refund" {
		card.CurrentBalance -= ccTransaction.Amount
		if card.CurrentBalance < 0 {
			card.CurrentBalance = 0
		}
	}

	if err := tx.Save(&card).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Database error updating card balance: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update card balance")
		return
	}

	tx.Commit()

	// Log credit card transaction activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleCreditCard,
		"CreditCardTransaction", ccTransaction.ID, "Created credit card transaction: "+ccTransaction.Description, nil)

	logger.Infof(ctx, "Successfully recorded credit card transaction for user: %s", userID)
	utilities.CreatedResponse(c, ccTransaction, "Transaction recorded successfully")
}

// GetCreditCardTransactions returns all transactions for a credit card
func GetCreditCardTransactions(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetCreditCardTransactions - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid credit card ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	// Verify card belongs to user
	logger.Debugf(ctx, "Verifying credit card ownership...")
	var card models.CreditCard
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching credit card: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	logger.Debugf(ctx, "Fetching credit card transactions...")
	var transactions []models.CreditCardTransaction
	if err := database.DB.WithContext(ctx).Where("card_id = ? AND user_id = ?", cardID, userID).
		Order("date DESC").Find(&transactions).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching transactions: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d transactions for user: %s", len(transactions), userID)
	utilities.SuccessResponse(c, transactions, "Transactions retrieved successfully")
}

// DeleteCreditCardTransaction deletes a credit card transaction
func DeleteCreditCardTransaction(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteCreditCardTransaction - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid credit card ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	transactionID, err := uuid.Parse(c.Param("transactionId"))
	if err != nil {
		logger.Warnf(ctx, "Invalid transaction ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	// Get the transaction
	logger.Debugf(ctx, "Fetching transaction to delete: %s", transactionID)
	var transaction models.CreditCardTransaction
	if err := database.DB.WithContext(ctx).Where("id = ? AND card_id = ? AND user_id = ?", transactionID, cardID, userID).First(&transaction).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching transaction: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Transaction not found")
		return
	}

	// Get the card
	logger.Debugf(ctx, "Fetching credit card...")
	var card models.CreditCard
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching credit card: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	// Start database transaction
	logger.Debugf(ctx, "Starting database transaction...")
	tx := database.DB.Begin()

	// Delete the linked main transaction first (if it exists)
	if transaction.TransactionID != uuid.Nil {
		if err := tx.Delete(&models.Transaction{}, "id = ?", transaction.TransactionID).Error; err != nil {
			tx.Rollback()
			logger.Errorf(ctx, "Database error deleting linked transaction: %v", err)
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete linked transaction")
			return
		}
	}

	// Reverse the balance change
	if transaction.Type == "purchase" || transaction.Type == "fee" || transaction.Type == "interest" {
		card.CurrentBalance -= transaction.Amount
		if card.CurrentBalance < 0 {
			card.CurrentBalance = 0
		}
	} else if transaction.Type == "refund" {
		card.CurrentBalance += transaction.Amount
	}

	if err := tx.Save(&card).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Database error updating card balance: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update card balance")
		return
	}

	// Delete the credit card transaction
	if err := tx.Delete(&transaction).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Database error deleting transaction: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete credit card transaction")
		return
	}

	tx.Commit()

	// Log credit card transaction deletion activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, models.ModuleCreditCard,
		"CreditCardTransaction", transaction.ID, "Deleted credit card transaction: "+transaction.Description, nil)

	logger.Infof(ctx, "Successfully deleted credit card transaction for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Transaction deleted successfully")
}

// RecordPayment records a payment for a credit card
func RecordPayment(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "RecordPayment - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid credit card ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	var paymentData struct {
		Amount      float64    `json:"amount" binding:"required,gt=0"`
		AccountID   string     `json:"accountId" binding:"required"`
		PaymentDate *time.Time `json:"paymentDate"`
		Description string     `json:"description"`
	}

	if err := c.ShouldBindJSON(&paymentData); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accountID, err := uuid.Parse(paymentData.AccountID)
	if err != nil {
		logger.Warnf(ctx, "Invalid account ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	// Get the credit card
	logger.Debugf(ctx, "Fetching credit card...")
	var card models.CreditCard
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching credit card: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	// Get the payment account
	logger.Debugf(ctx, "Fetching payment account...")
	var account models.Account
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching account: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Payment account not found")
		return
	}

	// Validate payment amount
	if paymentData.Amount > card.CurrentBalance {
		logger.Warnf(ctx, "Payment amount exceeds current balance")
		utilities.ErrorResponse(c, http.StatusBadRequest, "Payment amount exceeds current balance")
		return
	}

	if paymentData.Amount > account.Balance {
		logger.Warnf(ctx, "Insufficient funds in payment account")
		utilities.ErrorResponse(c, http.StatusBadRequest, "Insufficient funds in payment account")
		return
	}

	paymentDate := time.Now()
	if paymentData.PaymentDate != nil {
		paymentDate = *paymentData.PaymentDate
	}

	// Start transaction
	logger.Debugf(ctx, "Starting database transaction...")
	tx := database.DB.Begin()

	// Create transaction record for the expense
	tags := []string{"credit_card_payment"}
	transaction := models.Transaction{
		UserID:       userID,
		AccountID:    accountID,
		CategoryID:   "credit_card_payment",
		Amount:       paymentData.Amount,
		Type:         "expense",
		Date:         models.Date{Time: paymentDate},
		Description:  paymentData.Description,
		CreditCardID: &cardID,
		Tags:         tags,
	}

	if transaction.Description == "" {
		transaction.Description = "Credit card payment: " + card.Name
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Database error creating transaction: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	// Deduct from account balance
	account.Balance -= paymentData.Amount
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Database error updating account balance: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account balance")
		return
	}

	// Update card balance and payment info
	card.CurrentBalance -= paymentData.Amount
	if card.CurrentBalance < 0 {
		card.CurrentBalance = 0
	}
	card.LastPaymentDate = &paymentDate
	card.LastPaymentAmount = paymentData.Amount

	if err := tx.Save(&card).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Database error updating card balance: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update card balance")
		return
	}

	// Create payment record
	payment := models.CreditCardPayment{
		UserID:        userID,
		CardID:        cardID,
		AccountID:     accountID,
		Amount:        paymentData.Amount,
		PaymentDate:   models.Date{Time: paymentDate},
		Description:   paymentData.Description,
		TransactionID: transaction.ID,
	}

	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Database error creating payment record: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create payment record")
		return
	}

	// Create credit card transaction record for the payment
	ccTransaction := models.CreditCardTransaction{
		UserID:      userID,
		CardID:      cardID,
		Amount:      paymentData.Amount,
		Description: paymentData.Description,
		Date:        models.Date{Time: paymentDate},
		Type:        "payment",
	}

	if err := tx.Create(&ccTransaction).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Database error creating card transaction: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create card transaction")
		return
	}

	tx.Commit()

	// Log credit card payment activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleCreditCard,
		"CreditCardPayment", payment.ID, "Recorded credit card payment: "+card.Name, nil)

	response := map[string]interface{}{
		"card":    card,
		"payment": payment,
	}

	logger.Infof(ctx, "Successfully recorded credit card payment for user: %s", userID)
	utilities.SuccessResponse(c, response, "Payment recorded successfully")
}

// GetPayments returns all payments for a credit card
func GetPayments(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetPayments - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid credit card ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	// Verify card belongs to user
	logger.Debugf(ctx, "Verifying credit card ownership...")
	var card models.CreditCard
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching credit card: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	logger.Debugf(ctx, "Fetching payments...")
	var payments []models.CreditCardPayment
	if err := database.DB.WithContext(ctx).Where("card_id = ? AND user_id = ?", cardID, userID).
		Order("payment_date DESC").Find(&payments).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching payments: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch payments")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d payments for user: %s", len(payments), userID)
	utilities.SuccessResponse(c, payments, "Payments retrieved successfully")
}

// GetStatements returns statements for a credit card
func GetStatements(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetStatements - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "Invalid credit card ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	// Verify card belongs to user
	logger.Debugf(ctx, "Verifying credit card ownership...")
	var card models.CreditCard
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching credit card: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	logger.Debugf(ctx, "Fetching statements...")
	var statements []models.Statement
	if err := database.DB.WithContext(ctx).Where("card_id = ? AND user_id = ?", cardID, userID).
		Order("statement_date DESC").Find(&statements).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching statements: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch statements")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d statements for user: %s", len(statements), userID)
	utilities.SuccessResponse(c, statements, "Statements retrieved successfully")
}

// CreateStatement creates a new statement
func CreateStatement(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateStatement - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	var statement models.Statement
	if err := c.ShouldBindJSON(&statement); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	statement.UserID = userID

	// Verify card belongs to user
	logger.Debugf(ctx, "Verifying credit card ownership...")
	var card models.CreditCard
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", statement.CardID, userID).First(&card).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching credit card: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	logger.Debugf(ctx, "Creating statement in database...")
	if err := database.DB.WithContext(ctx).Create(&statement).Error; err != nil {
		logger.Errorf(ctx, "Database error creating statement: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create statement")
		return
	}

	// Log statement creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleCreditCard,
		"Statement", statement.ID, "Created credit card statement", nil)

	logger.Infof(ctx, "Successfully created statement for user: %s", userID)
	utilities.CreatedResponse(c, statement, "Statement created successfully")
}

// ListRewards returns rewards for the authenticated user
func ListRewards(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListRewards - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	query := database.DB.WithContext(ctx).Where("user_id = ?", userID)

	// Optional filter by card
	if cardID := c.Query("cardId"); cardID != "" {
		query = query.Where("card_id = ?", cardID)
	}

	// Optional filter by redeemed status
	if redeemed := c.Query("redeemed"); redeemed != "" {
		query = query.Where("redeemed = ?", redeemed == "true")
	}

	logger.Debugf(ctx, "Fetching rewards from database...")
	var rewards []models.Reward
	if err := query.Order("earned_date DESC").Find(&rewards).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching rewards: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch rewards")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d rewards for user: %s", len(rewards), userID)
	utilities.SuccessResponse(c, rewards, "Rewards retrieved successfully")
}

// RecordReward records a new reward
func RecordReward(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "RecordReward - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	var reward models.Reward
	if err := c.ShouldBindJSON(&reward); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	reward.UserID = userID

	// Verify card belongs to user
	logger.Debugf(ctx, "Verifying credit card ownership...")
	var card models.CreditCard
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", reward.CardID, userID).First(&card).Error; err != nil {
		logger.Errorf(ctx, "Database error fetching credit card: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	logger.Debugf(ctx, "Recording reward in database...")
	if err := database.DB.WithContext(ctx).Create(&reward).Error; err != nil {
		logger.Errorf(ctx, "Database error creating reward: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to record reward")
		return
	}

	// Log reward creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleCreditCard,
		"Reward", reward.ID, "Recorded credit card reward", nil)

	logger.Infof(ctx, "Successfully recorded reward for user: %s", userID)
	utilities.CreatedResponse(c, reward, "Reward recorded successfully")
}

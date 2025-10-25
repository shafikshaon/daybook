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

// ListCreditCards returns all credit cards for the authenticated user
func ListCreditCards(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var cards []models.CreditCard
	if err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&cards).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch credit cards")
		return
	}

	utilities.SuccessResponse(c, cards, "Credit cards retrieved successfully")
}

// GetCreditCard returns a specific credit card by ID
func GetCreditCard(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	var card models.CreditCard
	if err := database.DB.Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	utilities.SuccessResponse(c, card, "Credit card retrieved successfully")
}

// CreateCreditCard creates a new credit card
func CreateCreditCard(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var card models.CreditCard
	if err := c.ShouldBindJSON(&card); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	card.UserID = userID

	if err := database.DB.Create(&card).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create credit card")
		return
	}

	utilities.CreatedResponse(c, card, "Credit card created successfully")
}

// UpdateCreditCard updates an existing credit card
func UpdateCreditCard(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	var existingCard models.CreditCard
	if err := database.DB.Where("id = ? AND user_id = ?", cardID, userID).First(&existingCard).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	var updateData models.CreditCard
	if err := c.ShouldBindJSON(&updateData); err != nil {
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

	if err := database.DB.Save(&existingCard).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update credit card")
		return
	}

	utilities.SuccessResponse(c, existingCard, "Credit card updated successfully")
}

// DeleteCreditCard deletes a credit card
func DeleteCreditCard(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	var card models.CreditCard
	if err := database.DB.Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	// Soft delete
	if err := database.DB.Delete(&card).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete credit card")
		return
	}

	utilities.SuccessResponse(c, nil, "Credit card deleted successfully")
}

// RecordCreditCardTransaction records a new credit card transaction (purchase)
func RecordCreditCardTransaction(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	var ccTransaction models.CreditCardTransaction
	if err := c.ShouldBindJSON(&ccTransaction); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Verify card belongs to user
	var card models.CreditCard
	if err := database.DB.Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	ccTransaction.UserID = userID
	ccTransaction.CardID = cardID

	// Start database transaction
	tx := database.DB.Begin()

	// Create entry in main transactions table so it appears in transaction list
	// For credit card transactions, we don't have a specific account, so we use the credit card as the source
	mainTransaction := models.Transaction{
		UserID:       userID,
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
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction record")
		return
	}

	// Link the credit card transaction to the main transaction
	ccTransaction.TransactionID = mainTransaction.ID

	// Create credit card transaction record
	if err := tx.Create(&ccTransaction).Error; err != nil {
		tx.Rollback()
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
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update card balance")
		return
	}

	tx.Commit()

	utilities.CreatedResponse(c, ccTransaction, "Transaction recorded successfully")
}

// GetCreditCardTransactions returns all transactions for a credit card
func GetCreditCardTransactions(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	// Verify card belongs to user
	var card models.CreditCard
	if err := database.DB.Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	var transactions []models.CreditCardTransaction
	if err := database.DB.Where("card_id = ? AND user_id = ?", cardID, userID).
		Order("date DESC").Find(&transactions).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

	utilities.SuccessResponse(c, transactions, "Transactions retrieved successfully")
}

// DeleteCreditCardTransaction deletes a credit card transaction
func DeleteCreditCardTransaction(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	transactionID, err := uuid.Parse(c.Param("transactionId"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	// Get the transaction
	var transaction models.CreditCardTransaction
	if err := database.DB.Where("id = ? AND card_id = ? AND user_id = ?", transactionID, cardID, userID).First(&transaction).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Transaction not found")
		return
	}

	// Get the card
	var card models.CreditCard
	if err := database.DB.Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	// Start database transaction
	tx := database.DB.Begin()

	// Delete the linked main transaction first (if it exists)
	if transaction.TransactionID != uuid.Nil {
		if err := tx.Delete(&models.Transaction{}, "id = ?", transaction.TransactionID).Error; err != nil {
			tx.Rollback()
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
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update card balance")
		return
	}

	// Delete the credit card transaction
	if err := tx.Delete(&transaction).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete credit card transaction")
		return
	}

	tx.Commit()

	utilities.SuccessResponse(c, nil, "Transaction deleted successfully")
}

// RecordPayment records a payment for a credit card
func RecordPayment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
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
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accountID, err := uuid.Parse(paymentData.AccountID)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	// Get the credit card
	var card models.CreditCard
	if err := database.DB.Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	// Get the payment account
	var account models.Account
	if err := database.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Payment account not found")
		return
	}

	// Validate payment amount
	if paymentData.Amount > card.CurrentBalance {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Payment amount exceeds current balance")
		return
	}

	if paymentData.Amount > account.Balance {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Insufficient funds in payment account")
		return
	}

	paymentDate := time.Now()
	if paymentData.PaymentDate != nil {
		paymentDate = *paymentData.PaymentDate
	}

	// Start transaction
	tx := database.DB.Begin()

	// Create transaction record for the expense
	tags := []string{"credit_card_payment"}
	transaction := models.Transaction{
		UserID:       userID,
		AccountID:    accountID,
		CategoryID:   "credit_card_payment",
		Amount:       paymentData.Amount,
		Type:         "expense",
		Date:         paymentDate,
		Description:  paymentData.Description,
		CreditCardID: &cardID,
		Tags:         tags,
	}

	if transaction.Description == "" {
		transaction.Description = "Credit card payment: " + card.Name
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	// Deduct from account balance
	account.Balance -= paymentData.Amount
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
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
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update card balance")
		return
	}

	// Create payment record
	payment := models.CreditCardPayment{
		UserID:        userID,
		CardID:        cardID,
		AccountID:     accountID,
		Amount:        paymentData.Amount,
		PaymentDate:   paymentDate,
		Description:   paymentData.Description,
		TransactionID: transaction.ID,
	}

	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create payment record")
		return
	}

	// Create credit card transaction record for the payment
	ccTransaction := models.CreditCardTransaction{
		UserID:      userID,
		CardID:      cardID,
		Amount:      paymentData.Amount,
		Description: paymentData.Description,
		Date:        paymentDate,
		Type:        "payment",
	}

	if err := tx.Create(&ccTransaction).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create card transaction")
		return
	}

	tx.Commit()

	response := map[string]interface{}{
		"card":    card,
		"payment": payment,
	}

	utilities.SuccessResponse(c, response, "Payment recorded successfully")
}

// GetPayments returns all payments for a credit card
func GetPayments(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	// Verify card belongs to user
	var card models.CreditCard
	if err := database.DB.Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	var payments []models.CreditCardPayment
	if err := database.DB.Where("card_id = ? AND user_id = ?", cardID, userID).
		Order("payment_date DESC").Find(&payments).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch payments")
		return
	}

	utilities.SuccessResponse(c, payments, "Payments retrieved successfully")
}

// GetStatements returns statements for a credit card
func GetStatements(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	cardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	// Verify card belongs to user
	var card models.CreditCard
	if err := database.DB.Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	var statements []models.Statement
	if err := database.DB.Where("card_id = ? AND user_id = ?", cardID, userID).
		Order("statement_date DESC").Find(&statements).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch statements")
		return
	}

	utilities.SuccessResponse(c, statements, "Statements retrieved successfully")
}

// CreateStatement creates a new statement
func CreateStatement(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var statement models.Statement
	if err := c.ShouldBindJSON(&statement); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	statement.UserID = userID

	// Verify card belongs to user
	var card models.CreditCard
	if err := database.DB.Where("id = ? AND user_id = ?", statement.CardID, userID).First(&card).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	if err := database.DB.Create(&statement).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create statement")
		return
	}

	utilities.CreatedResponse(c, statement, "Statement created successfully")
}

// ListRewards returns rewards for the authenticated user
func ListRewards(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	query := database.DB.Where("user_id = ?", userID)

	// Optional filter by card
	if cardID := c.Query("cardId"); cardID != "" {
		query = query.Where("card_id = ?", cardID)
	}

	// Optional filter by redeemed status
	if redeemed := c.Query("redeemed"); redeemed != "" {
		query = query.Where("redeemed = ?", redeemed == "true")
	}

	var rewards []models.Reward
	if err := query.Order("earned_date DESC").Find(&rewards).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch rewards")
		return
	}

	utilities.SuccessResponse(c, rewards, "Rewards retrieved successfully")
}

// RecordReward records a new reward
func RecordReward(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var reward models.Reward
	if err := c.ShouldBindJSON(&reward); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	reward.UserID = userID

	// Verify card belongs to user
	var card models.CreditCard
	if err := database.DB.Where("id = ? AND user_id = ?", reward.CardID, userID).First(&card).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid credit card ID")
		return
	}

	if err := database.DB.Create(&reward).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to record reward")
		return
	}

	utilities.CreatedResponse(c, reward, "Reward recorded successfully")
}

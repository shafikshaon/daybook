package handlers

import (
	"net/http"

	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/services"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreditCardHandler struct {
	service services.CreditCardService
}

func NewCreditCardHandler(service services.CreditCardService) *CreditCardHandler {
	return &CreditCardHandler{
		service: service,
	}
}

// ListCreditCards returns all credit cards for the authenticated user
func (h *CreditCardHandler) ListCreditCards(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListCreditCards - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	cards, err := h.service.ListCreditCards(ctx, userID)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch credit cards")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d credit cards for user: %s", len(cards), userID)
	utilities.SuccessResponse(c, cards, "Credit cards retrieved successfully")
}

// GetCreditCard returns a specific credit card by ID
func (h *CreditCardHandler) GetCreditCard(c *gin.Context) {
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
	card, err := h.service.GetCreditCard(ctx, cardID, userID)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	logger.Infof(ctx, "Successfully retrieved credit card for user: %s", userID)
	utilities.SuccessResponse(c, card, "Credit card retrieved successfully")
}

// CreateCreditCard creates a new credit card
func (h *CreditCardHandler) CreateCreditCard(c *gin.Context) {
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
	createdCard, err := h.service.CreateCreditCard(ctx, &card)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create credit card")
		return
	}

	logger.Infof(ctx, "Successfully created credit card for user: %s", userID)
	utilities.CreatedResponse(c, createdCard, "Credit card created successfully")
}

// UpdateCreditCard updates an existing credit card
func (h *CreditCardHandler) UpdateCreditCard(c *gin.Context) {
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

	var updateData models.CreditCard
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedCard, err := h.service.UpdateCreditCard(ctx, cardID, userID, &updateData)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update credit card")
		return
	}

	logger.Infof(ctx, "Successfully updated credit card for user: %s", userID)
	utilities.SuccessResponse(c, updatedCard, "Credit card updated successfully")
}

// DeleteCreditCard deletes a credit card
func (h *CreditCardHandler) DeleteCreditCard(c *gin.Context) {
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

	logger.Debugf(ctx, "Deleting credit card: %s", cardID)

	if err := h.service.DeleteCreditCard(ctx, cardID, userID); err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete credit card")
		return
	}

	logger.Infof(ctx, "Successfully deleted credit card for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Credit card deleted successfully")
}

// RecordCreditCardTransaction records a credit card transaction
func (h *CreditCardHandler) RecordCreditCardTransaction(c *gin.Context) {
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

	var transactionRequest services.RecordTransactionRequest
	if err := c.ShouldBindJSON(&transactionRequest); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	transaction, err := h.service.RecordTransaction(ctx, cardID, userID, &transactionRequest)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Infof(ctx, "Successfully recorded credit card transaction for user: %s", userID)
	utilities.CreatedResponse(c, transaction, "Transaction recorded successfully")
}

// GetCreditCardTransactions returns all transactions for a credit card
func (h *CreditCardHandler) GetCreditCardTransactions(c *gin.Context) {
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

	transactions, err := h.service.GetTransactions(ctx, cardID, userID)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch transactions")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d transactions for user: %s", len(transactions), userID)
	utilities.SuccessResponse(c, transactions, "Transactions retrieved successfully")
}

// DeleteCreditCardTransaction deletes a credit card transaction
func (h *CreditCardHandler) DeleteCreditCardTransaction(c *gin.Context) {
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

	if err := h.service.DeleteTransaction(ctx, cardID, transactionID, userID); err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Infof(ctx, "Successfully deleted transaction for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Transaction deleted successfully")
}

// RecordPayment records a payment to a credit card
func (h *CreditCardHandler) RecordPayment(c *gin.Context) {
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

	var paymentRequest services.RecordPaymentRequest
	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		logger.Warnf(ctx, "Validation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.RecordPayment(ctx, cardID, userID, &paymentRequest)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Infof(ctx, "Successfully recorded credit card payment for user: %s", userID)
	utilities.SuccessResponse(c, result, "Payment recorded successfully")
}

// GetPayments returns all payments for a credit card
func (h *CreditCardHandler) GetPayments(c *gin.Context) {
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

	payments, err := h.service.GetPayments(ctx, cardID, userID)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch payments")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d payments for user: %s", len(payments), userID)
	utilities.SuccessResponse(c, payments, "Payments retrieved successfully")
}

// GetStatements returns all statements for a credit card
func (h *CreditCardHandler) GetStatements(c *gin.Context) {
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

	statements, err := h.service.GetStatements(ctx, cardID, userID)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch statements")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d statements for user: %s", len(statements), userID)
	utilities.SuccessResponse(c, statements, "Statements retrieved successfully")
}

// CreateStatement creates a new credit card statement
func (h *CreditCardHandler) CreateStatement(c *gin.Context) {
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

	createdStatement, err := h.service.CreateStatement(ctx, &statement)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Infof(ctx, "Successfully created statement for user: %s", userID)
	utilities.CreatedResponse(c, createdStatement, "Statement created successfully")
}

// ListRewards returns all rewards for the authenticated user
func (h *CreditCardHandler) ListRewards(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListRewards - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	rewards, err := h.service.ListRewards(ctx, userID)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch rewards")
		return
	}

	logger.Infof(ctx, "Successfully retrieved %d rewards for user: %s", len(rewards), userID)
	utilities.SuccessResponse(c, rewards, "Rewards retrieved successfully")
}

// RecordReward records a new reward
func (h *CreditCardHandler) RecordReward(c *gin.Context) {
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

	createdReward, err := h.service.RecordReward(ctx, &reward)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Infof(ctx, "Successfully recorded reward for user: %s", userID)
	utilities.CreatedResponse(c, createdReward, "Reward recorded successfully")
}

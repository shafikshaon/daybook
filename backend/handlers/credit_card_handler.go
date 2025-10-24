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
		PaymentDate *time.Time `json:"paymentDate"`
	}

	if err := c.ShouldBindJSON(&paymentData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var card models.CreditCard
	if err := database.DB.Where("id = ? AND user_id = ?", cardID, userID).First(&card).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Credit card not found")
		return
	}

	// Update card balance and payment info
	card.CurrentBalance -= paymentData.Amount
	if card.CurrentBalance < 0 {
		card.CurrentBalance = 0
	}

	if paymentData.PaymentDate != nil {
		card.LastPaymentDate = paymentData.PaymentDate
	} else {
		now := time.Now()
		card.LastPaymentDate = &now
	}
	card.LastPaymentAmount = paymentData.Amount

	if err := database.DB.Save(&card).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to record payment")
		return
	}

	utilities.SuccessResponse(c, card, "Payment recorded successfully")
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

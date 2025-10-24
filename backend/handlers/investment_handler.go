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

// ListInvestments returns all investments for the authenticated user
func ListInvestments(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	query := database.DB.Where("user_id = ?", userID)

	// Optional filter by portfolio
	if portfolioID := c.Query("portfolioId"); portfolioID != "" {
		query = query.Where("portfolio_id = ?", portfolioID)
	}

	// Optional filter by asset type
	if assetType := c.Query("assetType"); assetType != "" {
		query = query.Where("asset_type = ?", assetType)
	}

	var investments []models.Investment
	if err := query.Order("created_at DESC").Find(&investments).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch investments")
		return
	}

	utilities.SuccessResponse(c, investments, "Investments retrieved successfully")
}

// GetInvestment returns a specific investment by ID
func GetInvestment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	investmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid investment ID")
		return
	}

	var investment models.Investment
	if err := database.DB.Where("id = ? AND user_id = ?", investmentID, userID).First(&investment).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Investment not found")
		return
	}

	utilities.SuccessResponse(c, investment, "Investment retrieved successfully")
}

// CreateInvestment creates a new investment
func CreateInvestment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var investment models.Investment
	if err := c.ShouldBindJSON(&investment); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	investment.UserID = userID
	investment.LastUpdated = time.Now()

	// If portfolio is specified, verify it belongs to user
	if investment.PortfolioID != nil {
		var portfolio models.Portfolio
		if err := database.DB.Where("id = ? AND user_id = ?", *investment.PortfolioID, userID).First(&portfolio).Error; err != nil {
			utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid portfolio ID")
			return
		}
	}

	if err := database.DB.Create(&investment).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create investment")
		return
	}

	utilities.CreatedResponse(c, investment, "Investment created successfully")
}

// UpdateInvestment updates an existing investment
func UpdateInvestment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	investmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid investment ID")
		return
	}

	var existingInvestment models.Investment
	if err := database.DB.Where("id = ? AND user_id = ?", investmentID, userID).First(&existingInvestment).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Investment not found")
		return
	}

	var updateData models.Investment
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update allowed fields
	existingInvestment.PortfolioID = updateData.PortfolioID
	existingInvestment.Symbol = updateData.Symbol
	existingInvestment.Name = updateData.Name
	existingInvestment.AssetType = updateData.AssetType
	existingInvestment.Quantity = updateData.Quantity
	existingInvestment.CostBasis = updateData.CostBasis
	existingInvestment.CurrentPrice = updateData.CurrentPrice
	existingInvestment.PurchaseDate = updateData.PurchaseDate
	existingInvestment.Notes = updateData.Notes
	existingInvestment.LastUpdated = time.Now()

	if err := database.DB.Save(&existingInvestment).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update investment")
		return
	}

	utilities.SuccessResponse(c, existingInvestment, "Investment updated successfully")
}

// DeleteInvestment deletes an investment
func DeleteInvestment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	investmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid investment ID")
		return
	}

	var investment models.Investment
	if err := database.DB.Where("id = ? AND user_id = ?", investmentID, userID).First(&investment).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Investment not found")
		return
	}

	// Soft delete
	if err := database.DB.Delete(&investment).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete investment")
		return
	}

	utilities.SuccessResponse(c, nil, "Investment deleted successfully")
}

// BuyShares increases the quantity of shares for an investment
func BuyShares(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	investmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid investment ID")
		return
	}

	var buyData struct {
		Quantity float64 `json:"quantity" binding:"required,gt=0"`
		Price    float64 `json:"price" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&buyData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var investment models.Investment
	if err := database.DB.Where("id = ? AND user_id = ?", investmentID, userID).First(&investment).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Investment not found")
		return
	}

	// Calculate new cost basis (weighted average)
	totalCost := (investment.Quantity * investment.CostBasis) + (buyData.Quantity * buyData.Price)
	totalQuantity := investment.Quantity + buyData.Quantity
	newCostBasis := totalCost / totalQuantity

	// Update investment
	investment.Quantity = totalQuantity
	investment.CostBasis = newCostBasis
	investment.CurrentPrice = buyData.Price
	investment.LastUpdated = time.Now()

	if err := database.DB.Save(&investment).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to buy shares")
		return
	}

	utilities.SuccessResponse(c, investment, "Shares purchased successfully")
}

// SellShares decreases the quantity of shares for an investment
func SellShares(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	investmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid investment ID")
		return
	}

	var sellData struct {
		Quantity float64 `json:"quantity" binding:"required,gt=0"`
		Price    float64 `json:"price" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&sellData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var investment models.Investment
	if err := database.DB.Where("id = ? AND user_id = ?", investmentID, userID).First(&investment).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Investment not found")
		return
	}

	// Verify sufficient quantity
	if sellData.Quantity > investment.Quantity {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Insufficient shares to sell")
		return
	}

	// Calculate realized gain/loss
	saleProceeds := sellData.Quantity * sellData.Price
	costOfSharesSold := sellData.Quantity * investment.CostBasis
	realizedGainLoss := saleProceeds - costOfSharesSold

	// Update investment
	investment.Quantity -= sellData.Quantity
	investment.CurrentPrice = sellData.Price
	investment.RealizedGainLoss += realizedGainLoss
	investment.LastUpdated = time.Now()

	if err := database.DB.Save(&investment).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to sell shares")
		return
	}

	result := map[string]interface{}{
		"investment":       investment,
		"realizedGainLoss": realizedGainLoss,
	}

	utilities.SuccessResponse(c, result, "Shares sold successfully")
}

// ListPortfolios returns all portfolios for the authenticated user
func ListPortfolios(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var portfolios []models.Portfolio
	if err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&portfolios).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch portfolios")
		return
	}

	utilities.SuccessResponse(c, portfolios, "Portfolios retrieved successfully")
}

// CreatePortfolio creates a new portfolio
func CreatePortfolio(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var portfolio models.Portfolio
	if err := c.ShouldBindJSON(&portfolio); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	portfolio.UserID = userID

	if err := database.DB.Create(&portfolio).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create portfolio")
		return
	}

	utilities.CreatedResponse(c, portfolio, "Portfolio created successfully")
}

// ListDividends returns dividends for the authenticated user
func ListDividends(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	query := database.DB.Where("user_id = ?", userID)

	// Optional filter by investment
	if investmentID := c.Query("investmentId"); investmentID != "" {
		query = query.Where("investment_id = ?", investmentID)
	}

	var dividends []models.Dividend
	if err := query.Order("payment_date DESC").Find(&dividends).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch dividends")
		return
	}

	utilities.SuccessResponse(c, dividends, "Dividends retrieved successfully")
}

// RecordDividend records a new dividend payment
func RecordDividend(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var dividend models.Dividend
	if err := c.ShouldBindJSON(&dividend); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	dividend.UserID = userID

	// Verify investment belongs to user
	var investment models.Investment
	if err := database.DB.Where("id = ? AND user_id = ?", dividend.InvestmentID, userID).First(&investment).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid investment ID")
		return
	}

	if err := database.DB.Create(&dividend).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to record dividend")
		return
	}

	utilities.CreatedResponse(c, dividend, "Dividend recorded successfully")
}

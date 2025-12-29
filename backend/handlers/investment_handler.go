package handlers

import (
	"net/http"
	"strconv"
	"time"

	"daybook-backend/database"
	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
)

// ListInvestments returns all investments for the authenticated user
func ListInvestments(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "ListInvestments - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "ListInvestments - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	logger.Debugf(ctx, "ListInvestments - Fetching investments for user: %s", userID)

	query := database.DB.WithContext(ctx).Where("user_id = ?", userID)

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
		logger.Errorf(ctx, "ListInvestments - Failed to fetch investments: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch investments")
		return
	}

	logger.Infof(ctx, "ListInvestments - Successfully retrieved %d investments", len(investments))
	utilities.SuccessResponse(c, investments, "Investments retrieved successfully")
}

// GetInvestment returns a specific investment by ID
func GetInvestment(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "GetInvestment - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "GetInvestment - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	investmentIDStr := c.Param("id")
	investmentIDUint, err := strconv.ParseUint(investmentIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "GetInvestment - Invalid investment ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid investment ID")
		return
	}
	investmentID := uint(investmentIDUint)

	logger.Debugf(ctx, "GetInvestment - Fetching investment: %s", investmentID)

	var investment models.Investment
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", investmentID, userID).First(&investment).Error; err != nil {
		logger.Errorf(ctx, "GetInvestment - Investment not found: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Investment not found")
		return
	}

	logger.Infof(ctx, "GetInvestment - Successfully retrieved investment: %s", investmentID)
	utilities.SuccessResponse(c, investment, "Investment retrieved successfully")
}

// CreateInvestment creates a new investment
func CreateInvestment(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "CreateInvestment - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "CreateInvestment - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	var investment models.Investment
	if err := c.ShouldBindJSON(&investment); err != nil {
		logger.Warnf(ctx, "CreateInvestment - Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	investment.UserID = userID
	investment.LastUpdated = time.Now()

	// If portfolio is specified, verify it belongs to user
	if investment.PortfolioID != nil {
		logger.Debugf(ctx, "CreateInvestment - Verifying portfolio: %s", *investment.PortfolioID)
		var portfolio models.Portfolio
		if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", *investment.PortfolioID, userID).First(&portfolio).Error; err != nil {
			logger.Errorf(ctx, "CreateInvestment - Invalid portfolio ID: %v", err)
			utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid portfolio ID")
			return
		}
	}

	logger.Debugf(ctx, "CreateInvestment - Creating investment: %s", investment.Name)

	if err := database.DB.WithContext(ctx).Create(&investment).Error; err != nil {
		logger.Errorf(ctx, "CreateInvestment - Failed to create investment: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create investment")
		return
	}

	// Log investment creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, "investment",
		"Investment", investment.ID, "Created investment: "+investment.Name, nil)

	logger.Infof(ctx, "CreateInvestment - Successfully created investment: %s", investment.ID)
	utilities.CreatedResponse(c, investment, "Investment created successfully")
}

// UpdateInvestment updates an existing investment
func UpdateInvestment(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "UpdateInvestment - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "UpdateInvestment - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	investmentIDStr := c.Param("id")
	investmentIDUint, err := strconv.ParseUint(investmentIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "UpdateInvestment - Invalid investment ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid investment ID")
		return
	}
	investmentID := uint(investmentIDUint)

	logger.Debugf(ctx, "UpdateInvestment - Updating investment: %s", investmentID)

	var existingInvestment models.Investment
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", investmentID, userID).First(&existingInvestment).Error; err != nil {
		logger.Errorf(ctx, "UpdateInvestment - Investment not found: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Investment not found")
		return
	}

	var updateData models.Investment
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "UpdateInvestment - Invalid request body: %v", err)
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

	if err := database.DB.WithContext(ctx).Save(&existingInvestment).Error; err != nil {
		logger.Errorf(ctx, "UpdateInvestment - Failed to update investment: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update investment")
		return
	}

	// Log investment update activity
	utilities.LogEntityActivity(c, userID, models.ActionUpdate, "investment",
		"Investment", existingInvestment.ID, "Updated investment: "+existingInvestment.Name, nil)

	logger.Infof(ctx, "UpdateInvestment - Successfully updated investment: %s", investmentID)
	utilities.SuccessResponse(c, existingInvestment, "Investment updated successfully")
}

// DeleteInvestment deletes an investment
func DeleteInvestment(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "DeleteInvestment - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "DeleteInvestment - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	investmentIDStr := c.Param("id")
	investmentIDUint, err := strconv.ParseUint(investmentIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "DeleteInvestment - Invalid investment ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid investment ID")
		return
	}
	investmentID := uint(investmentIDUint)

	logger.Debugf(ctx, "DeleteInvestment - Deleting investment: %s", investmentID)

	var investment models.Investment
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", investmentID, userID).First(&investment).Error; err != nil {
		logger.Errorf(ctx, "DeleteInvestment - Investment not found: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Investment not found")
		return
	}

	// Soft delete
	if err := database.DB.WithContext(ctx).Delete(&investment).Error; err != nil {
		logger.Errorf(ctx, "DeleteInvestment - Failed to delete investment: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete investment")
		return
	}

	// Log investment deletion activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, "investment",
		"Investment", investment.ID, "Deleted investment: "+investment.Name, nil)

	logger.Infof(ctx, "DeleteInvestment - Successfully deleted investment: %s", investmentID)
	utilities.SuccessResponse(c, nil, "Investment deleted successfully")
}

// BuyShares increases the quantity of shares for an investment
func BuyShares(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "BuyShares - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "BuyShares - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	investmentIDStr := c.Param("id")
	investmentIDUint, err := strconv.ParseUint(investmentIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "BuyShares - Invalid investment ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid investment ID")
		return
	}
	investmentID := uint(investmentIDUint)

	var buyData struct {
		Quantity float64 `json:"quantity" binding:"required,gt=0"`
		Price    float64 `json:"price" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&buyData); err != nil {
		logger.Warnf(ctx, "BuyShares - Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "BuyShares - Buying %f shares at price %f for investment: %s", buyData.Quantity, buyData.Price, investmentID)

	var investment models.Investment
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", investmentID, userID).First(&investment).Error; err != nil {
		logger.Errorf(ctx, "BuyShares - Investment not found: %v", err)
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

	if err := database.DB.WithContext(ctx).Save(&investment).Error; err != nil {
		logger.Errorf(ctx, "BuyShares - Failed to update investment: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to buy shares")
		return
	}

	// Log share purchase activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, "investment",
		"InvestmentTransaction", investment.ID, "Bought shares of "+investment.Name, nil)

	logger.Infof(ctx, "BuyShares - Successfully bought shares for investment: %s", investmentID)
	utilities.SuccessResponse(c, investment, "Shares purchased successfully")
}

// SellShares decreases the quantity of shares for an investment
func SellShares(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "SellShares - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "SellShares - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	investmentIDStr := c.Param("id")
	investmentIDUint, err := strconv.ParseUint(investmentIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "SellShares - Invalid investment ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid investment ID")
		return
	}
	investmentID := uint(investmentIDUint)

	var sellData struct {
		Quantity float64 `json:"quantity" binding:"required,gt=0"`
		Price    float64 `json:"price" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&sellData); err != nil {
		logger.Warnf(ctx, "SellShares - Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "SellShares - Selling %f shares at price %f for investment: %s", sellData.Quantity, sellData.Price, investmentID)

	var investment models.Investment
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", investmentID, userID).First(&investment).Error; err != nil {
		logger.Errorf(ctx, "SellShares - Investment not found: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Investment not found")
		return
	}

	// Verify sufficient quantity
	if sellData.Quantity > investment.Quantity {
		logger.Warnf(ctx, "SellShares - Insufficient shares to sell")
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

	if err := database.DB.WithContext(ctx).Save(&investment).Error; err != nil {
		logger.Errorf(ctx, "SellShares - Failed to update investment: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to sell shares")
		return
	}

	// Log share sale activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, "investment",
		"InvestmentTransaction", investment.ID, "Sold shares of "+investment.Name, nil)

	logger.Infof(ctx, "SellShares - Successfully sold shares for investment: %s", investmentID)

	result := map[string]interface{}{
		"investment":       investment,
		"realizedGainLoss": realizedGainLoss,
	}

	utilities.SuccessResponse(c, result, "Shares sold successfully")
}

// ListPortfolios returns all portfolios for the authenticated user
func ListPortfolios(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "ListPortfolios - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "ListPortfolios - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	logger.Debugf(ctx, "ListPortfolios - Fetching portfolios for user: %s", userID)

	var portfolios []models.Portfolio
	if err := database.DB.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&portfolios).Error; err != nil {
		logger.Errorf(ctx, "ListPortfolios - Failed to fetch portfolios: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch portfolios")
		return
	}

	logger.Infof(ctx, "ListPortfolios - Successfully retrieved %d portfolios", len(portfolios))
	utilities.SuccessResponse(c, portfolios, "Portfolios retrieved successfully")
}

// CreatePortfolio creates a new portfolio
func CreatePortfolio(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "CreatePortfolio - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "CreatePortfolio - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	var portfolio models.Portfolio
	if err := c.ShouldBindJSON(&portfolio); err != nil {
		logger.Warnf(ctx, "CreatePortfolio - Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	portfolio.UserID = userID

	logger.Debugf(ctx, "CreatePortfolio - Creating portfolio: %s", portfolio.Name)

	if err := database.DB.WithContext(ctx).Create(&portfolio).Error; err != nil {
		logger.Errorf(ctx, "CreatePortfolio - Failed to create portfolio: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create portfolio")
		return
	}

	// Log portfolio creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, "investment",
		"Portfolio", portfolio.ID, "Created portfolio: "+portfolio.Name, nil)

	logger.Infof(ctx, "CreatePortfolio - Successfully created portfolio: %s", portfolio.ID)
	utilities.CreatedResponse(c, portfolio, "Portfolio created successfully")
}

// ListDividends returns dividends for the authenticated user
func ListDividends(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "ListDividends - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "ListDividends - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	logger.Debugf(ctx, "ListDividends - Fetching dividends for user: %s", userID)

	query := database.DB.WithContext(ctx).Where("user_id = ?", userID)

	// Optional filter by investment
	if investmentID := c.Query("investmentId"); investmentID != "" {
		query = query.Where("investment_id = ?", investmentID)
	}

	var dividends []models.Dividend
	if err := query.Order("payment_date DESC").Find(&dividends).Error; err != nil {
		logger.Errorf(ctx, "ListDividends - Failed to fetch dividends: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch dividends")
		return
	}

	logger.Infof(ctx, "ListDividends - Successfully retrieved %d dividends", len(dividends))
	utilities.SuccessResponse(c, dividends, "Dividends retrieved successfully")
}

// RecordDividend records a new dividend payment
func RecordDividend(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "RecordDividend - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "RecordDividend - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	var dividend models.Dividend
	if err := c.ShouldBindJSON(&dividend); err != nil {
		logger.Warnf(ctx, "RecordDividend - Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	dividend.UserID = userID

	// Verify investment belongs to user
	logger.Debugf(ctx, "RecordDividend - Verifying investment: %s", dividend.InvestmentID)
	var investment models.Investment
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", dividend.InvestmentID, userID).First(&investment).Error; err != nil {
		logger.Errorf(ctx, "RecordDividend - Invalid investment ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid investment ID")
		return
	}

	logger.Debugf(ctx, "RecordDividend - Recording dividend for investment: %s", investment.Name)

	if err := database.DB.WithContext(ctx).Create(&dividend).Error; err != nil {
		logger.Errorf(ctx, "RecordDividend - Failed to record dividend: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to record dividend")
		return
	}

	// Log dividend record activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, "investment",
		"Dividend", dividend.ID, "Recorded dividend from "+investment.Name, nil)

	logger.Infof(ctx, "RecordDividend - Successfully recorded dividend: %s", dividend.ID)
	utilities.CreatedResponse(c, dividend, "Dividend recorded successfully")
}

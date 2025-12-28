package handlers

import (
	"net/http"
	"strconv"
	"time"

	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReportHandler handles all report-related API endpoints
type ReportHandler struct {
	service services.ReportService
}

// NewReportHandler creates a new report handler instance
func NewReportHandler(service services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GetDashboardSummary returns a comprehensive dashboard with key metrics
// GET /api/v1/reports/dashboard
func (h *ReportHandler) GetDashboardSummary(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Errorf(ctx, "Failed to get user ID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	logger.Infof(ctx, "Fetching dashboard summary for user: %s", userID)

	summary, err := h.service.GetDashboardSummary(ctx, userID)
	if err != nil {
		logger.Errorf(ctx, "Failed to get dashboard summary: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate dashboard summary"})
		return
	}

	logger.Infof(ctx, "Dashboard summary retrieved successfully for user: %s", userID)
	c.JSON(http.StatusOK, summary)
}

// GetIncomeExpenseReport returns income vs expense analysis with trends
// GET /api/v1/reports/income-expense?startDate=2024-01-01&endDate=2024-12-31&groupBy=month
func (h *ReportHandler) GetIncomeExpenseReport(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Errorf(ctx, "Failed to get user ID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Parse query parameters
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	groupBy := c.DefaultQuery("groupBy", "month")

	// Validate and parse dates
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		logger.Errorf(ctx, "Invalid start date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		logger.Errorf(ctx, "Invalid end date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	// Validate groupBy
	validGroupBy := map[string]bool{"day": true, "week": true, "month": true, "year": true}
	if !validGroupBy[groupBy] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid groupBy. Must be one of: day, week, month, year"})
		return
	}

	logger.Infof(ctx, "Fetching income-expense report for user: %s, period: %s to %s, groupBy: %s",
		userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), groupBy)

	req := &services.TimeRangeRequest{
		StartDate: startDate,
		EndDate:   endDate,
		GroupBy:   groupBy,
	}

	report, err := h.service.GetIncomeExpenseReport(ctx, userID, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to get income-expense report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate income-expense report"})
		return
	}

	logger.Infof(ctx, "Income-expense report retrieved successfully for user: %s", userID)
	c.JSON(http.StatusOK, report)
}

// GetCategoryAnalysis returns category breakdown and top categories
// GET /api/v1/reports/category-analysis?startDate=2024-01-01&endDate=2024-12-31&type=expense&limit=10
func (h *ReportHandler) GetCategoryAnalysis(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Errorf(ctx, "Failed to get user ID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Parse query parameters
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	transactionType := c.DefaultQuery("type", "expense")
	limitStr := c.DefaultQuery("limit", "10")

	// Validate and parse dates
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		logger.Errorf(ctx, "Invalid start date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		logger.Errorf(ctx, "Invalid end date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	// Parse limit
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Validate transaction type
	validTypes := map[string]bool{"income": true, "expense": true}
	if !validTypes[transactionType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type. Must be 'income' or 'expense'"})
		return
	}

	logger.Infof(ctx, "Fetching category analysis for user: %s, type: %s, period: %s to %s, limit: %d",
		userID, transactionType, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), limit)

	req := &services.CategoryAnalysisRequest{
		TimeRangeRequest: services.TimeRangeRequest{
			StartDate: startDate,
			EndDate:   endDate,
			GroupBy:   "month", // Default grouping for category analysis
		},
		Type:  transactionType,
		Limit: limit,
	}

	report, err := h.service.GetCategoryAnalysis(ctx, userID, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to get category analysis: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate category analysis"})
		return
	}

	logger.Infof(ctx, "Category analysis retrieved successfully for user: %s", userID)
	c.JSON(http.StatusOK, report)
}

// GetAccountReport returns current balances for all accounts
// GET /api/v1/reports/accounts
func (h *ReportHandler) GetAccountReport(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Errorf(ctx, "Failed to get user ID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	logger.Infof(ctx, "Fetching account report for user: %s", userID)

	report, err := h.service.GetAccountReport(ctx, userID)
	if err != nil {
		logger.Errorf(ctx, "Failed to get account report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate account report"})
		return
	}

	logger.Infof(ctx, "Account report retrieved successfully for user: %s", userID)
	c.JSON(http.StatusOK, report)
}

// GetAccountBalanceHistory returns balance history for a specific account
// GET /api/v1/reports/accounts/:id/history?startDate=2024-01-01&endDate=2024-12-31&groupBy=month
func (h *ReportHandler) GetAccountBalanceHistory(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Errorf(ctx, "Failed to get user ID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Parse account ID from URL
	accountIDStr := c.Param("id")
	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		logger.Errorf(ctx, "Invalid account ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// Parse query parameters
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	groupBy := c.DefaultQuery("groupBy", "month")

	// Validate and parse dates
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		logger.Errorf(ctx, "Invalid start date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		logger.Errorf(ctx, "Invalid end date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	// Validate groupBy
	validGroupBy := map[string]bool{"day": true, "week": true, "month": true, "year": true}
	if !validGroupBy[groupBy] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid groupBy. Must be one of: day, week, month, year"})
		return
	}

	logger.Infof(ctx, "Fetching account balance history for user: %s, account: %s, period: %s to %s, groupBy: %s",
		userID, accountID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), groupBy)

	req := &services.TimeRangeRequest{
		StartDate: startDate,
		EndDate:   endDate,
		GroupBy:   groupBy,
	}

	report, err := h.service.GetAccountBalanceHistory(ctx, userID, accountID, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to get account balance history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate account balance history"})
		return
	}

	logger.Infof(ctx, "Account balance history retrieved successfully for user: %s, account: %s", userID, accountID)
	c.JSON(http.StatusOK, report)
}

// GetNetWorthReport returns net worth with trend analysis
// GET /api/v1/reports/net-worth?startDate=2024-01-01&endDate=2024-12-31&groupBy=month
func (h *ReportHandler) GetNetWorthReport(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Errorf(ctx, "Failed to get user ID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Parse query parameters (optional for net worth)
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	groupBy := c.DefaultQuery("groupBy", "month")

	var req *services.TimeRangeRequest
	if startDateStr != "" && endDateStr != "" {
		// Validate and parse dates
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			logger.Errorf(ctx, "Invalid start date: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			logger.Errorf(ctx, "Invalid end date: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
			return
		}

		// Validate groupBy
		validGroupBy := map[string]bool{"day": true, "week": true, "month": true, "year": true}
		if !validGroupBy[groupBy] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid groupBy. Must be one of: day, week, month, year"})
			return
		}

		req = &services.TimeRangeRequest{
			StartDate: startDate,
			EndDate:   endDate,
			GroupBy:   groupBy,
		}

		logger.Infof(ctx, "Fetching net worth report with trend for user: %s, period: %s to %s, groupBy: %s",
			userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), groupBy)
	} else {
		logger.Infof(ctx, "Fetching current net worth for user: %s", userID)
	}

	report, err := h.service.GetNetWorthReport(ctx, userID, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to get net worth report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate net worth report"})
		return
	}

	logger.Infof(ctx, "Net worth report retrieved successfully for user: %s", userID)
	c.JSON(http.StatusOK, report)
}

// GetBudgetReport returns budget performance for a specific month
// GET /api/v1/reports/budget?month=2024-01
func (h *ReportHandler) GetBudgetReport(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Errorf(ctx, "Failed to get user ID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Parse month parameter
	monthStr := c.Query("month")
	if monthStr == "" {
		// Default to current month
		now := time.Now()
		monthStr = now.Format("2006-01")
	}

	// Parse month
	month, err := time.Parse("2006-01", monthStr)
	if err != nil {
		logger.Errorf(ctx, "Invalid month format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month format. Use YYYY-MM"})
		return
	}

	logger.Infof(ctx, "Fetching budget report for user: %s, month: %s", userID, month.Format("2006-01"))

	report, err := h.service.GetBudgetReport(ctx, userID, month)
	if err != nil {
		logger.Errorf(ctx, "Failed to get budget report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate budget report"})
		return
	}

	logger.Infof(ctx, "Budget report retrieved successfully for user: %s", userID)
	c.JSON(http.StatusOK, report)
}

// GetCashFlowReport returns cash flow analysis with monthly breakdown
// GET /api/v1/reports/cash-flow?startDate=2024-01-01&endDate=2024-12-31&groupBy=month
func (h *ReportHandler) GetCashFlowReport(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Errorf(ctx, "Failed to get user ID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Parse query parameters
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	groupBy := c.DefaultQuery("groupBy", "month")

	// Validate and parse dates
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		logger.Errorf(ctx, "Invalid start date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		logger.Errorf(ctx, "Invalid end date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD"})
		return
	}

	// Validate groupBy
	validGroupBy := map[string]bool{"day": true, "week": true, "month": true, "year": true}
	if !validGroupBy[groupBy] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid groupBy. Must be one of: day, week, month, year"})
		return
	}

	logger.Infof(ctx, "Fetching cash flow report for user: %s, period: %s to %s, groupBy: %s",
		userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), groupBy)

	req := &services.TimeRangeRequest{
		StartDate: startDate,
		EndDate:   endDate,
		GroupBy:   groupBy,
	}

	report, err := h.service.GetCashFlowReport(ctx, userID, req)
	if err != nil {
		logger.Errorf(ctx, "Failed to get cash flow report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate cash flow report"})
		return
	}

	logger.Infof(ctx, "Cash flow report retrieved successfully for user: %s", userID)
	c.JSON(http.StatusOK, report)
}

// GetMonthlySummary returns comprehensive summary for a specific month
// GET /api/v1/reports/monthly-summary?month=2024-01
func (h *ReportHandler) GetMonthlySummary(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Errorf(ctx, "Failed to get user ID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Parse month parameter
	monthStr := c.Query("month")
	if monthStr == "" {
		// Default to current month
		now := time.Now()
		monthStr = now.Format("2006-01")
	}

	// Parse month
	month, err := time.Parse("2006-01", monthStr)
	if err != nil {
		logger.Errorf(ctx, "Invalid month format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month format. Use YYYY-MM"})
		return
	}

	logger.Infof(ctx, "Fetching monthly summary for user: %s, month: %s", userID, month.Format("2006-01"))

	report, err := h.service.GetMonthlySummary(ctx, userID, month)
	if err != nil {
		logger.Errorf(ctx, "Failed to get monthly summary: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate monthly summary"})
		return
	}

	logger.Infof(ctx, "Monthly summary retrieved successfully for user: %s", userID)
	c.JSON(http.StatusOK, report)
}

// GetYearlySummary returns comprehensive summary for a specific year
// GET /api/v1/reports/yearly-summary?year=2024
func (h *ReportHandler) GetYearlySummary(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Errorf(ctx, "Failed to get user ID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Parse year parameter
	yearStr := c.Query("year")
	if yearStr == "" {
		// Default to current year
		now := time.Now()
		yearStr = strconv.Itoa(now.Year())
	}

	// Parse year
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 1900 || year > 2100 {
		logger.Errorf(ctx, "Invalid year: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}

	logger.Infof(ctx, "Fetching yearly summary for user: %s, year: %d", userID, year)

	report, err := h.service.GetYearlySummary(ctx, userID, year)
	if err != nil {
		logger.Errorf(ctx, "Failed to get yearly summary: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate yearly summary"})
		return
	}

	logger.Infof(ctx, "Yearly summary retrieved successfully for user: %s", userID)
	c.JSON(http.StatusOK, report)
}

// GetPeriodComparison compares two time periods
// POST /api/v1/reports/comparison
func (h *ReportHandler) GetPeriodComparison(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Errorf(ctx, "Failed to get user ID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req services.ComparisonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf(ctx, "Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate dates
	if req.Period1Start.IsZero() || req.Period1End.IsZero() || req.Period2Start.IsZero() || req.Period2End.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All period dates are required"})
		return
	}

	if req.Period1Start.After(req.Period1End) || req.Period2Start.After(req.Period2End) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date must be before end date"})
		return
	}

	logger.Infof(ctx, "Fetching period comparison for user: %s, period1: %s to %s, period2: %s to %s",
		userID,
		req.Period1Start.Format("2006-01-02"), req.Period1End.Format("2006-01-02"),
		req.Period2Start.Format("2006-01-02"), req.Period2End.Format("2006-01-02"))

	report, err := h.service.GetPeriodComparison(ctx, userID, &req)
	if err != nil {
		logger.Errorf(ctx, "Failed to get period comparison: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate period comparison"})
		return
	}

	logger.Infof(ctx, "Period comparison retrieved successfully for user: %s", userID)
	c.JSON(http.StatusOK, report)
}

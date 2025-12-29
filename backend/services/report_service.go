package services

import (
	"context"
	"fmt"
	"time"

	"daybook-backend/repository"
)

// ReportService defines business logic for reports
type ReportService interface {
	// Dashboard Summary
	GetDashboardSummary(ctx context.Context, userID uint) (*DashboardSummary, error)

	// Income vs Expense
	GetIncomeExpenseReport(ctx context.Context, userID uint, req *TimeRangeRequest) (*IncomeExpenseReport, error)

	// Category Analysis
	GetCategoryAnalysis(ctx context.Context, userID uint, req *CategoryAnalysisRequest) (*CategoryAnalysisReport, error)

	// Account Reports
	GetAccountReport(ctx context.Context, userID uint) (*AccountReport, error)
	GetAccountBalanceHistory(ctx context.Context, userID uint, accountID uint, req *TimeRangeRequest) (*AccountHistoryReport, error)

	// Net Worth
	GetNetWorthReport(ctx context.Context, userID uint, req *TimeRangeRequest) (*NetWorthReport, error)

	// Budget Reports
	GetBudgetReport(ctx context.Context, userID uint, month time.Time) (*BudgetReport, error)

	// Cash Flow
	GetCashFlowReport(ctx context.Context, userID uint, req *TimeRangeRequest) (*CashFlowReport, error)

	// Period Summaries
	GetMonthlySummary(ctx context.Context, userID uint, month time.Time) (*PeriodSummaryReport, error)
	GetYearlySummary(ctx context.Context, userID uint, year int) (*PeriodSummaryReport, error)

	// Comparisons
	GetPeriodComparison(ctx context.Context, userID uint, req *ComparisonRequest) (*ComparisonReport, error)
}

// Request/Response structures

type TimeRangeRequest struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	GroupBy   string    `json:"groupBy"` // day, week, month, year
}

type CategoryAnalysisRequest struct {
	TimeRangeRequest
	Type  string `json:"type"`  // income or expense
	Limit int    `json:"limit"` // for top categories
}

type ComparisonRequest struct {
	Period1Start time.Time `json:"period1Start"`
	Period1End   time.Time `json:"period1End"`
	Period2Start time.Time `json:"period2Start"`
	Period2End   time.Time `json:"period2End"`
}

type DashboardSummary struct {
	CurrentMonth    *repository.PeriodSummary      `json:"currentMonth"`
	PreviousMonth   *repository.PeriodSummary      `json:"previousMonth"`
	NetWorth        *repository.NetWorth           `json:"netWorth"`
	TopExpenses     []repository.CategoryBreakdown `json:"topExpenses"`
	RecentTrend     []repository.TrendData         `json:"recentTrend"`
	AccountBalances []repository.AccountBalance    `json:"accountBalances"`
	BudgetSummary   *BudgetSummaryData             `json:"budgetSummary"`
}

type BudgetSummaryData struct {
	TotalBudget  float64 `json:"totalBudget"`
	TotalSpent   float64 `json:"totalSpent"`
	Remaining    float64 `json:"remaining"`
	UtilizedPerc float64 `json:"utilizedPerc"`
}

type IncomeExpenseReport struct {
	Summary *repository.IncomeExpenseSummary `json:"summary"`
	Trend   []repository.TrendData           `json:"trend"`
	Period  string                           `json:"period"`
}

type CategoryAnalysisReport struct {
	Type       string                         `json:"type"`
	Breakdown  []repository.CategoryBreakdown `json:"breakdown"`
	TopN       []repository.CategoryBreakdown `json:"topCategories"`
	Period     string                         `json:"period"`
	TotalCount int                            `json:"totalCount"`
}

type AccountReport struct {
	Balances   []repository.AccountBalance `json:"balances"`
	TotalAsset float64                     `json:"totalAsset"`
	TotalDebt  float64                     `json:"totalDebt"`
}

type AccountHistoryReport struct {
	AccountID   uint                        `json:"accountId"`
	AccountName string                      `json:"accountName"`
	History     []repository.BalanceHistory `json:"history"`
	StartDate   time.Time                   `json:"startDate"`
	EndDate     time.Time                   `json:"endDate"`
}

type NetWorthReport struct {
	Current *repository.NetWorth       `json:"current"`
	Trend   []repository.NetWorthTrend `json:"trend"`
}

type BudgetReport struct {
	Month       time.Time                      `json:"month"`
	Performance []repository.BudgetPerformance `json:"performance"`
	TotalBudget float64                        `json:"totalBudget"`
	TotalSpent  float64                        `json:"totalSpent"`
	OverallPerc float64                        `json:"overallPerc"`
	Status      string                         `json:"status"`
}

type CashFlowReport struct {
	Summary         *repository.CashFlowSummary  `json:"summary"`
	MonthlyCashFlow []repository.MonthlyCashFlow `json:"monthlyCashFlow"`
	Period          string                       `json:"period"`
}

type PeriodSummaryReport struct {
	Summary           *repository.PeriodSummary      `json:"summary"`
	CategoryBreakdown []repository.CategoryBreakdown `json:"categoryBreakdown"`
	TopExpenses       []repository.CategoryBreakdown `json:"topExpenses"`
}

type ComparisonReport struct {
	Comparison *repository.PeriodComparison `json:"comparison"`
	Insights   []string                     `json:"insights"`
}

type reportService struct {
	repo repository.ReportRepository
}

// NewReportService creates a new report service
func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

// GetDashboardSummary retrieves dashboard summary with key metrics
func (s *reportService) GetDashboardSummary(ctx context.Context, userID uint) (*DashboardSummary, error) {
	summary := &DashboardSummary{}

	now := time.Now()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	currentMonthEnd := currentMonthStart.AddDate(0, 1, 0).Add(-time.Second)

	previousMonthStart := currentMonthStart.AddDate(0, -1, 0)

	// Get current month summary
	currentMonth, err := s.repo.GetMonthlySummary(ctx, userID, now)
	if err == nil {
		summary.CurrentMonth = currentMonth
	}

	// Get previous month summary
	previousMonth, err := s.repo.GetMonthlySummary(ctx, userID, previousMonthStart)
	if err == nil {
		summary.PreviousMonth = previousMonth
	}

	// Get net worth
	netWorth, err := s.repo.GetNetWorth(ctx, userID)
	if err == nil {
		summary.NetWorth = netWorth
	}

	// Get top 5 expense categories for current month
	topExpenses, err := s.repo.GetTopCategories(ctx, userID, currentMonthStart, currentMonthEnd, "expense", 5)
	if err == nil {
		summary.TopExpenses = topExpenses
	}

	// Get last 6 months trend
	sixMonthsAgo := currentMonthStart.AddDate(0, -5, 0)
	trend, err := s.repo.GetIncomeExpenseTrend(ctx, userID, sixMonthsAgo, currentMonthEnd, "month")
	if err == nil {
		summary.RecentTrend = trend
	}

	// Get account balances
	balances, err := s.repo.GetAccountBalances(ctx, userID)
	if err == nil {
		summary.AccountBalances = balances
	}

	// Get budget summary for current month
	budgetPerf, err := s.repo.GetBudgetPerformance(ctx, userID, now)
	if err == nil {
		budgetSummary := &BudgetSummaryData{}
		for _, perf := range budgetPerf {
			budgetSummary.TotalBudget += perf.BudgetAmount
			budgetSummary.TotalSpent += perf.ActualAmount
		}
		budgetSummary.Remaining = budgetSummary.TotalBudget - budgetSummary.TotalSpent
		if budgetSummary.TotalBudget > 0 {
			budgetSummary.UtilizedPerc = (budgetSummary.TotalSpent / budgetSummary.TotalBudget) * 100
		}
		summary.BudgetSummary = budgetSummary
	}

	return summary, nil
}

// GetIncomeExpenseReport retrieves income vs expense report with trends
func (s *reportService) GetIncomeExpenseReport(ctx context.Context, userID uint, req *TimeRangeRequest) (*IncomeExpenseReport, error) {
	report := &IncomeExpenseReport{
		Period: fmt.Sprintf("%s to %s", req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02")),
	}

	// Get summary
	summary, err := s.repo.GetIncomeExpenseSummary(ctx, userID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	report.Summary = summary

	// Get trend
	groupBy := req.GroupBy
	if groupBy == "" {
		groupBy = "month" // default
	}

	trend, err := s.repo.GetIncomeExpenseTrend(ctx, userID, req.StartDate, req.EndDate, groupBy)
	if err != nil {
		return nil, err
	}
	report.Trend = trend

	return report, nil
}

// GetCategoryAnalysis retrieves category-wise analysis
func (s *reportService) GetCategoryAnalysis(ctx context.Context, userID uint, req *CategoryAnalysisRequest) (*CategoryAnalysisReport, error) {
	report := &CategoryAnalysisReport{
		Type:   req.Type,
		Period: fmt.Sprintf("%s to %s", req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02")),
	}

	// Get full breakdown
	breakdown, err := s.repo.GetCategoryBreakdown(ctx, userID, req.StartDate, req.EndDate, req.Type)
	if err != nil {
		return nil, err
	}
	report.Breakdown = breakdown
	report.TotalCount = len(breakdown)

	// Get top N categories
	limit := req.Limit
	if limit == 0 {
		limit = 10 // default top 10
	}

	topN, err := s.repo.GetTopCategories(ctx, userID, req.StartDate, req.EndDate, req.Type, limit)
	if err != nil {
		return nil, err
	}
	report.TopN = topN

	return report, nil
}

// GetAccountReport retrieves account balances report
func (s *reportService) GetAccountReport(ctx context.Context, userID uint) (*AccountReport, error) {
	report := &AccountReport{}

	balances, err := s.repo.GetAccountBalances(ctx, userID)
	if err != nil {
		return nil, err
	}
	report.Balances = balances

	// Calculate totals
	for _, balance := range balances {
		if balance.Balance >= 0 {
			report.TotalAsset += balance.Balance
		} else {
			report.TotalDebt += balance.Balance
		}
	}

	return report, nil
}

// GetAccountBalanceHistory retrieves balance history for an account
func (s *reportService) GetAccountBalanceHistory(ctx context.Context, userID uint, accountID uint, req *TimeRangeRequest) (*AccountHistoryReport, error) {
	report := &AccountHistoryReport{
		AccountID: accountID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	groupBy := req.GroupBy
	if groupBy == "" {
		groupBy = "month"
	}

	history, err := s.repo.GetAccountBalanceHistory(ctx, userID, accountID, req.StartDate, req.EndDate, groupBy)
	if err != nil {
		return nil, err
	}
	report.History = history

	// Get account name
	balances, err := s.repo.GetAccountBalances(ctx, userID)
	if err == nil {
		for _, balance := range balances {
			if balance.AccountID == accountID {
				report.AccountName = balance.AccountName
				break
			}
		}
	}

	return report, nil
}

// GetNetWorthReport retrieves net worth report with trend
func (s *reportService) GetNetWorthReport(ctx context.Context, userID uint, req *TimeRangeRequest) (*NetWorthReport, error) {
	report := &NetWorthReport{}

	// Get current net worth
	current, err := s.repo.GetNetWorth(ctx, userID)
	if err != nil {
		return nil, err
	}
	report.Current = current

	// Get trend
	groupBy := req.GroupBy
	if groupBy == "" {
		groupBy = "month"
	}

	trend, err := s.repo.GetNetWorthTrend(ctx, userID, req.StartDate, req.EndDate, groupBy)
	if err != nil {
		return nil, err
	}
	report.Trend = trend

	return report, nil
}

// GetBudgetReport retrieves budget performance report
func (s *reportService) GetBudgetReport(ctx context.Context, userID uint, month time.Time) (*BudgetReport, error) {
	report := &BudgetReport{
		Month: month,
	}

	performance, err := s.repo.GetBudgetPerformance(ctx, userID, month)
	if err != nil {
		return nil, err
	}
	report.Performance = performance

	// Calculate totals
	for _, perf := range performance {
		report.TotalBudget += perf.BudgetAmount
		report.TotalSpent += perf.ActualAmount
	}

	if report.TotalBudget > 0 {
		report.OverallPerc = (report.TotalSpent / report.TotalBudget) * 100
	}

	// Determine overall status
	if report.TotalSpent > report.TotalBudget {
		report.Status = "over"
	} else if report.TotalSpent >= report.TotalBudget*0.9 {
		report.Status = "on-track"
	} else {
		report.Status = "under"
	}

	return report, nil
}

// GetCashFlowReport retrieves cash flow report
func (s *reportService) GetCashFlowReport(ctx context.Context, userID uint, req *TimeRangeRequest) (*CashFlowReport, error) {
	report := &CashFlowReport{
		Period: fmt.Sprintf("%s to %s", req.StartDate.Format("2006-01-02"), req.EndDate.Format("2006-01-02")),
	}

	// Get summary
	summary, err := s.repo.GetCashFlowSummary(ctx, userID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}
	report.Summary = summary

	// Get monthly cash flow breakdown for the date range
	monthlyCashFlow, err := s.repo.GetCashFlowByDateRange(ctx, userID, req.StartDate, req.EndDate)
	if err == nil {
		report.MonthlyCashFlow = monthlyCashFlow
	}

	return report, nil
}

// GetMonthlySummary retrieves monthly summary report
func (s *reportService) GetMonthlySummary(ctx context.Context, userID uint, month time.Time) (*PeriodSummaryReport, error) {
	report := &PeriodSummaryReport{}

	summary, err := s.repo.GetMonthlySummary(ctx, userID, month)
	if err != nil {
		return nil, err
	}
	report.Summary = summary

	// Get category breakdown for the month
	startDate := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	breakdown, err := s.repo.GetCategoryBreakdown(ctx, userID, startDate, endDate, "expense")
	if err == nil {
		report.CategoryBreakdown = breakdown
	}

	topExpenses, err := s.repo.GetTopCategories(ctx, userID, startDate, endDate, "expense", 5)
	if err == nil {
		report.TopExpenses = topExpenses
	}

	return report, nil
}

// GetYearlySummary retrieves yearly summary report
func (s *reportService) GetYearlySummary(ctx context.Context, userID uint, year int) (*PeriodSummaryReport, error) {
	report := &PeriodSummaryReport{}

	summary, err := s.repo.GetYearlySummary(ctx, userID, year)
	if err != nil {
		return nil, err
	}
	report.Summary = summary

	// Get category breakdown for the year
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	breakdown, err := s.repo.GetCategoryBreakdown(ctx, userID, startDate, endDate, "expense")
	if err == nil {
		report.CategoryBreakdown = breakdown
	}

	topExpenses, err := s.repo.GetTopCategories(ctx, userID, startDate, endDate, "expense", 10)
	if err == nil {
		report.TopExpenses = topExpenses
	}

	return report, nil
}

// GetPeriodComparison retrieves period comparison report with insights
func (s *reportService) GetPeriodComparison(ctx context.Context, userID uint, req *ComparisonRequest) (*ComparisonReport, error) {
	report := &ComparisonReport{
		Insights: []string{},
	}

	comparison, err := s.repo.GetPeriodComparison(ctx, userID, req.Period1Start, req.Period1End, req.Period2Start, req.Period2End)
	if err != nil {
		return nil, err
	}
	report.Comparison = comparison

	// Generate insights
	if comparison.Changes.IncomeChange > 0 {
		report.Insights = append(report.Insights, fmt.Sprintf("Income increased by %.2f%% compared to previous period", comparison.Changes.IncomeChangePerc))
	} else if comparison.Changes.IncomeChange < 0 {
		report.Insights = append(report.Insights, fmt.Sprintf("Income decreased by %.2f%% compared to previous period", -comparison.Changes.IncomeChangePerc))
	}

	if comparison.Changes.ExpenseChange > 0 {
		report.Insights = append(report.Insights, fmt.Sprintf("Expenses increased by %.2f%% compared to previous period", comparison.Changes.ExpenseChangePerc))
	} else if comparison.Changes.ExpenseChange < 0 {
		report.Insights = append(report.Insights, fmt.Sprintf("Expenses decreased by %.2f%% compared to previous period", -comparison.Changes.ExpenseChangePerc))
	}

	if comparison.Changes.NetChange > 0 {
		report.Insights = append(report.Insights, "Net savings improved compared to previous period")
	} else if comparison.Changes.NetChange < 0 {
		report.Insights = append(report.Insights, "Net savings decreased compared to previous period")
	}

	return report, nil
}

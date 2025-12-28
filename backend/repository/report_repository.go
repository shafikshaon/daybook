package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReportRepository handles report data aggregation
type ReportRepository interface {
	// Income vs Expense
	GetIncomeExpenseSummary(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*IncomeExpenseSummary, error)
	GetIncomeExpenseTrend(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, groupBy string) ([]TrendData, error)

	// Category Analysis
	GetCategoryBreakdown(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, transactionType string) ([]CategoryBreakdown, error)
	GetTopCategories(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, transactionType string, limit int) ([]CategoryBreakdown, error)

	// Account Analysis
	GetAccountBalances(ctx context.Context, userID uuid.UUID) ([]AccountBalance, error)
	GetAccountBalanceHistory(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, startDate, endDate time.Time, groupBy string) ([]BalanceHistory, error)

	// Net Worth
	GetNetWorth(ctx context.Context, userID uuid.UUID) (*NetWorth, error)
	GetNetWorthTrend(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, groupBy string) ([]NetWorthTrend, error)

	// Budget Analysis
	GetBudgetPerformance(ctx context.Context, userID uuid.UUID, month time.Time) ([]BudgetPerformance, error)

	// Cash Flow
	GetCashFlowSummary(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*CashFlowSummary, error)
	GetCashFlowByMonth(ctx context.Context, userID uuid.UUID, year int) ([]MonthlyCashFlow, error)

	// Monthly/Yearly Summary
	GetMonthlySummary(ctx context.Context, userID uuid.UUID, month time.Time) (*PeriodSummary, error)
	GetYearlySummary(ctx context.Context, userID uuid.UUID, year int) (*PeriodSummary, error)

	// Comparison Reports
	GetPeriodComparison(ctx context.Context, userID uuid.UUID, period1Start, period1End, period2Start, period2End time.Time) (*PeriodComparison, error)
}

// Data structures for reports

type IncomeExpenseSummary struct {
	TotalIncome   float64 `json:"totalIncome"`
	TotalExpense  float64 `json:"totalExpense"`
	NetAmount     float64 `json:"netAmount"`
	TotalTransfer float64 `json:"totalTransfer"`
}

type TrendData struct {
	Period  string  `json:"period"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Net     float64 `json:"net"`
}

type CategoryBreakdown struct {
	CategoryID   string  `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	Amount       float64 `json:"amount"`
	Count        int     `json:"count"`
	Percentage   float64 `json:"percentage"`
}

type AccountBalance struct {
	AccountID   uuid.UUID `json:"accountId"`
	AccountName string    `json:"accountName"`
	Balance     float64   `json:"balance"`
	Type        string    `json:"type"`
}

type BalanceHistory struct {
	Period  string  `json:"period"`
	Balance float64 `json:"balance"`
}

type NetWorth struct {
	TotalAssets      float64 `json:"totalAssets"`
	TotalLiabilities float64 `json:"totalLiabilities"`
	NetWorth         float64 `json:"netWorth"`
}

type NetWorthTrend struct {
	Period           string  `json:"period"`
	TotalAssets      float64 `json:"totalAssets"`
	TotalLiabilities float64 `json:"totalLiabilities"`
	NetWorth         float64 `json:"netWorth"`
}

type BudgetPerformance struct {
	CategoryID   string  `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	BudgetAmount float64 `json:"budgetAmount"`
	ActualAmount float64 `json:"actualAmount"`
	Difference   float64 `json:"difference"`
	Percentage   float64 `json:"percentage"`
	Status       string  `json:"status"` // under, over, on-track
}

type CashFlowSummary struct {
	OpeningBalance float64 `json:"openingBalance"`
	TotalInflow    float64 `json:"totalInflow"`
	TotalOutflow   float64 `json:"totalOutflow"`
	ClosingBalance float64 `json:"closingBalance"`
	NetCashFlow    float64 `json:"netCashFlow"`
}

type MonthlyCashFlow struct {
	Month       string  `json:"month"`
	Inflow      float64 `json:"inflow"`
	Outflow     float64 `json:"outflow"`
	NetCashFlow float64 `json:"netCashFlow"`
}

type PeriodSummary struct {
	Period           string  `json:"period"`
	TotalIncome      float64 `json:"totalIncome"`
	TotalExpense     float64 `json:"totalExpense"`
	NetAmount        float64 `json:"netAmount"`
	TransactionCount int     `json:"transactionCount"`
	AvgDailySpending float64 `json:"avgDailySpending"`
}

type PeriodComparison struct {
	Period1 *PeriodSummary `json:"period1"`
	Period2 *PeriodSummary `json:"period2"`
	Changes *PeriodChanges `json:"changes"`
}

type PeriodChanges struct {
	IncomeChange      float64 `json:"incomeChange"`
	IncomeChangePerc  float64 `json:"incomeChangePerc"`
	ExpenseChange     float64 `json:"expenseChange"`
	ExpenseChangePerc float64 `json:"expenseChangePerc"`
	NetChange         float64 `json:"netChange"`
	NetChangePerc     float64 `json:"netChangePerc"`
}

type reportRepository struct {
	db *gorm.DB
}

// NewReportRepository creates a new report repository
func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

// GetIncomeExpenseSummary retrieves income vs expense summary
func (r *reportRepository) GetIncomeExpenseSummary(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*IncomeExpenseSummary, error) {
	var result IncomeExpenseSummary

	// Get total income
	r.db.WithContext(ctx).
		Table("transactions").
		Where("user_id = ? AND type = ? AND date >= ? AND date <= ? AND deleted_at IS NULL", userID, "income", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&result.TotalIncome)

	// Get total expense
	r.db.WithContext(ctx).
		Table("transactions").
		Where("user_id = ? AND type = ? AND date >= ? AND date <= ? AND deleted_at IS NULL", userID, "expense", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&result.TotalExpense)

	// Get total transfers
	r.db.WithContext(ctx).
		Table("transactions").
		Where("user_id = ? AND type = ? AND date >= ? AND date <= ? AND deleted_at IS NULL", userID, "transfer", startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&result.TotalTransfer)

	result.NetAmount = result.TotalIncome - result.TotalExpense

	return &result, nil
}

// GetIncomeExpenseTrend retrieves income vs expense trend over time
func (r *reportRepository) GetIncomeExpenseTrend(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, groupBy string) ([]TrendData, error) {
	var trends []TrendData

	// Determine SQL date grouping based on groupBy parameter
	var dateFormat string
	switch groupBy {
	case "day":
		dateFormat = "TO_CHAR(date, 'YYYY-MM-DD')"
	case "week":
		dateFormat = "TO_CHAR(DATE_TRUNC('week', date), 'YYYY-MM-DD')"
	case "month":
		dateFormat = "TO_CHAR(date, 'YYYY-MM')"
	case "year":
		dateFormat = "TO_CHAR(date, 'YYYY')"
	default:
		dateFormat = "TO_CHAR(date, 'YYYY-MM')" // default to month
	}

	query := `
		SELECT
			` + dateFormat + ` as period,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expense,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE -amount END), 0) as net
		FROM transactions
		WHERE user_id = ? AND date >= ? AND date <= ? AND deleted_at IS NULL AND type != 'tracking'
		GROUP BY period
		ORDER BY period ASC
	`

	err := r.db.WithContext(ctx).Raw(query, userID, startDate, endDate).Scan(&trends).Error
	return trends, err
}

// GetCategoryBreakdown retrieves spending/income by category
func (r *reportRepository) GetCategoryBreakdown(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, transactionType string) ([]CategoryBreakdown, error) {
	var breakdowns []CategoryBreakdown

	// First get the total for percentage calculation
	var total float64
	r.db.WithContext(ctx).
		Table("transactions").
		Where("user_id = ? AND type = ? AND date >= ? AND date <= ? AND deleted_at IS NULL", userID, transactionType, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total)

	query := `
		SELECT
			t.category_id as category_id,
			c.name as category_name,
			COALESCE(SUM(t.amount), 0) as amount,
			COUNT(*) as count,
			CASE WHEN ? > 0 THEN (COALESCE(SUM(t.amount), 0) / ? * 100) ELSE 0 END as percentage
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = ? AND t.type = ? AND t.date >= ? AND t.date <= ? AND t.deleted_at IS NULL
		GROUP BY t.category_id, c.name
		ORDER BY amount DESC
	`

	err := r.db.WithContext(ctx).Raw(query, total, total, userID, transactionType, startDate, endDate).Scan(&breakdowns).Error
	return breakdowns, err
}

// GetTopCategories retrieves top N categories by spending/income
func (r *reportRepository) GetTopCategories(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, transactionType string, limit int) ([]CategoryBreakdown, error) {
	var breakdowns []CategoryBreakdown

	var total float64
	r.db.WithContext(ctx).
		Table("transactions").
		Where("user_id = ? AND type = ? AND date >= ? AND date <= ? AND deleted_at IS NULL", userID, transactionType, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total)

	query := `
		SELECT
			t.category_id as category_id,
			c.name as category_name,
			COALESCE(SUM(t.amount), 0) as amount,
			COUNT(*) as count,
			CASE WHEN ? > 0 THEN (COALESCE(SUM(t.amount), 0) / ? * 100) ELSE 0 END as percentage
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = ? AND t.type = ? AND t.date >= ? AND t.date <= ? AND t.deleted_at IS NULL
		GROUP BY t.category_id, c.name
		ORDER BY amount DESC
		LIMIT ?
	`

	err := r.db.WithContext(ctx).Raw(query, total, total, userID, transactionType, startDate, endDate, limit).Scan(&breakdowns).Error
	return breakdowns, err
}

// GetAccountBalances retrieves current balances for all accounts
func (r *reportRepository) GetAccountBalances(ctx context.Context, userID uuid.UUID) ([]AccountBalance, error) {
	var balances []AccountBalance

	query := `
		SELECT
			a.id as account_id,
			a.name as account_name,
			a.balance,
			at.name as type
		FROM accounts a
		LEFT JOIN account_types at ON a.account_type_id = at.id
		WHERE a.user_id = ? AND a.deleted_at IS NULL
		ORDER BY a.balance DESC
	`

	err := r.db.WithContext(ctx).Raw(query, userID).Scan(&balances).Error
	return balances, err
}

// GetAccountBalanceHistory retrieves balance history for a specific account
func (r *reportRepository) GetAccountBalanceHistory(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, startDate, endDate time.Time, groupBy string) ([]BalanceHistory, error) {
	var history []BalanceHistory

	// This is a complex query that needs to calculate running balance
	// For simplicity, we'll calculate balance at end of each period
	var dateFormat string
	switch groupBy {
	case "day":
		dateFormat = "TO_CHAR(date, 'YYYY-MM-DD')"
	case "week":
		dateFormat = "TO_CHAR(DATE_TRUNC('week', date), 'YYYY-MM-DD')"
	case "month":
		dateFormat = "TO_CHAR(date, 'YYYY-MM')"
	default:
		dateFormat = "TO_CHAR(date, 'YYYY-MM')"
	}

	// Get account's opening balance (from account creation or start date)
	var openingBalance float64
	r.db.WithContext(ctx).Table("accounts").
		Where("id = ? AND user_id = ?", accountID, userID).
		Select("balance - COALESCE((SELECT SUM(CASE WHEN type = 'income' THEN amount WHEN type = 'expense' THEN -amount ELSE 0 END) FROM transactions WHERE account_id = ? AND deleted_at IS NULL), 0)", accountID).
		Scan(&openingBalance)

	query := `
		WITH balance_changes AS (
			SELECT
				` + dateFormat + ` as period,
				SUM(CASE WHEN type = 'income' THEN amount WHEN type = 'expense' THEN -amount ELSE 0 END) as change
			FROM transactions
			WHERE account_id = ? AND user_id = ? AND date >= ? AND date <= ? AND deleted_at IS NULL
			GROUP BY period
			ORDER BY period ASC
		)
		SELECT period, change as balance FROM balance_changes
	`

	var changes []BalanceHistory
	err := r.db.WithContext(ctx).Raw(query, accountID, userID, startDate, endDate).Scan(&changes).Error

	// Calculate running balance
	runningBalance := openingBalance
	for i := range changes {
		runningBalance += changes[i].Balance
		changes[i].Balance = runningBalance
		history = append(history, changes[i])
	}

	return history, err
}

// GetNetWorth calculates total net worth
func (r *reportRepository) GetNetWorth(ctx context.Context, userID uuid.UUID) (*NetWorth, error) {
	var netWorth NetWorth

	// Total assets (account balances + asset values + goal holdings)
	var accountBalance, assetValue, goalValue float64

	r.db.WithContext(ctx).Table("accounts").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Select("COALESCE(SUM(balance), 0)").
		Scan(&accountBalance)

	r.db.WithContext(ctx).Table("assets").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Select("COALESCE(SUM(purchase_price), 0)").
		Scan(&assetValue)

	r.db.WithContext(ctx).Table("goals").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Select("COALESCE(SUM(current_amount), 0)").
		Scan(&goalValue)

	netWorth.TotalAssets = accountBalance + assetValue + goalValue

	// Total liabilities (credit card balances + debts - lends)
	var creditCardBalance, debtBalance, lendBalance float64

	r.db.WithContext(ctx).Table("credit_cards").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Select("COALESCE(SUM(current_balance), 0)").
		Scan(&creditCardBalance)

	r.db.WithContext(ctx).Table("debt_records").
		Where("user_id = ? AND status IN ('active', 'partial') AND deleted_at IS NULL", userID).
		Select("COALESCE(SUM(remaining_amount), 0)").
		Scan(&debtBalance)

	r.db.WithContext(ctx).Table("lend_records").
		Where("user_id = ? AND status IN ('active', 'partial') AND deleted_at IS NULL", userID).
		Select("COALESCE(SUM(remaining_amount), 0)").
		Scan(&lendBalance)

	netWorth.TotalLiabilities = creditCardBalance + debtBalance - lendBalance
	netWorth.NetWorth = netWorth.TotalAssets - netWorth.TotalLiabilities

	return &netWorth, nil
}

// GetNetWorthTrend retrieves net worth trend over time
func (r *reportRepository) GetNetWorthTrend(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, groupBy string) ([]NetWorthTrend, error) {
	// This is a simplified version - in production you'd want to store snapshots
	// For now, we'll return current net worth for each period
	trends := []NetWorthTrend{}

	netWorth, err := r.GetNetWorth(ctx, userID)
	if err != nil {
		return trends, err
	}

	// Create trend entries (simplified - would need historical data for accurate trends)
	trend := NetWorthTrend{
		Period:           time.Now().Format("2006-01"),
		TotalAssets:      netWorth.TotalAssets,
		TotalLiabilities: netWorth.TotalLiabilities,
		NetWorth:         netWorth.NetWorth,
	}
	trends = append(trends, trend)

	return trends, nil
}

// GetBudgetPerformance retrieves budget vs actual spending
func (r *reportRepository) GetBudgetPerformance(ctx context.Context, userID uuid.UUID, month time.Time) ([]BudgetPerformance, error) {
	var performances []BudgetPerformance

	startDate := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	query := `
		SELECT
			b.category_id,
			c.name as category_name,
			b.amount as budget_amount,
			COALESCE(SUM(t.amount), 0) as actual_amount,
			b.amount - COALESCE(SUM(t.amount), 0) as difference,
			CASE WHEN b.amount > 0 THEN (COALESCE(SUM(t.amount), 0) / b.amount * 100) ELSE 0 END as percentage,
			CASE
				WHEN COALESCE(SUM(t.amount), 0) > b.amount THEN 'over'
				WHEN COALESCE(SUM(t.amount), 0) >= b.amount * 0.9 THEN 'on-track'
				ELSE 'under'
			END as status
		FROM budgets b
		LEFT JOIN categories c ON b.category_id = c.id
		LEFT JOIN transactions t ON t.category_id = b.category_id
			AND t.user_id = b.user_id
			AND t.type = 'expense'
			AND t.date >= ?
			AND t.date <= ?
			AND t.deleted_at IS NULL
		WHERE b.user_id = ? AND b.month = ? AND b.deleted_at IS NULL
		GROUP BY b.category_id, c.name, b.amount
		ORDER BY actual_amount DESC
	`

	err := r.db.WithContext(ctx).Raw(query, startDate, endDate, userID, month).Scan(&performances).Error
	return performances, err
}

// GetCashFlowSummary retrieves cash flow summary
func (r *reportRepository) GetCashFlowSummary(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*CashFlowSummary, error) {
	var summary CashFlowSummary

	// Get opening balance (sum of all accounts at start date)
	r.db.WithContext(ctx).Table("accounts").
		Where("user_id = ? AND created_at < ? AND deleted_at IS NULL", userID, startDate).
		Select("COALESCE(SUM(balance), 0)").
		Scan(&summary.OpeningBalance)

	// Get total inflow (income)
	r.db.WithContext(ctx).Table("transactions").
		Where("user_id = ? AND type = 'income' AND date >= ? AND date <= ? AND deleted_at IS NULL", userID, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&summary.TotalInflow)

	// Get total outflow (expense)
	r.db.WithContext(ctx).Table("transactions").
		Where("user_id = ? AND type = 'expense' AND date >= ? AND date <= ? AND deleted_at IS NULL", userID, startDate, endDate).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&summary.TotalOutflow)

	summary.NetCashFlow = summary.TotalInflow - summary.TotalOutflow
	summary.ClosingBalance = summary.OpeningBalance + summary.NetCashFlow

	return &summary, nil
}

// GetCashFlowByMonth retrieves monthly cash flow for a year
func (r *reportRepository) GetCashFlowByMonth(ctx context.Context, userID uuid.UUID, year int) ([]MonthlyCashFlow, error) {
	var cashFlows []MonthlyCashFlow

	query := `
		SELECT
			TO_CHAR(date, 'YYYY-MM') as month,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as inflow,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as outflow,
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount WHEN type = 'expense' THEN -amount ELSE 0 END), 0) as net_cash_flow
		FROM transactions
		WHERE user_id = ? AND EXTRACT(YEAR FROM date) = ? AND deleted_at IS NULL
		GROUP BY month
		ORDER BY month ASC
	`

	err := r.db.WithContext(ctx).Raw(query, userID, year).Scan(&cashFlows).Error
	return cashFlows, err
}

// GetMonthlySummary retrieves summary for a specific month
func (r *reportRepository) GetMonthlySummary(ctx context.Context, userID uuid.UUID, month time.Time) (*PeriodSummary, error) {
	startDate := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	return r.getPeriodSummary(ctx, userID, startDate, endDate, month.Format("2006-01"))
}

// GetYearlySummary retrieves summary for a specific year
func (r *reportRepository) GetYearlySummary(ctx context.Context, userID uuid.UUID, year int) (*PeriodSummary, error) {
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	return r.getPeriodSummary(ctx, userID, startDate, endDate, string(rune(year)))
}

// getPeriodSummary is a helper function to get period summary
func (r *reportRepository) getPeriodSummary(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, period string) (*PeriodSummary, error) {
	var summary PeriodSummary
	summary.Period = period

	query := `
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as total_expense,
			COUNT(*) as transaction_count
		FROM transactions
		WHERE user_id = ? AND date >= ? AND date <= ? AND deleted_at IS NULL AND type != 'tracking'
	`

	r.db.WithContext(ctx).Raw(query, userID, startDate, endDate).Scan(&summary)

	summary.NetAmount = summary.TotalIncome - summary.TotalExpense

	// Calculate average daily spending
	days := int(endDate.Sub(startDate).Hours()/24) + 1
	if days > 0 {
		summary.AvgDailySpending = summary.TotalExpense / float64(days)
	}

	return &summary, nil
}

// GetPeriodComparison compares two time periods
func (r *reportRepository) GetPeriodComparison(ctx context.Context, userID uuid.UUID, period1Start, period1End, period2Start, period2End time.Time) (*PeriodComparison, error) {
	comparison := &PeriodComparison{
		Changes: &PeriodChanges{},
	}

	// Get summaries for both periods
	period1, err := r.getPeriodSummary(ctx, userID, period1Start, period1End, "Period 1")
	if err != nil {
		return nil, err
	}
	comparison.Period1 = period1

	period2, err := r.getPeriodSummary(ctx, userID, period2Start, period2End, "Period 2")
	if err != nil {
		return nil, err
	}
	comparison.Period2 = period2

	// Calculate changes
	comparison.Changes.IncomeChange = period2.TotalIncome - period1.TotalIncome
	if period1.TotalIncome > 0 {
		comparison.Changes.IncomeChangePerc = (comparison.Changes.IncomeChange / period1.TotalIncome) * 100
	}

	comparison.Changes.ExpenseChange = period2.TotalExpense - period1.TotalExpense
	if period1.TotalExpense > 0 {
		comparison.Changes.ExpenseChangePerc = (comparison.Changes.ExpenseChange / period1.TotalExpense) * 100
	}

	comparison.Changes.NetChange = period2.NetAmount - period1.NetAmount
	if period1.NetAmount != 0 {
		comparison.Changes.NetChangePerc = (comparison.Changes.NetChange / period1.NetAmount) * 100
	}

	return comparison, nil
}

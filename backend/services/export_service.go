package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"time"

	"daybook-backend/models"
	"daybook-backend/repository"
)

// ExportService handles data export operations
type ExportService interface {
	ExportTransactionsCSV(ctx context.Context, userID uint, startDate, endDate time.Time, transactionType string) ([]byte, error)
	ExportTransactionsJSON(ctx context.Context, userID uint, startDate, endDate time.Time, transactionType string) ([]byte, error)
	ExportAccountsCSV(ctx context.Context, userID uint) ([]byte, error)
	ExportAccountsJSON(ctx context.Context, userID uint) ([]byte, error)
	ExportBudgetsCSV(ctx context.Context, userID uint) ([]byte, error)
	ExportBudgetsJSON(ctx context.Context, userID uint) ([]byte, error)
	ExportGoalsCSV(ctx context.Context, userID uint) ([]byte, error)
	ExportGoalsJSON(ctx context.Context, userID uint) ([]byte, error)
	ExportCategoriesCSV(ctx context.Context, userID uint) ([]byte, error)
	ExportCategoriesJSON(ctx context.Context, userID uint) ([]byte, error)
	ExportAssetsCSV(ctx context.Context, userID uint) ([]byte, error)
	ExportAssetsJSON(ctx context.Context, userID uint) ([]byte, error)
	ExportAllDataJSON(ctx context.Context, userID uint) ([]byte, error)
}

type exportService struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	budgetRepo      repository.BudgetRepository
	goalRepo        repository.GoalRepository
	categoryRepo    repository.CategoryRepository
	assetRepo       repository.AssetRepository
	creditCardRepo  repository.CreditCardRepository
	debtRepo        repository.DebtRepository
	lendRepo        repository.LendRepository
	activityLogger  ActivityLogService
}

// NewExportService creates a new export service
func NewExportService(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	budgetRepo repository.BudgetRepository,
	goalRepo repository.GoalRepository,
	categoryRepo repository.CategoryRepository,
	assetRepo repository.AssetRepository,
	creditCardRepo repository.CreditCardRepository,
	debtRepo repository.DebtRepository,
	lendRepo repository.LendRepository,
	activityLogger ActivityLogService,
) ExportService {
	return &exportService{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		budgetRepo:      budgetRepo,
		goalRepo:        goalRepo,
		categoryRepo:    categoryRepo,
		assetRepo:       assetRepo,
		creditCardRepo:  creditCardRepo,
		debtRepo:        debtRepo,
		lendRepo:        lendRepo,
		activityLogger:  activityLogger,
	}
}

// ExportTransactionsCSV exports transactions as CSV
func (s *exportService) ExportTransactionsCSV(ctx context.Context, userID uint, startDate, endDate time.Time, transactionType string) ([]byte, error) {
	// Use FindWithFilters to get transactions by date range
	filters := repository.TransactionFilters{
		StartDate:       &startDate,
		EndDate:         &endDate,
		IncludeTracking: true, // Include all transactions
	}
	if transactionType != "" {
		filters.Type = &transactionType
	}
	pagination := repository.PaginationParams{
		Page:  1,
		Limit: 100000, // Large number to get all transactions
	}

	result, err := s.transactionRepo.FindWithFilters(ctx, userID, filters, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	transactions := result.Transactions

	// Fetch categories and accounts for name lookups
	categories, err := s.categoryRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}
	categoryMap := make(map[uint]string)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat.Name
	}

	accounts, err := s.accountRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %w", err)
	}
	accountMap := make(map[uint]string)
	for _, acc := range accounts {
		accountMap[acc.ID] = acc.Name
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{"ID", "Date", "Type", "Amount", "Category", "Account", "To Account", "Description", "Tags", "Created At"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	for _, tx := range transactions {
		tags := ""
		if len(tx.Tags) > 0 {
			tagsBytes, _ := json.Marshal(tx.Tags)
			tags = string(tagsBytes)
		}

		categoryName := categoryMap[uint(tx.CategoryID)]
		accountName := accountMap[uint(tx.AccountID)]
		toAccountName := ""
		if tx.ToAccountID != nil {
			toAccountName = accountMap[uint(*tx.ToAccountID)]
		}

		record := []string{
			fmt.Sprintf("%d", tx.ID),
			tx.Date.Format("2006-01-02"),
			tx.Type,
			fmt.Sprintf("%.2f", tx.Amount),
			categoryName,
			accountName,
			toAccountName,
			tx.Description,
			tags,
			tx.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	// Log activity
	s.activityLogger.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      "export",
		Module:      "transactions",
		EntityType:  "csv",
		Description: fmt.Sprintf("Exported %d transactions to CSV", len(transactions)),
	})

	return buf.Bytes(), nil
}

// TransactionExport represents a transaction with category and account names for export
type TransactionExport struct {
	ID            uint        `json:"id"`
	Date          string      `json:"date"`
	Type          string      `json:"type"`
	Amount        float64     `json:"amount"`
	CategoryID    uint        `json:"categoryId"`
	CategoryName  string      `json:"categoryName"`
	AccountID     uint        `json:"accountId"`
	AccountName   string      `json:"accountName"`
	ToAccountID   *uint       `json:"toAccountId,omitempty"`
	ToAccountName string      `json:"toAccountName,omitempty"`
	Description   string      `json:"description"`
	Tags          []string    `json:"tags"`
	CreatedAt     string      `json:"createdAt"`
}

// ExportTransactionsJSON exports transactions as JSON
func (s *exportService) ExportTransactionsJSON(ctx context.Context, userID uint, startDate, endDate time.Time, transactionType string) ([]byte, error) {
	// Use FindWithFilters to get transactions by date range
	filters := repository.TransactionFilters{
		StartDate:       &startDate,
		EndDate:         &endDate,
		IncludeTracking: true, // Include all transactions
	}
	if transactionType != "" {
		filters.Type = &transactionType
	}
	pagination := repository.PaginationParams{
		Page:  1,
		Limit: 100000, // Large number to get all transactions
	}

	result, err := s.transactionRepo.FindWithFilters(ctx, userID, filters, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	transactions := result.Transactions

	// Fetch categories and accounts for name lookups
	categories, err := s.categoryRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}
	categoryMap := make(map[uint]string)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat.Name
	}

	accounts, err := s.accountRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %w", err)
	}
	accountMap := make(map[uint]string)
	for _, acc := range accounts {
		accountMap[acc.ID] = acc.Name
	}

	// Build enriched transaction export data
	exportData := make([]TransactionExport, len(transactions))
	for i, tx := range transactions {
		toAccountName := ""
		if tx.ToAccountID != nil {
			toAccountName = accountMap[uint(*tx.ToAccountID)]
		}

		exportData[i] = TransactionExport{
			ID:            tx.ID,
			Date:          tx.Date.Format("2006-01-02"),
			Type:          tx.Type,
			Amount:        tx.Amount,
			CategoryID:    uint(tx.CategoryID),
			CategoryName:  categoryMap[uint(tx.CategoryID)],
			AccountID:     uint(tx.AccountID),
			AccountName:   accountMap[uint(tx.AccountID)],
			ToAccountID:   tx.ToAccountID,
			ToAccountName: toAccountName,
			Description:   tx.Description,
			Tags:          tx.Tags,
			CreatedAt:     tx.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Log activity
	s.activityLogger.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      "export",
		Module:      "transactions",
		EntityType:  "json",
		Description: fmt.Sprintf("Exported %d transactions to JSON", len(transactions)),
	})

	return data, nil
}

// ExportAccountsCSV exports accounts as CSV
func (s *exportService) ExportAccountsCSV(ctx context.Context, userID uint) ([]byte, error) {
	accounts, err := s.accountRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{"ID", "Name", "Type", "Balance", "Currency", "Institution", "Account Number", "Description", "Created At"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	for _, acc := range accounts {
		record := []string{
			fmt.Sprintf("%d", acc.ID),
			acc.Name,
			acc.Type,
			fmt.Sprintf("%.2f", acc.Balance),
			acc.Currency,
			acc.Institution,
			acc.AccountNumber,
			acc.Description,
			acc.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	// Log activity
	s.activityLogger.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      "export",
		Module:      "accounts",
		EntityType:  "csv",
		Description: fmt.Sprintf("Exported %d accounts to CSV", len(accounts)),
	})

	return buf.Bytes(), nil
}

// ExportAccountsJSON exports accounts as JSON
func (s *exportService) ExportAccountsJSON(ctx context.Context, userID uint) ([]byte, error) {
	accounts, err := s.accountRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %w", err)
	}

	data, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Log activity
	s.activityLogger.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      "export",
		Module:      "accounts",
		EntityType:  "json",
		Description: fmt.Sprintf("Exported %d accounts to JSON", len(accounts)),
	})

	return data, nil
}

// ExportBudgetsCSV exports budgets as CSV
func (s *exportService) ExportBudgetsCSV(ctx context.Context, userID uint) ([]byte, error) {
	budgets, err := s.budgetRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch budgets: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{"ID", "Category ID", "Amount", "Period", "Rollover", "Alert Threshold", "Enabled", "Notes", "Created At"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	for _, budget := range budgets {
		record := []string{
			fmt.Sprintf("%d", budget.ID),
			fmt.Sprintf("%d", budget.CategoryID),
			fmt.Sprintf("%.2f", budget.Amount),
			budget.Period,
			fmt.Sprintf("%t", budget.Rollover),
			fmt.Sprintf("%.2f", budget.AlertThreshold),
			fmt.Sprintf("%t", budget.Enabled),
			budget.Notes,
			budget.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	// Log activity
	s.activityLogger.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      "export",
		Module:      "budgets",
		EntityType:  "csv",
		Description: fmt.Sprintf("Exported %d budgets to CSV", len(budgets)),
	})

	return buf.Bytes(), nil
}

// ExportBudgetsJSON exports budgets as JSON
func (s *exportService) ExportBudgetsJSON(ctx context.Context, userID uint) ([]byte, error) {
	budgets, err := s.budgetRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch budgets: %w", err)
	}

	data, err := json.MarshalIndent(budgets, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Log activity
	s.activityLogger.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      "export",
		Module:      "budgets",
		EntityType:  "json",
		Description: fmt.Sprintf("Exported %d budgets to JSON", len(budgets)),
	})

	return data, nil
}

// ExportGoalsCSV exports goals as CSV
func (s *exportService) ExportGoalsCSV(ctx context.Context, userID uint) ([]byte, error) {
	goals, err := s.goalRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch goals: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{"ID", "Name", "Target Amount", "Current Amount", "Target Date", "Status", "Created At"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	for _, goal := range goals {
		record := []string{
			fmt.Sprintf("%d", goal.ID),
			goal.Name,
			fmt.Sprintf("%.2f", goal.TargetAmount),
			fmt.Sprintf("%.2f", goal.CurrentAmount),
			goal.TargetDate.Format("2006-01-02"),
			goal.Status,
			goal.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	// Log activity
	s.activityLogger.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      "export",
		Module:      "goals",
		EntityType:  "csv",
		Description: fmt.Sprintf("Exported %d goals to CSV", len(goals)),
	})

	return buf.Bytes(), nil
}

// ExportGoalsJSON exports goals as JSON
func (s *exportService) ExportGoalsJSON(ctx context.Context, userID uint) ([]byte, error) {
	goals, err := s.goalRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch goals: %w", err)
	}

	data, err := json.MarshalIndent(goals, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Log activity
	s.activityLogger.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      "export",
		Module:      "goals",
		EntityType:  "json",
		Description: fmt.Sprintf("Exported %d goals to JSON", len(goals)),
	})

	return data, nil
}

// ExportCategoriesCSV exports categories as CSV
func (s *exportService) ExportCategoriesCSV(ctx context.Context, userID uint) ([]byte, error) {
	categories, err := s.categoryRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{"ID", "Name", "Type", "Icon", "Color", "Description", "Is Default", "Order"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	for _, cat := range categories {
		record := []string{
			fmt.Sprintf("%d", cat.ID),
			cat.Name,
			cat.Type,
			cat.Icon,
			cat.Color,
			cat.Description,
			fmt.Sprintf("%t", cat.IsDefault),
			fmt.Sprintf("%d", cat.Order),
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	// Log activity
	s.activityLogger.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      "export",
		Module:      "categories",
		EntityType:  "csv",
		Description: fmt.Sprintf("Exported %d categories to CSV", len(categories)),
	})

	return buf.Bytes(), nil
}

// ExportCategoriesJSON exports categories as JSON
func (s *exportService) ExportCategoriesJSON(ctx context.Context, userID uint) ([]byte, error) {
	categories, err := s.categoryRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	data, err := json.MarshalIndent(categories, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Log activity
	s.activityLogger.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      "export",
		Module:      "categories",
		EntityType:  "json",
		Description: fmt.Sprintf("Exported %d categories to JSON", len(categories)),
	})

	return data, nil
}

// ExportAssetsCSV exports assets as CSV
func (s *exportService) ExportAssetsCSV(ctx context.Context, userID uint) ([]byte, error) {
	assets, err := s.assetRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch assets: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{"ID", "Name", "Category", "Brand", "Model", "Serial Number", "Purchase Date", "Purchase Price", "Purchase Location", "Warranty Start", "Warranty End", "Warranty Provider", "Warranty Type", "Status", "Notes", "Created At"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	for _, asset := range assets {
		warrantyStart := ""
		if asset.WarrantyStartDate != nil {
			warrantyStart = asset.WarrantyStartDate.Time.Format("2006-01-02")
		}

		warrantyEnd := ""
		if asset.WarrantyEndDate != nil {
			warrantyEnd = asset.WarrantyEndDate.Time.Format("2006-01-02")
		}

		record := []string{
			fmt.Sprintf("%d", asset.ID),
			asset.Name,
			asset.Category,
			asset.Brand,
			asset.Model,
			asset.SerialNumber,
			asset.PurchaseDate.Time.Format("2006-01-02"),
			fmt.Sprintf("%.2f", asset.PurchasePrice),
			asset.PurchaseLocation,
			warrantyStart,
			warrantyEnd,
			asset.WarrantyProvider,
			asset.WarrantyType,
			asset.Status,
			asset.Notes,
			asset.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	// Log activity
	s.activityLogger.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      "export",
		Module:      "assets",
		EntityType:  "csv",
		Description: fmt.Sprintf("Exported %d assets to CSV", len(assets)),
	})

	return buf.Bytes(), nil
}

// ExportAssetsJSON exports assets as JSON
func (s *exportService) ExportAssetsJSON(ctx context.Context, userID uint) ([]byte, error) {
	assets, err := s.assetRepo.FindAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch assets: %w", err)
	}

	data, err := json.MarshalIndent(assets, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Log activity
	s.activityLogger.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      "export",
		Module:      "assets",
		EntityType:  "json",
		Description: fmt.Sprintf("Exported %d assets to JSON", len(assets)),
	})

	return data, nil
}

// ExportAllDataJSON exports all user data as a comprehensive JSON file
func (s *exportService) ExportAllDataJSON(ctx context.Context, userID uint) ([]byte, error) {
	// Fetch all data
	accounts, _ := s.accountRepo.FindAll(ctx, userID)
	categories, _ := s.categoryRepo.FindAll(ctx, userID)
	budgets, _ := s.budgetRepo.FindAll(ctx, userID)
	goals, _ := s.goalRepo.FindAll(ctx, userID)
	assets, _ := s.assetRepo.FindAll(ctx, userID)
	creditCards, _ := s.creditCardRepo.FindAll(ctx, userID)
	debts, _ := s.debtRepo.FindAll(ctx, userID)
	lends, _ := s.lendRepo.FindAll(ctx, userID)

	// Fetch transactions for the last 2 years
	endDate := time.Now()
	startDate := endDate.AddDate(-2, 0, 0)
	filters := repository.TransactionFilters{
		StartDate:       &startDate,
		EndDate:         &endDate,
		IncludeTracking: true,
	}
	pagination := repository.PaginationParams{
		Page:  1,
		Limit: 100000, // Large number to get all transactions
	}
	result, _ := s.transactionRepo.FindWithFilters(ctx, userID, filters, pagination)
	transactions := []models.Transaction{}
	if result != nil {
		transactions = result.Transactions
	}

	// Build comprehensive export structure
	exportData := map[string]interface{}{
		"export_date":    time.Now().Format(time.RFC3339),
		"export_version": "1.0",
		"user_id":        userID,
		"data": map[string]interface{}{
			"accounts":     accounts,
			"categories":   categories,
			"transactions": transactions,
			"budgets":      budgets,
			"goals":        goals,
			"assets":       assets,
			"credit_cards": creditCards,
			"debts":        debts,
			"lends":        lends,
		},
		"metadata": map[string]interface{}{
			"accounts_count":     len(accounts),
			"categories_count":   len(categories),
			"transactions_count": len(transactions),
			"budgets_count":      len(budgets),
			"goals_count":        len(goals),
			"assets_count":       len(assets),
			"credit_cards_count": len(creditCards),
			"debts_count":        len(debts),
			"lends_count":        len(lends),
			"date_range": map[string]string{
				"start": startDate.Format("2006-01-02"),
				"end":   endDate.Format("2006-01-02"),
			},
		},
	}

	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Log activity
	totalRecords := len(accounts) + len(categories) + len(transactions) + len(budgets) + len(goals) + len(assets) + len(creditCards) + len(debts) + len(lends)
	s.activityLogger.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      "export",
		Module:      "all_data",
		EntityType:  "json",
		Description: fmt.Sprintf("Exported all data (%d total records) to JSON", totalRecords),
	})

	return data, nil
}

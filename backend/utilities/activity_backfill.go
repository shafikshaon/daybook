package utilities

import (
	"encoding/json"
	"fmt"
	"time"

	"daybook-backend/database"
	"daybook-backend/models"

	"github.com/google/uuid"
)

// BackfillOptions contains configuration for backfilling activity logs
type BackfillOptions struct {
	UserID    *uuid.UUID // If nil, backfill for all users
	Module    string     // If empty, backfill for all modules
	DryRun    bool       // If true, only count records without creating logs
	StartDate *time.Time // Only backfill records created after this date
	EndDate   *time.Time // Only backfill records created before this date
	BatchSize int        // Number of records to process in each batch
}

// BackfillResult contains statistics about the backfill operation
type BackfillResult struct {
	Module       string `json:"module"`
	TotalRecords int64  `json:"total_records"`
	LogsCreated  int64  `json:"logs_created"`
	LogsSkipped  int64  `json:"logs_skipped"`
	Errors       int64  `json:"errors"`
}

// BackfillAllActivities backfills activity logs for all or specified modules
func BackfillAllActivities(options BackfillOptions) ([]BackfillResult, error) {
	results := []BackfillResult{}

	if options.BatchSize == 0 {
		options.BatchSize = 100 // Default batch size
	}

	// Determine which modules to backfill
	modules := []string{}
	if options.Module != "" {
		modules = []string{options.Module}
	} else {
		modules = []string{
			models.ModuleAccount,
			models.ModuleTransaction,
			models.ModuleBudget,
			models.ModuleCreditCard,
			models.ModuleDebt,
			models.ModuleLend,
			models.ModuleAsset,
			models.ModuleGoal,
		}
	}

	// Backfill each module
	for _, module := range modules {
		result, err := backfillModule(module, options)
		if err != nil {
			return results, fmt.Errorf("failed to backfill module %s: %w", module, err)
		}
		results = append(results, result)
	}

	return results, nil
}

// backfillModule backfills activity logs for a specific module
func backfillModule(module string, options BackfillOptions) (BackfillResult, error) {
	result := BackfillResult{Module: module}

	switch module {
	case models.ModuleAccount:
		return backfillAccounts(options)
	case models.ModuleTransaction:
		return backfillTransactions(options)
	case models.ModuleBudget:
		return backfillBudgets(options)
	case models.ModuleCreditCard:
		return backfillCreditCards(options)
	case models.ModuleDebt:
		return backfillDebts(options)
	case models.ModuleLend:
		return backfillLends(options)
	case models.ModuleAsset:
		return backfillAssets(options)
	case models.ModuleGoal:
		return backfillGoals(options)
	default:
		return result, fmt.Errorf("unsupported module: %s", module)
	}
}

// backfillAccounts creates activity logs for existing accounts
func backfillAccounts(options BackfillOptions) (BackfillResult, error) {
	result := BackfillResult{Module: models.ModuleAccount}

	query := database.DB.Model(&models.Account{})

	// Apply filters
	if options.UserID != nil {
		query = query.Where("user_id = ?", options.UserID)
	}
	if options.StartDate != nil {
		query = query.Where("created_at >= ?", options.StartDate)
	}
	if options.EndDate != nil {
		query = query.Where("created_at <= ?", options.EndDate)
	}

	// Count total records
	query.Count(&result.TotalRecords)

	if options.DryRun {
		result.LogsCreated = result.TotalRecords
		return result, nil
	}

	// Process in batches
	offset := 0
	for {
		var accounts []models.Account
		err := query.Limit(options.BatchSize).Offset(offset).Find(&accounts).Error
		if err != nil {
			return result, err
		}

		if len(accounts) == 0 {
			break
		}

		for _, account := range accounts {
			// Check if activity log already exists for this account creation
			if activityLogExists(account.UserID, models.ActionCreate, models.ModuleAccount, account.ID, account.CreatedAt) {
				result.LogsSkipped++
				continue
			}

			// Create backfilled activity log
			metadata := map[string]interface{}{
				"backfilled": true,
				"name":       account.Name,
				"type":       account.Type,
			}
			metadataJSON, _ := json.Marshal(metadata)

			activityLog := models.ActivityLog{
				ID:          uuid.New(),
				UserID:      account.UserID,
				Action:      models.ActionCreate,
				Module:      models.ModuleAccount,
				EntityType:  "Account",
				EntityID:    &account.ID,
				Description: fmt.Sprintf("Created account: %s (backfilled)", account.Name),
				Metadata:    string(metadataJSON),
				CreatedAt:   account.CreatedAt,
				UpdatedAt:   account.CreatedAt,
			}

			if err := database.DB.Create(&activityLog).Error; err != nil {
				result.Errors++
			} else {
				result.LogsCreated++
			}
		}

		offset += options.BatchSize
	}

	return result, nil
}

// backfillTransactions creates activity logs for existing transactions
func backfillTransactions(options BackfillOptions) (BackfillResult, error) {
	result := BackfillResult{Module: models.ModuleTransaction}

	query := database.DB.Model(&models.Transaction{})

	if options.UserID != nil {
		query = query.Where("user_id = ?", options.UserID)
	}
	if options.StartDate != nil {
		query = query.Where("created_at >= ?", options.StartDate)
	}
	if options.EndDate != nil {
		query = query.Where("created_at <= ?", options.EndDate)
	}

	query.Count(&result.TotalRecords)

	if options.DryRun {
		result.LogsCreated = result.TotalRecords
		return result, nil
	}

	offset := 0
	for {
		var transactions []models.Transaction
		err := query.Limit(options.BatchSize).Offset(offset).Find(&transactions).Error
		if err != nil {
			return result, err
		}

		if len(transactions) == 0 {
			break
		}

		for _, transaction := range transactions {
			if activityLogExists(transaction.UserID, models.ActionCreate, models.ModuleTransaction, transaction.ID, transaction.CreatedAt) {
				result.LogsSkipped++
				continue
			}

			metadata := map[string]interface{}{
				"backfilled": true,
				"amount":     transaction.Amount,
				"type":       transaction.Type,
			}
			metadataJSON, _ := json.Marshal(metadata)

			activityLog := models.ActivityLog{
				ID:          uuid.New(),
				UserID:      transaction.UserID,
				Action:      models.ActionCreate,
				Module:      models.ModuleTransaction,
				EntityType:  "Transaction",
				EntityID:    &transaction.ID,
				Description: fmt.Sprintf("Created transaction: %s (backfilled)", transaction.Description),
				Metadata:    string(metadataJSON),
				CreatedAt:   transaction.CreatedAt,
				UpdatedAt:   transaction.CreatedAt,
			}

			if err := database.DB.Create(&activityLog).Error; err != nil {
				result.Errors++
			} else {
				result.LogsCreated++
			}
		}

		offset += options.BatchSize
	}

	return result, nil
}

// backfillBudgets creates activity logs for existing budgets
func backfillBudgets(options BackfillOptions) (BackfillResult, error) {
	result := BackfillResult{Module: models.ModuleBudget}

	query := database.DB.Model(&models.Budget{})

	if options.UserID != nil {
		query = query.Where("user_id = ?", options.UserID)
	}
	if options.StartDate != nil {
		query = query.Where("created_at >= ?", options.StartDate)
	}
	if options.EndDate != nil {
		query = query.Where("created_at <= ?", options.EndDate)
	}

	query.Count(&result.TotalRecords)

	if options.DryRun {
		result.LogsCreated = result.TotalRecords
		return result, nil
	}

	offset := 0
	for {
		var budgets []models.Budget
		err := query.Limit(options.BatchSize).Offset(offset).Find(&budgets).Error
		if err != nil {
			return result, err
		}

		if len(budgets) == 0 {
			break
		}

		for _, budget := range budgets {
			if activityLogExists(budget.UserID, models.ActionCreate, models.ModuleBudget, budget.ID, budget.CreatedAt) {
				result.LogsSkipped++
				continue
			}

			metadata := map[string]interface{}{
				"backfilled": true,
				"amount":     budget.Amount,
				"period":     budget.Period,
				"categoryId": budget.CategoryID,
			}
			metadataJSON, _ := json.Marshal(metadata)

			activityLog := models.ActivityLog{
				ID:          uuid.New(),
				UserID:      budget.UserID,
				Action:      models.ActionCreate,
				Module:      models.ModuleBudget,
				EntityType:  "Budget",
				EntityID:    &budget.ID,
				Description: fmt.Sprintf("Created budget for category %s (backfilled)", budget.CategoryID),
				Metadata:    string(metadataJSON),
				CreatedAt:   budget.CreatedAt,
				UpdatedAt:   budget.CreatedAt,
			}

			if err := database.DB.Create(&activityLog).Error; err != nil {
				result.Errors++
			} else {
				result.LogsCreated++
			}
		}

		offset += options.BatchSize
	}

	return result, nil
}

// backfillCreditCards creates activity logs for existing credit cards
func backfillCreditCards(options BackfillOptions) (BackfillResult, error) {
	result := BackfillResult{Module: models.ModuleCreditCard}

	query := database.DB.Model(&models.CreditCard{})

	if options.UserID != nil {
		query = query.Where("user_id = ?", options.UserID)
	}
	if options.StartDate != nil {
		query = query.Where("created_at >= ?", options.StartDate)
	}
	if options.EndDate != nil {
		query = query.Where("created_at <= ?", options.EndDate)
	}

	query.Count(&result.TotalRecords)

	if options.DryRun {
		result.LogsCreated = result.TotalRecords
		return result, nil
	}

	offset := 0
	for {
		var creditCards []models.CreditCard
		err := query.Limit(options.BatchSize).Offset(offset).Find(&creditCards).Error
		if err != nil {
			return result, err
		}

		if len(creditCards) == 0 {
			break
		}

		for _, card := range creditCards {
			if activityLogExists(card.UserID, models.ActionCreate, models.ModuleCreditCard, card.ID, card.CreatedAt) {
				result.LogsSkipped++
				continue
			}

			metadata := map[string]interface{}{
				"backfilled": true,
				"name":       card.Name,
			}
			metadataJSON, _ := json.Marshal(metadata)

			activityLog := models.ActivityLog{
				ID:          uuid.New(),
				UserID:      card.UserID,
				Action:      models.ActionCreate,
				Module:      models.ModuleCreditCard,
				EntityType:  "CreditCard",
				EntityID:    &card.ID,
				Description: fmt.Sprintf("Created credit card: %s (backfilled)", card.Name),
				Metadata:    string(metadataJSON),
				CreatedAt:   card.CreatedAt,
				UpdatedAt:   card.CreatedAt,
			}

			if err := database.DB.Create(&activityLog).Error; err != nil {
				result.Errors++
			} else {
				result.LogsCreated++
			}
		}

		offset += options.BatchSize
	}

	return result, nil
}

// backfillDebts creates activity logs for existing debt records
func backfillDebts(options BackfillOptions) (BackfillResult, error) {
	result := BackfillResult{Module: models.ModuleDebt}

	query := database.DB.Model(&models.DebtRecord{})

	if options.UserID != nil {
		query = query.Where("user_id = ?", options.UserID)
	}
	if options.StartDate != nil {
		query = query.Where("created_at >= ?", options.StartDate)
	}
	if options.EndDate != nil {
		query = query.Where("created_at <= ?", options.EndDate)
	}

	query.Count(&result.TotalRecords)

	if options.DryRun {
		result.LogsCreated = result.TotalRecords
		return result, nil
	}

	offset := 0
	for {
		var debts []models.DebtRecord
		err := query.Limit(options.BatchSize).Offset(offset).Find(&debts).Error
		if err != nil {
			return result, err
		}

		if len(debts) == 0 {
			break
		}

		for _, debt := range debts {
			if activityLogExists(debt.UserID, models.ActionCreate, models.ModuleDebt, debt.ID, debt.CreatedAt) {
				result.LogsSkipped++
				continue
			}

			metadata := map[string]interface{}{
				"backfilled": true,
				"amount":     debt.OriginalAmount,
			}
			metadataJSON, _ := json.Marshal(metadata)

			activityLog := models.ActivityLog{
				ID:          uuid.New(),
				UserID:      debt.UserID,
				Action:      models.ActionCreate,
				Module:      models.ModuleDebt,
				EntityType:  "DebtRecord",
				EntityID:    &debt.ID,
				Description: fmt.Sprintf("Created debt record: %s (backfilled)", debt.CreditorName),
				Metadata:    string(metadataJSON),
				CreatedAt:   debt.CreatedAt,
				UpdatedAt:   debt.CreatedAt,
			}

			if err := database.DB.Create(&activityLog).Error; err != nil {
				result.Errors++
			} else {
				result.LogsCreated++
			}
		}

		offset += options.BatchSize
	}

	return result, nil
}

// backfillLends creates activity logs for existing lend records
func backfillLends(options BackfillOptions) (BackfillResult, error) {
	result := BackfillResult{Module: models.ModuleLend}

	query := database.DB.Model(&models.LendRecord{})

	if options.UserID != nil {
		query = query.Where("user_id = ?", options.UserID)
	}
	if options.StartDate != nil {
		query = query.Where("created_at >= ?", options.StartDate)
	}
	if options.EndDate != nil {
		query = query.Where("created_at <= ?", options.EndDate)
	}

	query.Count(&result.TotalRecords)

	if options.DryRun {
		result.LogsCreated = result.TotalRecords
		return result, nil
	}

	offset := 0
	for {
		var lends []models.LendRecord
		err := query.Limit(options.BatchSize).Offset(offset).Find(&lends).Error
		if err != nil {
			return result, err
		}

		if len(lends) == 0 {
			break
		}

		for _, lend := range lends {
			if activityLogExists(lend.UserID, models.ActionCreate, models.ModuleLend, lend.ID, lend.CreatedAt) {
				result.LogsSkipped++
				continue
			}

			metadata := map[string]interface{}{
				"backfilled": true,
				"amount":     lend.OriginalAmount,
			}
			metadataJSON, _ := json.Marshal(metadata)

			activityLog := models.ActivityLog{
				ID:          uuid.New(),
				UserID:      lend.UserID,
				Action:      models.ActionCreate,
				Module:      models.ModuleLend,
				EntityType:  "LendRecord",
				EntityID:    &lend.ID,
				Description: fmt.Sprintf("Created lend record: %s (backfilled)", lend.DebtorName),
				Metadata:    string(metadataJSON),
				CreatedAt:   lend.CreatedAt,
				UpdatedAt:   lend.CreatedAt,
			}

			if err := database.DB.Create(&activityLog).Error; err != nil {
				result.Errors++
			} else {
				result.LogsCreated++
			}
		}

		offset += options.BatchSize
	}

	return result, nil
}

// backfillAssets creates activity logs for existing assets
func backfillAssets(options BackfillOptions) (BackfillResult, error) {
	result := BackfillResult{Module: models.ModuleAsset}

	query := database.DB.Model(&models.Asset{})

	if options.UserID != nil {
		query = query.Where("user_id = ?", options.UserID)
	}
	if options.StartDate != nil {
		query = query.Where("created_at >= ?", options.StartDate)
	}
	if options.EndDate != nil {
		query = query.Where("created_at <= ?", options.EndDate)
	}

	query.Count(&result.TotalRecords)

	if options.DryRun {
		result.LogsCreated = result.TotalRecords
		return result, nil
	}

	offset := 0
	for {
		var assets []models.Asset
		err := query.Limit(options.BatchSize).Offset(offset).Find(&assets).Error
		if err != nil {
			return result, err
		}

		if len(assets) == 0 {
			break
		}

		for _, asset := range assets {
			if activityLogExists(asset.UserID, models.ActionCreate, models.ModuleAsset, asset.ID, asset.CreatedAt) {
				result.LogsSkipped++
				continue
			}

			metadata := map[string]interface{}{
				"backfilled":    true,
				"name":          asset.Name,
				"purchasePrice": asset.PurchasePrice,
			}
			metadataJSON, _ := json.Marshal(metadata)

			activityLog := models.ActivityLog{
				ID:          uuid.New(),
				UserID:      asset.UserID,
				Action:      models.ActionCreate,
				Module:      models.ModuleAsset,
				EntityType:  "Asset",
				EntityID:    &asset.ID,
				Description: fmt.Sprintf("Created asset: %s (backfilled)", asset.Name),
				Metadata:    string(metadataJSON),
				CreatedAt:   asset.CreatedAt,
				UpdatedAt:   asset.CreatedAt,
			}

			if err := database.DB.Create(&activityLog).Error; err != nil {
				result.Errors++
			} else {
				result.LogsCreated++
			}
		}

		offset += options.BatchSize
	}

	return result, nil
}

// backfillGoals creates activity logs for existing goals
func backfillGoals(options BackfillOptions) (BackfillResult, error) {
	result := BackfillResult{Module: models.ModuleGoal}

	query := database.DB.Model(&models.Goal{})

	if options.UserID != nil {
		query = query.Where("user_id = ?", options.UserID)
	}
	if options.StartDate != nil {
		query = query.Where("created_at >= ?", options.StartDate)
	}
	if options.EndDate != nil {
		query = query.Where("created_at <= ?", options.EndDate)
	}

	query.Count(&result.TotalRecords)

	if options.DryRun {
		result.LogsCreated = result.TotalRecords
		return result, nil
	}

	offset := 0
	for {
		var goals []models.Goal
		err := query.Limit(options.BatchSize).Offset(offset).Find(&goals).Error
		if err != nil {
			return result, err
		}

		if len(goals) == 0 {
			break
		}

		for _, goal := range goals {
			if activityLogExists(goal.UserID, models.ActionCreate, models.ModuleGoal, goal.ID, goal.CreatedAt) {
				result.LogsSkipped++
				continue
			}

			metadata := map[string]interface{}{
				"backfilled":    true,
				"name":          goal.Name,
				"target_amount": goal.TargetAmount,
			}
			metadataJSON, _ := json.Marshal(metadata)

			activityLog := models.ActivityLog{
				ID:          uuid.New(),
				UserID:      goal.UserID,
				Action:      models.ActionCreate,
				Module:      models.ModuleGoal,
				EntityType:  "Goal",
				EntityID:    &goal.ID,
				Description: fmt.Sprintf("Created goal: %s (backfilled)", goal.Name),
				Metadata:    string(metadataJSON),
				CreatedAt:   goal.CreatedAt,
				UpdatedAt:   goal.CreatedAt,
			}

			if err := database.DB.Create(&activityLog).Error; err != nil {
				result.Errors++
			} else {
				result.LogsCreated++
			}
		}

		offset += options.BatchSize
	}

	return result, nil
}

// activityLogExists checks if an activity log already exists for the given parameters
func activityLogExists(userID uuid.UUID, action, module string, entityID uuid.UUID, createdAt time.Time) bool {
	var count int64
	database.DB.Model(&models.ActivityLog{}).
		Where("user_id = ? AND action = ? AND module = ? AND entity_id = ? AND created_at = ?",
			userID, action, module, entityID, createdAt).
		Count(&count)
	return count > 0
}

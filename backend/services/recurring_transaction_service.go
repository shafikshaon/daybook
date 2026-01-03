package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"daybook-backend/models"
	"daybook-backend/repository"

	"gorm.io/gorm"
)

// RecurringTransactionService defines business logic for recurring transactions
type RecurringTransactionService interface {
	ListRecurringTransactions(ctx context.Context, userID uint) ([]models.RecurringTransaction, error)
	GetRecurringTransaction(ctx context.Context, id, userID uint) (*models.RecurringTransaction, error)
	CreateRecurringTransaction(ctx context.Context, recurringTransaction *models.RecurringTransaction) (*models.RecurringTransaction, error)
	UpdateRecurringTransaction(ctx context.Context, id, userID uint, updateData *models.RecurringTransaction) (*models.RecurringTransaction, error)
	DeleteRecurringTransaction(ctx context.Context, id, userID uint) error
	ProcessRecurringTransactions(ctx context.Context, userID uint) (*ProcessRecurringResult, error)
}

type recurringTransactionService struct {
	repo            repository.RecurringTransactionRepository
	accountRepo     repository.AccountRepository
	creditCardRepo  repository.CreditCardRepository
	txManager       repository.TransactionManager
	activityService ActivityLogService
}

// ProcessRecurringResult holds the result of processing recurring transactions
type ProcessRecurringResult struct {
	Created int    `json:"created"`
	Skipped int    `json:"skipped"`
	Errors  int    `json:"errors"`
	Message string `json:"message"`
}

// NewRecurringTransactionService creates a new recurring transaction service
func NewRecurringTransactionService(
	repo repository.RecurringTransactionRepository,
	accountRepo repository.AccountRepository,
	creditCardRepo repository.CreditCardRepository,
	txManager repository.TransactionManager,
	activityService ActivityLogService,
) RecurringTransactionService {
	return &recurringTransactionService{
		repo:            repo,
		accountRepo:     accountRepo,
		creditCardRepo:  creditCardRepo,
		txManager:       txManager,
		activityService: activityService,
	}
}

// ListRecurringTransactions retrieves all recurring transactions for a user
func (s *recurringTransactionService) ListRecurringTransactions(ctx context.Context, userID uint) ([]models.RecurringTransaction, error) {
	return s.repo.FindAll(ctx, userID)
}

// GetRecurringTransaction retrieves a specific recurring transaction
func (s *recurringTransactionService) GetRecurringTransaction(ctx context.Context, id, userID uint) (*models.RecurringTransaction, error) {
	return s.repo.FindByID(ctx, id, userID)
}

// CreateRecurringTransaction creates a new recurring transaction with validation
func (s *recurringTransactionService) CreateRecurringTransaction(ctx context.Context, recurringTransaction *models.RecurringTransaction) (*models.RecurringTransaction, error) {
	// Validate required UUID fields
	if recurringTransaction.TransactionTemplate.AccountID == 0 {
		return nil, errors.New("account ID is required")
	}

	// Validate transfer-specific requirements
	if recurringTransaction.TransactionTemplate.Type == "transfer" {
		if recurringTransaction.TransactionTemplate.ToAccountID == nil || *recurringTransaction.TransactionTemplate.ToAccountID == 0 {
			return nil, errors.New("to account ID is required for transfers")
		}
	}

	// Verify account belongs to user
	_, err := s.accountRepo.FindByID(ctx, recurringTransaction.TransactionTemplate.AccountID, recurringTransaction.UserID)
	if err != nil {
		return nil, errors.New("invalid account ID")
	}

	// Validate frequency
	validFrequencies := []string{"daily", "weekly", "biweekly", "monthly", "quarterly", "yearly"}
	isValidFrequency := false
	for _, freq := range validFrequencies {
		if recurringTransaction.Frequency == freq {
			isValidFrequency = true
			break
		}
	}
	if !isValidFrequency {
		return nil, errors.New("invalid frequency. Must be one of: daily, weekly, biweekly, monthly, quarterly, yearly")
	}

	// Set template fields to satisfy database constraints
	recurringTransaction.TransactionTemplate.Date = recurringTransaction.StartDate
	recurringTransaction.TransactionTemplate.UserID = recurringTransaction.UserID

	// Create the recurring transaction
	if err := s.repo.Create(ctx, recurringTransaction); err != nil {
		return nil, err
	}

	// Log activity
	entityID := recurringTransaction.ID
	s.activityService.LogActivity(ctx, ActivityLogParams{
		UserID:      recurringTransaction.UserID,
		Action:      models.ActionCreate,
		Module:      models.ModuleTransaction,
		EntityType:  "RecurringTransaction",
		EntityID:    &entityID,
		Description: "Created recurring transaction: " + recurringTransaction.TransactionTemplate.Description,
	})

	return recurringTransaction, nil
}

// UpdateRecurringTransaction updates an existing recurring transaction
func (s *recurringTransactionService) UpdateRecurringTransaction(ctx context.Context, id, userID uint, updateData *models.RecurringTransaction) (*models.RecurringTransaction, error) {
	// Fetch existing recurring transaction
	existing, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// Validate required UUID fields
	if updateData.TransactionTemplate.AccountID == 0 {
		return nil, errors.New("account ID is required")
	}

	// Validate transfer-specific requirements
	if updateData.TransactionTemplate.Type == "transfer" {
		if updateData.TransactionTemplate.ToAccountID == nil || *updateData.TransactionTemplate.ToAccountID == 0 {
			return nil, errors.New("to account ID is required for transfers")
		}
	}

	// Verify account belongs to user if changed
	if updateData.TransactionTemplate.AccountID != existing.TransactionTemplate.AccountID {
		_, err := s.accountRepo.FindByID(ctx, updateData.TransactionTemplate.AccountID, userID)
		if err != nil {
			return nil, errors.New("invalid account ID")
		}
	}

	// Update fields
	existing.TransactionTemplate = updateData.TransactionTemplate
	existing.Frequency = updateData.Frequency
	existing.StartDate = updateData.StartDate
	existing.EndDate = updateData.EndDate
	existing.Enabled = updateData.Enabled

	// Set template fields to satisfy database constraints
	existing.TransactionTemplate.Date = existing.StartDate
	existing.TransactionTemplate.UserID = existing.UserID

	// Save updates
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	// Log activity
	s.activityService.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      models.ActionUpdate,
		Module:      models.ModuleTransaction,
		EntityType:  "RecurringTransaction",
		EntityID:    &id,
		Description: "Updated recurring transaction: " + existing.TransactionTemplate.Description,
	})

	return existing, nil
}

// DeleteRecurringTransaction deletes a recurring transaction
func (s *recurringTransactionService) DeleteRecurringTransaction(ctx context.Context, id, userID uint) error {
	// Fetch the recurring transaction to get description for logging
	recurringTransaction, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}

	// Delete the recurring transaction
	if err := s.repo.Delete(ctx, id, userID); err != nil {
		return err
	}

	// Log activity
	s.activityService.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      models.ActionDelete,
		Module:      models.ModuleTransaction,
		EntityType:  "RecurringTransaction",
		EntityID:    &id,
		Description: "Deleted recurring transaction: " + recurringTransaction.TransactionTemplate.Description,
	})

	return nil
}

// ProcessRecurringTransactions generates missing transactions for all enabled recurring transactions
func (s *recurringTransactionService) ProcessRecurringTransactions(ctx context.Context, userID uint) (*ProcessRecurringResult, error) {
	// Get all enabled recurring transactions for the user
	recurringTransactions, err := s.repo.FindEnabled(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	result := &ProcessRecurringResult{
		Message: "Recurring transactions processed successfully",
	}

	// Process each recurring transaction within a database transaction
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		for _, recurring := range recurringTransactions {
			// Skip if start date is in the future
			if recurring.StartDate.Time.After(now) {
				result.Skipped++
				continue
			}

			// Skip if end date has passed
			if recurring.EndDate != nil && recurring.EndDate.Time.Before(now) {
				result.Skipped++
				continue
			}

			// Determine the start point for generating transactions
			startFrom := recurring.StartDate.Time
			if recurring.LastProcessed != nil {
				startFrom = *recurring.LastProcessed
			}

			// Generate transactions from startFrom to now
			transactionDates := calculateTransactionDates(startFrom, now, recurring.Frequency, recurring.StartDate.Time)

			for _, txnDate := range transactionDates {
				// Skip if transaction date is before start date
				if txnDate.Before(recurring.StartDate.Time) {
					continue
				}

				// Skip if transaction date is after end date
				if recurring.EndDate != nil && txnDate.After(recurring.EndDate.Time) {
					continue
				}

				// Check if transaction already exists for this date and recurring ID (duplicate prevention)
				// Use both date comparison and recurring_id to ensure uniqueness
				var existingCount int64
				err := tx.WithContext(ctx).Model(&models.Transaction{}).
					Where("user_id = ? AND recurring_id = ? AND date::date = ?::date",
						userID, recurring.ID, txnDate.Format("2006-01-02")).
					Count(&existingCount).Error
				if err != nil {
					result.Errors++
					continue
				}

				if existingCount > 0 {
					// Transaction already exists for this date - skip to prevent duplicate
					result.Skipped++
					continue
				}

				// Create the transaction from template
				transaction := models.Transaction{
					UserID:       userID,
					AccountID:    recurring.TransactionTemplate.AccountID,
					ToAccountID:  recurring.TransactionTemplate.ToAccountID,
					Type:         recurring.TransactionTemplate.Type,
					Amount:       recurring.TransactionTemplate.Amount,
					CategoryID:   recurring.TransactionTemplate.CategoryID,
					Date:         models.Date{Time: txnDate},
					Description:  recurring.TransactionTemplate.Description,
					Tags:         recurring.TransactionTemplate.Tags,
					CreditCardID: recurring.TransactionTemplate.CreditCardID,
					Attachments:  recurring.TransactionTemplate.Attachments,
					RecurringID:  &recurring.ID,
				}

				// Create the transaction
				if err := tx.Create(&transaction).Error; err != nil {
					// Check if it's a unique constraint violation (duplicate)
					errMsg := err.Error()
					if strings.Contains(errMsg, "unique constraint") || strings.Contains(errMsg, "duplicate key") {
						// This is a duplicate - skip it instead of counting as error
						result.Skipped++
					} else {
						result.Errors++
					}
					continue
				}

				// Update balances
				isCreditCardTransaction := transaction.CreditCardID != nil

				if isCreditCardTransaction {
					// Update credit card balance using direct tx query
					var creditCard models.CreditCard
					if err := tx.WithContext(ctx).Where("id = ? AND user_id = ?", *transaction.CreditCardID, userID).First(&creditCard).Error; err != nil {
						result.Errors++
						continue
					}

					if transaction.Type == "income" || transaction.Type == "expense" {
						creditCard.CurrentBalance += transaction.Amount
					}

					if err := tx.WithContext(ctx).Save(&creditCard).Error; err != nil {
						result.Errors++
						continue
					}
				} else {
					// Update account balance using direct tx query
					var account models.Account
					if err := tx.WithContext(ctx).Where("id = ? AND user_id = ?", transaction.AccountID, userID).First(&account).Error; err != nil {
						result.Errors++
						continue
					}

					if transaction.Type == "income" {
						account.Balance += transaction.Amount
					} else if transaction.Type == "expense" {
						account.Balance -= transaction.Amount
					} else if transaction.Type == "transfer" && transaction.ToAccountID != nil {
						account.Balance -= transaction.Amount
						// Update to-account balance
						tx.Model(&models.Account{}).Where("id = ?", *transaction.ToAccountID).
							UpdateColumn("balance", gorm.Expr("balance + ?", transaction.Amount))
					}

					if err := tx.WithContext(ctx).Save(&account).Error; err != nil {
						result.Errors++
						continue
					}
				}

				result.Created++
			}

			// Update LastProcessed date using UpdateColumn to bypass hooks and avoid template_date issues
			if err := tx.WithContext(ctx).Model(&models.RecurringTransaction{}).
				Where("id = ?", recurring.ID).
				UpdateColumn("last_processed", now).Error; err != nil {
				result.Errors++
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// calculateTransactionDates calculates all transaction dates between start and end based on frequency
func calculateTransactionDates(start, end time.Time, frequency string, originalStartDate time.Time) []time.Time {
	var dates []time.Time
	current := start

	// For the first run, include the start date if it's the original start date
	if start.Equal(originalStartDate) {
		dates = append(dates, start)
	}

	for {
		// Calculate next date based on frequency
		var next time.Time
		switch frequency {
		case "daily":
			next = current.AddDate(0, 0, 1)
		case "weekly":
			next = current.AddDate(0, 0, 7)
		case "biweekly":
			next = current.AddDate(0, 0, 14)
		case "monthly":
			next = current.AddDate(0, 1, 0)
		case "quarterly":
			next = current.AddDate(0, 3, 0)
		case "yearly":
			next = current.AddDate(1, 0, 0)
		default:
			return dates
		}

		if next.After(end) {
			break
		}

		dates = append(dates, next)
		current = next
	}

	return dates
}

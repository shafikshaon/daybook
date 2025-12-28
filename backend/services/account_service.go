package services

import (
	"context"
	"errors"

	"daybook-backend/models"
	"daybook-backend/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AccountService handles account business logic
type AccountService interface {
	// ListAccounts retrieves all accounts for a user
	ListAccounts(ctx context.Context, userID uuid.UUID) ([]models.Account, error)

	// GetAccount retrieves a specific account by ID
	GetAccount(ctx context.Context, accountID, userID uuid.UUID) (*models.Account, error)

	// CreateAccount creates a new account with optional opening balance transaction
	CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error)

	// UpdateAccount updates an existing account
	UpdateAccount(ctx context.Context, accountID, userID uuid.UUID, updateData *models.Account) (*models.Account, error)

	// DeleteAccount deletes an account
	DeleteAccount(ctx context.Context, accountID, userID uuid.UUID) error
}

type accountService struct {
	accountRepo    repository.AccountRepository
	categoryRepo   repository.CategoryRepository
	txManager      repository.TransactionManager
	activityLogger ActivityLogService
}

// NewAccountService creates a new account service
func NewAccountService(
	accountRepo repository.AccountRepository,
	categoryRepo repository.CategoryRepository,
	txManager repository.TransactionManager,
	activityLogger ActivityLogService,
) AccountService {
	return &accountService{
		accountRepo:    accountRepo,
		categoryRepo:   categoryRepo,
		txManager:      txManager,
		activityLogger: activityLogger,
	}
}

// ListAccounts retrieves all accounts
func (s *accountService) ListAccounts(ctx context.Context, userID uuid.UUID) ([]models.Account, error) {
	return s.accountRepo.FindAll(ctx, userID)
}

// GetAccount retrieves a specific account
func (s *accountService) GetAccount(ctx context.Context, accountID, userID uuid.UUID) (*models.Account, error) {
	return s.accountRepo.FindByID(ctx, accountID, userID)
}

// CreateAccount creates a new account with opening balance transaction
func (s *accountService) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	// Validate required fields
	if account.Name == "" {
		return nil, errors.New("account name is required")
	}
	if account.Type == "" {
		return nil, errors.New("account type is required")
	}

	// Use transaction for atomicity
	err := s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// Create the account using transactional repository
		accountRepoTx := s.accountRepo.WithTx(tx)
		if err := accountRepoTx.Create(ctx, account); err != nil {
			return err
		}

		// If there's an initial balance, create an opening balance transaction
		if account.InitialBalance > 0 {
			// Find or create "Opening Balance" category
			var openingBalanceCategory models.Category
			err := tx.WithContext(ctx).
				Where("user_id = ? AND name = ? AND type = ?", account.UserID, "Opening Balance", "income").
				First(&openingBalanceCategory).Error

			if err != nil {
				if err == gorm.ErrRecordNotFound {
					// Category doesn't exist, create it
					openingBalanceCategory = models.Category{
						UserID:      account.UserID,
						Name:        "Opening Balance",
						Type:        "income",
						Icon:        "🏦",
						Color:       "#10B981",
						Description: "Initial account balance",
						IsDefault:   true,
					}

					if err := tx.WithContext(ctx).Create(&openingBalanceCategory).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}

			// Create opening balance transaction
			transaction := models.Transaction{
				UserID:      account.UserID,
				AccountID:   account.ID,
				Type:        "income",
				CategoryID:  openingBalanceCategory.ID.String(),
				Amount:      account.InitialBalance,
				Date:        models.Date{Time: account.CreatedAt},
				Description: "Opening balance for " + account.Name,
			}

			if err := tx.WithContext(ctx).Create(&transaction).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		account.UserID,
		models.ActionCreate,
		models.ModuleAccount,
		"Account",
		account.ID,
		"Created account: "+account.Name,
		nil,
	)

	return account, nil
}

// UpdateAccount updates an existing account
func (s *accountService) UpdateAccount(ctx context.Context, accountID, userID uuid.UUID, updateData *models.Account) (*models.Account, error) {
	// Fetch existing account
	existing, err := s.accountRepo.FindByID(ctx, accountID, userID)
	if err != nil {
		return nil, err
	}

	// Update only allowed fields
	if updateData.Name != "" {
		existing.Name = updateData.Name
	}
	if updateData.Type != "" {
		existing.Type = updateData.Type
	}
	if updateData.Currency != "" {
		existing.Currency = updateData.Currency
	}
	existing.Description = updateData.Description
	existing.Institution = updateData.Institution
	existing.AccountNumber = updateData.AccountNumber
	existing.Active = updateData.Active

	// NOTE: InitialBalance and Balance are NOT updated
	// InitialBalance is set once at creation and never changes
	// Balance is updated automatically by transactions

	// Save updates
	if err := s.accountRepo.Update(ctx, existing); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionUpdate,
		models.ModuleAccount,
		"Account",
		existing.ID,
		"Updated account: "+existing.Name,
		nil,
	)

	return existing, nil
}

// DeleteAccount deletes an account
func (s *accountService) DeleteAccount(ctx context.Context, accountID, userID uuid.UUID) error {
	// Fetch the account to get its name
	account, err := s.accountRepo.FindByID(ctx, accountID, userID)
	if err != nil {
		return err
	}

	// TODO: Add validation to check if account has transactions
	// This should prevent deletion of accounts with transaction history

	// Delete the account
	if err := s.accountRepo.Delete(ctx, accountID, userID); err != nil {
		return err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionDelete,
		models.ModuleAccount,
		"Account",
		account.ID,
		"Deleted account: "+account.Name,
		nil,
	)

	return nil
}

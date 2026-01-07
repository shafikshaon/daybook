package services

import (
	"context"
	"errors"

	"daybook-backend/models"
	"daybook-backend/repository"

	"gorm.io/gorm"
)

// TransactionService defines business logic for transactions
type TransactionService interface {
	ListTransactions(ctx context.Context, userID uint, filters repository.TransactionFilters, pagination repository.PaginationParams) (*TransactionListResponse, error)
	GetTransaction(ctx context.Context, id, userID uint) (*TransactionResponse, error)
	CreateTransaction(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error)
	UpdateTransaction(ctx context.Context, id, userID uint, updateData *models.Transaction) (*models.Transaction, error)
	DeleteTransaction(ctx context.Context, id, userID uint) error
	BulkImportTransactions(ctx context.Context, userID uint, transactions []models.Transaction) (*BulkImportResult, error)
	GetTransactionStats(ctx context.Context, userID uint, filters repository.TransactionFilters) (*repository.TransactionStats, error)
}

type transactionService struct {
	repo            repository.TransactionRepository
	accountRepo     repository.AccountRepository
	creditCardRepo  repository.CreditCardRepository
	categoryRepo    repository.CategoryRepository
	txManager       repository.TransactionManager
	activityService ActivityLogService
}

// TransactionResponse holds enriched transaction data
type TransactionResponse struct {
	models.Transaction
	AccountName    *string `json:"accountName,omitempty"`
	CreditCardName *string `json:"creditCardName,omitempty"`
	ToAccountName  *string `json:"toAccountName,omitempty"`
}

// TransactionListResponse holds paginated transaction list with enrichment
type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Pagination   PaginationMetadata    `json:"pagination"`
}

// PaginationMetadata holds pagination information
type PaginationMetadata struct {
	CurrentPage int   `json:"currentPage"`
	Limit       int   `json:"limit"`
	TotalCount  int64 `json:"totalCount"`
	TotalPages  int   `json:"totalPages"`
	HasNext     bool  `json:"hasNext"`
	HasPrev     bool  `json:"hasPrev"`
}

// BulkImportResult holds the result of bulk import operation
type BulkImportResult struct {
	Imported int    `json:"imported"`
	Failed   int    `json:"failed"`
	Message  string `json:"message"`
}

// NewTransactionService creates a new transaction service
func NewTransactionService(
	repo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	creditCardRepo repository.CreditCardRepository,
	categoryRepo repository.CategoryRepository,
	txManager repository.TransactionManager,
	activityService ActivityLogService,
) TransactionService {
	return &transactionService{
		repo:            repo,
		accountRepo:     accountRepo,
		creditCardRepo:  creditCardRepo,
		categoryRepo:    categoryRepo,
		txManager:       txManager,
		activityService: activityService,
	}
}

// ListTransactions retrieves transactions with filtering, pagination, and enrichment
func (s *transactionService) ListTransactions(ctx context.Context, userID uint, filters repository.TransactionFilters, pagination repository.PaginationParams) (*TransactionListResponse, error) {
	// Fetch transactions
	result, err := s.repo.FindWithFilters(ctx, userID, filters, pagination)
	if err != nil {
		return nil, err
	}

	// Enrich with account/credit card names
	enrichedTransactions := make([]TransactionResponse, len(result.Transactions))
	for i, txn := range result.Transactions {
		enrichedTransactions[i] = TransactionResponse{Transaction: txn}

		// If transaction has a credit card, fetch credit card name
		if txn.CreditCardID != nil {
			card, err := s.creditCardRepo.FindByID(ctx, *txn.CreditCardID, userID)
			if err == nil {
				enrichedTransactions[i].CreditCardName = &card.Name
				enrichedTransactions[i].AccountName = &card.Name // Use credit card name as account name
			}
		} else {
			// For regular transactions, fetch the account name
			account, err := s.accountRepo.FindByID(ctx, txn.AccountID, userID)
			if err == nil {
				enrichedTransactions[i].AccountName = &account.Name
			}
		}

		// For transfers, fetch the destination account name
		if txn.ToAccountID != nil {
			toAccount, err := s.accountRepo.FindByID(ctx, *txn.ToAccountID, userID)
			if err == nil {
				enrichedTransactions[i].ToAccountName = &toAccount.Name
			}
		}
	}

	// Calculate pagination metadata
	totalPages := int((result.TotalCount + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	return &TransactionListResponse{
		Transactions: enrichedTransactions,
		Pagination: PaginationMetadata{
			CurrentPage: pagination.Page,
			Limit:       pagination.Limit,
			TotalCount:  result.TotalCount,
			TotalPages:  totalPages,
			HasNext:     pagination.Page < totalPages,
			HasPrev:     pagination.Page > 1,
		},
	}, nil
}

// GetTransaction retrieves a specific transaction with enrichment
func (s *transactionService) GetTransaction(ctx context.Context, id, userID uint) (*TransactionResponse, error) {
	transaction, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	response := &TransactionResponse{Transaction: *transaction}

	// If transaction has a credit card, fetch credit card name
	if transaction.CreditCardID != nil {
		card, err := s.creditCardRepo.FindByID(ctx, *transaction.CreditCardID, userID)
		if err == nil {
			response.CreditCardName = &card.Name
			response.AccountName = &card.Name
		}
	} else {
		// For regular transactions, fetch the account name
		account, err := s.accountRepo.FindByID(ctx, transaction.AccountID, userID)
		if err == nil {
			response.AccountName = &account.Name
		}
	}

	// For transfers, fetch the destination account name
	if transaction.ToAccountID != nil {
		toAccount, err := s.accountRepo.FindByID(ctx, *transaction.ToAccountID, userID)
		if err == nil {
			response.ToAccountName = &toAccount.Name
		}
	}

	return response, nil
}

// CreateTransaction creates a new transaction and updates account/credit card balances
func (s *transactionService) CreateTransaction(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error) {
	// Validate date
	if transaction.Date.IsZero() {
		return nil, errors.New("date is required")
	}

	// For transfers, automatically get or create Transfer category if categoryId is 0 or not provided
	if transaction.Type == "transfer" && transaction.CategoryID == 0 {
		// Try to find existing Transfer category
		categories, err := s.categoryRepo.FindAll(ctx, transaction.UserID)
		if err == nil {
			for _, cat := range categories {
				if cat.Name == "Transfer" && cat.Type == "transfer" {
					transaction.CategoryID = models.FlexibleUint(cat.ID)
					break
				}
			}
		}

		// If still not found, create Transfer category
		if transaction.CategoryID == 0 {
			transferCategory := &models.Category{
				UserID:      transaction.UserID,
				Name:        "Transfer",
				Type:        "transfer",
				Icon:        "🔄",
				Color:       "#8B5CF6",
				IsDefault:   true,
				Description: "Transfer between accounts",
			}
			if err := s.categoryRepo.Create(ctx, transferCategory); err == nil {
				transaction.CategoryID = models.FlexibleUint(transferCategory.ID)
			}
		}
	}

	// Determine if this is a credit card transaction or account transaction
	isCreditCardTransaction := transaction.CreditCardID != nil

	// Verify ownership
	if isCreditCardTransaction {
		_, err := s.creditCardRepo.FindByID(ctx, *transaction.CreditCardID, transaction.UserID)
		if err != nil {
			return nil, errors.New("invalid credit card ID")
		}
	} else {
		_, err := s.accountRepo.FindByID(ctx, transaction.AccountID, transaction.UserID)
		if err != nil {
			return nil, errors.New("invalid account ID")
		}

		// For transfers, verify the destination account
		if transaction.Type == "transfer" && transaction.ToAccountID != nil {
			_, err := s.accountRepo.FindByID(ctx, *transaction.ToAccountID, transaction.UserID)
			if err != nil {
				return nil, errors.New("invalid destination account ID")
			}
		}
	}

	// Execute within a transaction
	err := s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// Create the transaction
		repoTx := s.repo.WithTx(tx)
		if err := repoTx.Create(ctx, transaction); err != nil {
			return err
		}

		// Update balance
		if isCreditCardTransaction {
			// Update credit card balance
			cardRepoTx := s.creditCardRepo.WithTx(tx)
			creditCard, err := cardRepoTx.(repository.CreditCardRepository).FindByID(ctx, *transaction.CreditCardID, transaction.UserID)
			if err != nil {
				return err
			}

			if transaction.Type == "income" {
				creditCard.CurrentBalance -= transaction.Amount // Income reduces credit card debt
			} else if transaction.Type == "expense" {
				creditCard.CurrentBalance += transaction.Amount // Expense increases credit card debt
			}

			if err := cardRepoTx.Update(ctx, creditCard); err != nil {
				return err
			}
		} else {
			// Update account balance
			accountRepoTx := s.accountRepo.WithTx(tx)
			account, err := accountRepoTx.FindByID(ctx, transaction.AccountID, transaction.UserID)
			if err != nil {
				return err
			}

			if transaction.Type == "income" {
				account.Balance += transaction.Amount
			} else if transaction.Type == "expense" {
				account.Balance -= transaction.Amount
			} else if transaction.Type == "transfer" && transaction.ToAccountID != nil {
				// Deduct from source account
				account.Balance -= transaction.Amount
				// Add to destination account
				tx.Model(&models.Account{}).Where("id = ?", *transaction.ToAccountID).
					UpdateColumn("balance", gorm.Expr("balance + ?", transaction.Amount))
			}

			if err := accountRepoTx.Update(ctx, account); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Log activity
	entityID := transaction.ID
	s.activityService.LogActivity(ctx, ActivityLogParams{
		UserID:      transaction.UserID,
		Action:      models.ActionCreate,
		Module:      models.ModuleTransaction,
		EntityType:  "Transaction",
		EntityID:    &entityID,
		Description: "Created transaction: " + transaction.Description,
	})

	return transaction, nil
}

// UpdateTransaction updates a transaction and recalculates balances
func (s *transactionService) UpdateTransaction(ctx context.Context, id, userID uint, updateData *models.Transaction) (*models.Transaction, error) {
	// Fetch existing transaction
	existing, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// Validate date
	if updateData.Date.IsZero() {
		return nil, errors.New("date is required")
	}

	// Store old values for balance reversal
	oldAmount := existing.Amount
	oldType := existing.Type
	oldAccountID := existing.AccountID
	oldToAccountID := existing.ToAccountID
	oldCreditCardID := existing.CreditCardID

	// Determine transaction types
	wasCreditCardTransaction := oldCreditCardID != nil
	isCreditCardTransaction := updateData.CreditCardID != nil

	// Execute within a transaction
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// Revert old balance changes first
		if wasCreditCardTransaction {
			cardRepoTx := s.creditCardRepo.WithTx(tx)
			creditCard, err := cardRepoTx.(repository.CreditCardRepository).FindByID(ctx, *oldCreditCardID, userID)
			if err != nil {
				return err
			}

			// Reverse the old transaction
			if oldType == "income" {
				creditCard.CurrentBalance += oldAmount // Reverse income
			} else if oldType == "expense" {
				creditCard.CurrentBalance -= oldAmount // Reverse expense
			}

			if err := cardRepoTx.Update(ctx, creditCard); err != nil {
				return err
			}
		} else {
			accountRepoTx := s.accountRepo.WithTx(tx)
			account, err := accountRepoTx.FindByID(ctx, oldAccountID, userID)
			if err != nil {
				return err
			}

			// Reverse the old transaction
			if oldType == "income" {
				account.Balance -= oldAmount
			} else if oldType == "expense" {
				account.Balance += oldAmount
			} else if oldType == "transfer" && oldToAccountID != nil {
				// Reverse transfer
				account.Balance += oldAmount
				tx.Model(&models.Account{}).Where("id = ?", *oldToAccountID).
					UpdateColumn("balance", gorm.Expr("balance - ?", oldAmount))
			}

			if err := accountRepoTx.Update(ctx, account); err != nil {
				return err
			}
		}

		// Update the transaction record
		existing.AccountID = updateData.AccountID
		existing.ToAccountID = updateData.ToAccountID
		existing.Type = updateData.Type
		existing.Amount = updateData.Amount
		existing.CategoryID = updateData.CategoryID
		existing.Date = updateData.Date
		existing.Description = updateData.Description
		existing.Tags = updateData.Tags
		existing.CreditCardID = updateData.CreditCardID
		existing.Attachments = updateData.Attachments

		repoTx := s.repo.WithTx(tx)
		if err := repoTx.Update(ctx, existing); err != nil {
			return err
		}

		// Apply new balance changes
		if isCreditCardTransaction {
			cardRepoTx := s.creditCardRepo.WithTx(tx)
			creditCard, err := cardRepoTx.(repository.CreditCardRepository).FindByID(ctx, *updateData.CreditCardID, userID)
			if err != nil {
				return err
			}

			if updateData.Type == "income" {
				creditCard.CurrentBalance -= updateData.Amount
			} else if updateData.Type == "expense" {
				creditCard.CurrentBalance += updateData.Amount
			}

			if err := cardRepoTx.Update(ctx, creditCard); err != nil {
				return err
			}
		} else {
			accountRepoTx := s.accountRepo.WithTx(tx)
			account, err := accountRepoTx.FindByID(ctx, updateData.AccountID, userID)
			if err != nil {
				return err
			}

			if updateData.Type == "income" {
				account.Balance += updateData.Amount
			} else if updateData.Type == "expense" {
				account.Balance -= updateData.Amount
			} else if updateData.Type == "transfer" && updateData.ToAccountID != nil {
				account.Balance -= updateData.Amount
				tx.Model(&models.Account{}).Where("id = ?", *updateData.ToAccountID).
					UpdateColumn("balance", gorm.Expr("balance + ?", updateData.Amount))
			}

			if err := accountRepoTx.Update(ctx, account); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Log activity
	s.activityService.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      models.ActionUpdate,
		Module:      models.ModuleTransaction,
		EntityType:  "Transaction",
		EntityID:    &id,
		Description: "Updated transaction: " + existing.Description,
	})

	return existing, nil
}

// DeleteTransaction deletes a transaction and reverses balance changes
func (s *transactionService) DeleteTransaction(ctx context.Context, id, userID uint) error {
	// Fetch the transaction
	transaction, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}

	// Execute within a transaction
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// Reverse balance changes
		isCreditCardTransaction := transaction.CreditCardID != nil

		if isCreditCardTransaction {
			cardRepoTx := s.creditCardRepo.WithTx(tx)
			creditCard, err := cardRepoTx.(repository.CreditCardRepository).FindByID(ctx, *transaction.CreditCardID, userID)
			if err != nil {
				return err
			}

			// Reverse the transaction
			if transaction.Type == "income" {
				creditCard.CurrentBalance += transaction.Amount
			} else if transaction.Type == "expense" {
				creditCard.CurrentBalance -= transaction.Amount
			}

			if err := cardRepoTx.Update(ctx, creditCard); err != nil {
				return err
			}
		} else {
			accountRepoTx := s.accountRepo.WithTx(tx)
			account, err := accountRepoTx.FindByID(ctx, transaction.AccountID, userID)
			if err != nil {
				return err
			}

			// Reverse the transaction
			if transaction.Type == "income" {
				account.Balance -= transaction.Amount
			} else if transaction.Type == "expense" {
				account.Balance += transaction.Amount
			} else if transaction.Type == "transfer" && transaction.ToAccountID != nil {
				// Reverse transfer
				account.Balance += transaction.Amount
				tx.Model(&models.Account{}).Where("id = ?", *transaction.ToAccountID).
					UpdateColumn("balance", gorm.Expr("balance - ?", transaction.Amount))
			}

			if err := accountRepoTx.Update(ctx, account); err != nil {
				return err
			}
		}

		// Delete the transaction
		repoTx := s.repo.WithTx(tx)
		if err := repoTx.Delete(ctx, id, userID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Log activity
	s.activityService.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      models.ActionDelete,
		Module:      models.ModuleTransaction,
		EntityType:  "Transaction",
		EntityID:    &id,
		Description: "Deleted transaction: " + transaction.Description,
	})

	return nil
}

// BulkImportTransactions imports multiple transactions
func (s *transactionService) BulkImportTransactions(ctx context.Context, userID uint, transactions []models.Transaction) (*BulkImportResult, error) {
	result := &BulkImportResult{
		Message: "Bulk import completed",
	}

	for i := range transactions {
		transactions[i].UserID = userID

		// Create each transaction (this handles balance updates)
		if _, err := s.CreateTransaction(ctx, &transactions[i]); err != nil {
			result.Failed++
		} else {
			result.Imported++
		}
	}

	return result, nil
}

// GetTransactionStats calculates transaction statistics
func (s *transactionService) GetTransactionStats(ctx context.Context, userID uint, filters repository.TransactionFilters) (*repository.TransactionStats, error) {
	return s.repo.CalculateStats(ctx, userID, filters)
}

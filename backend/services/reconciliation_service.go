package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"daybook-backend/models"
	"daybook-backend/repository"

	"gorm.io/gorm"
)

// CreateReconciliationRequest represents the request body for creating a reconciliation
type CreateReconciliationRequest struct {
	AccountID          uint      `json:"accountId" binding:"required"`
	ReconciliationDate time.Time `json:"reconciliationDate" binding:"required"`
	StatementBalance   float64   `json:"statementBalance" binding:"required"`
	Notes              string    `json:"notes"`
	TransactionIDs     []uint    `json:"transactionIds"` // Optional: specific transactions to reconcile
}

// ReconciliationService handles reconciliation business logic
type ReconciliationService interface {
	// ListReconciliations retrieves all reconciliations with optional filters
	ListReconciliations(ctx context.Context, userID uint, filters repository.ReconciliationFilters) ([]models.Reconciliation, error)

	// GetReconciliation retrieves a specific reconciliation by ID
	GetReconciliation(ctx context.Context, reconciliationID, userID uint) (*models.Reconciliation, error)

	// CreateReconciliation creates a new reconciliation
	CreateReconciliation(ctx context.Context, userID uint, req *CreateReconciliationRequest) (*models.Reconciliation, error)

	// UpdateReconciliation updates an existing reconciliation
	UpdateReconciliation(ctx context.Context, reconciliationID, userID uint, updateData *models.Reconciliation) (*models.Reconciliation, error)

	// DeleteReconciliation deletes a reconciliation
	DeleteReconciliation(ctx context.Context, reconciliationID, userID uint) error

	// GetUnreconciledTransactions retrieves unreconciled transactions for an account
	GetUnreconciledTransactions(ctx context.Context, accountID, userID uint) ([]models.Transaction, error)

	// GetStats calculates reconciliation statistics for an account
	GetStats(ctx context.Context, accountID, userID uint) (*repository.ReconciliationStats, error)
}

type reconciliationService struct {
	repo           repository.ReconciliationRepository
	accountRepo    repository.AccountRepository
	txManager      repository.TransactionManager
	activityLogger ActivityLogService
}

// NewReconciliationService creates a new reconciliation service
func NewReconciliationService(
	repo repository.ReconciliationRepository,
	accountRepo repository.AccountRepository,
	txManager repository.TransactionManager,
	activityLogger ActivityLogService,
) ReconciliationService {
	return &reconciliationService{
		repo:           repo,
		accountRepo:    accountRepo,
		txManager:      txManager,
		activityLogger: activityLogger,
	}
}

// ListReconciliations retrieves reconciliations with optional filters
func (s *reconciliationService) ListReconciliations(ctx context.Context, userID uint, filters repository.ReconciliationFilters) ([]models.Reconciliation, error) {
	return s.repo.FindWithFilters(ctx, userID, filters)
}

// GetReconciliation retrieves a specific reconciliation
func (s *reconciliationService) GetReconciliation(ctx context.Context, reconciliationID, userID uint) (*models.Reconciliation, error) {
	return s.repo.FindByIDWithPreloads(ctx, reconciliationID, userID)
}

// CreateReconciliation creates a new reconciliation
func (s *reconciliationService) CreateReconciliation(ctx context.Context, userID uint, req *CreateReconciliationRequest) (*models.Reconciliation, error) {
	// Verify account belongs to user
	account, err := s.accountRepo.FindByID(ctx, req.AccountID, userID)
	if err != nil {
		return nil, errors.New("account not found")
	}

	// Create reconciliation within a transaction
	var reconciliation models.Reconciliation
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// Create reconciliation record
		reconciliation = models.Reconciliation{
			UserID:             userID,
			AccountID:          req.AccountID,
			ReconciliationDate: req.ReconciliationDate,
			StatementBalance:   req.StatementBalance,
			BookBalance:        account.Balance,
			Notes:              req.Notes,
		}

		// Use transactional repository
		reconciliationRepoTx := s.repo.WithTx(tx)
		if err := reconciliationRepoTx.Create(ctx, &reconciliation); err != nil {
			return err
		}

		// If specific transactions are provided, link them to this reconciliation
		if len(req.TransactionIDs) > 0 {
			// Verify all transactions belong to this account and user
			for _, transactionID := range req.TransactionIDs {
				var transaction models.Transaction
				if err := tx.WithContext(ctx).
					Where("id = ? AND user_id = ? AND account_id = ?", transactionID, userID, req.AccountID).
					First(&transaction).Error; err != nil {
					return fmt.Errorf("invalid transaction ID: %d", transactionID)
				}
			}

			// Link transactions using repository
			reconciliationRepoTx := s.repo.WithTx(tx)
			if err := reconciliationRepoTx.(repository.ReconciliationRepository).LinkTransactions(ctx, reconciliation.ID, req.TransactionIDs); err != nil {
				return err
			}
		}

		// Update account's reconciliation status
		accountRepoTx := s.accountRepo.WithTx(tx)
		if reconciliation.Status == models.ReconciliationCompleted {
			account.LastReconciled = &reconciliation.ReconciliationDate
			account.ReconciliationDifference = 0
		} else {
			account.ReconciliationDifference = reconciliation.Difference
		}

		if err := accountRepoTx.Update(ctx, account); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Reload reconciliation with relationships
	reloaded, err := s.repo.FindByIDWithPreloads(ctx, reconciliation.ID, userID)
	if err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionCreate,
		models.ModuleReconcile,
		"Reconciliation",
		reconciliation.ID,
		"Created reconciliation for account",
		nil,
	)

	return reloaded, nil
}

// UpdateReconciliation updates an existing reconciliation
func (s *reconciliationService) UpdateReconciliation(ctx context.Context, reconciliationID, userID uint, updateData *models.Reconciliation) (*models.Reconciliation, error) {
	// Fetch existing reconciliation
	existing, err := s.repo.FindByID(ctx, reconciliationID, userID)
	if err != nil {
		return nil, errors.New("reconciliation not found")
	}

	// Update allowed fields
	existing.ReconciliationDate = updateData.ReconciliationDate
	existing.StatementBalance = updateData.StatementBalance
	existing.Notes = updateData.Notes
	if updateData.Status != "" {
		existing.Status = updateData.Status
	}

	// Save updates
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	// Update account's last reconciled date if reconciliation is completed
	if existing.Status == models.ReconciliationCompleted {
		account, err := s.accountRepo.FindByID(ctx, existing.AccountID, userID)
		if err == nil {
			account.LastReconciled = &existing.ReconciliationDate
			account.ReconciliationDifference = 0
			s.accountRepo.Update(ctx, account)
		}
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionUpdate,
		models.ModuleReconcile,
		"Reconciliation",
		existing.ID,
		"Updated reconciliation",
		nil,
	)

	return existing, nil
}

// DeleteReconciliation deletes a reconciliation
func (s *reconciliationService) DeleteReconciliation(ctx context.Context, reconciliationID, userID uint) error {
	// Fetch the reconciliation to get its details
	reconciliation, err := s.repo.FindByID(ctx, reconciliationID, userID)
	if err != nil {
		return errors.New("reconciliation not found")
	}

	// Delete within a transaction
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// Use transactional repository
		reconciliationRepoTx := s.repo.WithTx(tx)

		// Unlink all transactions
		if err := reconciliationRepoTx.(repository.ReconciliationRepository).UnlinkAllTransactions(ctx, reconciliationID); err != nil {
			return err
		}

		// Delete reconciliation
		if err := reconciliationRepoTx.Delete(ctx, reconciliationID, userID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionDelete,
		models.ModuleReconcile,
		"Reconciliation",
		reconciliation.ID,
		"Deleted reconciliation",
		nil,
	)

	return nil
}

// GetUnreconciledTransactions retrieves unreconciled transactions for an account
func (s *reconciliationService) GetUnreconciledTransactions(ctx context.Context, accountID, userID uint) ([]models.Transaction, error) {
	// Verify account belongs to user
	_, err := s.accountRepo.FindByID(ctx, accountID, userID)
	if err != nil {
		return nil, errors.New("account not found")
	}

	return s.repo.GetUnreconciledTransactions(ctx, accountID, userID)
}

// GetStats calculates reconciliation statistics for an account
func (s *reconciliationService) GetStats(ctx context.Context, accountID, userID uint) (*repository.ReconciliationStats, error) {
	// Verify account belongs to user
	_, err := s.accountRepo.FindByID(ctx, accountID, userID)
	if err != nil {
		return nil, errors.New("account not found")
	}

	return s.repo.GetStats(ctx, accountID, userID)
}

package services

import (
	"context"
	"errors"

	"daybook-backend/models"
	"daybook-backend/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LendResponse represents enriched lend with account name
type LendResponse struct {
	models.LendRecord
	AccountName *string `json:"accountName,omitempty"`
}

// LendPaymentResponse represents enriched payment with account name
type LendPaymentResponse struct {
	models.LendPayment
	AccountName string `json:"accountName"`
}

// LendService handles lend business logic
type LendService interface {
	// ListLends retrieves all lends with optional filters
	ListLends(ctx context.Context, userID uuid.UUID, filters repository.LendFilters) ([]LendResponse, error)

	// GetLend retrieves a specific lend by ID
	GetLend(ctx context.Context, lendID, userID uuid.UUID) (*LendResponse, error)

	// CreateLend creates a new lend
	CreateLend(ctx context.Context, lend *models.LendRecord) (*models.LendRecord, error)

	// UpdateLend updates an existing lend
	UpdateLend(ctx context.Context, lendID, userID uuid.UUID, updateData map[string]interface{}) (*models.LendRecord, error)

	// DeleteLend deletes a lend
	DeleteLend(ctx context.Context, lendID, userID uuid.UUID) error

	// RecordPayment records a payment received for a lend
	RecordPayment(ctx context.Context, lendID, userID uuid.UUID, payment *models.LendPayment) (*models.LendPayment, *models.LendRecord, error)

	// ListPayments retrieves all payments for a specific lend
	ListPayments(ctx context.Context, lendID, userID uuid.UUID) ([]LendPaymentResponse, error)
}

type lendService struct {
	repo           repository.LendRepository
	accountRepo    repository.AccountRepository
	txManager      repository.TransactionManager
	activityLogger ActivityLogService
}

// NewLendService creates a new lend service
func NewLendService(
	repo repository.LendRepository,
	accountRepo repository.AccountRepository,
	txManager repository.TransactionManager,
	activityLogger ActivityLogService,
) LendService {
	return &lendService{
		repo:           repo,
		accountRepo:    accountRepo,
		txManager:      txManager,
		activityLogger: activityLogger,
	}
}

// ListLends retrieves lends with optional filters
func (s *lendService) ListLends(ctx context.Context, userID uuid.UUID, filters repository.LendFilters) ([]LendResponse, error) {
	lends, err := s.repo.FindWithFilters(ctx, userID, filters)
	if err != nil {
		return nil, err
	}

	// Enrich with account names
	enrichedLends := make([]LendResponse, len(lends))
	for i, lend := range lends {
		enrichedLends[i] = LendResponse{LendRecord: lend}

		if lend.AccountID != nil {
			accountName, err := s.repo.GetAccountName(ctx, *lend.AccountID)
			if err == nil {
				enrichedLends[i].AccountName = &accountName
			}
		}
	}

	return enrichedLends, nil
}

// GetLend retrieves a specific lend
func (s *lendService) GetLend(ctx context.Context, lendID, userID uuid.UUID) (*LendResponse, error) {
	lend, err := s.repo.FindByID(ctx, lendID, userID)
	if err != nil {
		return nil, errors.New("lend not found")
	}

	response := &LendResponse{LendRecord: *lend}

	if lend.AccountID != nil {
		accountName, err := s.repo.GetAccountName(ctx, *lend.AccountID)
		if err == nil {
			response.AccountName = &accountName
		}
	}

	return response, nil
}

// CreateLend creates a new lend
func (s *lendService) CreateLend(ctx context.Context, lend *models.LendRecord) (*models.LendRecord, error) {
	// Validate lent date
	if lend.LentDate.IsZero() {
		return nil, errors.New("lent date is required")
	}

	// If account is specified, verify it belongs to user and create transaction
	if lend.AccountID != nil {
		// Verify account belongs to user
		_, err := s.accountRepo.FindByID(ctx, *lend.AccountID, lend.UserID)
		if err != nil {
			return nil, errors.New("invalid account ID")
		}

		// Create lend within a transaction
		err = s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
			// Create the lend record
			lendRepoTx := s.repo.WithTx(tx)
			if err := lendRepoTx.Create(ctx, lend); err != nil {
				return err
			}

			// If not initial lend, create transaction and update account balance
			if !lend.IsInitial {
				// Create transaction record for the lent money
				transaction := models.Transaction{
					UserID:      lend.UserID,
					AccountID:   *lend.AccountID,
					Type:        "expense",
					Amount:      lend.OriginalAmount,
					CategoryID:  "lend",
					Date:        lend.LentDate,
					Description: "Lent to " + lend.DebtorName,
				}

				if err := tx.WithContext(ctx).Create(&transaction).Error; err != nil {
					return err
				}

				// Update account balance
				var account models.Account
				if err := tx.WithContext(ctx).Where("id = ?", *lend.AccountID).First(&account).Error; err != nil {
					return err
				}

				account.Balance -= lend.OriginalAmount
				if err := tx.WithContext(ctx).Save(&account).Error; err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	} else {
		// No account specified, just create the lend record
		if err := s.repo.Create(ctx, lend); err != nil {
			return nil, err
		}
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		lend.UserID,
		models.ActionCreate,
		models.ModuleLend,
		"Lend",
		lend.ID,
		"Created lend: "+lend.DebtorName,
		nil,
	)

	return lend, nil
}

// UpdateLend updates an existing lend
func (s *lendService) UpdateLend(ctx context.Context, lendID, userID uuid.UUID, updateData map[string]interface{}) (*models.LendRecord, error) {
	// Fetch existing lend
	existing, err := s.repo.FindByID(ctx, lendID, userID)
	if err != nil {
		return nil, errors.New("lend not found")
	}

	// Prevent updating certain fields
	delete(updateData, "id")
	delete(updateData, "userId")
	delete(updateData, "originalAmount")
	delete(updateData, "remainingAmount")
	delete(updateData, "createdAt")

	// Update using map to allow partial updates
	if err := s.repo.Query(ctx, userID).Model(existing).Updates(updateData).Error; err != nil {
		return nil, err
	}

	// Reload the updated lend
	updated, err := s.repo.FindByID(ctx, lendID, userID)
	if err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionUpdate,
		models.ModuleLend,
		"Lend",
		existing.ID,
		"Updated lend: "+existing.DebtorName,
		nil,
	)

	return updated, nil
}

// DeleteLend deletes a lend
func (s *lendService) DeleteLend(ctx context.Context, lendID, userID uuid.UUID) error {
	// Fetch the lend to get its details
	lend, err := s.repo.FindByID(ctx, lendID, userID)
	if err != nil {
		return errors.New("lend not found")
	}

	// Delete the lend
	if err := s.repo.Delete(ctx, lendID, userID); err != nil {
		return err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionDelete,
		models.ModuleLend,
		"Lend",
		lend.ID,
		"Deleted lend: "+lend.DebtorName,
		nil,
	)

	return nil
}

// RecordPayment records a payment received for a lend
func (s *lendService) RecordPayment(ctx context.Context, lendID, userID uuid.UUID, payment *models.LendPayment) (*models.LendPayment, *models.LendRecord, error) {
	// Validate payment date
	if payment.PaymentDate.IsZero() {
		return nil, nil, errors.New("payment date is required")
	}

	payment.UserID = userID
	payment.LendID = lendID

	// Verify account belongs to user
	account, err := s.accountRepo.FindByID(ctx, payment.AccountID, userID)
	if err != nil {
		return nil, nil, errors.New("invalid account ID")
	}

	var updatedLend *models.LendRecord

	// Record payment within a transaction
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// Get lend record
		lendRepoTx := s.repo.WithTx(tx)
		lend, err := lendRepoTx.FindByID(ctx, lendID, userID)
		if err != nil {
			return errors.New("lend not found")
		}

		// Validate payment amount doesn't exceed remaining amount
		if payment.Amount > lend.RemainingAmount {
			return errors.New("payment amount exceeds remaining lend")
		}

		// Create payment record
		lendRepoTx2 := s.repo.WithTx(tx)
		if err := lendRepoTx2.(repository.LendRepository).CreatePayment(ctx, payment); err != nil {
			return err
		}

		// Update lend remaining amount and status
		lend.RemainingAmount -= payment.Amount
		if lend.RemainingAmount == 0 {
			lend.Status = "fully_received"
		} else if lend.RemainingAmount < lend.OriginalAmount {
			lend.Status = "partially_received"
		}

		if err := lendRepoTx.Update(ctx, lend); err != nil {
			return err
		}

		// Create transaction record for the payment received
		description := "Payment from " + lend.DebtorName
		if payment.Description != "" {
			description = payment.Description
		}

		transaction := models.Transaction{
			UserID:      userID,
			AccountID:   payment.AccountID,
			Type:        "income",
			Amount:      payment.Amount,
			CategoryID:  "lend_return",
			Date:        payment.PaymentDate,
			Description: description,
		}

		if err := tx.WithContext(ctx).Create(&transaction).Error; err != nil {
			return err
		}

		// Update account balance
		accountRepoTx := s.accountRepo.WithTx(tx)
		account.Balance += payment.Amount
		if err := accountRepoTx.Update(ctx, account); err != nil {
			return err
		}

		updatedLend = lend
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionCreate,
		models.ModuleLend,
		"LendPayment",
		payment.ID,
		"Recorded lend payment: "+updatedLend.DebtorName,
		nil,
	)

	return payment, updatedLend, nil
}

// ListPayments retrieves all payments for a specific lend
func (s *lendService) ListPayments(ctx context.Context, lendID, userID uuid.UUID) ([]LendPaymentResponse, error) {
	// Verify lend belongs to user
	_, err := s.repo.FindByID(ctx, lendID, userID)
	if err != nil {
		return nil, errors.New("lend not found")
	}

	payments, err := s.repo.FindPaymentsByLend(ctx, lendID, userID)
	if err != nil {
		return nil, err
	}

	// Enrich with account names
	enrichedPayments := make([]LendPaymentResponse, len(payments))
	for i, payment := range payments {
		enrichedPayments[i] = LendPaymentResponse{LendPayment: payment}

		accountName, err := s.repo.GetAccountName(ctx, payment.AccountID)
		if err == nil {
			enrichedPayments[i].AccountName = accountName
		}
	}

	return enrichedPayments, nil
}

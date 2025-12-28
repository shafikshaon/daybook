package services

import (
	"context"
	"errors"

	"daybook-backend/models"
	"daybook-backend/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DebtResponse represents enriched debt with account name
type DebtResponse struct {
	models.DebtRecord
	AccountName *string `json:"accountName,omitempty"`
}

// DebtPaymentResponse represents enriched payment with account name
type DebtPaymentResponse struct {
	models.DebtPayment
	AccountName string `json:"accountName"`
}

// DebtService handles debt business logic
type DebtService interface {
	// ListDebts retrieves all debts with optional filters
	ListDebts(ctx context.Context, userID uuid.UUID, filters repository.DebtFilters) ([]DebtResponse, error)

	// GetDebt retrieves a specific debt by ID
	GetDebt(ctx context.Context, debtID, userID uuid.UUID) (*DebtResponse, error)

	// CreateDebt creates a new debt
	CreateDebt(ctx context.Context, debt *models.DebtRecord) (*models.DebtRecord, error)

	// UpdateDebt updates an existing debt
	UpdateDebt(ctx context.Context, debtID, userID uuid.UUID, updateData map[string]interface{}) (*models.DebtRecord, error)

	// DeleteDebt deletes a debt
	DeleteDebt(ctx context.Context, debtID, userID uuid.UUID) error

	// RecordPayment records a payment towards a debt
	RecordPayment(ctx context.Context, debtID, userID uuid.UUID, payment *models.DebtPayment) (*models.DebtPayment, *models.DebtRecord, error)

	// ListPayments retrieves all payments for a specific debt
	ListPayments(ctx context.Context, debtID, userID uuid.UUID) ([]DebtPaymentResponse, error)
}

type debtService struct {
	repo           repository.DebtRepository
	accountRepo    repository.AccountRepository
	txManager      repository.TransactionManager
	activityLogger ActivityLogService
}

// NewDebtService creates a new debt service
func NewDebtService(
	repo repository.DebtRepository,
	accountRepo repository.AccountRepository,
	txManager repository.TransactionManager,
	activityLogger ActivityLogService,
) DebtService {
	return &debtService{
		repo:           repo,
		accountRepo:    accountRepo,
		txManager:      txManager,
		activityLogger: activityLogger,
	}
}

// ListDebts retrieves debts with optional filters
func (s *debtService) ListDebts(ctx context.Context, userID uuid.UUID, filters repository.DebtFilters) ([]DebtResponse, error) {
	debts, err := s.repo.FindWithFilters(ctx, userID, filters)
	if err != nil {
		return nil, err
	}

	// Enrich with account names
	enrichedDebts := make([]DebtResponse, len(debts))
	for i, debt := range debts {
		enrichedDebts[i] = DebtResponse{DebtRecord: debt}

		if debt.AccountID != nil {
			accountName, err := s.repo.GetAccountName(ctx, *debt.AccountID)
			if err == nil {
				enrichedDebts[i].AccountName = &accountName
			}
		}
	}

	return enrichedDebts, nil
}

// GetDebt retrieves a specific debt
func (s *debtService) GetDebt(ctx context.Context, debtID, userID uuid.UUID) (*DebtResponse, error) {
	debt, err := s.repo.FindByID(ctx, debtID, userID)
	if err != nil {
		return nil, errors.New("debt not found")
	}

	response := &DebtResponse{DebtRecord: *debt}

	if debt.AccountID != nil {
		accountName, err := s.repo.GetAccountName(ctx, *debt.AccountID)
		if err == nil {
			response.AccountName = &accountName
		}
	}

	return response, nil
}

// CreateDebt creates a new debt
func (s *debtService) CreateDebt(ctx context.Context, debt *models.DebtRecord) (*models.DebtRecord, error) {
	// Validate borrowed date
	if debt.BorrowedDate.IsZero() {
		return nil, errors.New("borrowed date is required")
	}

	// If account is specified, verify it belongs to user and create transaction
	if debt.AccountID != nil {
		// Verify account belongs to user
		_, err := s.accountRepo.FindByID(ctx, *debt.AccountID, debt.UserID)
		if err != nil {
			return nil, errors.New("invalid account ID")
		}

		// Create debt within a transaction
		err = s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
			// Create the debt record
			debtRepoTx := s.repo.WithTx(tx)
			if err := debtRepoTx.Create(ctx, debt); err != nil {
				return err
			}

			// If not initial debt, create transaction and update account balance
			if !debt.IsInitial {
				// Create transaction record for the borrowed money
				transaction := models.Transaction{
					UserID:      debt.UserID,
					AccountID:   *debt.AccountID,
					Type:        "income",
					Amount:      debt.OriginalAmount,
					CategoryID:  "debt",
					Date:        debt.BorrowedDate,
					Description: "Borrowed from " + debt.CreditorName,
				}

				if err := tx.WithContext(ctx).Create(&transaction).Error; err != nil {
					return err
				}

				// Update account balance
				var account models.Account
				if err := tx.WithContext(ctx).Where("id = ?", *debt.AccountID).First(&account).Error; err != nil {
					return err
				}

				account.Balance += debt.OriginalAmount
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
		// No account specified, just create the debt record
		if err := s.repo.Create(ctx, debt); err != nil {
			return nil, err
		}
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		debt.UserID,
		models.ActionCreate,
		models.ModuleDebt,
		"Debt",
		debt.ID,
		"Created debt: "+debt.CreditorName,
		nil,
	)

	return debt, nil
}

// UpdateDebt updates an existing debt
func (s *debtService) UpdateDebt(ctx context.Context, debtID, userID uuid.UUID, updateData map[string]interface{}) (*models.DebtRecord, error) {
	// Fetch existing debt
	existing, err := s.repo.FindByID(ctx, debtID, userID)
	if err != nil {
		return nil, errors.New("debt not found")
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

	// Reload the updated debt
	updated, err := s.repo.FindByID(ctx, debtID, userID)
	if err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionUpdate,
		models.ModuleDebt,
		"Debt",
		existing.ID,
		"Updated debt: "+existing.CreditorName,
		nil,
	)

	return updated, nil
}

// DeleteDebt deletes a debt
func (s *debtService) DeleteDebt(ctx context.Context, debtID, userID uuid.UUID) error {
	// Fetch the debt to get its details
	debt, err := s.repo.FindByID(ctx, debtID, userID)
	if err != nil {
		return errors.New("debt not found")
	}

	// Delete the debt
	if err := s.repo.Delete(ctx, debtID, userID); err != nil {
		return err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionDelete,
		models.ModuleDebt,
		"Debt",
		debt.ID,
		"Deleted debt: "+debt.CreditorName,
		nil,
	)

	return nil
}

// RecordPayment records a payment towards a debt
func (s *debtService) RecordPayment(ctx context.Context, debtID, userID uuid.UUID, payment *models.DebtPayment) (*models.DebtPayment, *models.DebtRecord, error) {
	// Validate payment date
	if payment.PaymentDate.IsZero() {
		return nil, nil, errors.New("payment date is required")
	}

	payment.UserID = userID
	payment.DebtID = debtID

	// Verify account belongs to user
	account, err := s.accountRepo.FindByID(ctx, payment.AccountID, userID)
	if err != nil {
		return nil, nil, errors.New("invalid account ID")
	}

	var updatedDebt *models.DebtRecord

	// Record payment within a transaction
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// Get debt record
		debtRepoTx := s.repo.WithTx(tx)
		debt, err := debtRepoTx.FindByID(ctx, debtID, userID)
		if err != nil {
			return errors.New("debt not found")
		}

		// Validate payment amount doesn't exceed remaining amount
		if payment.Amount > debt.RemainingAmount {
			return errors.New("payment amount exceeds remaining debt")
		}

		// Validate account has sufficient balance
		if account.Balance < payment.Amount {
			return errors.New("insufficient account balance")
		}

		// Create payment record
		debtRepoTx2 := s.repo.WithTx(tx)
		if err := debtRepoTx2.(repository.DebtRepository).CreatePayment(ctx, payment); err != nil {
			return err
		}

		// Update debt remaining amount and status
		debt.RemainingAmount -= payment.Amount
		if debt.RemainingAmount == 0 {
			debt.Status = "fully_paid"
		} else if debt.RemainingAmount < debt.OriginalAmount {
			debt.Status = "partially_paid"
		}

		if err := debtRepoTx.Update(ctx, debt); err != nil {
			return err
		}

		// Create transaction record for the payment
		description := "Payment to " + debt.CreditorName
		if payment.Description != "" {
			description = payment.Description
		}

		transaction := models.Transaction{
			UserID:      userID,
			AccountID:   payment.AccountID,
			Type:        "expense",
			Amount:      payment.Amount,
			CategoryID:  "debt_payment",
			Date:        payment.PaymentDate,
			Description: description,
		}

		if err := tx.WithContext(ctx).Create(&transaction).Error; err != nil {
			return err
		}

		// Update account balance
		accountRepoTx := s.accountRepo.WithTx(tx)
		account.Balance -= payment.Amount
		if err := accountRepoTx.Update(ctx, account); err != nil {
			return err
		}

		updatedDebt = debt
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
		models.ModuleDebt,
		"DebtPayment",
		payment.ID,
		"Recorded debt payment: "+updatedDebt.CreditorName,
		nil,
	)

	return payment, updatedDebt, nil
}

// ListPayments retrieves all payments for a specific debt
func (s *debtService) ListPayments(ctx context.Context, debtID, userID uuid.UUID) ([]DebtPaymentResponse, error) {
	// Verify debt belongs to user
	_, err := s.repo.FindByID(ctx, debtID, userID)
	if err != nil {
		return nil, errors.New("debt not found")
	}

	payments, err := s.repo.FindPaymentsByDebt(ctx, debtID, userID)
	if err != nil {
		return nil, err
	}

	// Enrich with account names
	enrichedPayments := make([]DebtPaymentResponse, len(payments))
	for i, payment := range payments {
		enrichedPayments[i] = DebtPaymentResponse{DebtPayment: payment}

		accountName, err := s.repo.GetAccountName(ctx, payment.AccountID)
		if err == nil {
			enrichedPayments[i].AccountName = accountName
		}
	}

	return enrichedPayments, nil
}

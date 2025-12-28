package services

import (
	"context"
	"errors"
	"time"

	"daybook-backend/models"
	"daybook-backend/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RecordTransactionRequest represents the request to record a credit card transaction
type RecordTransactionRequest struct {
	models.CreditCardTransaction
}

// RecordPaymentRequest represents the request to record a credit card payment
type RecordPaymentRequest struct {
	Amount      float64    `json:"amount" binding:"required,gt=0"`
	AccountID   uuid.UUID  `json:"accountId" binding:"required"`
	PaymentDate *time.Time `json:"paymentDate"`
	Description string     `json:"description"`
}

// RecordPaymentResponse represents the response after recording a payment
type RecordPaymentResponse struct {
	Card    models.CreditCard        `json:"card"`
	Payment models.CreditCardPayment `json:"payment"`
}

// CreditCardService handles credit card business logic
type CreditCardService interface {
	// ListCreditCards retrieves all credit cards
	ListCreditCards(ctx context.Context, userID uuid.UUID) ([]models.CreditCard, error)

	// GetCreditCard retrieves a specific credit card by ID
	GetCreditCard(ctx context.Context, cardID, userID uuid.UUID) (*models.CreditCard, error)

	// CreateCreditCard creates a new credit card
	CreateCreditCard(ctx context.Context, card *models.CreditCard) (*models.CreditCard, error)

	// UpdateCreditCard updates an existing credit card
	UpdateCreditCard(ctx context.Context, cardID, userID uuid.UUID, updateData *models.CreditCard) (*models.CreditCard, error)

	// DeleteCreditCard deletes a credit card
	DeleteCreditCard(ctx context.Context, cardID, userID uuid.UUID) error

	// RecordTransaction records a credit card transaction
	RecordTransaction(ctx context.Context, cardID, userID uuid.UUID, req *RecordTransactionRequest) (*models.CreditCardTransaction, error)

	// GetTransactions retrieves all transactions for a credit card
	GetTransactions(ctx context.Context, cardID, userID uuid.UUID) ([]models.CreditCardTransaction, error)

	// DeleteTransaction deletes a credit card transaction
	DeleteTransaction(ctx context.Context, cardID, transactionID, userID uuid.UUID) error

	// RecordPayment records a credit card payment
	RecordPayment(ctx context.Context, cardID, userID uuid.UUID, req *RecordPaymentRequest) (*RecordPaymentResponse, error)

	// GetPayments retrieves all payments for a credit card
	GetPayments(ctx context.Context, cardID, userID uuid.UUID) ([]models.CreditCardPayment, error)

	// GetStatements retrieves all statements for a credit card
	GetStatements(ctx context.Context, cardID, userID uuid.UUID) ([]models.Statement, error)

	// CreateStatement creates a new statement
	CreateStatement(ctx context.Context, statement *models.Statement) (*models.Statement, error)

	// ListRewards retrieves all rewards for a user
	ListRewards(ctx context.Context, userID uuid.UUID) ([]models.Reward, error)

	// RecordReward records a new reward
	RecordReward(ctx context.Context, reward *models.Reward) (*models.Reward, error)
}

type creditCardService struct {
	repo           repository.CreditCardRepository
	accountRepo    repository.AccountRepository
	txManager      repository.TransactionManager
	activityLogger ActivityLogService
}

// NewCreditCardService creates a new credit card service
func NewCreditCardService(
	repo repository.CreditCardRepository,
	accountRepo repository.AccountRepository,
	txManager repository.TransactionManager,
	activityLogger ActivityLogService,
) CreditCardService {
	return &creditCardService{
		repo:           repo,
		accountRepo:    accountRepo,
		txManager:      txManager,
		activityLogger: activityLogger,
	}
}

// ListCreditCards retrieves all credit cards
func (s *creditCardService) ListCreditCards(ctx context.Context, userID uuid.UUID) ([]models.CreditCard, error) {
	return s.repo.FindAll(ctx, userID)
}

// GetCreditCard retrieves a specific credit card
func (s *creditCardService) GetCreditCard(ctx context.Context, cardID, userID uuid.UUID) (*models.CreditCard, error) {
	card, err := s.repo.FindByID(ctx, cardID, userID)
	if err != nil {
		return nil, errors.New("credit card not found")
	}
	return card, nil
}

// CreateCreditCard creates a new credit card
func (s *creditCardService) CreateCreditCard(ctx context.Context, card *models.CreditCard) (*models.CreditCard, error) {
	if err := s.repo.Create(ctx, card); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		card.UserID,
		models.ActionCreate,
		models.ModuleCreditCard,
		"CreditCard",
		card.ID,
		"Created credit card: "+card.Name,
		nil,
	)

	return card, nil
}

// UpdateCreditCard updates an existing credit card
func (s *creditCardService) UpdateCreditCard(ctx context.Context, cardID, userID uuid.UUID, updateData *models.CreditCard) (*models.CreditCard, error) {
	// Fetch existing card
	existing, err := s.repo.FindByID(ctx, cardID, userID)
	if err != nil {
		return nil, errors.New("credit card not found")
	}

	// Update allowed fields
	existing.Name = updateData.Name
	existing.LastFourDigits = updateData.LastFourDigits
	existing.CardNetwork = updateData.CardNetwork
	existing.CreditLimit = updateData.CreditLimit
	existing.APR = updateData.APR
	existing.DueDate = updateData.DueDate
	existing.StatementDate = updateData.StatementDate
	existing.MinimumPayment = updateData.MinimumPayment
	existing.RewardsProgram = updateData.RewardsProgram
	existing.Active = updateData.Active
	existing.Notes = updateData.Notes

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionUpdate,
		models.ModuleCreditCard,
		"CreditCard",
		existing.ID,
		"Updated credit card: "+existing.Name,
		nil,
	)

	return existing, nil
}

// DeleteCreditCard deletes a credit card
func (s *creditCardService) DeleteCreditCard(ctx context.Context, cardID, userID uuid.UUID) error {
	// Fetch the card to get its details
	card, err := s.repo.FindByID(ctx, cardID, userID)
	if err != nil {
		return errors.New("credit card not found")
	}

	// Delete the card
	if err := s.repo.Delete(ctx, cardID, userID); err != nil {
		return err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionDelete,
		models.ModuleCreditCard,
		"CreditCard",
		card.ID,
		"Deleted credit card: "+card.Name,
		nil,
	)

	return nil
}

// RecordTransaction records a credit card transaction
func (s *creditCardService) RecordTransaction(ctx context.Context, cardID, userID uuid.UUID, req *RecordTransactionRequest) (*models.CreditCardTransaction, error) {
	// Verify card belongs to user
	card, err := s.repo.FindByID(ctx, cardID, userID)
	if err != nil {
		return nil, errors.New("credit card not found")
	}

	req.CreditCardTransaction.UserID = userID
	req.CreditCardTransaction.CardID = cardID

	var createdTransaction *models.CreditCardTransaction

	// Perform all operations within a transaction
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// Create entry in main transactions table
		mainTransaction := models.Transaction{
			UserID:       userID,
			AccountID:    cardID, // Set account_id to credit card ID
			Type:         "expense",
			Amount:       req.Amount,
			Date:         req.Date,
			Description:  req.Description,
			CategoryID:   req.CategoryID,
			CreditCardID: &cardID,
			Tags:         req.Tags,
			Attachments:  req.Attachments,
		}

		// For refunds, it's income
		if req.Type == "refund" {
			mainTransaction.Type = "income"
		}

		if err := tx.WithContext(ctx).Create(&mainTransaction).Error; err != nil {
			return err
		}

		// Link the credit card transaction to the main transaction
		req.TransactionID = mainTransaction.ID

		// Create credit card transaction record
		cardRepoTx := s.repo.WithTx(tx)
		if err := cardRepoTx.(repository.CreditCardRepository).CreateTransaction(ctx, &req.CreditCardTransaction); err != nil {
			return err
		}

		// Update card balance based on transaction type
		if req.Type == "purchase" || req.Type == "fee" || req.Type == "interest" {
			card.CurrentBalance += req.Amount
		} else if req.Type == "refund" {
			card.CurrentBalance -= req.Amount
			if card.CurrentBalance < 0 {
				card.CurrentBalance = 0
			}
		}

		if err := cardRepoTx.Update(ctx, card); err != nil {
			return err
		}

		createdTransaction = &req.CreditCardTransaction
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionCreate,
		models.ModuleCreditCard,
		"CreditCardTransaction",
		createdTransaction.ID,
		"Created credit card transaction: "+createdTransaction.Description,
		nil,
	)

	return createdTransaction, nil
}

// GetTransactions retrieves all transactions for a credit card
func (s *creditCardService) GetTransactions(ctx context.Context, cardID, userID uuid.UUID) ([]models.CreditCardTransaction, error) {
	// Verify card belongs to user
	_, err := s.repo.FindByID(ctx, cardID, userID)
	if err != nil {
		return nil, errors.New("credit card not found")
	}

	return s.repo.FindTransactionsByCard(ctx, cardID, userID)
}

// DeleteTransaction deletes a credit card transaction
func (s *creditCardService) DeleteTransaction(ctx context.Context, cardID, transactionID, userID uuid.UUID) error {
	// Get the transaction
	transaction, err := s.repo.FindTransactionByID(ctx, transactionID, userID)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Verify transaction belongs to the card
	if transaction.CardID != cardID {
		return errors.New("transaction does not belong to this card")
	}

	// Get the card
	card, err := s.repo.FindByID(ctx, cardID, userID)
	if err != nil {
		return errors.New("credit card not found")
	}

	// Perform all operations within a transaction
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// Reverse the balance update
		if transaction.Type == "purchase" || transaction.Type == "fee" || transaction.Type == "interest" {
			card.CurrentBalance -= transaction.Amount
			if card.CurrentBalance < 0 {
				card.CurrentBalance = 0
			}
		} else if transaction.Type == "refund" {
			card.CurrentBalance += transaction.Amount
		}

		cardRepoTx := s.repo.WithTx(tx)
		if err := cardRepoTx.Update(ctx, card); err != nil {
			return err
		}

		// Delete main transaction if it exists
		if transaction.TransactionID != uuid.Nil {
			if err := tx.WithContext(ctx).Delete(&models.Transaction{}, transaction.TransactionID).Error; err != nil {
				return err
			}
		}

		// Delete the credit card transaction
		if err := cardRepoTx.(repository.CreditCardRepository).DeleteTransaction(ctx, transactionID, userID); err != nil {
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
		models.ModuleCreditCard,
		"CreditCardTransaction",
		transactionID,
		"Deleted credit card transaction",
		nil,
	)

	return nil
}

// RecordPayment records a credit card payment
func (s *creditCardService) RecordPayment(ctx context.Context, cardID, userID uuid.UUID, req *RecordPaymentRequest) (*RecordPaymentResponse, error) {
	// Get the credit card
	card, err := s.repo.FindByID(ctx, cardID, userID)
	if err != nil {
		return nil, errors.New("credit card not found")
	}

	// Get the payment account
	account, err := s.accountRepo.FindByID(ctx, req.AccountID, userID)
	if err != nil {
		return nil, errors.New("payment account not found")
	}

	// Validate payment amount
	if req.Amount > card.CurrentBalance {
		return nil, errors.New("payment amount exceeds current balance")
	}

	if req.Amount > account.Balance {
		return nil, errors.New("insufficient funds in payment account")
	}

	paymentDate := time.Now()
	if req.PaymentDate != nil {
		paymentDate = *req.PaymentDate
	}

	var response RecordPaymentResponse

	// Perform all operations within a transaction
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// Create transaction record for the expense
		tags := []string{"credit_card_payment"}
		transaction := models.Transaction{
			UserID:       userID,
			AccountID:    req.AccountID,
			CategoryID:   "credit_card_payment",
			Amount:       req.Amount,
			Type:         "expense",
			Date:         models.Date{Time: paymentDate},
			Description:  req.Description,
			CreditCardID: &cardID,
			Tags:         tags,
		}

		if transaction.Description == "" {
			transaction.Description = "Credit card payment: " + card.Name
		}

		if err := tx.WithContext(ctx).Create(&transaction).Error; err != nil {
			return err
		}

		// Deduct from account balance
		accountRepoTx := s.accountRepo.WithTx(tx)
		account.Balance -= req.Amount
		if err := accountRepoTx.Update(ctx, account); err != nil {
			return err
		}

		// Update card balance and payment info
		cardRepoTx := s.repo.WithTx(tx)
		card.CurrentBalance -= req.Amount
		if card.CurrentBalance < 0 {
			card.CurrentBalance = 0
		}
		card.LastPaymentDate = &paymentDate
		card.LastPaymentAmount = req.Amount

		if err := cardRepoTx.Update(ctx, card); err != nil {
			return err
		}

		// Create payment record
		payment := models.CreditCardPayment{
			UserID:        userID,
			CardID:        cardID,
			AccountID:     req.AccountID,
			Amount:        req.Amount,
			PaymentDate:   models.Date{Time: paymentDate},
			Description:   req.Description,
			TransactionID: transaction.ID,
		}

		if err := cardRepoTx.(repository.CreditCardRepository).CreatePayment(ctx, &payment); err != nil {
			return err
		}

		// Create credit card transaction record for the payment
		ccTransaction := models.CreditCardTransaction{
			UserID:      userID,
			CardID:      cardID,
			Amount:      req.Amount,
			Description: req.Description,
			Date:        models.Date{Time: paymentDate},
			Type:        "payment",
		}

		if err := cardRepoTx.(repository.CreditCardRepository).CreateTransaction(ctx, &ccTransaction); err != nil {
			return err
		}

		response.Card = *card
		response.Payment = payment

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionCreate,
		models.ModuleCreditCard,
		"CreditCardPayment",
		response.Payment.ID,
		"Recorded credit card payment: "+card.Name,
		nil,
	)

	return &response, nil
}

// GetPayments retrieves all payments for a credit card
func (s *creditCardService) GetPayments(ctx context.Context, cardID, userID uuid.UUID) ([]models.CreditCardPayment, error) {
	// Verify card belongs to user
	_, err := s.repo.FindByID(ctx, cardID, userID)
	if err != nil {
		return nil, errors.New("credit card not found")
	}

	return s.repo.FindPaymentsByCard(ctx, cardID, userID)
}

// GetStatements retrieves all statements for a credit card
func (s *creditCardService) GetStatements(ctx context.Context, cardID, userID uuid.UUID) ([]models.Statement, error) {
	// Verify card belongs to user
	_, err := s.repo.FindByID(ctx, cardID, userID)
	if err != nil {
		return nil, errors.New("credit card not found")
	}

	return s.repo.FindStatementsByCard(ctx, cardID, userID)
}

// CreateStatement creates a new statement
func (s *creditCardService) CreateStatement(ctx context.Context, statement *models.Statement) (*models.Statement, error) {
	// Verify card belongs to user
	_, err := s.repo.FindByID(ctx, statement.CardID, statement.UserID)
	if err != nil {
		return nil, errors.New("credit card not found")
	}

	if err := s.repo.CreateStatement(ctx, statement); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		statement.UserID,
		models.ActionCreate,
		models.ModuleCreditCard,
		"Statement",
		statement.ID,
		"Created statement",
		nil,
	)

	return statement, nil
}

// ListRewards retrieves all rewards for a user
func (s *creditCardService) ListRewards(ctx context.Context, userID uuid.UUID) ([]models.Reward, error) {
	return s.repo.FindRewardsByUser(ctx, userID)
}

// RecordReward records a new reward
func (s *creditCardService) RecordReward(ctx context.Context, reward *models.Reward) (*models.Reward, error) {
	// Verify card belongs to user
	_, err := s.repo.FindByID(ctx, reward.CardID, reward.UserID)
	if err != nil {
		return nil, errors.New("credit card not found")
	}

	if err := s.repo.CreateReward(ctx, reward); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		reward.UserID,
		models.ActionCreate,
		models.ModuleCreditCard,
		"Reward",
		reward.ID,
		"Recorded reward: "+reward.Description,
		nil,
	)

	return reward, nil
}

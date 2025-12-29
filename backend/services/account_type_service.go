package services

import (
	"context"
	"errors"

	"daybook-backend/models"
	"daybook-backend/repository"
	"daybook-backend/utilities"
)

// AccountTypeService handles account type business logic
type AccountTypeService interface {
	// ListAccountTypes retrieves all account types for a user
	ListAccountTypes(ctx context.Context, userID uint) ([]models.AccountType, error)

	// GetAccountType retrieves a specific account type by ID
	GetAccountType(ctx context.Context, accountTypeID, userID uint) (*models.AccountType, error)

	// CreateAccountType creates a new account type
	CreateAccountType(ctx context.Context, accountType *models.AccountType) (*models.AccountType, error)

	// UpdateAccountType updates an existing account type
	UpdateAccountType(ctx context.Context, accountTypeID, userID uint, updateData *models.AccountType) (*models.AccountType, error)

	// DeleteAccountType deletes an account type
	DeleteAccountType(ctx context.Context, accountTypeID, userID uint) error
}

type accountTypeService struct {
	repo           repository.AccountTypeRepository
	activityLogger ActivityLogService
}

// NewAccountTypeService creates a new account type service
func NewAccountTypeService(
	repo repository.AccountTypeRepository,
	activityLogger ActivityLogService,
) AccountTypeService {
	return &accountTypeService{
		repo:           repo,
		activityLogger: activityLogger,
	}
}

// ListAccountTypes retrieves all account types
func (s *accountTypeService) ListAccountTypes(ctx context.Context, userID uint) ([]models.AccountType, error) {
	return s.repo.FindAll(ctx, userID)
}

// GetAccountType retrieves a specific account type
func (s *accountTypeService) GetAccountType(ctx context.Context, accountTypeID, userID uint) (*models.AccountType, error) {
	return s.repo.FindByID(ctx, accountTypeID, userID)
}

// CreateAccountType creates a new account type
func (s *accountTypeService) CreateAccountType(ctx context.Context, accountType *models.AccountType) (*models.AccountType, error) {
	// Validate required fields
	if accountType.Name == "" {
		return nil, errors.New("account type name is required")
	}

	// Create the account type
	if err := s.repo.Create(ctx, accountType); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		accountType.UserID,
		models.ActionCreate,
		models.ModuleAccount,
		"AccountType",
		accountType.ID,
		"Created account type: "+accountType.Name,
		nil,
	)

	return accountType, nil
}

// UpdateAccountType updates an existing account type
func (s *accountTypeService) UpdateAccountType(ctx context.Context, accountTypeID, userID uint, updateData *models.AccountType) (*models.AccountType, error) {
	// Fetch existing account type
	existing, err := s.repo.FindByID(ctx, accountTypeID, userID)
	if err != nil {
		return nil, err
	}

	// Update allowed fields
	if updateData.Name != "" {
		existing.Name = updateData.Name
	}
	existing.Icon = updateData.Icon
	existing.Description = updateData.Description
	existing.Active = updateData.Active
	existing.SortOrder = updateData.SortOrder

	// Save updates
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionUpdate,
		models.ModuleAccount,
		"AccountType",
		existing.ID,
		"Updated account type: "+existing.Name,
		nil,
	)

	return existing, nil
}

// DeleteAccountType deletes an account type
func (s *accountTypeService) DeleteAccountType(ctx context.Context, accountTypeID, userID uint) error {
	// Fetch the account type to get its name
	accountType, err := s.repo.FindByID(ctx, accountTypeID, userID)
	if err != nil {
		return err
	}

	// Check if any accounts are using this type
	// The account type is stored in lowercase with underscores (e.g., "digital_wallet")
	typeValue := utilities.ToSnakeCase(accountType.Name)

	accountCount, err := s.repo.CountAccountsUsingType(ctx, userID, typeValue)
	if err != nil {
		return err
	}
	if accountCount > 0 {
		return errors.New("cannot delete account type that is in use by existing accounts")
	}

	// Delete the account type
	if err := s.repo.Delete(ctx, accountTypeID, userID); err != nil {
		return err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionDelete,
		models.ModuleAccount,
		"AccountType",
		accountType.ID,
		"Deleted account type: "+accountType.Name,
		nil,
	)

	return nil
}

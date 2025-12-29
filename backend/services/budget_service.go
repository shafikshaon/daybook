package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"daybook-backend/models"
	"daybook-backend/repository"
)

// BudgetProgress represents budget spending progress
type BudgetProgress struct {
	Budget         *models.Budget `json:"budget"`
	TotalSpent     float64        `json:"totalSpent"`
	Remaining      float64        `json:"remaining"`
	PercentageUsed float64        `json:"percentageUsed"`
	StartDate      time.Time      `json:"startDate"`
	EndDate        time.Time      `json:"endDate"`
	IsOverBudget   bool           `json:"isOverBudget"`
	AlertTriggered bool           `json:"alertTriggered"`
}

// BudgetService handles budget business logic
type BudgetService interface {
	// ListBudgets retrieves all budgets for a user with optional filters
	ListBudgets(ctx context.Context, userID uint, filters repository.BudgetFilters) ([]models.Budget, error)

	// GetBudget retrieves a specific budget by ID
	GetBudget(ctx context.Context, budgetID, userID uint) (*models.Budget, error)

	// CreateBudget creates a new budget
	CreateBudget(ctx context.Context, budget *models.Budget) (*models.Budget, error)

	// UpdateBudget updates an existing budget
	UpdateBudget(ctx context.Context, budgetID, userID uint, updateData *models.Budget) (*models.Budget, error)

	// DeleteBudget deletes a budget
	DeleteBudget(ctx context.Context, budgetID, userID uint) error

	// GetBudgetProgress calculates spending progress for a budget
	GetBudgetProgress(ctx context.Context, budgetID, userID uint) (*BudgetProgress, error)
}

type budgetService struct {
	repo           repository.BudgetRepository
	activityLogger ActivityLogService
}

// NewBudgetService creates a new budget service
func NewBudgetService(
	repo repository.BudgetRepository,
	activityLogger ActivityLogService,
) BudgetService {
	return &budgetService{
		repo:           repo,
		activityLogger: activityLogger,
	}
}

// ListBudgets retrieves budgets with optional filters
func (s *budgetService) ListBudgets(ctx context.Context, userID uint, filters repository.BudgetFilters) ([]models.Budget, error) {
	return s.repo.FindWithFilters(ctx, userID, filters)
}

// GetBudget retrieves a specific budget
func (s *budgetService) GetBudget(ctx context.Context, budgetID, userID uint) (*models.Budget, error) {
	return s.repo.FindByID(ctx, budgetID, userID)
}

// CreateBudget creates a new budget
func (s *budgetService) CreateBudget(ctx context.Context, budget *models.Budget) (*models.Budget, error) {
	// Validate required fields
	if budget.CategoryID == 0 {
		return nil, errors.New("category ID is required")
	}
	if budget.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if budget.Period == "" {
		return nil, errors.New("period is required")
	}

	// Validate custom period dates
	if budget.Period == "custom" {
		if budget.CustomStartDate == nil || budget.CustomEndDate == nil {
			return nil, errors.New("custom period requires start and end dates")
		}
	}

	// Create the budget
	if err := s.repo.Create(ctx, budget); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		budget.UserID,
		models.ActionCreate,
		models.ModuleBudget,
		"Budget",
		budget.ID,
		fmt.Sprintf("Created budget for category %d", budget.CategoryID),
		nil,
	)

	return budget, nil
}

// UpdateBudget updates an existing budget
func (s *budgetService) UpdateBudget(ctx context.Context, budgetID, userID uint, updateData *models.Budget) (*models.Budget, error) {
	// Fetch existing budget
	existing, err := s.repo.FindByID(ctx, budgetID, userID)
	if err != nil {
		return nil, err
	}

	// Update allowed fields
	existing.CategoryID = updateData.CategoryID
	existing.Amount = updateData.Amount
	existing.Period = updateData.Period
	existing.CustomStartDate = updateData.CustomStartDate
	existing.CustomEndDate = updateData.CustomEndDate
	existing.Rollover = updateData.Rollover
	existing.AlertThreshold = updateData.AlertThreshold
	existing.Enabled = updateData.Enabled
	existing.Notes = updateData.Notes

	// Save updates
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionUpdate,
		models.ModuleBudget,
		"Budget",
		existing.ID,
		fmt.Sprintf("Updated budget for category %d", existing.CategoryID),
		nil,
	)

	return existing, nil
}

// DeleteBudget deletes a budget
func (s *budgetService) DeleteBudget(ctx context.Context, budgetID, userID uint) error {
	// Fetch the budget to get its details
	budget, err := s.repo.FindByID(ctx, budgetID, userID)
	if err != nil {
		return err
	}

	// Delete the budget
	if err := s.repo.Delete(ctx, budgetID, userID); err != nil {
		return err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionDelete,
		models.ModuleBudget,
		"Budget",
		budget.ID,
		fmt.Sprintf("Deleted budget for category %d", budget.CategoryID),
		nil,
	)

	return nil
}

// GetBudgetProgress calculates spending progress for a budget
func (s *budgetService) GetBudgetProgress(ctx context.Context, budgetID, userID uint) (*BudgetProgress, error) {
	// Fetch the budget
	budget, err := s.repo.FindByID(ctx, budgetID, userID)
	if err != nil {
		return nil, err
	}

	// Calculate date range based on period
	var startDate, endDate time.Time
	now := time.Now().UTC()

	switch budget.Period {
	case "weekly":
		// Start of current week (Sunday) in UTC
		startDate = now.AddDate(0, 0, -int(now.Weekday()))
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 0, 7)

	case "monthly":
		// Start of current month in UTC
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 1, 0)

	case "quarterly":
		// Start of current quarter in UTC
		currentMonth := int(now.Month())
		quarterStartMonth := ((currentMonth-1)/3)*3 + 1
		startDate = time.Date(now.Year(), time.Month(quarterStartMonth), 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 3, 0)

	case "yearly":
		// Start of current year in UTC
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(1, 0, 0)

	case "custom":
		if budget.CustomStartDate != nil && budget.CustomEndDate != nil {
			startDate = *budget.CustomStartDate
			endDate = *budget.CustomEndDate
		} else {
			return nil, errors.New("custom budget dates not set")
		}

	default:
		return nil, errors.New("invalid budget period")
	}

	// Calculate total spending
	totalSpent, err := s.repo.CalculateTotalSpent(ctx, userID, budget.CategoryID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Calculate progress
	progress := &BudgetProgress{
		Budget:         budget,
		TotalSpent:     totalSpent,
		Remaining:      budget.Amount - totalSpent,
		PercentageUsed: (totalSpent / budget.Amount) * 100,
		StartDate:      startDate,
		EndDate:        endDate,
		IsOverBudget:   totalSpent > budget.Amount,
		AlertTriggered: (totalSpent / budget.Amount * 100) >= budget.AlertThreshold,
	}

	return progress, nil
}

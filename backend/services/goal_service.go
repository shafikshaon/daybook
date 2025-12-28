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

// AddHoldingRequest represents the request to add a holding to a goal
type AddHoldingRequest struct {
	models.GoalHolding
	AccountID  *uuid.UUID `json:"accountId"`
	IsExisting bool       `json:"isExisting"`
}

// RemoveHoldingRequest represents the request to remove a holding
type RemoveHoldingRequest struct {
	AccountID    uuid.UUID   `json:"accountId" binding:"required"`
	CurrentValue float64     `json:"currentValue" binding:"required,gt=0"`
	Date         models.Date `json:"date"`
	Notes        string      `json:"notes"`
}

// AddHoldingResponse represents the response after adding a holding
type AddHoldingResponse struct {
	Holding      models.GoalHolding      `json:"holding"`
	Contribution models.GoalContribution `json:"contribution"`
	Transaction  models.Transaction      `json:"transaction"`
	Goal         models.Goal             `json:"goal"`
}

// RemoveHoldingResponse represents the response after removing a holding
type RemoveHoldingResponse struct {
	Holding      models.GoalHolding      `json:"holding"`
	Transaction  models.Transaction      `json:"transaction"`
	Contribution models.GoalContribution `json:"contribution"`
}

// GoalService handles goal business logic
type GoalService interface {
	// ListGoals retrieves all goals with optional filters
	ListGoals(ctx context.Context, userID uuid.UUID, filters repository.GoalFilters) ([]models.Goal, error)

	// GetGoal retrieves a specific goal by ID
	GetGoal(ctx context.Context, goalID, userID uuid.UUID) (*models.Goal, error)

	// CreateGoal creates a new goal
	CreateGoal(ctx context.Context, goal *models.Goal) (*models.Goal, error)

	// UpdateGoal updates an existing goal
	UpdateGoal(ctx context.Context, goalID, userID uuid.UUID, updateData *models.Goal) (*models.Goal, error)

	// DeleteGoal deletes a goal
	DeleteGoal(ctx context.Context, goalID, userID uuid.UUID) error

	// AddHolding adds a new holding to a goal
	AddHolding(ctx context.Context, goalID, userID uuid.UUID, req *AddHoldingRequest) (*AddHoldingResponse, error)

	// UpdateHolding updates a holding
	UpdateHolding(ctx context.Context, holdingID, userID uuid.UUID, updateData *models.GoalHolding) (*models.GoalHolding, error)

	// RemoveHolding removes/liquidates a holding
	RemoveHolding(ctx context.Context, holdingID, userID uuid.UUID, req *RemoveHoldingRequest) (*RemoveHoldingResponse, error)

	// GetHoldingTypes returns all available holding types
	GetHoldingTypes(ctx context.Context) map[string]interface{}
}

type goalService struct {
	repo           repository.GoalRepository
	accountRepo    repository.AccountRepository
	txManager      repository.TransactionManager
	activityLogger ActivityLogService
}

// NewGoalService creates a new goal service
func NewGoalService(
	repo repository.GoalRepository,
	accountRepo repository.AccountRepository,
	txManager repository.TransactionManager,
	activityLogger ActivityLogService,
) GoalService {
	return &goalService{
		repo:           repo,
		accountRepo:    accountRepo,
		txManager:      txManager,
		activityLogger: activityLogger,
	}
}

// ListGoals retrieves goals with optional filters
func (s *goalService) ListGoals(ctx context.Context, userID uuid.UUID, filters repository.GoalFilters) ([]models.Goal, error) {
	goals, err := s.repo.FindWithFilters(ctx, userID, filters)
	if err != nil {
		return nil, err
	}

	// Update current amount for each goal
	for i := range goals {
		s.repo.UpdateCurrentAmount(ctx, goals[i].ID)
	}

	// Reload goals to get updated amounts
	return s.repo.FindWithFilters(ctx, userID, filters)
}

// GetGoal retrieves a specific goal
func (s *goalService) GetGoal(ctx context.Context, goalID, userID uuid.UUID) (*models.Goal, error) {
	_, err := s.repo.FindByIDWithPreloads(ctx, goalID, userID)
	if err != nil {
		return nil, errors.New("goal not found")
	}

	// Update current amount
	s.repo.UpdateCurrentAmount(ctx, goalID)

	// Reload goal to get updated amount
	return s.repo.FindByIDWithPreloads(ctx, goalID, userID)
}

// CreateGoal creates a new goal
func (s *goalService) CreateGoal(ctx context.Context, goal *models.Goal) (*models.Goal, error) {
	goal.Status = models.GoalStatusActive
	goal.CurrentAmount = 0

	if err := s.repo.Create(ctx, goal); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		goal.UserID,
		models.ActionCreate,
		models.ModuleGoal,
		"Goal",
		goal.ID,
		"Created goal: "+goal.Name,
		nil,
	)

	return goal, nil
}

// UpdateGoal updates an existing goal
func (s *goalService) UpdateGoal(ctx context.Context, goalID, userID uuid.UUID, updateData *models.Goal) (*models.Goal, error) {
	// Fetch existing goal
	existing, err := s.repo.FindByID(ctx, goalID, userID)
	if err != nil {
		return nil, errors.New("goal not found")
	}

	// Update allowed fields
	existing.Name = updateData.Name
	existing.Description = updateData.Description
	existing.Icon = updateData.Icon
	existing.Color = updateData.Color
	existing.Category = updateData.Category
	existing.Priority = updateData.Priority
	existing.TargetAmount = updateData.TargetAmount
	existing.TargetDate = updateData.TargetDate
	existing.MonthlyContribution = updateData.MonthlyContribution
	existing.Status = updateData.Status

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionUpdate,
		models.ModuleGoal,
		"Goal",
		existing.ID,
		"Updated goal: "+existing.Name,
		nil,
	)

	return existing, nil
}

// DeleteGoal deletes a goal
func (s *goalService) DeleteGoal(ctx context.Context, goalID, userID uuid.UUID) error {
	// Fetch the goal to get its details
	goal, err := s.repo.FindByID(ctx, goalID, userID)
	if err != nil {
		return errors.New("goal not found")
	}

	// Delete the goal
	if err := s.repo.Delete(ctx, goalID, userID); err != nil {
		return err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionDelete,
		models.ModuleGoal,
		"Goal",
		goal.ID,
		"Deleted goal: "+goal.Name,
		nil,
	)

	return nil
}

// AddHolding adds a new holding to a goal
func (s *goalService) AddHolding(ctx context.Context, goalID, userID uuid.UUID, req *AddHoldingRequest) (*AddHoldingResponse, error) {
	// Verify goal belongs to user
	goal, err := s.repo.FindByID(ctx, goalID, userID)
	if err != nil {
		return nil, errors.New("goal not found")
	}

	// For existing investments, we don't need to verify account or check balance
	var account *models.Account
	if !req.IsExisting {
		// Verify account belongs to user
		if req.AccountID == nil || *req.AccountID == uuid.Nil {
			return nil, errors.New("account ID is required for new investments")
		}

		account, err = s.accountRepo.FindByID(ctx, *req.AccountID, userID)
		if err != nil {
			return nil, errors.New("invalid account ID")
		}

		// Check sufficient balance
		if account.Balance < req.Amount {
			return nil, errors.New("insufficient account balance")
		}
	}

	req.GoalHolding.UserID = userID
	req.GoalHolding.GoalID = goalID
	req.GoalHolding.Status = models.HoldingStatusActive

	// Always set currentValue to amount initially
	if req.CurrentValue == 0 {
		req.CurrentValue = req.Amount
	}

	// For market instruments, calculate value based on quantity and price
	req.GoalHolding.UpdateMarketValue()

	var response AddHoldingResponse

	// Perform all operations within a transaction
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// Create holding
		goalRepoTx := s.repo.WithTx(tx)
		if err := goalRepoTx.(repository.GoalRepository).CreateHolding(ctx, &req.GoalHolding); err != nil {
			return err
		}

		// Create transaction record
		var transaction models.Transaction
		if req.IsExisting {
			// For existing investments, create a special "tracking" type transaction
			var trackingAccountID uuid.UUID
			if req.AccountID != nil {
				trackingAccountID = *req.AccountID
			}

			transaction = models.Transaction{
				UserID:      userID,
				AccountID:   trackingAccountID,
				Type:        "tracking",
				Amount:      req.Amount,
				CategoryID:  "goal_external_holding",
				Date:        req.PurchaseDate,
				Description: "External " + req.Name + " tracked for " + goal.Name,
				Tags:        []string{"goal", "holding", "external", "tracking", "hidden"},
			}
		} else {
			// For new investments, create normal transaction
			transaction = models.Transaction{
				UserID:      userID,
				AccountID:   *req.AccountID,
				Type:        "expense",
				Amount:      req.Amount,
				CategoryID:  "goal_holding_added",
				Date:        req.PurchaseDate,
				Description: "Added to " + goal.Name + ": " + req.Name,
				Tags:        []string{"goal", "holding"},
			}
		}

		if err := tx.WithContext(ctx).Create(&transaction).Error; err != nil {
			return err
		}

		req.TransactionID = transaction.ID

		// Update account balance only for new investments
		if !req.IsExisting {
			accountRepoTx := s.accountRepo.WithTx(tx)
			account.Balance -= req.Amount
			if err := accountRepoTx.Update(ctx, account); err != nil {
				return err
			}
		}

		// Create contribution record
		contributionNotes := "Added " + req.Name
		if req.IsExisting {
			contributionNotes = "External holding: " + req.Name
		}

		contribution := models.GoalContribution{
			UserID:        userID,
			GoalID:        goalID,
			HoldingID:     &req.GoalHolding.ID,
			Type:          models.ContributionTypeContribution,
			Amount:        req.Amount,
			Date:          req.PurchaseDate,
			Notes:         contributionNotes,
			TransactionID: transaction.ID,
		}

		if err := goalRepoTx.(repository.GoalRepository).CreateContribution(ctx, &contribution); err != nil {
			return err
		}

		// Update goal metadata
		goal.LastContribution = req.Amount
		goal.LastContributionDate = &req.PurchaseDate

		if err := goalRepoTx.Update(ctx, goal); err != nil {
			return err
		}

		// Save response data
		response.Holding = req.GoalHolding
		response.Contribution = contribution
		response.Transaction = transaction

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Update current amount after transaction commits
	s.repo.UpdateCurrentAmount(ctx, goalID)

	// Reload goal with all relations to return complete data
	updatedGoal, err := s.repo.FindByIDWithPreloads(ctx, goalID, userID)
	if err == nil {
		response.Goal = *updatedGoal
	} else {
		response.Goal = *goal
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionCreate,
		models.ModuleGoal,
		"GoalHolding",
		req.GoalHolding.ID,
		"Added holding "+req.Name+" to goal "+goal.Name,
		nil,
	)

	return &response, nil
}

// UpdateHolding updates a holding
func (s *goalService) UpdateHolding(ctx context.Context, holdingID, userID uuid.UUID, updateData *models.GoalHolding) (*models.GoalHolding, error) {
	existingHolding, err := s.repo.FindHoldingByID(ctx, holdingID, userID)
	if err != nil {
		return nil, errors.New("holding not found")
	}

	// Update core fields
	if updateData.Name != "" {
		existingHolding.Name = updateData.Name
	}
	if updateData.CurrentValue > 0 {
		existingHolding.CurrentValue = updateData.CurrentValue
	}
	if updateData.Status != "" {
		existingHolding.Status = updateData.Status
	}
	if !updateData.PurchaseDate.IsZero() {
		existingHolding.PurchaseDate = updateData.PurchaseDate
	}

	// Update market instrument fields
	if updateData.Symbol != nil {
		existingHolding.Symbol = updateData.Symbol
	}
	if updateData.Quantity != nil {
		existingHolding.Quantity = updateData.Quantity
	}
	if updateData.CostBasis != nil {
		existingHolding.CostBasis = updateData.CostBasis
	}
	if updateData.CurrentPrice != nil {
		existingHolding.CurrentPrice = updateData.CurrentPrice
	}

	// Update bank product fields
	if updateData.Institution != nil {
		existingHolding.Institution = updateData.Institution
	}
	if updateData.InterestRate != nil {
		existingHolding.InterestRate = updateData.InterestRate
	}
	if updateData.TenureMonths != nil {
		existingHolding.TenureMonths = updateData.TenureMonths
	}
	if updateData.MaturityDate != nil {
		existingHolding.MaturityDate = updateData.MaturityDate
	}
	if updateData.MaturityAmount != nil {
		existingHolding.MaturityAmount = updateData.MaturityAmount
	}

	// Update DPS field
	if updateData.MonthlyDeposit != nil {
		existingHolding.MonthlyDeposit = updateData.MonthlyDeposit
	}

	// Recalculate market value
	existingHolding.UpdateMarketValue()

	if err := s.repo.UpdateHolding(ctx, existingHolding); err != nil {
		return nil, err
	}

	// Update goal's current amount
	s.repo.UpdateCurrentAmount(ctx, existingHolding.GoalID)

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionUpdate,
		models.ModuleGoal,
		"GoalHolding",
		existingHolding.ID,
		"Updated holding: "+existingHolding.Name,
		nil,
	)

	return existingHolding, nil
}

// RemoveHolding removes/liquidates a holding
func (s *goalService) RemoveHolding(ctx context.Context, holdingID, userID uuid.UUID, req *RemoveHoldingRequest) (*RemoveHoldingResponse, error) {
	holding, err := s.repo.FindHoldingByID(ctx, holdingID, userID)
	if err != nil {
		return nil, errors.New("holding not found")
	}

	// Verify account
	account, err := s.accountRepo.FindByID(ctx, req.AccountID, userID)
	if err != nil {
		return nil, errors.New("invalid account ID")
	}

	if req.Date.IsZero() {
		req.Date = models.Date{Time: time.Now()}
	}

	var response RemoveHoldingResponse

	// Perform all operations within a transaction
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// Mark holding as sold/closed
		holding.Status = models.HoldingStatusSold
		holding.CurrentValue = req.CurrentValue

		goalRepoTx := s.repo.WithTx(tx)
		if err := goalRepoTx.(repository.GoalRepository).UpdateHolding(ctx, holding); err != nil {
			return err
		}

		// Create transaction (income as money returns)
		transaction := models.Transaction{
			UserID:      userID,
			AccountID:   req.AccountID,
			Type:        "income",
			Amount:      req.CurrentValue,
			CategoryID:  "goal_holding_removed",
			Date:        req.Date,
			Description: "Sold/Closed: " + holding.Name,
			Tags:        []string{"goal", "holding", "liquidation"},
		}

		if err := tx.WithContext(ctx).Create(&transaction).Error; err != nil {
			return err
		}

		// Credit account
		accountRepoTx := s.accountRepo.WithTx(tx)
		account.Balance += req.CurrentValue
		if err := accountRepoTx.Update(ctx, account); err != nil {
			return err
		}

		// Create contribution record
		contribution := models.GoalContribution{
			UserID:        userID,
			GoalID:        holding.GoalID,
			HoldingID:     &holding.ID,
			Type:          models.ContributionTypeWithdrawal,
			Amount:        req.CurrentValue,
			Date:          req.Date,
			Notes:         req.Notes,
			TransactionID: transaction.ID,
		}

		if err := goalRepoTx.(repository.GoalRepository).CreateContribution(ctx, &contribution); err != nil {
			return err
		}

		// Update goal's current amount
		if err := goalRepoTx.(repository.GoalRepository).UpdateCurrentAmount(ctx, holding.GoalID); err != nil {
			return err
		}

		// Save response data
		response.Holding = *holding
		response.Transaction = transaction
		response.Contribution = contribution

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionDelete,
		models.ModuleGoal,
		"GoalHolding",
		holding.ID,
		"Removed holding: "+holding.Name,
		nil,
	)

	return &response, nil
}

// GetHoldingTypes returns all available holding types
func (s *goalService) GetHoldingTypes(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"Savings": []map[string]string{
			{"value": "savings", "label": "Savings", "icon": "💰"},
			{"value": "fixed_deposit", "label": "Fixed Deposit", "icon": "🏦"},
			{"value": "dps", "label": "DPS (Deposit Pension Scheme)", "icon": "📅"},
			{"value": "recurring_deposit", "label": "Recurring Deposit", "icon": "🔄"},
			{"value": "savings_bond", "label": "Savings Bond", "icon": "🎫"},
			{"value": "ppf", "label": "PPF (Public Provident Fund)", "icon": "🏛️"},
			{"value": "nsc", "label": "NSC (National Savings Certificate)", "icon": "📄"},
		},
		"Investments": []map[string]string{
			{"value": "stocks", "label": "Stocks", "icon": "📈"},
			{"value": "mutual_fund", "label": "Mutual Fund", "icon": "🏛️"},
			{"value": "etf", "label": "ETF", "icon": "📊"},
			{"value": "index_fund", "label": "Index Fund", "icon": "📉"},
			{"value": "bonds", "label": "Bonds", "icon": "📜"},
			{"value": "cryptocurrency", "label": "Cryptocurrency", "icon": "₿"},
		},
		"Alternatives": []map[string]string{
			{"value": "real_estate", "label": "Real Estate", "icon": "🏢"},
			{"value": "reit", "label": "REIT", "icon": "🏗️"},
			{"value": "gold", "label": "Gold", "icon": "🥇"},
			{"value": "commodities", "label": "Commodities", "icon": "🛢️"},
		},
		"Retirement": []map[string]string{
			{"value": "pension_fund", "label": "Pension Fund", "icon": "👴"},
			{"value": "retirement_401k", "label": "401(k) / Retirement", "icon": "🏦"},
			{"value": "provident_fund", "label": "Provident Fund (EPF)", "icon": "💼"},
		},
		"Insurance": []map[string]string{
			{"value": "life_insurance", "label": "Life Insurance", "icon": "🛡️"},
			{"value": "ulip", "label": "ULIP", "icon": "🔗"},
		},
		"Other": []map[string]string{
			{"value": "custom", "label": "Custom Investment", "icon": "💎"},
		},
	}
}

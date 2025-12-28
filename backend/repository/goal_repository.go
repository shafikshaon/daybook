package repository

import (
	"context"

	"daybook-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GoalFilters represents query filters for goals
type GoalFilters struct {
	Status   *string
	Category *string
	Priority *string
}

// GoalRepository handles goal data access
type GoalRepository interface {
	BaseRepository[models.Goal]

	// FindWithFilters retrieves goals with optional filters and preloads
	FindWithFilters(ctx context.Context, userID uuid.UUID, filters GoalFilters) ([]models.Goal, error)

	// FindByIDWithPreloads retrieves a goal with all relationships
	FindByIDWithPreloads(ctx context.Context, goalID, userID uuid.UUID) (*models.Goal, error)

	// UpdateCurrentAmount recalculates current amount from all holdings
	UpdateCurrentAmount(ctx context.Context, goalID uuid.UUID) error

	// Holding operations
	CreateHolding(ctx context.Context, holding *models.GoalHolding) error
	FindHoldingByID(ctx context.Context, holdingID, userID uuid.UUID) (*models.GoalHolding, error)
	UpdateHolding(ctx context.Context, holding *models.GoalHolding) error

	// Contribution operations
	CreateContribution(ctx context.Context, contribution *models.GoalContribution) error
}

type goalRepository struct {
	*GormBaseRepository[models.Goal]
}

// NewGoalRepository creates a new goal repository
func NewGoalRepository(db *gorm.DB) GoalRepository {
	return &goalRepository{
		GormBaseRepository: NewGormBaseRepository[models.Goal](db),
	}
}

// FindWithFilters retrieves goals with optional filters and preloads
func (r *goalRepository) FindWithFilters(ctx context.Context, userID uuid.UUID, filters GoalFilters) ([]models.Goal, error) {
	var goals []models.Goal

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.Category != nil {
		query = query.Where("category = ?", *filters.Category)
	}
	if filters.Priority != nil {
		query = query.Where("priority = ?", *filters.Priority)
	}

	err := query.Order("name DESC, created_at DESC").
		Preload("Holdings").
		Preload("Contributions").
		Find(&goals).Error

	return goals, err
}

// FindByIDWithPreloads retrieves a goal with all relationships
func (r *goalRepository) FindByIDWithPreloads(ctx context.Context, goalID, userID uuid.UUID) (*models.Goal, error) {
	var goal models.Goal
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", goalID, userID).
		Preload("Holdings").
		Preload("Contributions").
		First(&goal).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

// UpdateCurrentAmount recalculates current amount from all holdings
func (r *goalRepository) UpdateCurrentAmount(ctx context.Context, goalID uuid.UUID) error {
	var total float64

	// Calculate sum of all active/matured holdings for this goal
	r.db.WithContext(ctx).Model(&models.GoalHolding{}).
		Where("goal_id = ? AND status IN ?", goalID, []string{"active", "matured"}).
		Select("COALESCE(SUM(current_value), 0)").
		Scan(&total)

	// Update the goal's current amount
	return r.db.WithContext(ctx).
		Model(&models.Goal{}).
		Where("id = ?", goalID).
		Update("current_amount", total).Error
}

// CreateHolding creates a new goal holding
func (r *goalRepository) CreateHolding(ctx context.Context, holding *models.GoalHolding) error {
	return r.db.WithContext(ctx).Create(holding).Error
}

// FindHoldingByID retrieves a specific holding
func (r *goalRepository) FindHoldingByID(ctx context.Context, holdingID, userID uuid.UUID) (*models.GoalHolding, error) {
	var holding models.GoalHolding
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", holdingID, userID).
		First(&holding).Error
	if err != nil {
		return nil, err
	}
	return &holding, nil
}

// UpdateHolding updates a goal holding
func (r *goalRepository) UpdateHolding(ctx context.Context, holding *models.GoalHolding) error {
	return r.db.WithContext(ctx).Save(holding).Error
}

// CreateContribution creates a new goal contribution
func (r *goalRepository) CreateContribution(ctx context.Context, contribution *models.GoalContribution) error {
	return r.db.WithContext(ctx).Create(contribution).Error
}

package repository

import (
	"context"
	"daybook-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ActivityLogRepository handles activity log database operations
type ActivityLogRepository interface {
	BaseRepository[models.ActivityLog]

	// FindByModule retrieves activity logs filtered by module
	FindByModule(ctx context.Context, userID uuid.UUID, module string, limit int) ([]models.ActivityLog, error)

	// FindByAction retrieves activity logs filtered by action
	FindByAction(ctx context.Context, userID uuid.UUID, action string, limit int) ([]models.ActivityLog, error)

	// FindByEntity retrieves activity logs for a specific entity
	FindByEntity(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID) ([]models.ActivityLog, error)

	// FindRecent retrieves the most recent activity logs for a user
	FindRecent(ctx context.Context, userID uuid.UUID, limit int) ([]models.ActivityLog, error)
}

type activityLogRepository struct {
	*GormBaseRepository[models.ActivityLog]
}

// NewActivityLogRepository creates a new activity log repository
func NewActivityLogRepository(db *gorm.DB) ActivityLogRepository {
	return &activityLogRepository{
		GormBaseRepository: NewGormBaseRepository[models.ActivityLog](db),
	}
}

// FindByModule retrieves activity logs for a specific module
func (r *activityLogRepository) FindByModule(ctx context.Context, userID uuid.UUID, module string, limit int) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := r.Query(ctx, userID).
		Where("module = ?", module).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// FindByAction retrieves activity logs for a specific action
func (r *activityLogRepository) FindByAction(ctx context.Context, userID uuid.UUID, action string, limit int) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := r.Query(ctx, userID).
		Where("action = ?", action).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// FindByEntity retrieves activity logs for a specific entity
func (r *activityLogRepository) FindByEntity(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := r.Query(ctx, userID).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// FindRecent retrieves the most recent activity logs
func (r *activityLogRepository) FindRecent(ctx context.Context, userID uuid.UUID, limit int) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := r.Query(ctx, userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

package repository

import (
	"context"
	"daybook-backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SettingsRepository handles settings database operations
type SettingsRepository interface {
	BaseRepository[models.Settings]

	// FindByUserID retrieves settings for a specific user
	// Settings are unique per user (one settings record per user)
	FindByUserID(ctx context.Context, userID uuid.UUID) (*models.Settings, error)

	// CreateOrUpdate creates new settings or updates existing ones
	// This is useful for settings which should have exactly one record per user
	CreateOrUpdate(ctx context.Context, settings *models.Settings) error
}

type settingsRepository struct {
	*GormBaseRepository[models.Settings]
}

// NewSettingsRepository creates a new settings repository
func NewSettingsRepository(db *gorm.DB) SettingsRepository {
	return &settingsRepository{
		GormBaseRepository: NewGormBaseRepository[models.Settings](db),
	}
}

// FindByUserID retrieves settings for a user
func (r *settingsRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.Settings, error) {
	var settings models.Settings
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

// CreateOrUpdate creates new settings or updates existing ones
func (r *settingsRepository) CreateOrUpdate(ctx context.Context, settings *models.Settings) error {
	var existing models.Settings
	err := r.db.WithContext(ctx).
		Where("user_id = ?", settings.UserID).
		First(&existing).Error

	// If not found, create new
	if err == gorm.ErrRecordNotFound {
		return r.Create(ctx, settings)
	}

	// If error occurred, return it
	if err != nil {
		return err
	}

	// Update existing settings
	settings.ID = existing.ID
	return r.Update(ctx, settings)
}

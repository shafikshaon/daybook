package services

import (
	"context"
	"daybook-backend/models"
	"daybook-backend/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SettingsService handles settings business logic
type SettingsService interface {
	// GetSettings retrieves user settings, creating default settings if none exist
	GetSettings(ctx context.Context, userID uuid.UUID) (*models.Settings, error)

	// UpdateSettings updates user settings
	UpdateSettings(ctx context.Context, userID uuid.UUID, settings *models.Settings) (*models.Settings, error)
}

type settingsService struct {
	repo           repository.SettingsRepository
	activityLogger ActivityLogService
}

// NewSettingsService creates a new settings service
func NewSettingsService(
	repo repository.SettingsRepository,
	activityLogger ActivityLogService,
) SettingsService {
	return &settingsService{
		repo:           repo,
		activityLogger: activityLogger,
	}
}

// GetSettings retrieves user settings, creating defaults if needed
func (s *settingsService) GetSettings(ctx context.Context, userID uuid.UUID) (*models.Settings, error) {
	settings, err := s.repo.FindByUserID(ctx, userID)

	// If settings don't exist, create default settings
	if err == gorm.ErrRecordNotFound {
		settings = &models.Settings{
			UserID:         userID,
			Currency:       "BDT",
			DarkMode:       false,
			DateFormat:     "MM/DD/YYYY",
			FirstDayOfWeek: 0,
			Language:       "en",
			Notifications: &models.Notifications{
				Push:         true,
				Email:        true,
				BudgetAlerts: true,
			},
		}

		if err := s.repo.Create(ctx, settings); err != nil {
			return nil, err
		}

		// Log activity for creating default settings
		s.activityLogger.LogEntityActivity(
			ctx,
			userID,
			models.ActionCreate,
			models.ModuleSettings,
			"Settings",
			settings.ID,
			"Created default user settings",
			nil,
		)

		return settings, nil
	}

	return settings, err
}

// UpdateSettings updates user settings
func (s *settingsService) UpdateSettings(ctx context.Context, userID uuid.UUID, updateData *models.Settings) (*models.Settings, error) {
	// Ensure the userID is set correctly
	updateData.UserID = userID

	// Create or update settings
	if err := s.repo.CreateOrUpdate(ctx, updateData); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionUpdate,
		models.ModuleSettings,
		"Settings",
		updateData.ID,
		"Updated user settings",
		nil,
	)

	return updateData, nil
}

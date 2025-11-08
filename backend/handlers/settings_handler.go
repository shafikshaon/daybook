package handlers

import (
	"net/http"

	"daybook-backend/database"
	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
)

// GetSettings returns the settings for the authenticated user
func GetSettings(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetSettings - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Fetching settings for user: %s", userID)
	var settings models.Settings
	if err := database.DB.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error; err != nil {
		// If settings don't exist, create default settings
		logger.Debugf(ctx, "Settings not found, creating default settings for user: %s", userID)
		settings = models.Settings{
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

		if err := database.DB.WithContext(ctx).Create(&settings).Error; err != nil {
			logger.Errorf(ctx, "Failed to create default settings: %v", err)
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create settings")
			return
		}
		logger.Infof(ctx, "Default settings created for user: %s", userID)
	}

	logger.Infof(ctx, "Settings retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, settings, "Settings retrieved successfully")
}

// UpdateSettings updates the settings for the authenticated user
func UpdateSettings(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateSettings - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Parsing settings update request for user: %s", userID)
	var updateData models.Settings
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "Fetching existing settings for user: %s", userID)
	var settings models.Settings
	result := database.DB.WithContext(ctx).Where("user_id = ?", userID).First(&settings)

	if result.Error != nil {
		// If settings don't exist, create new settings
		logger.Debugf(ctx, "Settings not found, creating new settings for user: %s", userID)
		settings = models.Settings{
			UserID:         userID,
			Currency:       updateData.Currency,
			DarkMode:       updateData.DarkMode,
			DateFormat:     updateData.DateFormat,
			FirstDayOfWeek: updateData.FirstDayOfWeek,
			Language:       updateData.Language,
			Notifications:  updateData.Notifications,
		}

		if err := database.DB.WithContext(ctx).Create(&settings).Error; err != nil {
			logger.Errorf(ctx, "Failed to create settings: %v", err)
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create settings")
			return
		}
		logger.Infof(ctx, "Settings created successfully for user: %s", userID)
	} else {
		// Update existing settings
		logger.Debugf(ctx, "Updating existing settings for user: %s", userID)
		settings.Currency = updateData.Currency
		settings.DarkMode = updateData.DarkMode
		settings.DateFormat = updateData.DateFormat
		settings.FirstDayOfWeek = updateData.FirstDayOfWeek
		settings.Language = updateData.Language

		// Update notifications if provided
		if updateData.Notifications != nil {
			settings.Notifications = updateData.Notifications
		}

		if err := database.DB.WithContext(ctx).Save(&settings).Error; err != nil {
			logger.Errorf(ctx, "Failed to update settings: %v", err)
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update settings")
			return
		}
		logger.Infof(ctx, "Settings updated successfully for user: %s", userID)
	}

	// Log settings update activity
	utilities.LogEntityActivity(c, userID, models.ActionUpdate, models.ModuleSettings,
		"Settings", settings.ID, "Updated user settings", nil)

	logger.Infof(ctx, "Settings operation completed successfully for user: %s", userID)
	utilities.SuccessResponse(c, settings, "Settings updated successfully")
}

package handlers

import (
	"net/http"

	"daybook-backend/database"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
)

// GetSettings returns the settings for the authenticated user
func GetSettings(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var settings models.Settings
	if err := database.DB.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		// If settings don't exist, create default settings
		settings = models.Settings{
			UserID:         userID,
			Currency:       "USD",
			DarkMode:       false,
			DateFormat:     "MM/DD/YYYY",
			FirstDayOfWeek: 0,
			Language:       "en",
			Notifications: &models.Notifications{
				Push:          true,
				Email:         true,
				BudgetAlerts:  true,
				BillReminders: true,
			},
		}

		if err := database.DB.Create(&settings).Error; err != nil {
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create settings")
			return
		}
	}

	utilities.SuccessResponse(c, settings, "Settings retrieved successfully")
}

// UpdateSettings updates the settings for the authenticated user
func UpdateSettings(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var updateData models.Settings
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var settings models.Settings
	result := database.DB.Where("user_id = ?", userID).First(&settings)

	if result.Error != nil {
		// If settings don't exist, create new settings
		settings = models.Settings{
			UserID:         userID,
			Currency:       updateData.Currency,
			DarkMode:       updateData.DarkMode,
			DateFormat:     updateData.DateFormat,
			FirstDayOfWeek: updateData.FirstDayOfWeek,
			Language:       updateData.Language,
			Notifications:  updateData.Notifications,
		}

		if err := database.DB.Create(&settings).Error; err != nil {
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create settings")
			return
		}
	} else {
		// Update existing settings
		settings.Currency = updateData.Currency
		settings.DarkMode = updateData.DarkMode
		settings.DateFormat = updateData.DateFormat
		settings.FirstDayOfWeek = updateData.FirstDayOfWeek
		settings.Language = updateData.Language

		// Update notifications if provided
		if updateData.Notifications != nil {
			settings.Notifications = updateData.Notifications
		}

		if err := database.DB.Save(&settings).Error; err != nil {
			utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update settings")
			return
		}
	}

	utilities.SuccessResponse(c, settings, "Settings updated successfully")
}

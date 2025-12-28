package handlers

import (
	"net/http"

	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/services"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
)

// SettingsHandler handles settings-related HTTP requests
type SettingsHandler struct {
	service services.SettingsService
}

// NewSettingsHandler creates a new settings handler
func NewSettingsHandler(service services.SettingsService) *SettingsHandler {
	return &SettingsHandler{service: service}
}

// GetSettings returns the settings for the authenticated user
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetSettings - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Fetching settings for user: %s", userID)
	settings, err := h.service.GetSettings(ctx, userID)
	if err != nil {
		logger.Errorf(ctx, "Failed to get settings: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to get settings")
		return
	}

	logger.Infof(ctx, "Settings retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, settings, "Settings retrieved successfully")
}

// UpdateSettings updates the settings for the authenticated user
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
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

	logger.Debugf(ctx, "Updating settings for user: %s", userID)
	settings, err := h.service.UpdateSettings(ctx, userID, &updateData)
	if err != nil {
		logger.Errorf(ctx, "Failed to update settings: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update settings")
		return
	}

	logger.Infof(ctx, "Settings updated successfully for user: %s", userID)
	utilities.SuccessResponse(c, settings, "Settings updated successfully")
}

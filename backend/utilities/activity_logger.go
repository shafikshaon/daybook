package utilities

import (
	"encoding/json"
	"log"

	"daybook-backend/database"
	"daybook-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ActivityLogParams contains parameters for logging an activity
type ActivityLogParams struct {
	UserID      uuid.UUID
	Action      string
	Module      string
	EntityType  string
	EntityID    *uuid.UUID
	Description string
	Metadata    map[string]interface{}
}

// LogActivity logs a user activity to the database
func LogActivity(c *gin.Context, params ActivityLogParams) {
	// Get IP address and user agent from context
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Convert metadata to JSON string
	var metadataJSON string
	if params.Metadata != nil {
		if jsonBytes, err := json.Marshal(params.Metadata); err == nil {
			metadataJSON = string(jsonBytes)
		}
	}

	activityLog := models.ActivityLog{
		UserID:      params.UserID,
		Action:      params.Action,
		Module:      params.Module,
		EntityType:  params.EntityType,
		EntityID:    params.EntityID,
		Description: params.Description,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Metadata:    metadataJSON,
	}

	// Log to database asynchronously to avoid blocking the main request
	go func() {
		if err := database.DB.Create(&activityLog).Error; err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()
}

// LogAuthActivity logs authentication-related activities
func LogAuthActivity(c *gin.Context, userID uuid.UUID, action string, description string) {
	LogActivity(c, ActivityLogParams{
		UserID:      userID,
		Action:      action,
		Module:      models.ModuleAuth,
		EntityType:  "User",
		EntityID:    &userID,
		Description: description,
	})
}

// LogEntityActivity logs entity-related activities (create, update, delete)
func LogEntityActivity(c *gin.Context, userID uuid.UUID, action string, module string, entityType string, entityID uuid.UUID, description string, metadata map[string]interface{}) {
	LogActivity(c, ActivityLogParams{
		UserID:      userID,
		Action:      action,
		Module:      module,
		EntityType:  entityType,
		EntityID:    &entityID,
		Description: description,
		Metadata:    metadata,
	})
}

// LogViewActivity logs view activities
func LogViewActivity(c *gin.Context, userID uuid.UUID, module string, entityType string, entityID *uuid.UUID, description string) {
	LogActivity(c, ActivityLogParams{
		UserID:      userID,
		Action:      models.ActionView,
		Module:      module,
		EntityType:  entityType,
		EntityID:    entityID,
		Description: description,
	})
}

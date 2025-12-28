package services

import (
	"context"
	"daybook-backend/models"
	"daybook-backend/repository"
	"encoding/json"

	"github.com/google/uuid"
)

// ActivityLogService handles activity logging business logic
// This replaces the utilities/activity_logger.go and fixes the race condition
type ActivityLogService interface {
	// LogActivity logs a generic activity
	LogActivity(ctx context.Context, params ActivityLogParams) error

	// LogEntityActivity logs an activity for a specific entity
	LogEntityActivity(ctx context.Context, userID uuid.UUID, action, module, entityType string, entityID uuid.UUID, description string, metadata map[string]interface{}) error

	// LogAuthActivity logs authentication-related activities
	LogAuthActivity(ctx context.Context, userID uuid.UUID, action, description string) error

	// GetRecentLogs retrieves recent activity logs for a user
	GetRecentLogs(ctx context.Context, userID uuid.UUID, limit int) ([]models.ActivityLog, error)

	// GetLogsByModule retrieves logs filtered by module
	GetLogsByModule(ctx context.Context, userID uuid.UUID, module string, limit int) ([]models.ActivityLog, error)
}

// ActivityLogParams contains parameters for logging an activity
type ActivityLogParams struct {
	UserID      uuid.UUID
	Action      string
	Module      string
	EntityType  string
	EntityID    *uuid.UUID
	Description string
	IPAddress   string
	UserAgent   string
	Metadata    map[string]interface{}
}

type activityLogService struct {
	repo repository.ActivityLogRepository
}

// NewActivityLogService creates a new activity log service
func NewActivityLogService(repo repository.ActivityLogRepository) ActivityLogService {
	return &activityLogService{repo: repo}
}

// LogActivity logs a generic activity
func (s *activityLogService) LogActivity(ctx context.Context, params ActivityLogParams) error {
	// Convert metadata to JSON string
	var metadataJSON string
	if params.Metadata != nil && len(params.Metadata) > 0 {
		if jsonBytes, err := json.Marshal(params.Metadata); err == nil {
			metadataJSON = string(jsonBytes)
		} else {
			metadataJSON = "null"
		}
	} else {
		metadataJSON = "null"
	}

	activityLog := &models.ActivityLog{
		UserID:      params.UserID,
		Action:      params.Action,
		Module:      params.Module,
		EntityType:  params.EntityType,
		EntityID:    params.EntityID,
		Description: params.Description,
		IPAddress:   params.IPAddress,
		UserAgent:   params.UserAgent,
		Metadata:    metadataJSON,
	}

	// Log asynchronously using goroutine
	// Note: In production, consider using a message queue for better reliability
	// Using background context to avoid cancellation
	go func() {
		bgCtx := context.Background()
		_ = s.repo.Create(bgCtx, activityLog)
	}()

	return nil
}

// LogEntityActivity logs an activity for a specific entity
func (s *activityLogService) LogEntityActivity(ctx context.Context, userID uuid.UUID, action, module, entityType string, entityID uuid.UUID, description string, metadata map[string]interface{}) error {
	return s.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      action,
		Module:      module,
		EntityType:  entityType,
		EntityID:    &entityID,
		Description: description,
		Metadata:    metadata,
	})
}

// LogAuthActivity logs authentication-related activities
func (s *activityLogService) LogAuthActivity(ctx context.Context, userID uuid.UUID, action, description string) error {
	return s.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      action,
		Module:      models.ModuleAuth,
		EntityType:  "User",
		EntityID:    &userID,
		Description: description,
	})
}

// GetRecentLogs retrieves recent activity logs for a user
func (s *activityLogService) GetRecentLogs(ctx context.Context, userID uuid.UUID, limit int) ([]models.ActivityLog, error) {
	return s.repo.FindRecent(ctx, userID, limit)
}

// GetLogsByModule retrieves logs filtered by module
func (s *activityLogService) GetLogsByModule(ctx context.Context, userID uuid.UUID, module string, limit int) ([]models.ActivityLog, error) {
	return s.repo.FindByModule(ctx, userID, module, limit)
}

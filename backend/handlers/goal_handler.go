package handlers

import (
	"net/http"

	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/repository"
	"daybook-backend/services"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GoalHandler struct {
	service services.GoalService
}

func NewGoalHandler(service services.GoalService) *GoalHandler {
	return &GoalHandler{
		service: service,
	}
}

// ListGoals returns all goals for the authenticated user
func (h *GoalHandler) ListGoals(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "ListGoals - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "ListGoals - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	logger.Debugf(ctx, "ListGoals - Fetching goals for user: %s", userID)

	// Build filters
	filters := repository.GoalFilters{}
	if status := c.Query("status"); status != "" {
		filters.Status = &status
	}
	if category := c.Query("category"); category != "" {
		filters.Category = &category
	}
	if priority := c.Query("priority"); priority != "" {
		filters.Priority = &priority
	}

	goals, err := h.service.ListGoals(ctx, userID, filters)
	if err != nil {
		logger.Errorf(ctx, "ListGoals - Failed to fetch goals: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch goals")
		return
	}

	logger.Debugf(ctx, "ListGoals - Retrieved %d goals", len(goals))
	logger.Infof(ctx, "ListGoals - Successfully retrieved goals")
	utilities.SuccessResponse(c, goals, "Goals retrieved successfully")
}

// GetGoal returns a specific goal by ID with all details
func (h *GoalHandler) GetGoal(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "GetGoal - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "GetGoal - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "GetGoal - Invalid goal ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid goal ID")
		return
	}

	logger.Debugf(ctx, "GetGoal - Fetching goal: %s", goalID)

	goal, err := h.service.GetGoal(ctx, goalID, userID)
	if err != nil {
		logger.Errorf(ctx, "GetGoal - Failed to fetch goal: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Goal not found")
		return
	}

	logger.Infof(ctx, "GetGoal - Successfully retrieved goal: %s", goalID)
	utilities.SuccessResponse(c, goal, "Goal retrieved successfully")
}

// CreateGoal creates a new goal
func (h *GoalHandler) CreateGoal(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "CreateGoal - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "CreateGoal - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	var goal models.Goal
	if err := c.ShouldBindJSON(&goal); err != nil {
		logger.Warnf(ctx, "CreateGoal - Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "CreateGoal - Creating goal: %s", goal.Name)

	goal.UserID = userID

	createdGoal, err := h.service.CreateGoal(ctx, &goal)
	if err != nil {
		logger.Errorf(ctx, "CreateGoal - Failed to create goal: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create goal")
		return
	}

	logger.Infof(ctx, "CreateGoal - Successfully created goal: %s", createdGoal.ID)
	utilities.CreatedResponse(c, createdGoal, "Goal created successfully")
}

// UpdateGoal updates an existing goal
func (h *GoalHandler) UpdateGoal(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "UpdateGoal - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "UpdateGoal - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "UpdateGoal - Invalid goal ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid goal ID")
		return
	}

	logger.Debugf(ctx, "UpdateGoal - Updating goal: %s", goalID)

	var updateData models.Goal
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "UpdateGoal - Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedGoal, err := h.service.UpdateGoal(ctx, goalID, userID, &updateData)
	if err != nil {
		logger.Errorf(ctx, "UpdateGoal - Failed to update goal: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update goal")
		return
	}

	logger.Infof(ctx, "UpdateGoal - Successfully updated goal: %s", goalID)
	utilities.SuccessResponse(c, updatedGoal, "Goal updated successfully")
}

// DeleteGoal deletes a goal
func (h *GoalHandler) DeleteGoal(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "DeleteGoal - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "DeleteGoal - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "DeleteGoal - Invalid goal ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid goal ID")
		return
	}

	logger.Debugf(ctx, "DeleteGoal - Deleting goal: %s", goalID)

	if err := h.service.DeleteGoal(ctx, goalID, userID); err != nil {
		logger.Errorf(ctx, "DeleteGoal - Failed to delete goal: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete goal")
		return
	}

	logger.Infof(ctx, "DeleteGoal - Successfully deleted goal: %s", goalID)
	utilities.SuccessResponse(c, nil, "Goal deleted successfully")
}

// AddHolding adds a new holding to a goal
func (h *GoalHandler) AddHolding(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "AddHolding - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "AddHolding - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	goalID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		logger.Warnf(ctx, "AddHolding - Invalid goal ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid goal ID")
		return
	}

	logger.Debugf(ctx, "AddHolding - Adding holding to goal: %s", goalID)

	var holdingRequest services.AddHoldingRequest
	if err := c.ShouldBindJSON(&holdingRequest); err != nil {
		logger.Warnf(ctx, "AddHolding - Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.AddHolding(ctx, goalID, userID, &holdingRequest)
	if err != nil {
		logger.Errorf(ctx, "AddHolding - Failed to add holding: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Infof(ctx, "AddHolding - Successfully added holding to goal: %s", goalID)
	utilities.CreatedResponse(c, result, "Holding added successfully")
}

// UpdateHolding updates a holding (e.g., update stock price)
func (h *GoalHandler) UpdateHolding(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "UpdateHolding - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "UpdateHolding - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	holdingID, err := uuid.Parse(c.Param("holdingId"))
	if err != nil {
		logger.Warnf(ctx, "UpdateHolding - Invalid holding ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid holding ID")
		return
	}

	logger.Debugf(ctx, "UpdateHolding - Updating holding: %s", holdingID)

	var updateData models.GoalHolding
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "UpdateHolding - Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedHolding, err := h.service.UpdateHolding(ctx, holdingID, userID, &updateData)
	if err != nil {
		logger.Errorf(ctx, "UpdateHolding - Failed to update holding: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update holding")
		return
	}

	logger.Infof(ctx, "UpdateHolding - Successfully updated holding: %s", holdingID)
	utilities.SuccessResponse(c, updatedHolding, "Holding updated successfully")
}

// RemoveHolding removes/liquidates a holding
func (h *GoalHandler) RemoveHolding(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "RemoveHolding - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "RemoveHolding - Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx = middleware.GetContextWithUserID(c)
	holdingID, err := uuid.Parse(c.Param("holdingId"))
	if err != nil {
		logger.Warnf(ctx, "RemoveHolding - Invalid holding ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid holding ID")
		return
	}

	logger.Debugf(ctx, "RemoveHolding - Removing holding: %s", holdingID)

	var removeRequest services.RemoveHoldingRequest
	if err := c.ShouldBindJSON(&removeRequest); err != nil {
		logger.Warnf(ctx, "RemoveHolding - Invalid request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.RemoveHolding(ctx, holdingID, userID, &removeRequest)
	if err != nil {
		logger.Errorf(ctx, "RemoveHolding - Failed to remove holding: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Infof(ctx, "RemoveHolding - Successfully removed holding: %s", holdingID)
	utilities.SuccessResponse(c, result, "Holding removed successfully")
}

// GetHoldingTypes returns all available holding types
func (h *GoalHandler) GetHoldingTypes(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "GetHoldingTypes - Entry")

	holdingTypes := h.service.GetHoldingTypes(ctx)

	logger.Infof(ctx, "GetHoldingTypes - Successfully retrieved holding types")
	utilities.SuccessResponse(c, holdingTypes, "Holding types retrieved successfully")
}

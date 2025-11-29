package handlers

import (
	"net/http"
	"strconv"
	"time"

	"daybook-backend/database"
	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListActivityLogs returns all activity logs for the authenticated user with filtering
func ListActivityLogs(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListActivityLogs - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	// Filter parameters
	module := c.Query("module")
	action := c.Query("action")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// Build query
	logger.Debugf(ctx, "Building query for activity logs with filters")
	query := database.DB.WithContext(ctx).Where("user_id = ?", userID)

	if module != "" {
		query = query.Where("module = ?", module)
	}

	if action != "" {
		query = query.Where("action = ?", action)
	}

	if startDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			query = query.Where("created_at >= ?", start)
		}
	}

	if endDate != "" {
		end, err := time.Parse("2006-01-02", endDate)
		if err == nil {
			// Add one day to include the entire end date
			end = end.Add(24 * time.Hour)
			query = query.Where("created_at < ?", end)
		}
	}

	// Get total count
	logger.Debugf(ctx, "Counting activity logs")
	var total int64
	if err := query.Model(&models.ActivityLog{}).Count(&total).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to count activity logs")
		return
	}

	// Get paginated results
	logger.Debugf(ctx, "Fetching paginated activity logs")
	var activityLogs []models.ActivityLog
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&activityLogs).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch activity logs")
		return
	}

	response := map[string]interface{}{
		"data":       activityLogs,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	}

	logger.Infof(ctx, "Activity logs retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, response, "Activity logs retrieved successfully")
}

// GetActivityLog returns a specific activity log by ID
func GetActivityLog(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetActivityLog - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid log ID")
		return
	}

	logger.Debugf(ctx, "Fetching activity log with ID: %s", logID)
	var activityLog models.ActivityLog
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", logID, userID).First(&activityLog).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Activity log not found")
		return
	}

	logger.Infof(ctx, "Activity log retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, activityLog, "Activity log retrieved successfully")
}

// GetActivitySummary returns summary statistics of user activities
func GetActivitySummary(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetActivitySummary - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Date range filter
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	startDate := time.Now().AddDate(0, 0, -days)

	// Total activities count
	logger.Debugf(ctx, "Counting total activities")
	var totalCount int64
	if err := database.DB.WithContext(ctx).Model(&models.ActivityLog{}).
		Where("user_id = ? AND created_at >= ?", userID, startDate).
		Count(&totalCount).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to get total count")
		return
	}

	// Count by action
	type ActionCount struct {
		Action string `json:"action"`
		Count  int64  `json:"count"`
	}
	logger.Debugf(ctx, "Fetching action counts")
	var actionCounts []ActionCount
	if err := database.DB.WithContext(ctx).Model(&models.ActivityLog{}).
		Select("action, COUNT(*) as count").
		Where("user_id = ? AND created_at >= ?", userID, startDate).
		Group("action").
		Order("count DESC").
		Scan(&actionCounts).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to get action counts")
		return
	}

	// Count by module
	type ModuleCount struct {
		Module string `json:"module"`
		Count  int64  `json:"count"`
	}
	logger.Debugf(ctx, "Fetching module counts")
	var moduleCounts []ModuleCount
	if err := database.DB.WithContext(ctx).Model(&models.ActivityLog{}).
		Select("module, COUNT(*) as count").
		Where("user_id = ? AND created_at >= ?", userID, startDate).
		Group("module").
		Order("count DESC").
		Scan(&moduleCounts).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to get module counts")
		return
	}

	// Activities per day (last 7 days)
	type DailyActivity struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	logger.Debugf(ctx, "Fetching daily activities")
	var dailyActivities []DailyActivity
	last7Days := time.Now().AddDate(0, 0, -7)
	if err := database.DB.WithContext(ctx).Model(&models.ActivityLog{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("user_id = ? AND created_at >= ?", userID, last7Days).
		Group("DATE(created_at)").
		Order("date DESC").
		Scan(&dailyActivities).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to get daily activities")
		return
	}

	// Most active modules
	logger.Debugf(ctx, "Fetching top modules")
	var topModules []ModuleCount
	if err := database.DB.WithContext(ctx).Model(&models.ActivityLog{}).
		Select("module, COUNT(*) as count").
		Where("user_id = ? AND created_at >= ?", userID, startDate).
		Group("module").
		Order("count DESC").
		Limit(5).
		Scan(&topModules).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to get top modules")
		return
	}

	// Recent activities
	logger.Debugf(ctx, "Fetching recent activities")
	var recentActivities []models.ActivityLog
	if err := database.DB.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(10).
		Find(&recentActivities).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to get recent activities")
		return
	}

	summary := map[string]interface{}{
		"totalActivities":  totalCount,
		"dateRange":        days,
		"actionCounts":     actionCounts,
		"moduleCounts":     moduleCounts,
		"dailyActivities":  dailyActivities,
		"topModules":       topModules,
		"recentActivities": recentActivities,
	}

	logger.Infof(ctx, "Activity summary retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, summary, "Activity summary retrieved successfully")
}

// DeleteOldActivityLogs deletes activity logs older than specified days (admin function)
func DeleteOldActivityLogs(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteOldActivityLogs - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	days, _ := strconv.Atoi(c.DefaultQuery("days", "90"))
	cutoffDate := time.Now().AddDate(0, 0, -days)

	logger.Debugf(ctx, "Deleting activity logs older than %d days", days)
	result := database.DB.WithContext(ctx).Where("user_id = ? AND created_at < ?", userID, cutoffDate).
		Delete(&models.ActivityLog{})

	if result.Error != nil {
		logger.Errorf(ctx, "Database operation failed: %v", result.Error)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete old logs")
		return
	}

	logger.Infof(ctx, "Old activity logs deleted successfully for user: %s (deleted: %d)", userID, result.RowsAffected)
	utilities.SuccessResponse(c, map[string]interface{}{
		"deletedCount": result.RowsAffected,
		"cutoffDate":   cutoffDate,
	}, "Old activity logs deleted successfully")
}

// BackfillActivityLogs generates activity logs for existing historical data
func BackfillActivityLogs(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "BackfillActivityLogs - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse request body for backfill options
	var requestBody struct {
		Module    string  `json:"module"`    // Optional: specific module to backfill
		DryRun    bool    `json:"dryRun"`    // Optional: if true, only count without creating logs
		StartDate *string `json:"startDate"` // Optional: only backfill records after this date (YYYY-MM-DD)
		EndDate   *string `json:"endDate"`   // Optional: only backfill records before this date (YYYY-MM-DD)
		BatchSize int     `json:"batchSize"` // Optional: number of records to process in each batch
		AllUsers  bool    `json:"allUsers"`  // Optional: admin only - backfill for all users
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		// If no body provided, use defaults
		requestBody.DryRun = false
	}

	// Build backfill options
	options := utilities.BackfillOptions{
		UserID:    &userID,
		Module:    requestBody.Module,
		DryRun:    requestBody.DryRun,
		BatchSize: requestBody.BatchSize,
	}

	// Handle all users option (you may want to add admin check here)
	if requestBody.AllUsers {
		options.UserID = nil
	}

	// Parse dates if provided
	if requestBody.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *requestBody.StartDate)
		if err != nil {
			utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid start date format. Use YYYY-MM-DD")
			return
		}
		options.StartDate = &startDate
	}

	if requestBody.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *requestBody.EndDate)
		if err != nil {
			utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid end date format. Use YYYY-MM-DD")
			return
		}
		options.EndDate = &endDate
	}

	// Execute backfill
	results, err := utilities.BackfillAllActivities(options)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to backfill activity logs: "+err.Error())
		return
	}

	// Calculate totals
	var totalRecords, totalCreated, totalSkipped, totalErrors int64
	for _, result := range results {
		totalRecords += result.TotalRecords
		totalCreated += result.LogsCreated
		totalSkipped += result.LogsSkipped
		totalErrors += result.Errors
	}

	message := "Activity logs backfilled successfully"
	if requestBody.DryRun {
		message = "Dry run completed - no logs were created"
	}

	logger.Infof(ctx, "Activity logs backfill completed for user: %s (created: %d, skipped: %d, errors: %d)", userID, totalCreated, totalSkipped, totalErrors)
	utilities.SuccessResponse(c, map[string]interface{}{
		"summary": map[string]interface{}{
			"totalRecords": totalRecords,
			"logsCreated":  totalCreated,
			"logsSkipped":  totalSkipped,
			"errors":       totalErrors,
			"dryRun":       requestBody.DryRun,
		},
		"details": results,
	}, message)
}

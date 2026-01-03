package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	customLogger "daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/services"

	"github.com/gin-gonic/gin"
)

type BackupHandler struct {
	service services.BackupService
}

func NewBackupHandler(service services.BackupService) *BackupHandler {
	return &BackupHandler{service: service}
}

// CreateBackup initiates a database backup
// POST /api/v1/backups
func (h *BackupHandler) CreateBackup(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := c.Request.Context()
	backup, err := h.service.CreateBackup(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Backup initiated successfully",
		"data":    backup,
	})
}

// ListBackups retrieves all backups for the authenticated user
// GET /api/v1/backups
func (h *BackupHandler) ListBackups(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := c.Request.Context()
	backups, err := h.service.ListBackups(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": backups})
}

// GetBackup retrieves a specific backup
// GET /api/v1/backups/:id
func (h *BackupHandler) GetBackup(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	backupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup ID"})
		return
	}

	ctx := c.Request.Context()
	backup, err := h.service.GetBackup(ctx, userID, uint(backupID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": backup})
}

// DownloadBackup downloads a backup file
// GET /api/v1/backups/:id/download
func (h *BackupHandler) DownloadBackup(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	customLogger.Infof(ctx, "DownloadBackup - Entry")

	userID := c.GetUint("userID")
	if userID == 0 {
		customLogger.Warnf(ctx, "Unauthorized access")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	backupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		customLogger.Warnf(ctx, "Invalid backup ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup ID"})
		return
	}

	customLogger.Debugf(ctx, "Fetching backup ID=%d for user=%d", backupID, userID)
	backup, err := h.service.GetBackup(ctx, userID, uint(backupID))
	if err != nil {
		customLogger.Errorf(ctx, "Backup not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	customLogger.Debugf(ctx, "Backup found: status=%s, fileName=%s", backup.Status, backup.FileName)

	if backup.Status != "completed" {
		customLogger.Warnf(ctx, "Backup is not completed: status=%s", backup.Status)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Backup is not ready for download"})
		return
	}

	filePath := h.service.GetBackupFilePath(backup)
	customLogger.Debugf(ctx, "Backup file path: %s", filePath)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		customLogger.Errorf(ctx, "Backup file does not exist at path: %s", filePath)
		c.JSON(http.StatusNotFound, gin.H{"error": "Backup file not found"})
		return
	}

	customLogger.Infof(ctx, "Serving backup file: %s (size: %d bytes)", backup.FileName, backup.FileSize)

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", backup.FileName))
	c.Header("Content-Type", "application/sql")
	c.Header("Content-Length", fmt.Sprintf("%d", backup.FileSize))

	c.File(filePath)
}

// DeleteBackup deletes a backup file and its record
// DELETE /api/v1/backups/:id
func (h *BackupHandler) DeleteBackup(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	backupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup ID"})
		return
	}

	ctx := c.Request.Context()
	if err := h.service.DeleteBackup(ctx, userID, uint(backupID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Backup deleted successfully"})
}

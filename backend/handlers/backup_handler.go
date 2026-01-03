package handlers

import (
	"fmt"
	"net/http"
	"strconv"

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

	if backup.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Backup is not ready for download"})
		return
	}

	filePath := h.service.GetBackupFilePath(backup)

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

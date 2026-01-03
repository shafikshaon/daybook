package services

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"daybook-backend/config"
	customLogger "daybook-backend/logger"
	"daybook-backend/models"
	"daybook-backend/repository"
)

// BackupService handles database backup operations
type BackupService interface {
	CreateBackup(ctx context.Context, userID uint) (*models.Backup, error)
	ListBackups(ctx context.Context, userID uint) ([]models.Backup, error)
	GetBackup(ctx context.Context, userID uint, backupID uint) (*models.Backup, error)
	DeleteBackup(ctx context.Context, userID uint, backupID uint) error
	GetBackupFilePath(backup *models.Backup) string
}

type backupService struct {
	backupRepo     repository.BackupRepository
	activityLogger ActivityLogService
	config         *config.Config
}

// NewBackupService creates a new backup service
func NewBackupService(
	backupRepo repository.BackupRepository,
	activityLogger ActivityLogService,
	cfg *config.Config,
) BackupService {
	return &backupService{
		backupRepo:     backupRepo,
		activityLogger: activityLogger,
		config:         cfg,
	}
}

// CreateBackup creates a new database backup
func (s *backupService) CreateBackup(ctx context.Context, userID uint) (*models.Backup, error) {
	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("backup_%s.sql", timestamp)

	// Ensure backup directory exists
	backupDir := s.getBackupDirectory()
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	filePath := filepath.Join(backupDir, fileName)

	// Create backup record with pending status
	backup := &models.Backup{
		UserID:   userID,
		FileName: fileName,
		FilePath: filePath,
		FileSize: 0,
		Status:   "pending",
	}

	if err := s.backupRepo.Create(ctx, backup); err != nil {
		return nil, fmt.Errorf("failed to create backup record: %w", err)
	}

	// Perform the actual backup
	go s.performBackup(backup.ID, filePath)

	// Log activity
	s.activityLogger.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      "create",
		Module:      "backup",
		EntityType:  "database",
		Description: fmt.Sprintf("Database backup initiated: %s", fileName),
	})

	return backup, nil
}

// performBackup executes pg_dump to create the backup file
func (s *backupService) performBackup(backupID uint, filePath string) {
	// Create a context with timeout (5 minutes max for backup)
	ctx, cancel := context.WithTimeout(customLogger.CreateContext(""), 5*time.Minute)
	defer cancel()

	// Find pg_dump executable - check common locations for macOS Homebrew
	pgDumpPath, err := s.findPgDump()
	if err != nil {
		errorMsg := fmt.Sprintf("pg_dump command not found: %v. Please install PostgreSQL client tools.", err)
		s.backupRepo.UpdateStatus(context.Background(), backupID, "failed", errorMsg)
		customLogger.Errorf(ctx, "[BACKUP] Failed to find pg_dump for backup ID=%d: %v", backupID, err)
		return
	}
	customLogger.Infof(ctx, "[BACKUP] Found pg_dump at: %s for backup ID=%d", pgDumpPath, backupID)

	// Create the backup file
	file, err := os.Create(filePath)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to create file: %v", err)
		s.backupRepo.UpdateStatus(context.Background(), backupID, "failed", errorMsg)
		customLogger.Errorf(ctx, "[BACKUP] Failed to create backup file for ID=%d: %v", backupID, err)
		return
	}
	defer file.Close()

	// Prepare pg_dump command
	dbConfig := s.config.Database

	// Log the backup attempt (don't log password)
	customLogger.Infof(ctx, "[BACKUP] Starting backup ID=%d, connecting to %s:%s/%s as %s",
		backupID, dbConfig.Host, dbConfig.Port, dbConfig.DBName, dbConfig.User)

	pgDumpCmd := exec.CommandContext(ctx,
		pgDumpPath,
		"-h", dbConfig.Host,
		"-p", dbConfig.Port,
		"-U", dbConfig.User,
		"-d", dbConfig.DBName,
		"--no-password",
		"--clean",
		"--if-exists",
		"--no-owner",
		"--no-privileges",
	)

	// Capture stderr for better error reporting
	var stderrBuf bytes.Buffer

	// Set PGPASSWORD environment variable (use dbConfig.Password)
	pgDumpCmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbConfig.Password))
	pgDumpCmd.Stdout = file
	pgDumpCmd.Stderr = &stderrBuf

	// Execute pg_dump
	customLogger.Infof(ctx, "[BACKUP] Executing pg_dump for backup ID=%d (timeout: 5 minutes)", backupID)
	startTime := time.Now()

	if err := pgDumpCmd.Run(); err != nil {
		os.Remove(filePath) // Clean up failed backup file
		stderrOutput := stderrBuf.String()

		// Check if it was a timeout
		if ctx.Err() == context.DeadlineExceeded {
			errorMsg := "Backup timeout: operation took longer than 5 minutes"
			s.backupRepo.UpdateStatus(context.Background(), backupID, "failed", errorMsg)
			customLogger.Errorf(ctx, "[BACKUP] Timeout for backup ID=%d after 5 minutes", backupID)
			return
		}

		errorMsg := fmt.Sprintf("pg_dump failed: %v", err)
		if stderrOutput != "" {
			errorMsg = fmt.Sprintf("pg_dump failed: %v - %s", err, stderrOutput)
		}
		s.backupRepo.UpdateStatus(context.Background(), backupID, "failed", errorMsg)
		customLogger.Errorf(ctx, "[BACKUP] pg_dump failed for backup ID=%d: %v, stderr: %s", backupID, err, stderrOutput)
		return
	}

	duration := time.Since(startTime)
	customLogger.Infof(ctx, "[BACKUP] pg_dump completed in %v for backup ID=%d", duration, backupID)

	// Get file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to get file info: %v", err)
		s.backupRepo.UpdateStatus(context.Background(), backupID, "failed", errorMsg)
		customLogger.Errorf(ctx, "[BACKUP] Failed to get file info for backup ID=%d: %v", backupID, err)
		return
	}

	customLogger.Infof(ctx, "[BACKUP] Backup ID=%d completed successfully, size=%d bytes (%.2f MB)",
		backupID, fileInfo.Size(), float64(fileInfo.Size())/(1024*1024))

	// Update the backup status and file size using UpdateStatus
	if err := s.backupRepo.UpdateStatus(context.Background(), backupID, "completed", ""); err != nil {
		customLogger.Errorf(ctx, "[BACKUP] Failed to update backup ID=%d status: %v", backupID, err)
		return
	}

	// Update file size separately using raw SQL to avoid userID constraint
	updateCtx := context.Background()
	if err := s.updateBackupFileSize(updateCtx, backupID, fileInfo.Size()); err != nil {
		customLogger.Warnf(ctx, "[BACKUP] Failed to update file size for backup ID=%d: %v", backupID, err)
		// Don't return here - the backup is still marked as completed
	}

	customLogger.Infof(ctx, "[BACKUP] Backup ID=%d status updated to completed", backupID)
}

// updateBackupFileSize updates just the file size field for a backup
func (s *backupService) updateBackupFileSize(ctx context.Context, backupID uint, fileSize int64) error {
	return s.backupRepo.UpdateFileSize(ctx, backupID, fileSize)
}

// findPgDump locates the pg_dump executable
func (s *backupService) findPgDump() (string, error) {
	// First, try to find it in PATH (works for both Linux and macOS)
	path, err := exec.LookPath("pg_dump")
	if err == nil {
		return path, nil
	}

	// Common locations for different operating systems
	commonPaths := []string{
		// Linux (EC2, Ubuntu, Debian, etc.)
		"/usr/bin/pg_dump",
		"/usr/local/bin/pg_dump",
		"/usr/pgsql-16/bin/pg_dump",
		"/usr/pgsql-15/bin/pg_dump",
		"/usr/pgsql-14/bin/pg_dump",
		"/usr/pgsql-13/bin/pg_dump",
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("pg_dump not found in PATH or common locations. Install with: sudo apt-get install postgresql-client (Ubuntu/Debian) or brew install postgresql (macOS)")
}

// ListBackups retrieves all backups for a user
func (s *backupService) ListBackups(ctx context.Context, userID uint) ([]models.Backup, error) {
	backups, err := s.backupRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch backups: %w", err)
	}

	customLogger.Debugf(ctx, "[BACKUP] Found %d backups for user %d", len(backups), userID)

	// Verify file existence and update sizes
	for i := range backups {
		if backups[i].Status == "completed" {
			customLogger.Debugf(ctx, "[BACKUP] Checking file for backup ID=%d at path: %s", backups[i].ID, backups[i].FilePath)
			if fileInfo, err := os.Stat(backups[i].FilePath); err == nil {
				backups[i].FileSize = fileInfo.Size()
				customLogger.Debugf(ctx, "[BACKUP] File exists, size: %d bytes", fileInfo.Size())
			} else {
				customLogger.Warnf(ctx, "[BACKUP] File not found for backup ID=%d: %v", backups[i].ID, err)
			}
		}
	}

	return backups, nil
}

// GetBackup retrieves a specific backup
func (s *backupService) GetBackup(ctx context.Context, userID uint, backupID uint) (*models.Backup, error) {
	backup, err := s.backupRepo.FindByID(ctx, backupID, userID)
	if err != nil {
		return nil, fmt.Errorf("backup not found: %w", err)
	}

	return backup, nil
}

// DeleteBackup deletes a backup file and its record
func (s *backupService) DeleteBackup(ctx context.Context, userID uint, backupID uint) error {
	backup, err := s.GetBackup(ctx, userID, backupID)
	if err != nil {
		return err
	}

	// Delete the file from filesystem
	if err := os.Remove(backup.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete backup file: %w", err)
	}

	// Delete the database record
	if err := s.backupRepo.Delete(ctx, backupID, userID); err != nil {
		return fmt.Errorf("failed to delete backup record: %w", err)
	}

	// Log activity
	s.activityLogger.LogActivity(ctx, ActivityLogParams{
		UserID:      userID,
		Action:      "delete",
		Module:      "backup",
		EntityType:  "database",
		Description: fmt.Sprintf("Deleted database backup: %s", backup.FileName),
	})

	return nil
}

// GetBackupFilePath returns the full path to a backup file
func (s *backupService) GetBackupFilePath(backup *models.Backup) string {
	// If the path is already absolute, return it as is
	if filepath.IsAbs(backup.FilePath) {
		return backup.FilePath
	}

	// If it's a relative path, convert it to absolute
	// This handles legacy backups created before the path fix
	cwd, err := os.Getwd()
	if err != nil {
		// Fallback to the stored path if we can't get cwd
		return backup.FilePath
	}

	return filepath.Join(cwd, backup.FilePath)
}

// getBackupDirectory returns the directory path for storing backups
func (s *backupService) getBackupDirectory() string {
	// Check if custom backup directory is set in environment
	if dir := os.Getenv("BACKUP_DIR"); dir != "" {
		return dir
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		// Fallback to relative path if we can't get cwd
		return "./backups"
	}

	// Return absolute path to backups directory
	return filepath.Join(cwd, "backups")
}

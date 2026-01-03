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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Find pg_dump executable - check common locations for macOS Homebrew
	pgDumpPath, err := s.findPgDump()
	if err != nil {
		errorMsg := fmt.Sprintf("pg_dump command not found: %v. Please install PostgreSQL client tools.", err)
		s.backupRepo.UpdateStatus(context.Background(), backupID, "failed", errorMsg)
		fmt.Printf("[BACKUP ERROR] %s\n", errorMsg)
		return
	}
	fmt.Printf("[BACKUP] Found pg_dump at: %s\n", pgDumpPath)

	// Create the backup file
	file, err := os.Create(filePath)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to create file: %v", err)
		s.backupRepo.UpdateStatus(context.Background(), backupID, "failed", errorMsg)
		fmt.Printf("[BACKUP ERROR] %s\n", errorMsg)
		return
	}
	defer file.Close()

	// Prepare pg_dump command
	dbConfig := s.config.Database

	// Log the backup attempt (don't log password)
	fmt.Printf("[BACKUP] Starting backup ID=%d, connecting to %s:%s/%s as %s\n",
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
	fmt.Printf("[BACKUP] Executing pg_dump for backup ID=%d (timeout: 5 minutes)\n", backupID)
	startTime := time.Now()

	if err := pgDumpCmd.Run(); err != nil {
		os.Remove(filePath) // Clean up failed backup file
		stderrOutput := stderrBuf.String()

		// Check if it was a timeout
		if ctx.Err() == context.DeadlineExceeded {
			errorMsg := "Backup timeout: operation took longer than 5 minutes"
			s.backupRepo.UpdateStatus(context.Background(), backupID, "failed", errorMsg)
			fmt.Printf("[BACKUP ERROR] %s\n", errorMsg)
			return
		}

		errorMsg := fmt.Sprintf("pg_dump failed: %v", err)
		if stderrOutput != "" {
			errorMsg = fmt.Sprintf("pg_dump failed: %v - %s", err, stderrOutput)
		}
		s.backupRepo.UpdateStatus(context.Background(), backupID, "failed", errorMsg)
		fmt.Printf("[BACKUP ERROR] %s\n", errorMsg)
		return
	}

	duration := time.Since(startTime)
	fmt.Printf("[BACKUP] pg_dump completed in %v\n", duration)

	// Get file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		errorMsg := fmt.Sprintf("Failed to get file info: %v", err)
		s.backupRepo.UpdateStatus(context.Background(), backupID, "failed", errorMsg)
		fmt.Printf("[BACKUP ERROR] %s\n", errorMsg)
		return
	}

	fmt.Printf("[BACKUP] Backup ID=%d completed successfully, size=%d bytes (%.2f MB)\n",
		backupID, fileInfo.Size(), float64(fileInfo.Size())/(1024*1024))

	// Get the backup to update
	backup, err := s.backupRepo.FindByID(context.Background(), backupID, 0) // userID=0 for system operation
	if err == nil {
		backup.FileSize = fileInfo.Size()
		backup.Status = "completed"
		s.backupRepo.Update(context.Background(), backup)
		fmt.Printf("[BACKUP] Backup ID=%d status updated to completed\n", backupID)
	} else {
		fmt.Printf("[BACKUP ERROR] Failed to update backup ID=%d: %v\n", backupID, err)
	}
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

	// Verify file existence and update sizes
	for i := range backups {
		if backups[i].Status == "completed" {
			if fileInfo, err := os.Stat(backups[i].FilePath); err == nil {
				backups[i].FileSize = fileInfo.Size()
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
	return backup.FilePath
}

// getBackupDirectory returns the directory path for storing backups
func (s *backupService) getBackupDirectory() string {
	// Check if custom backup directory is set in environment
	if dir := os.Getenv("BACKUP_DIR"); dir != "" {
		return dir
	}
	// Default to ./backups in the current working directory
	return "./backups"
}

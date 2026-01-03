package services

import (
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
	ctx := context.Background()

	// Create the backup file
	file, err := os.Create(filePath)
	if err != nil {
		s.backupRepo.UpdateStatus(ctx, backupID, "failed", fmt.Sprintf("Failed to create file: %v", err))
		return
	}
	defer file.Close()

	// Prepare pg_dump command
	dbConfig := s.config.Database
	pgDumpCmd := exec.Command(
		"pg_dump",
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

	// Set PGPASSWORD environment variable
	pgDumpCmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", dbConfig.Password))
	pgDumpCmd.Stdout = file
	pgDumpCmd.Stderr = os.Stderr

	// Execute pg_dump
	if err := pgDumpCmd.Run(); err != nil {
		os.Remove(filePath) // Clean up failed backup file
		s.backupRepo.UpdateStatus(ctx, backupID, "failed", fmt.Sprintf("pg_dump failed: %v", err))
		return
	}

	// Get file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		s.backupRepo.UpdateStatus(ctx, backupID, "failed", fmt.Sprintf("Failed to get file info: %v", err))
		return
	}

	// Get the backup to update
	backup, err := s.backupRepo.FindByID(ctx, backupID, 0) // userID=0 for system operation
	if err == nil {
		backup.FileSize = fileInfo.Size()
		backup.Status = "completed"
		s.backupRepo.Update(ctx, backup)
	}
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

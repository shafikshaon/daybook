package repository

import (
	"context"

	"daybook-backend/models"

	"gorm.io/gorm"
)

// BackupRepository handles backup data access
type BackupRepository interface {
	BaseRepository[models.Backup]

	// FindByUserID retrieves all backups for a user
	FindByUserID(ctx context.Context, userID uint) ([]models.Backup, error)

	// FindByFileName retrieves a backup by filename and user ID
	FindByFileName(ctx context.Context, userID uint, fileName string) (*models.Backup, error)

	// UpdateStatus updates the status of a backup
	UpdateStatus(ctx context.Context, id uint, status string, errorMsg string) error
}

type backupRepository struct {
	*GormBaseRepository[models.Backup]
}

// NewBackupRepository creates a new backup repository
func NewBackupRepository(db *gorm.DB) BackupRepository {
	return &backupRepository{
		GormBaseRepository: NewGormBaseRepository[models.Backup](db),
	}
}

// FindByUserID retrieves all backups for a user
func (r *backupRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Backup, error) {
	var backups []models.Backup
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&backups).Error
	return backups, err
}

// FindByFileName retrieves a backup by filename and user ID
func (r *backupRepository) FindByFileName(ctx context.Context, userID uint, fileName string) (*models.Backup, error) {
	var backup models.Backup
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND file_name = ?", userID, fileName).
		First(&backup).Error
	if err != nil {
		return nil, err
	}
	return &backup, nil
}

// UpdateStatus updates the status of a backup
func (r *backupRepository) UpdateStatus(ctx context.Context, id uint, status string, errorMsg string) error {
	return r.db.WithContext(ctx).Model(&models.Backup{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        status,
			"error_message": errorMsg,
		}).Error
}

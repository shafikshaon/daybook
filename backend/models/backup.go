package models

import (
	"time"

	"gorm.io/gorm"
)

// Backup represents a database backup with metadata
type Backup struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint           `gorm:"not null;index" json:"userId"`
	FileName     string         `gorm:"not null" json:"fileName"`         // e.g., backup_20240103_152030.sql
	FilePath     string         `gorm:"not null" json:"filePath"`         // Full path to backup file
	FileSize     int64          `gorm:"not null" json:"fileSize"`         // Size in bytes
	Status       string         `gorm:"not null;index" json:"status"`     // pending, completed, failed
	ErrorMessage string         `json:"errorMessage,omitempty"`           // Error message if failed
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b *Backup) BeforeCreate(tx *gorm.DB) error {
	// Set default status to pending if not set
	if b.Status == "" {
		b.Status = "pending"
	}
	return nil
}

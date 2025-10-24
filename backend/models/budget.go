package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Budget struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	CategoryID      string         `gorm:"not null;index" json:"categoryId" binding:"required"`
	Amount          float64        `gorm:"not null" json:"amount" binding:"required,gt=0"`
	Period          string         `gorm:"not null" json:"period" binding:"required"` // weekly, monthly, quarterly, yearly, custom
	CustomStartDate *time.Time     `json:"customStartDate"`
	CustomEndDate   *time.Time     `json:"customEndDate"`
	Rollover        bool           `gorm:"default:false" json:"rollover"`    // Rollover unused budget to next period
	AlertThreshold  float64        `gorm:"default:80" json:"alertThreshold"` // Alert when % of budget is reached
	Enabled         bool           `gorm:"default:true" json:"enabled"`
	Notes           string         `json:"notes"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b *Budget) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

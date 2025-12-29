package models

import (
	"time"

	"gorm.io/gorm"
)

type Budget struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint           `gorm:"not null;index" json:"userId"`
	CategoryID      uint           `gorm:"not null;index" json:"categoryId" binding:"required"`
	Amount          float64        `gorm:"not null" json:"amount" binding:"required,gt=0"`
	Period          string         `gorm:"not null" json:"period" binding:"required"` // weekly, monthly, quarterly, yearly, custom
	CustomStartDate *time.Time     `gorm:"type:timestamptz" json:"customStartDate"`
	CustomEndDate   *time.Time     `gorm:"type:timestamptz" json:"customEndDate"`
	Rollover        bool           `gorm:"default:false" json:"rollover"`    // Rollover unused budget to next period
	AlertThreshold  float64        `gorm:"default:80" json:"alertThreshold"` // Alert when % of budget is reached
	Enabled         bool           `gorm:"default:true" json:"enabled"`
	Notes           string         `json:"notes"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b *Budget) BeforeCreate(tx *gorm.DB) error {
	return nil
}

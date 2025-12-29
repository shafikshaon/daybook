package models

import (
	"time"

	"gorm.io/gorm"
)

type Settings struct {
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint           `gorm:"uniqueIndex;not null" json:"userId"`
	Currency       string         `gorm:"default:'BDT'" json:"currency"`
	DarkMode       bool           `gorm:"default:false" json:"darkMode"`
	DateFormat     string         `gorm:"default:'MM/DD/YYYY'" json:"dateFormat"`
	FirstDayOfWeek int            `gorm:"default:0" json:"firstDayOfWeek"` // 0 = Sunday
	Language       string         `gorm:"default:'en'" json:"language"`
	Notifications  *Notifications `gorm:"embedded;embeddedPrefix:notif_" json:"notifications"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type Notifications struct {
	Push         bool `json:"push"`
	Email        bool `json:"email"`
	BudgetAlerts bool `json:"budgetAlerts"`
}

func (s *Settings) BeforeCreate(tx *gorm.DB) error {
	return nil
}

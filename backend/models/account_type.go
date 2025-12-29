package models

import (
	"time"

	"gorm.io/gorm"
)

type AccountType struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"userId"` // All types are user-specific
	Name        string         `gorm:"not null" json:"name" binding:"required"`
	Icon        string         `json:"icon"`
	Description string         `json:"description"`
	Active      bool           `gorm:"default:true" json:"active"`
	SortOrder   int            `gorm:"default:0" json:"sortOrder"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (at *AccountType) BeforeCreate(tx *gorm.DB) error {
	return nil
}

// SeedDefaultAccountTypes creates default account types for a new user
func SeedDefaultAccountTypes(tx *gorm.DB, userID uint) error {
	// Define default types
	defaultTypes := []AccountType{
		{
			UserID:      userID,
			Name:        "Cash",
			Icon:        "💵",
			Description: "Physical cash",
			Active:      true,
			SortOrder:   1,
		},
		{
			UserID:      userID,
			Name:        "Bank",
			Icon:        "🏦",
			Description: "Bank accounts",
			Active:      true,
			SortOrder:   2,
		},
		{
			UserID:      userID,
			Name:        "Digital Wallet",
			Icon:        "📱",
			Description: "Digital payment services",
			Active:      true,
			SortOrder:   3,
		},
		{
			UserID:      userID,
			Name:        "Other",
			Icon:        "📋",
			Description: "Other account types",
			Active:      true,
			SortOrder:   4,
		},
	}

	// Create all account types
	for i := range defaultTypes {
		if err := tx.Create(&defaultTypes[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

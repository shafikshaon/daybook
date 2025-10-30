package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountType struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"` // All types are user-specific
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
	if at.ID == uuid.Nil {
		at.ID = uuid.New()
	}
	return nil
}

// SeedDefaultAccountTypes creates default account types for a new user
func SeedDefaultAccountTypes(tx *gorm.DB, userID uuid.UUID) error {
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

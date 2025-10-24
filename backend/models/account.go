package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID                       uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID                   uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	Name                     string         `gorm:"not null" json:"name" binding:"required"`
	Type                     string         `gorm:"not null" json:"type" binding:"required"` // cash, checking, savings, credit_card, etc
	Balance                  float64        `gorm:"default:0" json:"balance"`
	Currency                 string         `gorm:"default:'USD'" json:"currency"`
	Description              string         `json:"description"`
	Institution              string         `json:"institution"`
	AccountNumber            string         `json:"accountNumber"`
	LastReconciled           *time.Time     `json:"lastReconciled"`
	ReconciliationDifference float64        `gorm:"default:0" json:"reconciliationDifference"`
	Active                   bool           `gorm:"default:true" json:"active"`
	CreatedAt                time.Time      `json:"createdAt"`
	UpdatedAt                time.Time      `json:"updatedAt"`
	DeletedAt                gorm.DeletedAt `gorm:"index" json:"-"`
}

func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

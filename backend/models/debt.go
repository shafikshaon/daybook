package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DebtRecord represents money borrowed from someone
type DebtRecord struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	CreditorName    string         `gorm:"not null" json:"creditorName" binding:"required"` // Person/entity we owe
	OriginalAmount  float64        `gorm:"not null" json:"originalAmount" binding:"required,gt=0"`
	RemainingAmount float64        `gorm:"not null" json:"remainingAmount"`
	AccountID       *uuid.UUID     `gorm:"type:uuid;index" json:"accountId"`        // Account affected (if any)
	Status          string         `gorm:"not null;index" json:"status"`            // active, partially_paid, fully_paid
	BorrowedDate    Date           `gorm:"not null;index" json:"borrowedDate"`
	DueDate         *time.Time     `json:"dueDate"`
	InterestRate    *float64       `json:"interestRate"` // Annual interest rate in percentage
	Description     string         `json:"description"`
	IsInitial       bool           `gorm:"default:false" json:"isInitial"` // True if this is a pre-existing debt
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *DebtRecord) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	// Set remaining amount to original amount if not set
	if d.RemainingAmount == 0 {
		d.RemainingAmount = d.OriginalAmount
	}
	// Set status to active if not set
	if d.Status == "" {
		d.Status = "active"
	}
	return nil
}

// DebtPayment represents a payment made towards a debt
type DebtPayment struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	DebtID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"debtId"`
	AccountID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"accountId"`      // Account from which payment is made
	Amount      float64        `gorm:"not null" json:"amount" binding:"required,gt=0"`
	PaymentDate Date           `gorm:"not null;index" json:"paymentDate"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (dp *DebtPayment) BeforeCreate(tx *gorm.DB) error {
	if dp.ID == uuid.Nil {
		dp.ID = uuid.New()
	}
	return nil
}

package models

import (
	"time"

	"gorm.io/gorm"
)

// LendRecord represents money lent to someone
type LendRecord struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint           `gorm:"not null;index" json:"userId"`
	DebtorName      string         `gorm:"not null" json:"debtorName" binding:"required"` // Person/entity who owes us
	OriginalAmount  float64        `gorm:"not null" json:"originalAmount" binding:"required,gt=0"`
	RemainingAmount float64        `gorm:"not null" json:"remainingAmount"`
	AccountID       *uint          `gorm:"index" json:"accountId"`       // Account affected (if any)
	Status          string         `gorm:"not null;index" json:"status"` // active, partially_received, fully_received
	LentDate        Date           `gorm:"not null;index" json:"lentDate"`
	DueDate         *Date          `gorm:"type:timestamptz" json:"dueDate"`
	InterestRate    *float64       `json:"interestRate"` // Annual interest rate in percentage
	Description     string         `json:"description"`
	IsInitial       bool           `gorm:"default:false" json:"isInitial"` // True if this is a pre-existing lend
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (l *LendRecord) BeforeCreate(tx *gorm.DB) error {
	// Set remaining amount to original amount if not set
	if l.RemainingAmount == 0 {
		l.RemainingAmount = l.OriginalAmount
	}
	// Set status to active if not set
	if l.Status == "" {
		l.Status = "active"
	}
	return nil
}

// LendPayment represents a payment received for a lend
type LendPayment struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"userId"`
	LendID      uint           `gorm:"not null;index" json:"lendId"`
	AccountID   uint           `gorm:"not null;index" json:"accountId"` // Account to which payment is received
	Amount      float64        `gorm:"not null" json:"amount" binding:"required,gt=0"`
	PaymentDate Date           `gorm:"not null;index" json:"paymentDate"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (lp *LendPayment) BeforeCreate(tx *gorm.DB) error {
	return nil
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bill struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	Name           string         `gorm:"not null" json:"name" binding:"required"`
	Category       string         `gorm:"not null" json:"category" binding:"required"` // utilities, subscriptions, insurance, etc
	Amount         float64        `gorm:"not null" json:"amount" binding:"required,gt=0"`
	Frequency      string         `gorm:"not null" json:"frequency" binding:"required"` // weekly, biweekly, monthly, quarterly, etc
	StartDate      time.Time      `gorm:"not null" json:"startDate"`
	DueDay         int            `json:"dueDay"` // Day of month for monthly bills
	LastPaidDate   *time.Time     `json:"lastPaidDate"`
	LastPaidAmount float64        `gorm:"default:0" json:"lastPaidAmount"`
	AutoPay        bool           `gorm:"default:false" json:"autoPay"`
	ReminderDays   int            `gorm:"default:3" json:"reminderDays"` // Days before due date to remind
	Active         bool           `gorm:"default:true" json:"active"`
	Notes          string         `json:"notes"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (b *Bill) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type BillPayment struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	BillID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"billId"`
	Amount      float64        `gorm:"not null" json:"amount" binding:"required,gt=0"`
	PaymentDate time.Time      `gorm:"not null;index" json:"paymentDate"`
	AccountID   *uuid.UUID     `gorm:"type:uuid" json:"accountId"`
	Notes       string         `json:"notes"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (bp *BillPayment) BeforeCreate(tx *gorm.DB) error {
	if bp.ID == uuid.Nil {
		bp.ID = uuid.New()
	}
	return nil
}

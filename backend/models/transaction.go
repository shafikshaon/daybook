package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	AccountID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"accountId"`
	ToAccountID      *uuid.UUID     `gorm:"type:uuid;index" json:"toAccountId"`      // For transfers
	Type             string         `gorm:"not null" json:"type" binding:"required"` // income, expense, transfer
	Amount           float64        `gorm:"not null" json:"amount" binding:"required,gt=0"`
	CategoryID       string         `gorm:"not null;index" json:"categoryId" binding:"required"`
	Date             time.Time      `gorm:"not null;index" json:"date" binding:"required"`
	Description      string         `json:"description"`
	Tags             []string       `gorm:"type:jsonb;serializer:json" json:"tags"`
	SavingsGoalID    *uuid.UUID     `gorm:"type:uuid" json:"savingsGoalId"`
	RecurringID      *uuid.UUID     `gorm:"type:uuid" json:"recurringId"`
	CreditCardID     *uuid.UUID     `gorm:"type:uuid" json:"creditCardId"`
	Attachments      []string       `gorm:"type:jsonb;serializer:json" json:"attachments"`
	Reconciled       bool           `gorm:"default:false;index" json:"reconciled"`
	ReconciliationID *uuid.UUID     `gorm:"type:uuid" json:"reconciliationId"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type RecurringTransaction struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID              uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	TransactionTemplate Transaction    `gorm:"embedded;embeddedPrefix:template_" json:"transactionTemplate"`
	Frequency           string         `gorm:"not null" json:"frequency"` // daily, weekly, biweekly, monthly, quarterly, yearly
	StartDate           time.Time      `gorm:"not null" json:"startDate"`
	EndDate             *time.Time     `json:"endDate"`
	LastProcessed       *time.Time     `json:"lastProcessed"`
	Enabled             bool           `gorm:"default:true" json:"enabled"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

func (rt *RecurringTransaction) BeforeCreate(tx *gorm.DB) error {
	if rt.ID == uuid.Nil {
		rt.ID = uuid.New()
	}
	return nil
}

type Tag struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	Name      string         `gorm:"not null" json:"name" binding:"required"`
	Color     string         `json:"color"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

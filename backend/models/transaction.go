package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID               uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           uint           `gorm:"not null;index" json:"userId"`
	AccountID        uint           `gorm:"not null;index" json:"accountId"`
	ToAccountID      *uint          `gorm:"index" json:"toAccountId"`                // For transfers
	Type             string         `gorm:"not null" json:"type" binding:"required"` // income, expense, transfer
	Amount           float64        `gorm:"not null" json:"amount" binding:"required,gt=0"`
	CategoryID       uint           `gorm:"default:0;index" json:"categoryId"`
	Date             Date           `gorm:"not null;index" json:"date"`
	Description      string         `json:"description"`
	Tags             []string       `gorm:"type:jsonb;serializer:json" json:"tags"`
	SavingsGoalID    *uint          `gorm:"index" json:"savingsGoalId"`
	FixedDepositID   *uint          `gorm:"index" json:"fixedDepositId"`
	InvestmentID     *uint          `gorm:"index" json:"investmentId"`
	RecurringID      *uint          `json:"recurringId"`
	CreditCardID     *uint          `json:"creditCardId"`
	Attachments      []string       `gorm:"type:jsonb;serializer:json" json:"attachments"`
	Reconciled       bool           `gorm:"default:false;index" json:"reconciled"`
	ReconciliationID *uint          `json:"reconciliationId"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	return nil
}

type RecurringTransaction struct {
	ID                  uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID              uint           `gorm:"not null;index" json:"userId"`
	TransactionTemplate Transaction    `gorm:"embedded;embeddedPrefix:template_" json:"transactionTemplate"`
	Frequency           string         `gorm:"not null" json:"frequency"` // daily, weekly, biweekly, monthly, quarterly, yearly
	StartDate           Date           `gorm:"not null" json:"startDate"`
	EndDate             *Date          `json:"endDate"`
	LastProcessed       *time.Time     `gorm:"type:timestamptz" json:"lastProcessed"`
	Enabled             bool           `gorm:"default:true" json:"enabled"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

func (rt *RecurringTransaction) BeforeCreate(tx *gorm.DB) error {
	return nil
}

type Tag struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"userId"`
	Name      string         `gorm:"not null" json:"name" binding:"required"`
	Color     string         `json:"color"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	return nil
}

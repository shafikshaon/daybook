package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReconciliationStatus represents the status of a reconciliation
type ReconciliationStatus string

const (
	ReconciliationPending     ReconciliationStatus = "pending"
	ReconciliationCompleted   ReconciliationStatus = "completed"
	ReconciliationDiscrepancy ReconciliationStatus = "discrepancy"
)

// Reconciliation represents an account reconciliation record
type Reconciliation struct {
	ID                 uuid.UUID            `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID             uuid.UUID            `gorm:"type:uuid;not null;index" json:"userId"`
	AccountID          uuid.UUID            `gorm:"type:uuid;not null;index" json:"accountId"`
	ReconciliationDate time.Time            `gorm:"not null" json:"reconciliationDate" binding:"required"`
	StatementBalance   float64              `gorm:"type:decimal(15,2);not null" json:"statementBalance" binding:"required"`
	BookBalance        float64              `gorm:"type:decimal(15,2);not null" json:"bookBalance"`
	Difference         float64              `gorm:"type:decimal(15,2);not null" json:"difference"`
	Notes              string               `gorm:"type:text" json:"notes"`
	Status             ReconciliationStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	CreatedAt          time.Time            `json:"createdAt"`
	UpdatedAt          time.Time            `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt       `gorm:"index" json:"-"`

	// Relationships
	Account      Account                     `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Transactions []ReconciliationTransaction `gorm:"foreignKey:ReconciliationID" json:"transactions,omitempty"`
}

// ReconciliationTransaction links transactions to reconciliation records
type ReconciliationTransaction struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ReconciliationID uuid.UUID `gorm:"type:uuid;not null;index" json:"reconciliationId"`
	TransactionID    uuid.UUID `gorm:"type:uuid;not null;index" json:"transactionId"`
	CreatedAt        time.Time `json:"createdAt"`

	// Relationships
	Reconciliation Reconciliation `gorm:"foreignKey:ReconciliationID" json:"-"`
	Transaction    Transaction    `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`
}

// BeforeCreate sets the UUID and calculates difference
func (r *Reconciliation) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}

	// Calculate difference between statement and book balance
	r.Difference = r.StatementBalance - r.BookBalance

	// Determine status based on difference
	if r.Difference == 0 {
		r.Status = ReconciliationCompleted
	} else {
		r.Status = ReconciliationDiscrepancy
	}

	return nil
}

// BeforeUpdate recalculates difference
func (r *Reconciliation) BeforeUpdate(tx *gorm.DB) error {
	// Recalculate difference
	r.Difference = r.StatementBalance - r.BookBalance

	// Update status based on difference
	if r.Difference == 0 {
		r.Status = ReconciliationCompleted
	} else if r.Status != ReconciliationPending {
		r.Status = ReconciliationDiscrepancy
	}

	return nil
}

// BeforeCreate sets the UUID for ReconciliationTransaction
func (rt *ReconciliationTransaction) BeforeCreate(tx *gorm.DB) error {
	if rt.ID == uuid.Nil {
		rt.ID = uuid.New()
	}
	return nil
}

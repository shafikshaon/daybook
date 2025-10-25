package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FixedDeposit struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID               uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	Institution          string         `gorm:"not null" json:"institution" binding:"required"`
	AccountNumber        string         `json:"accountNumber"`
	Principal            float64        `gorm:"not null" json:"principal" binding:"required,gt=0"`
	InterestRate         float64        `gorm:"not null" json:"interestRate" binding:"required,gt=0"`
	TenureMonths         int            `gorm:"not null" json:"tenureMonths" binding:"required,gt=0"`
	Compounding          string         `gorm:"not null;default:'monthly'" json:"compounding"` // simple, daily, monthly, quarterly, semi-annually, annually
	StartDate            time.Time      `gorm:"not null" json:"startDate"`
	MaturityDate         time.Time      `gorm:"not null" json:"maturityDate"`
	MaturityAmount       float64        `json:"maturityAmount"`
	ActualMaturityAmount float64        `json:"actualMaturityAmount"`
	Withdrawn            bool           `gorm:"default:false" json:"withdrawn"`
	WithdrawnDate        *time.Time     `json:"withdrawnDate"`
	AutoRenew            bool           `gorm:"default:false" json:"autoRenew"`
	Notes                string         `json:"notes"`
	Attachments          []string       `gorm:"type:text[]" json:"attachments"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

func (fd *FixedDeposit) BeforeCreate(tx *gorm.DB) error {
	if fd.ID == uuid.Nil {
		fd.ID = uuid.New()
	}
	return nil
}

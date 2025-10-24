package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Investment struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	PortfolioID      *uuid.UUID     `gorm:"type:uuid;index" json:"portfolioId"`
	Symbol           string         `gorm:"not null" json:"symbol" binding:"required"`
	Name             string         `gorm:"not null" json:"name" binding:"required"`
	AssetType        string         `gorm:"not null" json:"assetType" binding:"required"` // stocks, bonds, mutual_funds, etf, crypto, etc
	Quantity         float64        `gorm:"not null" json:"quantity" binding:"required,gt=0"`
	CostBasis        float64        `gorm:"not null" json:"costBasis" binding:"required,gt=0"`
	CurrentPrice     float64        `gorm:"not null" json:"currentPrice" binding:"required,gt=0"`
	PurchaseDate     time.Time      `json:"purchaseDate"`
	LastUpdated      time.Time      `json:"lastUpdated"`
	RealizedGainLoss float64        `gorm:"default:0" json:"realizedGainLoss"`
	Notes            string         `json:"notes"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (i *Investment) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}

type Portfolio struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	Name        string         `gorm:"not null" json:"name" binding:"required"`
	Description string         `json:"description"`
	Type        string         `json:"type"` // retirement, taxable, education, etc
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *Portfolio) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type Dividend struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	InvestmentID uuid.UUID      `gorm:"type:uuid;not null;index" json:"investmentId"`
	Amount       float64        `gorm:"not null" json:"amount" binding:"required,gt=0"`
	PaymentDate  time.Time      `gorm:"not null" json:"paymentDate"`
	Reinvested   bool           `gorm:"default:false" json:"reinvested"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *Dividend) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

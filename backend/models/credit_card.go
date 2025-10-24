package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreditCard struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID            uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	Name              string         `gorm:"not null" json:"name" binding:"required"`
	LastFourDigits    string         `json:"lastFourDigits"`
	CardNetwork       string         `json:"cardNetwork"` // Visa, Mastercard, Amex, etc
	CreditLimit       float64        `gorm:"not null" json:"creditLimit" binding:"required,gt=0"`
	CurrentBalance    float64        `gorm:"default:0" json:"currentBalance"`
	APR               float64        `json:"apr"`
	DueDate           *time.Time     `json:"dueDate"`
	StatementDate     *time.Time     `json:"statementDate"`
	MinimumPayment    float64        `gorm:"default:0" json:"minimumPayment"`
	LastPaymentDate   *time.Time     `json:"lastPaymentDate"`
	LastPaymentAmount float64        `gorm:"default:0" json:"lastPaymentAmount"`
	RewardsProgram    string         `json:"rewardsProgram"`
	Active            bool           `gorm:"default:true" json:"active"`
	Notes             string         `json:"notes"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (cc *CreditCard) BeforeCreate(tx *gorm.DB) error {
	if cc.ID == uuid.Nil {
		cc.ID = uuid.New()
	}
	return nil
}

type Statement struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	CardID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"cardId"`
	StatementDate   time.Time      `gorm:"not null" json:"statementDate"`
	DueDate         time.Time      `gorm:"not null" json:"dueDate"`
	OpeningBalance  float64        `json:"openingBalance"`
	ClosingBalance  float64        `json:"closingBalance"`
	MinimumPayment  float64        `json:"minimumPayment"`
	TotalCharges    float64        `json:"totalCharges"`
	TotalPayments   float64        `json:"totalPayments"`
	InterestCharged float64        `json:"interestCharged"`
	Paid            bool           `gorm:"default:false" json:"paid"`
	PaidDate        *time.Time     `json:"paidDate"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Statement) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type Reward struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	CardID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"cardId"`
	Type        string         `json:"type"` // cashback, points, miles
	Amount      float64        `json:"amount"`
	Description string         `json:"description"`
	EarnedDate  time.Time      `gorm:"not null" json:"earnedDate"`
	Redeemed    bool           `gorm:"default:false" json:"redeemed"`
	RedeemedAt  *time.Time     `json:"redeemedAt"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *Reward) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

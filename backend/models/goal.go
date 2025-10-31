package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Goal represents a financial objective (replaces SavingsGoal)
type Goal struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`

	// Basic Info
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`

	// Category & Priority
	Category string `json:"category"` // emergency_fund, vacation, retirement, home, education, car, wedding, business, other
	Priority string `json:"priority"` // high, medium, low

	// Financial Targets
	TargetAmount        float64    `gorm:"not null" json:"targetAmount"`
	CurrentAmount       float64    `gorm:"default:0" json:"currentAmount"`
	TargetDate          *time.Time `json:"targetDate"`
	MonthlyContribution float64    `json:"monthlyContribution"`

	// Status
	Status       string     `gorm:"default:active" json:"status"` // active, achieved, paused, archived
	Achieved     bool       `gorm:"default:false" json:"achieved"`
	AchievedDate *time.Time `json:"achievedDate"`

	// Tracking
	LastContribution     float64    `json:"lastContribution"`
	LastContributionDate *time.Time `json:"lastContributionDate"`

	// Relationships
	Holdings      []GoalHolding      `gorm:"foreignKey:GoalID" json:"holdings,omitempty"`
	Contributions []GoalContribution `gorm:"foreignKey:GoalID" json:"contributions,omitempty"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// GoalHolding represents an investment/savings instrument toward a goal
// This replaces: Investment, FixedDeposit, and individual savings contributions
type GoalHolding struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	GoalID uuid.UUID `gorm:"type:uuid;not null;index" json:"goalId"`

	// Core Fields
	Name         string    `gorm:"not null" json:"name"`
	Type         string    `gorm:"not null;index" json:"type"`   // savings, fixed_deposit, dps, recurring_deposit, stocks, mutual_fund, etf, bonds, crypto, real_estate, gold, pension_fund, ulip, ppf, nsc, custom
	Status       string    `gorm:"default:active" json:"status"` // active, matured, sold, closed, withdrawn
	PurchaseDate time.Time `gorm:"not null" json:"purchaseDate"`
	Amount       float64   `gorm:"not null" json:"amount"` // Initial investment amount
	CurrentValue float64   `json:"currentValue"`           // Current market value

	// Common Fields (for bank products)
	Institution    *string    `json:"institution"`    // Bank/Fund house name
	AccountNumber  *string    `json:"accountNumber"`  // Account/Policy number
	InterestRate   *float64   `json:"interestRate"`   // Annual interest rate
	MaturityDate   *time.Time `json:"maturityDate"`   // When it matures
	MaturityAmount *float64   `json:"maturityAmount"` // Expected maturity value
	TenureMonths   *int       `json:"tenureMonths"`   // Duration in months

	// For market instruments (stocks, mutual funds, ETF, crypto)
	Symbol       *string  `json:"symbol"`       // Ticker symbol (AAPL, VTSAX, BTC)
	Quantity     *float64 `json:"quantity"`     // Number of shares/units
	CostBasis    *float64 `json:"costBasis"`    // Price per unit when purchased
	CurrentPrice *float64 `json:"currentPrice"` // Current market price per unit

	// For DPS/Recurring Deposits
	MonthlyDeposit *float64 `json:"monthlyDeposit"` // Monthly contribution amount

	// Additional metadata stored as JSON
	Details map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"details"`

	// Transaction link (the transaction that created this holding)
	TransactionID uuid.UUID `gorm:"type:uuid" json:"transactionId"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// GoalContribution tracks all contributions/withdrawals/gains for a goal
type GoalContribution struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"userId"`
	GoalID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"goalId"`
	HoldingID *uuid.UUID `gorm:"type:uuid" json:"holdingId"` // Link to specific holding

	Type   string    `gorm:"not null" json:"type"` // contribution, withdrawal, dividend, interest, appreciation, depreciation, maturity
	Amount float64   `gorm:"not null" json:"amount"`
	Date   time.Time `gorm:"not null;index" json:"date"`
	Notes  string    `json:"notes"`

	// Transaction link (creates a transaction for account balance)
	TransactionID uuid.UUID `gorm:"type:uuid;not null" json:"transactionId"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hooks
func (g *Goal) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

func (h *GoalHolding) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

func (c *GoalContribution) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// Helper methods

// CalculateProgress returns the completion percentage
func (g *Goal) CalculateProgress() float64 {
	if g.TargetAmount <= 0 {
		return 0
	}
	progress := (g.CurrentAmount / g.TargetAmount) * 100
	if progress > 100 {
		return 100
	}
	return progress
}

// UpdateCurrentAmount recalculates current amount from all holdings
func (g *Goal) UpdateCurrentAmount(db *gorm.DB) error {
	var total float64
	db.Model(&GoalHolding{}).
		Where("goal_id = ? AND status IN ?", g.ID, []string{"active", "matured"}).
		Select("COALESCE(SUM(current_value), 0)").
		Scan(&total)

	g.CurrentAmount = total
	return db.Save(g).Error
}

// CalculateGainLoss returns the gain/loss for a holding
func (h *GoalHolding) CalculateGainLoss() (float64, float64) {
	gainLoss := h.CurrentValue - h.Amount
	gainLossPercent := float64(0)
	if h.Amount > 0 {
		gainLossPercent = (gainLoss / h.Amount) * 100
	}
	return gainLoss, gainLossPercent
}

// UpdateMarketValue updates current value for market instruments
func (h *GoalHolding) UpdateMarketValue() {
	if h.Quantity != nil && h.CurrentPrice != nil {
		h.CurrentValue = *h.Quantity * *h.CurrentPrice
	} else if h.CurrentValue == 0 {
		// If no market value set, use initial amount
		h.CurrentValue = h.Amount
	}
}

// Holding type constants
const (
	// Traditional Savings
	HoldingTypeSavings          = "savings"
	HoldingTypeFixedDeposit     = "fixed_deposit"
	HoldingTypeDPS              = "dps"
	HoldingTypeRecurringDeposit = "recurring_deposit"
	HoldingTypeSavingsBond      = "savings_bond"
	HoldingTypePPF              = "ppf"
	HoldingTypeNSC              = "nsc"

	// Market Investments
	HoldingTypeStocks     = "stocks"
	HoldingTypeMutualFund = "mutual_fund"
	HoldingTypeETF        = "etf"
	HoldingTypeIndexFund  = "index_fund"
	HoldingTypeBonds      = "bonds"
	HoldingTypeCrypto     = "cryptocurrency"

	// Alternative Investments
	HoldingTypeRealEstate  = "real_estate"
	HoldingTypeREIT        = "reit"
	HoldingTypeGold        = "gold"
	HoldingTypeCommodities = "commodities"

	// Retirement & Pension
	HoldingTypePensionFund    = "pension_fund"
	HoldingTypeRetirement401k = "retirement_401k"
	HoldingTypeProvidentFund  = "provident_fund"

	// Insurance
	HoldingTypeLifeInsurance = "life_insurance"
	HoldingTypeULIP          = "ulip"

	// Other
	HoldingTypeCustom = "custom"
)

// Goal categories
const (
	GoalCategoryEmergencyFund = "emergency_fund"
	GoalCategoryVacation      = "vacation"
	GoalCategoryRetirement    = "retirement"
	GoalCategoryHome          = "home"
	GoalCategoryEducation     = "education"
	GoalCategoryCar           = "car"
	GoalCategoryWedding       = "wedding"
	GoalCategoryBusiness      = "business"
	GoalCategoryOther         = "other"
)

// Goal priorities
const (
	GoalPriorityHigh   = "high"
	GoalPriorityMedium = "medium"
	GoalPriorityLow    = "low"
)

// Goal statuses
const (
	GoalStatusActive   = "active"
	GoalStatusAchieved = "achieved"
	GoalStatusPaused   = "paused"
	GoalStatusArchived = "archived"
)

// Holding statuses
const (
	HoldingStatusActive    = "active"
	HoldingStatusMatured   = "matured"
	HoldingStatusSold      = "sold"
	HoldingStatusClosed    = "closed"
	HoldingStatusWithdrawn = "withdrawn"
)

// Contribution types
const (
	ContributionTypeContribution = "contribution"
	ContributionTypeWithdrawal   = "withdrawal"
	ContributionTypeDividend     = "dividend"
	ContributionTypeInterest     = "interest"
	ContributionTypeAppreciation = "appreciation"
	ContributionTypeDepreciation = "depreciation"
	ContributionTypeMaturity     = "maturity"
)

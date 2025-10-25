package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SavingsGoal struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID               uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	Name                 string         `gorm:"not null" json:"name" binding:"required"`
	Description          string         `json:"description"`
	TargetAmount         float64        `gorm:"not null" json:"targetAmount" binding:"required,gt=0"`
	CurrentAmount        float64        `gorm:"default:0" json:"currentAmount"`
	TargetDate           *time.Time     `json:"targetDate"`
	MonthlyContribution  float64        `gorm:"default:0" json:"monthlyContribution"`
	Category             string         `json:"category"` // emergency, vacation, purchase, etc
	Priority             string         `json:"priority"` // high, medium, low
	Achieved             bool           `gorm:"default:false" json:"achieved"`
	AchievedDate         *time.Time     `json:"achievedDate"`
	Archived             bool           `gorm:"default:false" json:"archived"`
	ArchivedDate         *time.Time     `json:"archivedDate"`
	LastContribution     float64        `gorm:"default:0" json:"lastContribution"`
	LastContributionDate *time.Time     `json:"lastContributionDate"`
	Attachments          []string       `gorm:"type:text[]" json:"attachments"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

func (sg *SavingsGoal) BeforeCreate(tx *gorm.DB) error {
	if sg.ID == uuid.Nil {
		sg.ID = uuid.New()
	}
	return nil
}

type SavingsContribution struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	GoalID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"goalId"`
	Amount      float64        `gorm:"not null" json:"amount" binding:"required,gt=0"`
	Date        time.Time      `gorm:"not null;index" json:"date"`
	Notes       string         `json:"notes"`
	Attachments []string       `gorm:"type:text[]" json:"attachments"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (sc *SavingsContribution) BeforeCreate(tx *gorm.DB) error {
	if sc.ID == uuid.Nil {
		sc.ID = uuid.New()
	}
	return nil
}

type AutomatedRule struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	GoalID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"goalId"`
	RuleType   string         `gorm:"not null" json:"ruleType"` // percentage_of_income, fixed_amount, round_up
	Amount     float64        `json:"amount"`
	Percentage float64        `json:"percentage"`
	Frequency  string         `json:"frequency"` // daily, weekly, monthly
	Enabled    bool           `gorm:"default:true" json:"enabled"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ar *AutomatedRule) BeforeCreate(tx *gorm.DB) error {
	if ar.ID == uuid.Nil {
		ar.ID = uuid.New()
	}
	return nil
}

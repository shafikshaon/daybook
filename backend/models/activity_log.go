package models

import (
	"time"

	"gorm.io/gorm"
)

type ActivityLog struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"userId"`
	Action      string         `gorm:"not null" json:"action"`
	Module      string         `gorm:"not null;index" json:"module"` // e.g., "transaction", "account", "budget", etc.
	EntityType  string         `json:"entityType"`                   // e.g., "Transaction", "Account", "Budget"
	EntityID    *uint          `json:"entityId"`                     // ID of the affected entity
	Description string         `json:"description"`                  // Human-readable description
	IPAddress   string         `json:"ipAddress"`
	UserAgent   string         `json:"userAgent"`
	Metadata    string         `gorm:"type:jsonb" json:"metadata"` // Additional data in JSON format
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (a *ActivityLog) BeforeCreate(tx *gorm.DB) error {
	return nil
}

// Action types constants
const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionView   = "view"
	ActionLogin  = "login"
	ActionLogout = "logout"
	ActionExport = "export"
	ActionImport = "import"
)

// Module constants
const (
	ModuleAuth        = "auth"
	ModuleAccount     = "account"
	ModuleCategory    = "category"
	ModuleTransaction = "transaction"
	ModuleBudget      = "budget"
	ModuleCreditCard  = "credit_card"
	ModuleDebt        = "debt"
	ModuleLend        = "lend"
	ModuleAsset       = "asset"
	ModuleGoal        = "goal"
	ModuleSettings    = "settings"
	ModuleReport      = "report"
	ModuleReconcile   = "reconciliation"
)

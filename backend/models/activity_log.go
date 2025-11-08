package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityLog struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	Action      string         `gorm:"not null" json:"action"`
	Module      string         `gorm:"not null;index" json:"module"` // e.g., "transaction", "account", "budget", etc.
	EntityType  string         `json:"entityType"`                   // e.g., "Transaction", "Account", "Budget"
	EntityID    *uuid.UUID     `gorm:"type:uuid" json:"entityId"`    // ID of the affected entity
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
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
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
	ModuleAuth         = "auth"
	ModuleAccount      = "account"
	ModuleTransaction  = "transaction"
	ModuleBudget       = "budget"
	ModuleCreditCard   = "credit_card"
	ModuleDebt         = "debt"
	ModuleLend         = "lend"
	ModuleAsset        = "asset"
	ModuleGoal         = "goal"
	ModuleSettings     = "settings"
	ModuleReport       = "report"
	ModuleReconcile    = "reconciliation"
)

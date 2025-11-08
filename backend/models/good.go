package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Good represents a purchased item with warranty tracking
type Good struct {
	ID                uuid.UUID        `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID            uuid.UUID        `gorm:"type:uuid;not null;index" json:"userId"`
	Name              string           `gorm:"not null" json:"name" binding:"required"`
	Description       string           `json:"description"`
	Category          string           `gorm:"index" json:"category"` // Electronics, Appliances, Furniture, etc.
	Brand             string           `json:"brand"`
	Model             string           `json:"model"`
	SerialNumber      string           `json:"serialNumber"`
	PurchaseDate      Date             `gorm:"not null;index" json:"purchaseDate"`
	PurchasePrice     float64          `gorm:"not null" json:"purchasePrice" binding:"required,gt=0"`
	PurchaseLocation  string           `json:"purchaseLocation"` // Store/website name
	WarrantyStartDate *time.Time       `json:"warrantyStartDate"`
	WarrantyEndDate   *time.Time       `json:"warrantyEndDate"`
	WarrantyProvider  string           `json:"warrantyProvider"` // Manufacturer, Retailer, Extended warranty provider
	WarrantyType      string           `json:"warrantyType"`     // manufacturer, extended, lifetime
	Status            string           `gorm:"not null;index" json:"status"` // active, archived, sold, disposed
	Notes             string           `json:"notes"`
	CreatedAt         time.Time        `json:"createdAt"`
	UpdatedAt         time.Time        `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt   `gorm:"index" json:"-"`

	// Relationships
	Attachments       []GoodAttachment `gorm:"foreignKey:GoodID" json:"attachments,omitempty"`
	ServiceRecords    []ServiceRecord  `gorm:"foreignKey:GoodID" json:"serviceRecords,omitempty"`
}

func (g *Good) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	// Set default status to active if not set
	if g.Status == "" {
		g.Status = "active"
	}
	return nil
}

// ServiceRecord represents a service/repair record for a good
type ServiceRecord struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	GoodID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"goodId"`
	ServiceDate     Date           `gorm:"not null;index" json:"serviceDate"`
	ServiceType     string         `gorm:"not null" json:"serviceType"` // repair, maintenance, inspection, replacement
	ServiceProvider string         `json:"serviceProvider"` // Company/person who provided service
	Cost            float64        `gorm:"not null" json:"cost" binding:"required,gte=0"`
	Description     string         `json:"description"`
	Notes           string         `json:"notes"`
	WarrantyCovered bool           `gorm:"default:false" json:"warrantyCovered"` // Was it covered by warranty?
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (sr *ServiceRecord) BeforeCreate(tx *gorm.DB) error {
	if sr.ID == uuid.Nil {
		sr.ID = uuid.New()
	}
	return nil
}

// GoodAttachment represents an attachment (photo, receipt, warranty document) for a good
type GoodAttachment struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	GoodID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"goodId"`
	FileName     string         `gorm:"not null" json:"fileName"`         // Unique filename on server
	OriginalName string         `gorm:"not null" json:"originalName"`     // Original filename
	FilePath     string         `gorm:"not null" json:"filePath"`         // Server file path
	FileURL      string         `gorm:"not null" json:"fileUrl"`          // API access URL
	FileSize     int64          `json:"fileSize"`
	MimeType     string         `json:"mimeType"`
	AttachmentType string       `gorm:"index" json:"attachmentType"` // photo, receipt, warranty_document, manual, other
	Description  string         `json:"description"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ga *GoodAttachment) BeforeCreate(tx *gorm.DB) error {
	if ga.ID == uuid.Nil {
		ga.ID = uuid.New()
	}
	return nil
}

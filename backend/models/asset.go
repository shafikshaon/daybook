package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Asset represents a purchased item with warranty tracking
type Asset struct {
	ID                uuid.UUID         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID            uuid.UUID         `gorm:"type:uuid;not null;index" json:"userId"`
	Name              string            `gorm:"not null" json:"name" binding:"required"`
	Description       string            `json:"description"`
	Category          string            `gorm:"index" json:"category"` // Electronics, Appliances, Furniture, etc.
	Brand             string            `json:"brand"`
	Model             string            `json:"model"`
	SerialNumber      string            `json:"serialNumber"`
	PurchaseDate      Date              `gorm:"not null;index" json:"purchaseDate"`
	PurchasePrice     float64           `gorm:"not null" json:"purchasePrice" binding:"required,gt=0"`
	PurchaseLocation  string            `json:"purchaseLocation"` // Store/website name
	WarrantyStartDate *Date             `json:"warrantyStartDate"`
	WarrantyEndDate   *Date             `json:"warrantyEndDate"`
	WarrantyProvider  string            `json:"warrantyProvider"` // Manufacturer, Retailer, Extended warranty provider
	WarrantyType      string            `json:"warrantyType"`     // manufacturer, extended, lifetime
	Status            string            `gorm:"not null;index" json:"status"` // active, archived, sold, disposed
	Notes             string            `json:"notes"`
	CreatedAt         time.Time         `json:"createdAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt    `gorm:"index" json:"-"`

	// Relationships
	Attachments       []AssetAttachment `gorm:"foreignKey:AssetID" json:"attachments,omitempty"`
	ServiceRecords    []ServiceRecord   `gorm:"foreignKey:AssetID" json:"serviceRecords,omitempty"`
}

func (a *Asset) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	// Set default status to active if not set
	if a.Status == "" {
		a.Status = "active"
	}
	return nil
}

// ServiceRecord represents a service/repair record for an asset
type ServiceRecord struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	AssetID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"assetId"`
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

// AssetAttachment represents an attachment (photo, receipt, warranty document) for an asset
type AssetAttachment struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	AssetID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"assetId"`
	FileName       string         `gorm:"not null" json:"fileName"`         // Unique filename on server
	OriginalName   string         `gorm:"not null" json:"originalName"`     // Original filename
	FilePath       string         `gorm:"not null" json:"filePath"`         // Server file path
	FileURL        string         `gorm:"not null" json:"fileUrl"`          // API access URL
	FileSize       int64          `json:"fileSize"`
	MimeType       string         `json:"mimeType"`
	AttachmentType string         `gorm:"index" json:"attachmentType"` // photo, receipt, warranty_document, manual, other
	Description    string         `json:"description"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (aa *AssetAttachment) BeforeCreate(tx *gorm.DB) error {
	if aa.ID == uuid.Nil {
		aa.ID = uuid.New()
	}
	return nil
}

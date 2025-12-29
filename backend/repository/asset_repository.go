package repository

import (
	"context"

	"daybook-backend/models"

	"gorm.io/gorm"
)

// AssetFilters represents query filters for assets
type AssetFilters struct {
	Status   *string
	Category *string
}

// AssetStatsResponse represents asset statistics
type AssetStatsResponse struct {
	TotalAssets        int64   `json:"totalAssets"`
	ActiveAssets       int64   `json:"activeAssets"`
	ArchivedAssets     int64   `json:"archivedAssets"`
	TotalValue         float64 `json:"totalValue"`
	TotalServiceCost   float64 `json:"totalServiceCost"`
	AverageAssetValue  float64 `json:"averageAssetValue"`
	AssetsWithWarranty int64   `json:"assetsWithWarranty"`
}

// AssetRepository handles asset data access
type AssetRepository interface {
	BaseRepository[models.Asset]

	// FindWithFilters retrieves assets with optional filters and preloads
	FindWithFilters(ctx context.Context, userID uint, filters AssetFilters) ([]models.Asset, error)

	// FindByIDWithPreloads retrieves an asset with all relationships
	FindByIDWithPreloads(ctx context.Context, assetID, userID uint) (*models.Asset, error)

	// GetStats calculates asset statistics
	GetStats(ctx context.Context, userID uint) (*AssetStatsResponse, error)

	// Service Record operations
	CreateServiceRecord(ctx context.Context, record *models.ServiceRecord) error
	FindServiceRecordsByAsset(ctx context.Context, assetID, userID uint) ([]models.ServiceRecord, error)
	DeleteServiceRecord(ctx context.Context, serviceID, userID uint) error

	// Attachment operations
	CreateAttachment(ctx context.Context, attachment *models.AssetAttachment) error
	FindAttachmentsByAsset(ctx context.Context, assetID, userID uint) ([]models.AssetAttachment, error)
	FindAttachmentByID(ctx context.Context, attachmentID, userID uint) (*models.AssetAttachment, error)
	DeleteAttachment(ctx context.Context, attachmentID, userID uint) error
}

type assetRepository struct {
	*GormBaseRepository[models.Asset]
}

// NewAssetRepository creates a new asset repository
func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{
		GormBaseRepository: NewGormBaseRepository[models.Asset](db),
	}
}

// FindWithFilters retrieves assets with optional filters and preloads
func (r *assetRepository) FindWithFilters(ctx context.Context, userID uint, filters AssetFilters) ([]models.Asset, error) {
	var assets []models.Asset

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.Category != nil {
		query = query.Where("category = ?", *filters.Category)
	}

	err := query.Order("purchase_date DESC, created_at DESC").
		Preload("Attachments").
		Preload("ServiceRecords").
		Find(&assets).Error
	return assets, err
}

// FindByIDWithPreloads retrieves an asset with all relationships
func (r *assetRepository) FindByIDWithPreloads(ctx context.Context, assetID, userID uint) (*models.Asset, error) {
	var asset models.Asset
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", assetID, userID).
		Preload("Attachments").
		Preload("ServiceRecords", "deleted_at IS NULL", func(db *gorm.DB) *gorm.DB {
			return db.Order("service_date DESC")
		}).
		First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// GetStats calculates asset statistics
func (r *assetRepository) GetStats(ctx context.Context, userID uint) (*AssetStatsResponse, error) {
	stats := &AssetStatsResponse{}

	// Total assets
	r.db.WithContext(ctx).Model(&models.Asset{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalAssets)

	// Active assets
	r.db.WithContext(ctx).Model(&models.Asset{}).
		Where("user_id = ? AND status = ?", userID, "active").
		Count(&stats.ActiveAssets)

	// Archived assets
	r.db.WithContext(ctx).Model(&models.Asset{}).
		Where("user_id = ? AND status = ?", userID, "archived").
		Count(&stats.ArchivedAssets)

	// Total value
	var totalValueResult struct {
		TotalValue float64
	}
	r.db.WithContext(ctx).Model(&models.Asset{}).
		Select("COALESCE(SUM(purchase_price), 0) as total_value").
		Where("user_id = ? AND status = ?", userID, "active").
		Scan(&totalValueResult)
	stats.TotalValue = totalValueResult.TotalValue

	// Total service cost
	var totalServiceCostResult struct {
		TotalServiceCost float64
	}
	r.db.WithContext(ctx).Model(&models.ServiceRecord{}).
		Select("COALESCE(SUM(cost), 0) as total_service_cost").
		Where("user_id = ?", userID).
		Scan(&totalServiceCostResult)
	stats.TotalServiceCost = totalServiceCostResult.TotalServiceCost

	// Average asset value
	if stats.TotalAssets > 0 {
		stats.AverageAssetValue = stats.TotalValue / float64(stats.TotalAssets)
	}

	// Assets with warranty
	r.db.WithContext(ctx).Model(&models.Asset{}).
		Where("user_id = ? AND warranty_end_date IS NOT NULL AND warranty_end_date >= CURRENT_DATE", userID).
		Count(&stats.AssetsWithWarranty)

	return stats, nil
}

// CreateServiceRecord creates a new service record
func (r *assetRepository) CreateServiceRecord(ctx context.Context, record *models.ServiceRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

// FindServiceRecordsByAsset retrieves all service records for an asset
func (r *assetRepository) FindServiceRecordsByAsset(ctx context.Context, assetID, userID uint) ([]models.ServiceRecord, error) {
	var records []models.ServiceRecord
	err := r.db.WithContext(ctx).
		Where("asset_id = ? AND user_id = ?", assetID, userID).
		Order("service_date DESC, created_at DESC").
		Find(&records).Error
	return records, err
}

// DeleteServiceRecord deletes a service record
func (r *assetRepository) DeleteServiceRecord(ctx context.Context, serviceID, userID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", serviceID, userID).
		Delete(&models.ServiceRecord{}).Error
}

// CreateAttachment creates a new asset attachment
func (r *assetRepository) CreateAttachment(ctx context.Context, attachment *models.AssetAttachment) error {
	return r.db.WithContext(ctx).Create(attachment).Error
}

// FindAttachmentsByAsset retrieves all attachments for an asset
func (r *assetRepository) FindAttachmentsByAsset(ctx context.Context, assetID, userID uint) ([]models.AssetAttachment, error) {
	var attachments []models.AssetAttachment
	err := r.db.WithContext(ctx).
		Where("asset_id = ? AND user_id = ?", assetID, userID).
		Order("created_at DESC").
		Find(&attachments).Error
	return attachments, err
}

// FindAttachmentByID retrieves a specific attachment
func (r *assetRepository) FindAttachmentByID(ctx context.Context, attachmentID, userID uint) (*models.AssetAttachment, error) {
	var attachment models.AssetAttachment
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", attachmentID, userID).
		First(&attachment).Error
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}

// DeleteAttachment deletes an attachment
func (r *assetRepository) DeleteAttachment(ctx context.Context, attachmentID, userID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", attachmentID, userID).
		Delete(&models.AssetAttachment{}).Error
}

package services

import (
	"context"
	"errors"
	"os"
	"time"

	"daybook-backend/models"
	"daybook-backend/repository"

	"github.com/google/uuid"
)

// AssetWithStats includes an asset with calculated statistics
type AssetWithStats struct {
	models.Asset
	WarrantyDaysTotal     *int    `json:"warrantyDaysTotal"`
	WarrantyDaysPassed    *int    `json:"warrantyDaysPassed"`
	WarrantyDaysRemaining *int    `json:"warrantyDaysRemaining"`
	WarrantyStatus        string  `json:"warrantyStatus"` // active, expired, no_warranty
	DaysOwned             int     `json:"daysOwned"`
	PricePerDay           float64 `json:"pricePerDay"`
	TotalServiceCost      float64 `json:"totalServiceCost"`
	ServiceCount          int     `json:"serviceCount"`
	TotalCost             float64 `json:"totalCost"` // Purchase price + service costs
}

// AssetService handles asset business logic
type AssetService interface {
	// ListAssets retrieves all assets with optional filters
	ListAssets(ctx context.Context, userID uuid.UUID, filters repository.AssetFilters) ([]AssetWithStats, error)

	// GetAsset retrieves a specific asset by ID
	GetAsset(ctx context.Context, assetID, userID uuid.UUID) (*AssetWithStats, error)

	// CreateAsset creates a new asset
	CreateAsset(ctx context.Context, asset *models.Asset) (*models.Asset, error)

	// UpdateAsset updates an existing asset
	UpdateAsset(ctx context.Context, assetID, userID uuid.UUID, updateData *models.Asset) (*models.Asset, error)

	// DeleteAsset deletes an asset
	DeleteAsset(ctx context.Context, assetID, userID uuid.UUID) error

	// GetStats calculates asset statistics
	GetStats(ctx context.Context, userID uuid.UUID) (*repository.AssetStatsResponse, error)

	// Service Record operations
	CreateServiceRecord(ctx context.Context, record *models.ServiceRecord) (*models.ServiceRecord, error)
	ListServiceRecords(ctx context.Context, assetID, userID uuid.UUID) ([]models.ServiceRecord, error)
	DeleteServiceRecord(ctx context.Context, serviceID, userID uuid.UUID) error

	// Attachment operations
	AddAttachment(ctx context.Context, attachment *models.AssetAttachment) (*models.AssetAttachment, error)
	ListAttachments(ctx context.Context, assetID, userID uuid.UUID) ([]models.AssetAttachment, error)
	DeleteAttachment(ctx context.Context, attachmentID, userID uuid.UUID) error
}

type assetService struct {
	repo           repository.AssetRepository
	activityLogger ActivityLogService
}

// NewAssetService creates a new asset service
func NewAssetService(
	repo repository.AssetRepository,
	activityLogger ActivityLogService,
) AssetService {
	return &assetService{
		repo:           repo,
		activityLogger: activityLogger,
	}
}

// ListAssets retrieves assets with optional filters
func (s *assetService) ListAssets(ctx context.Context, userID uuid.UUID, filters repository.AssetFilters) ([]AssetWithStats, error) {
	assets, err := s.repo.FindWithFilters(ctx, userID, filters)
	if err != nil {
		return nil, err
	}

	// Enrich with statistics
	enrichedAssets := make([]AssetWithStats, len(assets))
	for i, asset := range assets {
		enrichedAssets[i] = calculateAssetStats(asset)
	}

	return enrichedAssets, nil
}

// GetAsset retrieves a specific asset
func (s *assetService) GetAsset(ctx context.Context, assetID, userID uuid.UUID) (*AssetWithStats, error) {
	asset, err := s.repo.FindByIDWithPreloads(ctx, assetID, userID)
	if err != nil {
		return nil, errors.New("asset not found")
	}

	enriched := calculateAssetStats(*asset)
	return &enriched, nil
}

// CreateAsset creates a new asset
func (s *assetService) CreateAsset(ctx context.Context, asset *models.Asset) (*models.Asset, error) {
	if err := s.repo.Create(ctx, asset); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		asset.UserID,
		models.ActionCreate,
		models.ModuleAsset,
		"Asset",
		asset.ID,
		"Created asset: "+asset.Name,
		nil,
	)

	return asset, nil
}

// UpdateAsset updates an existing asset
func (s *assetService) UpdateAsset(ctx context.Context, assetID, userID uuid.UUID, updateData *models.Asset) (*models.Asset, error) {
	// Fetch existing asset
	existing, err := s.repo.FindByID(ctx, assetID, userID)
	if err != nil {
		return nil, errors.New("asset not found")
	}

	// Update allowed fields
	existing.Name = updateData.Name
	existing.Description = updateData.Description
	existing.Category = updateData.Category
	existing.Brand = updateData.Brand
	existing.Model = updateData.Model
	existing.SerialNumber = updateData.SerialNumber
	existing.PurchaseDate = updateData.PurchaseDate
	existing.PurchasePrice = updateData.PurchasePrice
	existing.PurchaseLocation = updateData.PurchaseLocation
	existing.WarrantyStartDate = updateData.WarrantyStartDate
	existing.WarrantyEndDate = updateData.WarrantyEndDate
	existing.WarrantyProvider = updateData.WarrantyProvider
	existing.WarrantyType = updateData.WarrantyType
	existing.Status = updateData.Status
	existing.Notes = updateData.Notes

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionUpdate,
		models.ModuleAsset,
		"Asset",
		existing.ID,
		"Updated asset: "+existing.Name,
		nil,
	)

	return existing, nil
}

// DeleteAsset deletes an asset
func (s *assetService) DeleteAsset(ctx context.Context, assetID, userID uuid.UUID) error {
	// Fetch the asset to get its details
	asset, err := s.repo.FindByID(ctx, assetID, userID)
	if err != nil {
		return errors.New("asset not found")
	}

	// Delete the asset
	if err := s.repo.Delete(ctx, assetID, userID); err != nil {
		return err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionDelete,
		models.ModuleAsset,
		"Asset",
		asset.ID,
		"Deleted asset: "+asset.Name,
		nil,
	)

	return nil
}

// GetStats calculates asset statistics
func (s *assetService) GetStats(ctx context.Context, userID uuid.UUID) (*repository.AssetStatsResponse, error) {
	return s.repo.GetStats(ctx, userID)
}

// CreateServiceRecord creates a new service record
func (s *assetService) CreateServiceRecord(ctx context.Context, record *models.ServiceRecord) (*models.ServiceRecord, error) {
	// Verify asset exists and belongs to user
	_, err := s.repo.FindByID(ctx, record.AssetID, record.UserID)
	if err != nil {
		return nil, errors.New("asset not found")
	}

	if err := s.repo.CreateServiceRecord(ctx, record); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		record.UserID,
		models.ActionCreate,
		models.ModuleAsset,
		"ServiceRecord",
		record.ID,
		"Created service record for asset",
		nil,
	)

	return record, nil
}

// ListServiceRecords retrieves all service records for an asset
func (s *assetService) ListServiceRecords(ctx context.Context, assetID, userID uuid.UUID) ([]models.ServiceRecord, error) {
	// Verify asset exists and belongs to user
	_, err := s.repo.FindByID(ctx, assetID, userID)
	if err != nil {
		return nil, errors.New("asset not found")
	}

	return s.repo.FindServiceRecordsByAsset(ctx, assetID, userID)
}

// DeleteServiceRecord deletes a service record
func (s *assetService) DeleteServiceRecord(ctx context.Context, serviceID, userID uuid.UUID) error {
	if err := s.repo.DeleteServiceRecord(ctx, serviceID, userID); err != nil {
		return err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionDelete,
		models.ModuleAsset,
		"ServiceRecord",
		serviceID,
		"Deleted service record",
		nil,
	)

	return nil
}

// AddAttachment adds an attachment to an asset
func (s *assetService) AddAttachment(ctx context.Context, attachment *models.AssetAttachment) (*models.AssetAttachment, error) {
	// Verify asset exists and belongs to user
	_, err := s.repo.FindByID(ctx, attachment.AssetID, attachment.UserID)
	if err != nil {
		return nil, errors.New("asset not found")
	}

	if err := s.repo.CreateAttachment(ctx, attachment); err != nil {
		return nil, err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		attachment.UserID,
		models.ActionCreate,
		models.ModuleAsset,
		"AssetAttachment",
		attachment.ID,
		"Added attachment to asset",
		nil,
	)

	return attachment, nil
}

// ListAttachments retrieves all attachments for an asset
func (s *assetService) ListAttachments(ctx context.Context, assetID, userID uuid.UUID) ([]models.AssetAttachment, error) {
	// Verify asset exists and belongs to user
	_, err := s.repo.FindByID(ctx, assetID, userID)
	if err != nil {
		return nil, errors.New("asset not found")
	}

	return s.repo.FindAttachmentsByAsset(ctx, assetID, userID)
}

// DeleteAttachment deletes an attachment
func (s *assetService) DeleteAttachment(ctx context.Context, attachmentID, userID uuid.UUID) error {
	// Get attachment to retrieve file path
	attachment, err := s.repo.FindAttachmentByID(ctx, attachmentID, userID)
	if err != nil {
		return errors.New("attachment not found")
	}

	// Delete the file from disk
	if err := os.Remove(attachment.FilePath); err != nil {
		// Log error but continue with database deletion
	}

	// Delete from database
	if err := s.repo.DeleteAttachment(ctx, attachmentID, userID); err != nil {
		return err
	}

	// Log activity
	s.activityLogger.LogEntityActivity(
		ctx,
		userID,
		models.ActionDelete,
		models.ModuleAsset,
		"AssetAttachment",
		attachmentID,
		"Deleted attachment from asset",
		nil,
	)

	return nil
}

// calculateAssetStats calculates statistics for an asset
func calculateAssetStats(asset models.Asset) AssetWithStats {
	stats := AssetWithStats{
		Asset:          asset,
		WarrantyStatus: "no_warranty",
	}

	// Calculate warranty statistics
	if asset.WarrantyStartDate != nil && asset.WarrantyEndDate != nil {
		warrantyStart := asset.WarrantyStartDate.Time
		warrantyEnd := asset.WarrantyEndDate.Time
		now := time.Now()

		totalDays := int(warrantyEnd.Sub(warrantyStart).Hours() / 24)
		daysPassed := int(now.Sub(warrantyStart).Hours() / 24)
		daysRemaining := int(warrantyEnd.Sub(now).Hours() / 24)

		stats.WarrantyDaysTotal = &totalDays
		stats.WarrantyDaysPassed = &daysPassed
		stats.WarrantyDaysRemaining = &daysRemaining

		if daysRemaining > 0 {
			stats.WarrantyStatus = "active"
		} else {
			stats.WarrantyStatus = "expired"
		}
	}

	// Calculate ownership statistics
	purchaseDate := asset.PurchaseDate.Time
	stats.DaysOwned = int(time.Since(purchaseDate).Hours() / 24)
	if stats.DaysOwned > 0 {
		stats.PricePerDay = asset.PurchasePrice / float64(stats.DaysOwned)
	}

	// Calculate service costs
	stats.ServiceCount = len(asset.ServiceRecords)
	for _, record := range asset.ServiceRecords {
		stats.TotalServiceCost += record.Cost
	}
	stats.TotalCost = asset.PurchasePrice + stats.TotalServiceCost

	return stats
}

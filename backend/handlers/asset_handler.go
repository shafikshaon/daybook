package handlers

import (
	"net/http"
	"time"

	"daybook-backend/database"
	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AssetWithStats includes an asset with calculated statistics
type AssetWithStats struct {
	models.Asset
	WarrantyDaysTotal     *int     `json:"warrantyDaysTotal"`
	WarrantyDaysPassed    *int     `json:"warrantyDaysPassed"`
	WarrantyDaysRemaining *int     `json:"warrantyDaysRemaining"`
	WarrantyStatus        string   `json:"warrantyStatus"` // active, expired, no_warranty
	DaysOwned             int      `json:"daysOwned"`
	PricePerDay           float64  `json:"pricePerDay"`
	TotalServiceCost      float64  `json:"totalServiceCost"`
	ServiceCount          int      `json:"serviceCount"`
	TotalCost             float64  `json:"totalCost"` // Purchase price + service costs
}

// ListAssets returns all assets for the authenticated user
func ListAssets(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListAssets - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Fetching assets for user")
	query := database.DB.WithContext(ctx).Where("user_id = ?", userID)

	// Apply filters
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	var assets []models.Asset
	if err := query.Order("purchase_date DESC, created_at DESC").
		Preload("Attachments").
		Preload("ServiceRecords").
		Find(&assets).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch assets")
		return
	}

	// Enrich with statistics
	enrichedAssets := make([]AssetWithStats, len(assets))
	for i, asset := range assets {
		enrichedAssets[i] = calculateAssetStats(asset)
	}

	logger.Infof(ctx, "Assets retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, enrichedAssets, "Assets retrieved successfully")
}

// GetAsset returns a specific asset by ID with statistics
func GetAsset(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetAsset - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	logger.Debugf(ctx, "Fetching asset with ID: %s", assetID)
	var asset models.Asset
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", assetID, userID).
		Preload("Attachments").
		Preload("ServiceRecords", "deleted_at IS NULL", func(db *gorm.DB) *gorm.DB {
			return db.Order("service_date DESC")
		}).
		First(&asset).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Asset not found")
		return
	}

	response := calculateAssetStats(asset)
	logger.Infof(ctx, "Asset retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, response, "Asset retrieved successfully")
}

// CreateAsset creates a new asset record
func CreateAsset(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateAsset - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var asset models.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate purchase date is provided
	if asset.PurchaseDate.IsZero() {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Purchase date is required")
		return
	}

	asset.UserID = userID

	logger.Debugf(ctx, "Creating asset: %s", asset.Name)
	if err := database.DB.WithContext(ctx).Create(&asset).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create asset")
		return
	}

	// Log asset creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleAsset,
		"Asset", asset.ID, "Created asset: "+asset.Name, nil)

	logger.Infof(ctx, "Asset created successfully for user: %s", userID)
	utilities.CreatedResponse(c, asset, "Asset created successfully")
}

// UpdateAsset updates an asset record
func UpdateAsset(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateAsset - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	logger.Debugf(ctx, "Fetching existing asset with ID: %s", assetID)
	var existingAsset models.Asset
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", assetID, userID).First(&existingAsset).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Asset not found")
		return
	}

	var updateData models.Asset
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Preserve protected fields from the existing asset
	updateData.ID = existingAsset.ID
	updateData.UserID = existingAsset.UserID
	updateData.CreatedAt = existingAsset.CreatedAt

	logger.Debugf(ctx, "Updating asset: %s", existingAsset.Name)
	if err := database.DB.WithContext(ctx).Model(&existingAsset).Updates(&updateData).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update asset")
		return
	}

	// Reload the asset to get the updated data with all calculated fields
	logger.Debugf(ctx, "Reloading asset data")
	if err := database.DB.WithContext(ctx).Where("id = ?", assetID).
		Preload("Attachments").
		Preload("ServiceRecords").
		First(&existingAsset).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to reload asset")
		return
	}

	// Log asset update activity
	utilities.LogEntityActivity(c, userID, models.ActionUpdate, models.ModuleAsset,
		"Asset", existingAsset.ID, "Updated asset: "+existingAsset.Name, nil)

	logger.Infof(ctx, "Asset updated successfully for user: %s", userID)
	utilities.SuccessResponse(c, existingAsset, "Asset updated successfully")
}

// DeleteAsset soft deletes an asset record
func DeleteAsset(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteAsset - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	logger.Debugf(ctx, "Fetching asset with ID: %s for deletion", assetID)
	var asset models.Asset
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", assetID, userID).First(&asset).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Asset not found")
		return
	}

	// Start transaction to delete asset and related records
	logger.Debugf(ctx, "Starting transaction for asset deletion")
	tx := database.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Soft delete attachments
	logger.Debugf(ctx, "Deleting asset attachments")
	if err := tx.Where("asset_id = ?", assetID).Delete(&models.AssetAttachment{}).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete attachments")
		return
	}

	// Soft delete service records
	logger.Debugf(ctx, "Deleting service records")
	if err := tx.Where("asset_id = ?", assetID).Delete(&models.ServiceRecord{}).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete service records")
		return
	}

	// Soft delete the asset
	logger.Debugf(ctx, "Deleting asset")
	if err := tx.Delete(&asset).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete asset")
		return
	}

	if err := tx.Commit().Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	// Log asset deletion activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, models.ModuleAsset,
		"Asset", asset.ID, "Deleted asset: "+asset.Name, nil)

	logger.Infof(ctx, "Asset deleted successfully for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Asset deleted successfully")
}

// CreateServiceRecord adds a service record for an asset
func CreateServiceRecord(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateServiceRecord - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	// Verify asset belongs to user
	logger.Debugf(ctx, "Verifying asset ownership for asset ID: %s", assetID)
	var asset models.Asset
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", assetID, userID).First(&asset).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Asset not found")
		return
	}

	var serviceRecord models.ServiceRecord
	if err := c.ShouldBindJSON(&serviceRecord); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate service date is provided
	if serviceRecord.ServiceDate.IsZero() {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Service date is required")
		return
	}

	serviceRecord.UserID = userID
	serviceRecord.AssetID = assetID

	logger.Debugf(ctx, "Creating service record for asset: %s", asset.Name)
	if err := database.DB.WithContext(ctx).Create(&serviceRecord).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create service record")
		return
	}

	// Log service record creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleAsset,
		"ServiceRecord", serviceRecord.ID, "Created service record for asset: "+asset.Name, nil)

	logger.Infof(ctx, "Service record created successfully for user: %s", userID)
	utilities.CreatedResponse(c, serviceRecord, "Service record created successfully")
}

// ListServiceRecords returns all service records for a specific asset
func ListServiceRecords(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListServiceRecords - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	// Verify asset belongs to user
	logger.Debugf(ctx, "Verifying asset ownership for asset ID: %s", assetID)
	var asset models.Asset
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", assetID, userID).First(&asset).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Asset not found")
		return
	}

	logger.Debugf(ctx, "Fetching service records for asset")
	var serviceRecords []models.ServiceRecord
	if err := database.DB.WithContext(ctx).Where("asset_id = ? AND user_id = ?", assetID, userID).
		Order("service_date DESC, created_at DESC").
		Find(&serviceRecords).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch service records")
		return
	}

	logger.Infof(ctx, "Service records retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, serviceRecords, "Service records retrieved successfully")
}

// DeleteServiceRecord deletes a service record
func DeleteServiceRecord(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteServiceRecord - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	serviceID, err := uuid.Parse(c.Param("serviceId"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid service record ID")
		return
	}

	logger.Debugf(ctx, "Fetching service record with ID: %s for deletion", serviceID)
	var serviceRecord models.ServiceRecord
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", serviceID, userID).First(&serviceRecord).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Service record not found")
		return
	}

	logger.Debugf(ctx, "Deleting service record")
	if err := database.DB.WithContext(ctx).Delete(&serviceRecord).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete service record")
		return
	}

	// Log service record deletion activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, models.ModuleAsset,
		"ServiceRecord", serviceRecord.ID, "Deleted service record", nil)

	logger.Infof(ctx, "Service record deleted successfully for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Service record deleted successfully")
}

// AddAttachment links an uploaded file to an asset
func AddAttachment(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "AddAttachment - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	// Verify asset belongs to user
	logger.Debugf(ctx, "Verifying asset ownership for asset ID: %s", assetID)
	var asset models.Asset
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", assetID, userID).First(&asset).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Asset not found")
		return
	}

	var attachment models.AssetAttachment
	if err := c.ShouldBindJSON(&attachment); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	attachment.UserID = userID
	attachment.AssetID = assetID

	logger.Debugf(ctx, "Creating attachment for asset: %s", asset.Name)
	if err := database.DB.WithContext(ctx).Create(&attachment).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to add attachment")
		return
	}

	// Log attachment addition activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, models.ModuleAsset,
		"AssetAttachment", attachment.ID, "Added attachment to asset: "+asset.Name, nil)

	logger.Infof(ctx, "Attachment added successfully for user: %s", userID)
	utilities.CreatedResponse(c, attachment, "Attachment added successfully")
}

// ListAttachments returns all attachments for a specific asset
func ListAttachments(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListAttachments - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	// Verify asset belongs to user
	logger.Debugf(ctx, "Verifying asset ownership for asset ID: %s", assetID)
	var asset models.Asset
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", assetID, userID).First(&asset).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Asset not found")
		return
	}

	logger.Debugf(ctx, "Fetching attachments for asset")
	var attachments []models.AssetAttachment
	if err := database.DB.WithContext(ctx).Where("asset_id = ? AND user_id = ?", assetID, userID).
		Order("created_at DESC").
		Find(&attachments).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch attachments")
		return
	}

	logger.Infof(ctx, "Attachments retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, attachments, "Attachments retrieved successfully")
}

// DeleteAttachment deletes an attachment
func DeleteAttachment(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteAttachment - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	attachmentID, err := uuid.Parse(c.Param("attachmentId"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid attachment ID")
		return
	}

	logger.Debugf(ctx, "Fetching attachment with ID: %s for deletion", attachmentID)
	var attachment models.AssetAttachment
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", attachmentID, userID).First(&attachment).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Attachment not found")
		return
	}

	logger.Debugf(ctx, "Deleting attachment")
	if err := database.DB.WithContext(ctx).Delete(&attachment).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete attachment")
		return
	}

	// Log attachment deletion activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, models.ModuleAsset,
		"AssetAttachment", attachment.ID, "Deleted attachment", nil)

	logger.Infof(ctx, "Attachment deleted successfully for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Attachment deleted successfully")
}

// GetAssetsStats returns summary statistics for all assets
func GetAssetsStats(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetAssetsStats - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Fetching assets for statistics calculation")
	var assets []models.Asset
	if err := database.DB.WithContext(ctx).Where("user_id = ?", userID).
		Preload("ServiceRecords").
		Find(&assets).Error; err != nil {
		logger.Errorf(ctx, "Database operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch assets")
		return
	}

	stats := struct {
		TotalGoods            int     `json:"totalGoods"`
		ActiveGoods           int     `json:"activeGoods"`
		TotalValue            float64 `json:"totalValue"`
		TotalServiceCost      float64 `json:"totalServiceCost"`
		GoodsUnderWarranty    int     `json:"goodsUnderWarranty"`
		GoodsWarrantyExpiring int     `json:"goodsWarrantyExpiring"` // Expiring in next 30 days
	}{}

	now := time.Now()
	// Truncate to day precision for date-only comparison
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	thirtyDaysFromNow := nowDate.AddDate(0, 0, 30)

	stats.TotalGoods = len(assets)

	for _, asset := range assets {
		if asset.Status == "active" {
			stats.ActiveGoods++
			stats.TotalValue += asset.PurchasePrice
		}

		// Calculate service costs
		for _, service := range asset.ServiceRecords {
			stats.TotalServiceCost += service.Cost
		}

		// Check warranty status
		if asset.WarrantyEndDate != nil {
			endDate := asset.WarrantyEndDate.Time
			endDateTrunc := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())

			// Warranty is active if today is on or before the end date
			if !nowDate.After(endDateTrunc) {
				stats.GoodsUnderWarranty++
				if endDateTrunc.Before(thirtyDaysFromNow) || endDateTrunc.Equal(thirtyDaysFromNow) {
					stats.GoodsWarrantyExpiring++
				}
			}
		}
	}

	logger.Infof(ctx, "Asset statistics retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, stats, "Statistics retrieved successfully")
}

// Helper function to calculate statistics for an asset
func calculateAssetStats(asset models.Asset) AssetWithStats {
	stats := AssetWithStats{
		Asset: asset,
	}

	now := time.Now()

	// Calculate warranty statistics
	if asset.WarrantyStartDate != nil && asset.WarrantyEndDate != nil {
		startDate := asset.WarrantyStartDate.Time
		endDate := asset.WarrantyEndDate.Time

		// Truncate to day precision for date-only comparison
		nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		startDateTrunc := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
		endDateTrunc := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())

		totalDays := int(endDateTrunc.Sub(startDateTrunc).Hours() / 24)
		daysPassed := int(nowDate.Sub(startDateTrunc).Hours() / 24)
		daysRemaining := int(endDateTrunc.Sub(nowDate).Hours() / 24)

		if daysPassed < 0 {
			daysPassed = 0
		}
		if daysRemaining < 0 {
			daysRemaining = 0
		}

		stats.WarrantyDaysTotal = &totalDays
		stats.WarrantyDaysPassed = &daysPassed
		stats.WarrantyDaysRemaining = &daysRemaining

		// Warranty is active if today is on or before the end date
		if nowDate.After(endDateTrunc) {
			stats.WarrantyStatus = "expired"
		} else {
			stats.WarrantyStatus = "active"
		}
	} else {
		stats.WarrantyStatus = "no_warranty"
	}

	// Calculate days owned (using day precision from purchase date to current date)
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	purchaseDate := asset.PurchaseDate.Time
	purchaseDateTrunc := time.Date(purchaseDate.Year(), purchaseDate.Month(), purchaseDate.Day(), 0, 0, 0, 0, purchaseDate.Location())

	stats.DaysOwned = int(nowDate.Sub(purchaseDateTrunc).Hours() / 24)
	if stats.DaysOwned < 1 {
		stats.DaysOwned = 1 // Avoid division by zero
	}

	// Calculate total service cost
	for _, service := range asset.ServiceRecords {
		stats.TotalServiceCost += service.Cost
		stats.ServiceCount++
	}

	// Calculate total cost
	stats.TotalCost = asset.PurchasePrice + stats.TotalServiceCost

	// Calculate price per day
	stats.PricePerDay = stats.TotalCost / float64(stats.DaysOwned)

	return stats
}

package handlers

import (
	"net/http"
	"time"

	"daybook-backend/database"
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
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	query := database.DB.Where("user_id = ?", userID)

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
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch assets")
		return
	}

	// Enrich with statistics
	enrichedAssets := make([]AssetWithStats, len(assets))
	for i, asset := range assets {
		enrichedAssets[i] = calculateAssetStats(asset)
	}

	utilities.SuccessResponse(c, enrichedAssets, "Assets retrieved successfully")
}

// GetAsset returns a specific asset by ID with statistics
func GetAsset(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	var asset models.Asset
	if err := database.DB.Where("id = ? AND user_id = ?", assetID, userID).
		Preload("Attachments").
		Preload("ServiceRecords", "deleted_at IS NULL", func(db *gorm.DB) *gorm.DB {
			return db.Order("service_date DESC")
		}).
		First(&asset).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Asset not found")
		return
	}

	response := calculateAssetStats(asset)
	utilities.SuccessResponse(c, response, "Asset retrieved successfully")
}

// CreateAsset creates a new asset record
func CreateAsset(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
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

	if err := database.DB.Create(&asset).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create asset")
		return
	}

	utilities.CreatedResponse(c, asset, "Asset created successfully")
}

// UpdateAsset updates an asset record
func UpdateAsset(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	var existingGood models.Asset
	if err := database.DB.Where("id = ? AND user_id = ?", assetID, userID).First(&existingGood).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Asset not found")
		return
	}

	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Prevent updating certain fields
	delete(updateData, "id")
	delete(updateData, "userId")
	delete(updateData, "createdAt")

	if err := database.DB.Model(&existingGood).Updates(updateData).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update asset")
		return
	}

	utilities.SuccessResponse(c, existingGood, "Asset updated successfully")
}

// DeleteAsset soft deletes an asset record
func DeleteAsset(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	var asset models.Asset
	if err := database.DB.Where("id = ? AND user_id = ?", assetID, userID).First(&asset).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Asset not found")
		return
	}

	// Start transaction to delete asset and related records
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Soft delete attachments
	if err := tx.Where("asset_id = ?", assetID).Delete(&models.AssetAttachment{}).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete attachments")
		return
	}

	// Soft delete service records
	if err := tx.Where("asset_id = ?", assetID).Delete(&models.ServiceRecord{}).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete service records")
		return
	}

	// Soft delete the asset
	if err := tx.Delete(&asset).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete asset")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	utilities.SuccessResponse(c, nil, "Asset deleted successfully")
}

// CreateServiceRecord adds a service record for an asset
func CreateServiceRecord(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	// Verify asset belongs to user
	var asset models.Asset
	if err := database.DB.Where("id = ? AND user_id = ?", assetID, userID).First(&asset).Error; err != nil {
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

	if err := database.DB.Create(&serviceRecord).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create service record")
		return
	}

	utilities.CreatedResponse(c, serviceRecord, "Service record created successfully")
}

// ListServiceRecords returns all service records for a specific asset
func ListServiceRecords(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	// Verify asset belongs to user
	var asset models.Asset
	if err := database.DB.Where("id = ? AND user_id = ?", assetID, userID).First(&asset).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Asset not found")
		return
	}

	var serviceRecords []models.ServiceRecord
	if err := database.DB.Where("asset_id = ? AND user_id = ?", assetID, userID).
		Order("service_date DESC, created_at DESC").
		Find(&serviceRecords).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch service records")
		return
	}

	utilities.SuccessResponse(c, serviceRecords, "Service records retrieved successfully")
}

// DeleteServiceRecord deletes a service record
func DeleteServiceRecord(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	serviceID, err := uuid.Parse(c.Param("serviceId"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid service record ID")
		return
	}

	var serviceRecord models.ServiceRecord
	if err := database.DB.Where("id = ? AND user_id = ?", serviceID, userID).First(&serviceRecord).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Service record not found")
		return
	}

	if err := database.DB.Delete(&serviceRecord).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete service record")
		return
	}

	utilities.SuccessResponse(c, nil, "Service record deleted successfully")
}

// AddAttachment links an uploaded file to an asset
func AddAttachment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	// Verify asset belongs to user
	var asset models.Asset
	if err := database.DB.Where("id = ? AND user_id = ?", assetID, userID).First(&asset).Error; err != nil {
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

	if err := database.DB.Create(&attachment).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to add attachment")
		return
	}

	utilities.CreatedResponse(c, attachment, "Attachment added successfully")
}

// ListAttachments returns all attachments for a specific asset
func ListAttachments(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	assetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	// Verify asset belongs to user
	var asset models.Asset
	if err := database.DB.Where("id = ? AND user_id = ?", assetID, userID).First(&asset).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Asset not found")
		return
	}

	var attachments []models.AssetAttachment
	if err := database.DB.Where("asset_id = ? AND user_id = ?", assetID, userID).
		Order("created_at DESC").
		Find(&attachments).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch attachments")
		return
	}

	utilities.SuccessResponse(c, attachments, "Attachments retrieved successfully")
}

// DeleteAttachment deletes an attachment
func DeleteAttachment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	attachmentID, err := uuid.Parse(c.Param("attachmentId"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid attachment ID")
		return
	}

	var attachment models.AssetAttachment
	if err := database.DB.Where("id = ? AND user_id = ?", attachmentID, userID).First(&attachment).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Attachment not found")
		return
	}

	if err := database.DB.Delete(&attachment).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete attachment")
		return
	}

	utilities.SuccessResponse(c, nil, "Attachment deleted successfully")
}

// GetAssetsStats returns summary statistics for all assets
func GetAssetsStats(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var assets []models.Asset
	if err := database.DB.Where("user_id = ?", userID).
		Preload("ServiceRecords").
		Find(&assets).Error; err != nil {
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

	// Calculate days owned
	stats.DaysOwned = int(now.Sub(asset.PurchaseDate.Time).Hours() / 24)
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

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

// GoodWithStats includes a good with calculated statistics
type GoodWithStats struct {
	models.Good
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

// ListGoods returns all goods for the authenticated user
func ListGoods(c *gin.Context) {
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

	var goods []models.Good
	if err := query.Order("purchase_date DESC, created_at DESC").
		Preload("Attachments").
		Preload("ServiceRecords").
		Find(&goods).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch goods")
		return
	}

	// Enrich with statistics
	enrichedGoods := make([]GoodWithStats, len(goods))
	for i, good := range goods {
		enrichedGoods[i] = calculateGoodStats(good)
	}

	utilities.SuccessResponse(c, enrichedGoods, "Goods retrieved successfully")
}

// GetGood returns a specific good by ID with statistics
func GetGood(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid good ID")
		return
	}

	var good models.Good
	if err := database.DB.Where("id = ? AND user_id = ?", goodID, userID).
		Preload("Attachments").
		Preload("ServiceRecords", "deleted_at IS NULL", func(db *gorm.DB) *gorm.DB {
			return db.Order("service_date DESC")
		}).
		First(&good).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Good not found")
		return
	}

	response := calculateGoodStats(good)
	utilities.SuccessResponse(c, response, "Good retrieved successfully")
}

// CreateGood creates a new good record
func CreateGood(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var good models.Good
	if err := c.ShouldBindJSON(&good); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate purchase date is provided
	if good.PurchaseDate.IsZero() {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Purchase date is required")
		return
	}

	good.UserID = userID

	if err := database.DB.Create(&good).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create good")
		return
	}

	utilities.CreatedResponse(c, good, "Good created successfully")
}

// UpdateGood updates a good record
func UpdateGood(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid good ID")
		return
	}

	var existingGood models.Good
	if err := database.DB.Where("id = ? AND user_id = ?", goodID, userID).First(&existingGood).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Good not found")
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
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update good")
		return
	}

	utilities.SuccessResponse(c, existingGood, "Good updated successfully")
}

// DeleteGood soft deletes a good record
func DeleteGood(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid good ID")
		return
	}

	var good models.Good
	if err := database.DB.Where("id = ? AND user_id = ?", goodID, userID).First(&good).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Good not found")
		return
	}

	// Start transaction to delete good and related records
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Soft delete attachments
	if err := tx.Where("good_id = ?", goodID).Delete(&models.GoodAttachment{}).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete attachments")
		return
	}

	// Soft delete service records
	if err := tx.Where("good_id = ?", goodID).Delete(&models.ServiceRecord{}).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete service records")
		return
	}

	// Soft delete the good
	if err := tx.Delete(&good).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete good")
		return
	}

	if err := tx.Commit().Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	utilities.SuccessResponse(c, nil, "Good deleted successfully")
}

// CreateServiceRecord adds a service record for a good
func CreateServiceRecord(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid good ID")
		return
	}

	// Verify good belongs to user
	var good models.Good
	if err := database.DB.Where("id = ? AND user_id = ?", goodID, userID).First(&good).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Good not found")
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
	serviceRecord.GoodID = goodID

	if err := database.DB.Create(&serviceRecord).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create service record")
		return
	}

	utilities.CreatedResponse(c, serviceRecord, "Service record created successfully")
}

// ListServiceRecords returns all service records for a specific good
func ListServiceRecords(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid good ID")
		return
	}

	// Verify good belongs to user
	var good models.Good
	if err := database.DB.Where("id = ? AND user_id = ?", goodID, userID).First(&good).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Good not found")
		return
	}

	var serviceRecords []models.ServiceRecord
	if err := database.DB.Where("good_id = ? AND user_id = ?", goodID, userID).
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

// AddAttachment links an uploaded file to a good
func AddAttachment(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid good ID")
		return
	}

	// Verify good belongs to user
	var good models.Good
	if err := database.DB.Where("id = ? AND user_id = ?", goodID, userID).First(&good).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Good not found")
		return
	}

	var attachment models.GoodAttachment
	if err := c.ShouldBindJSON(&attachment); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	attachment.UserID = userID
	attachment.GoodID = goodID

	if err := database.DB.Create(&attachment).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to add attachment")
		return
	}

	utilities.CreatedResponse(c, attachment, "Attachment added successfully")
}

// ListAttachments returns all attachments for a specific good
func ListAttachments(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	goodID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid good ID")
		return
	}

	// Verify good belongs to user
	var good models.Good
	if err := database.DB.Where("id = ? AND user_id = ?", goodID, userID).First(&good).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Good not found")
		return
	}

	var attachments []models.GoodAttachment
	if err := database.DB.Where("good_id = ? AND user_id = ?", goodID, userID).
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

	var attachment models.GoodAttachment
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

// GetGoodsStats returns summary statistics for all goods
func GetGoodsStats(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var goods []models.Good
	if err := database.DB.Where("user_id = ?", userID).
		Preload("ServiceRecords").
		Find(&goods).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch goods")
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
	thirtyDaysFromNow := now.AddDate(0, 0, 30)

	stats.TotalGoods = len(goods)

	for _, good := range goods {
		if good.Status == "active" {
			stats.ActiveGoods++
			stats.TotalValue += good.PurchasePrice
		}

		// Calculate service costs
		for _, service := range good.ServiceRecords {
			stats.TotalServiceCost += service.Cost
		}

		// Check warranty status
		if good.WarrantyEndDate != nil {
			if good.WarrantyEndDate.After(now) {
				stats.GoodsUnderWarranty++
				if good.WarrantyEndDate.Before(thirtyDaysFromNow) {
					stats.GoodsWarrantyExpiring++
				}
			}
		}
	}

	utilities.SuccessResponse(c, stats, "Statistics retrieved successfully")
}

// Helper function to calculate statistics for a good
func calculateGoodStats(good models.Good) GoodWithStats {
	stats := GoodWithStats{
		Good: good,
	}

	now := time.Now()

	// Calculate warranty statistics
	if good.WarrantyStartDate != nil && good.WarrantyEndDate != nil {
		startDate := *good.WarrantyStartDate
		endDate := *good.WarrantyEndDate

		totalDays := int(endDate.Sub(startDate).Hours() / 24)
		daysPassed := int(now.Sub(startDate).Hours() / 24)
		daysRemaining := int(endDate.Sub(now).Hours() / 24)

		if daysPassed < 0 {
			daysPassed = 0
		}
		if daysRemaining < 0 {
			daysRemaining = 0
		}

		stats.WarrantyDaysTotal = &totalDays
		stats.WarrantyDaysPassed = &daysPassed
		stats.WarrantyDaysRemaining = &daysRemaining

		if now.After(endDate) {
			stats.WarrantyStatus = "expired"
		} else {
			stats.WarrantyStatus = "active"
		}
	} else {
		stats.WarrantyStatus = "no_warranty"
	}

	// Calculate days owned
	stats.DaysOwned = int(now.Sub(good.PurchaseDate.Time).Hours() / 24)
	if stats.DaysOwned < 1 {
		stats.DaysOwned = 1 // Avoid division by zero
	}

	// Calculate total service cost
	for _, service := range good.ServiceRecords {
		stats.TotalServiceCost += service.Cost
		stats.ServiceCount++
	}

	// Calculate total cost
	stats.TotalCost = good.PurchasePrice + stats.TotalServiceCost

	// Calculate price per day
	stats.PricePerDay = stats.TotalCost / float64(stats.DaysOwned)

	return stats
}

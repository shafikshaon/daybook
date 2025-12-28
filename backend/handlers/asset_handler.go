package handlers

import (
	"net/http"

	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/repository"
	"daybook-backend/services"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssetHandler struct {
	service services.AssetService
}

func NewAssetHandler(service services.AssetService) *AssetHandler {
	return &AssetHandler{
		service: service,
	}
}

// ListAssets returns all assets for the authenticated user
func (h *AssetHandler) ListAssets(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListAssets - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Build filters
	filters := repository.AssetFilters{}
	if status := c.Query("status"); status != "" {
		filters.Status = &status
	}
	if category := c.Query("category"); category != "" {
		filters.Category = &category
	}

	logger.Debugf(ctx, "Fetching assets for user")
	assets, err := h.service.ListAssets(ctx, userID, filters)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch assets")
		return
	}

	logger.Infof(ctx, "Assets retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, assets, "Assets retrieved successfully")
}

// GetAsset returns a specific asset by ID with statistics
func (h *AssetHandler) GetAsset(c *gin.Context) {
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
	asset, err := h.service.GetAsset(ctx, assetID, userID)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Asset not found")
		return
	}

	logger.Infof(ctx, "Asset retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, asset, "Asset retrieved successfully")
}

// CreateAsset creates a new asset record
func (h *AssetHandler) CreateAsset(c *gin.Context) {
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
	createdAsset, err := h.service.CreateAsset(ctx, &asset)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create asset")
		return
	}

	logger.Infof(ctx, "Asset created successfully for user: %s", userID)
	utilities.CreatedResponse(c, createdAsset, "Asset created successfully")
}

// UpdateAsset updates an asset record
func (h *AssetHandler) UpdateAsset(c *gin.Context) {
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

	var updateData models.Asset
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "Updating asset with ID: %s", assetID)
	updatedAsset, err := h.service.UpdateAsset(ctx, assetID, userID, &updateData)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update asset")
		return
	}

	logger.Infof(ctx, "Asset updated successfully for user: %s", userID)
	utilities.SuccessResponse(c, updatedAsset, "Asset updated successfully")
}

// DeleteAsset soft deletes an asset record
func (h *AssetHandler) DeleteAsset(c *gin.Context) {
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

	logger.Debugf(ctx, "Deleting asset with ID: %s", assetID)
	if err := h.service.DeleteAsset(ctx, assetID, userID); err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete asset")
		return
	}

	logger.Infof(ctx, "Asset deleted successfully for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Asset deleted successfully")
}

// CreateServiceRecord adds a service record for an asset
func (h *AssetHandler) CreateServiceRecord(c *gin.Context) {
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

	logger.Debugf(ctx, "Creating service record for asset ID: %s", assetID)
	createdRecord, err := h.service.CreateServiceRecord(ctx, &serviceRecord)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create service record")
		return
	}

	logger.Infof(ctx, "Service record created successfully for user: %s", userID)
	utilities.CreatedResponse(c, createdRecord, "Service record created successfully")
}

// ListServiceRecords returns all service records for a specific asset
func (h *AssetHandler) ListServiceRecords(c *gin.Context) {
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

	logger.Debugf(ctx, "Fetching service records for asset ID: %s", assetID)
	serviceRecords, err := h.service.ListServiceRecords(ctx, assetID, userID)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Failed to fetch service records")
		return
	}

	logger.Infof(ctx, "Service records retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, serviceRecords, "Service records retrieved successfully")
}

// DeleteServiceRecord deletes a service record
func (h *AssetHandler) DeleteServiceRecord(c *gin.Context) {
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

	logger.Debugf(ctx, "Deleting service record with ID: %s", serviceID)
	if err := h.service.DeleteServiceRecord(ctx, serviceID, userID); err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete service record")
		return
	}

	logger.Infof(ctx, "Service record deleted successfully for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Service record deleted successfully")
}

// AddAttachment links an uploaded file to an asset
func (h *AssetHandler) AddAttachment(c *gin.Context) {
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

	var attachment models.AssetAttachment
	if err := c.ShouldBindJSON(&attachment); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	attachment.UserID = userID
	attachment.AssetID = assetID

	logger.Debugf(ctx, "Adding attachment for asset ID: %s", assetID)
	createdAttachment, err := h.service.AddAttachment(ctx, &attachment)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to add attachment")
		return
	}

	logger.Infof(ctx, "Attachment added successfully for user: %s", userID)
	utilities.CreatedResponse(c, createdAttachment, "Attachment added successfully")
}

// ListAttachments returns all attachments for a specific asset
func (h *AssetHandler) ListAttachments(c *gin.Context) {
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

	logger.Debugf(ctx, "Fetching attachments for asset ID: %s", assetID)
	attachments, err := h.service.ListAttachments(ctx, assetID, userID)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusNotFound, "Failed to fetch attachments")
		return
	}

	logger.Infof(ctx, "Attachments retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, attachments, "Attachments retrieved successfully")
}

// DeleteAttachment deletes an attachment
func (h *AssetHandler) DeleteAttachment(c *gin.Context) {
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

	logger.Debugf(ctx, "Deleting attachment with ID: %s", attachmentID)
	if err := h.service.DeleteAttachment(ctx, attachmentID, userID); err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete attachment")
		return
	}

	logger.Infof(ctx, "Attachment deleted successfully for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Attachment deleted successfully")
}

// GetAssetsStats returns summary statistics for all assets
func (h *AssetHandler) GetAssetsStats(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetAssetsStats - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Fetching asset statistics for user")
	stats, err := h.service.GetStats(ctx, userID)
	if err != nil {
		logger.Errorf(ctx, "Service operation failed: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch statistics")
		return
	}

	logger.Infof(ctx, "Asset statistics retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, stats, "Statistics retrieved successfully")
}

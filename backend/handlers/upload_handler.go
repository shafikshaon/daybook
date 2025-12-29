package handlers

import (
	"fmt"
	"net/http"

	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/services"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
)

// UploadHandler handles file upload-related HTTP requests
type UploadHandler struct {
	service services.UploadService
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(service services.UploadService) *UploadHandler {
	return &UploadHandler{service: service}
}

// UploadFiles handles multiple file uploads
func (h *UploadHandler) UploadFiles(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UploadFiles - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	// Parse multipart form with max memory of 32 MB
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		logger.Warnf(ctx, "Failed to parse multipart form: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "File too large or invalid form data")
		return
	}

	form := c.Request.MultipartForm
	files := form.File["files"]

	if len(files) == 0 {
		logger.Warnf(ctx, "No files provided in request")
		utilities.ErrorResponse(c, http.StatusBadRequest, "No files provided")
		return
	}

	logger.Debugf(ctx, "Processing %d files...", len(files))

	// Upload files using service
	uploadedFiles, errors, err := h.service.UploadFiles(userID, files)
	if err != nil {
		logger.Errorf(ctx, "Failed to upload files: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Prepare response
	response := gin.H{
		"files":         uploadedFiles,
		"uploadedCount": len(uploadedFiles),
		"totalFiles":    len(files),
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	if len(uploadedFiles) == 0 {
		logger.Warnf(ctx, "No files were uploaded successfully")
		utilities.ErrorResponse(c, http.StatusBadRequest, "No files were uploaded successfully")
		return
	}

	logger.Infof(ctx, "Successfully uploaded %d files for user: %s", len(uploadedFiles), userID)
	utilities.SuccessResponse(c, response, "Files uploaded successfully")
}

// UploadSingleFile handles single file upload
func (h *UploadHandler) UploadSingleFile(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UploadSingleFile - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Processing request for user: %s", userID)

	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		logger.Warnf(ctx, "No file provided: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "No file provided")
		return
	}

	logger.Debugf(ctx, "Uploading file: %s", fileHeader.Filename)

	// Upload file using service
	response, err := h.service.UploadFile(userID, fileHeader)
	if err != nil {
		logger.Errorf(ctx, "Failed to upload file: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Successfully uploaded file for user: %s", userID)
	utilities.CreatedResponse(c, response, "File uploaded successfully")
}

// ServeUploadedFile serves the uploaded files
func (h *UploadHandler) ServeUploadedFile(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ServeUploadedFile - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	requestedUserID := c.Param("userId")
	filename := c.Param("filename")

	logger.Debugf(ctx, "Processing request for user: %s, requested user: %s, file: %s", userID, requestedUserID, filename)

	// Verify user can only access their own files
	if fmt.Sprintf("%d", userID) != requestedUserID {
		logger.Warnf(ctx, "Access denied: user %s tried to access files of user %s", userID, requestedUserID)
		utilities.ErrorResponse(c, http.StatusForbidden, "Access denied")
		return
	}

	// Get file path from service (includes security validation)
	filePath, err := h.service.GetFilePath(userID, filename)
	if err != nil {
		logger.Errorf(ctx, "Failed to get file path: %v", err)
		if err.Error() == "file not found" {
			utilities.ErrorResponse(c, http.StatusNotFound, "File not found")
		} else {
			utilities.ErrorResponse(c, http.StatusForbidden, err.Error())
		}
		return
	}

	logger.Infof(ctx, "Serving file for user: %s", userID)
	c.File(filePath)
}

// DeleteFile deletes an uploaded file
func (h *UploadHandler) DeleteFile(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteFile - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	filename := c.Param("filename")

	logger.Debugf(ctx, "Processing request for user: %s, file: %s", userID, filename)

	// Delete file using service
	if err := h.service.DeleteFile(userID, filename); err != nil {
		logger.Errorf(ctx, "Failed to delete file: %v", err)
		if err.Error() == "file not found" {
			utilities.ErrorResponse(c, http.StatusNotFound, "File not found")
		} else {
			utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	logger.Infof(ctx, "Successfully deleted file for user: %s", userID)
	utilities.SuccessResponse(c, gin.H{
		"filename": filename,
	}, "File deleted successfully")
}

// GetFileInfo returns information about an uploaded file
func (h *UploadHandler) GetFileInfo(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetFileInfo - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	filename := c.Param("filename")

	logger.Debugf(ctx, "Processing request for user: %s, file: %s", userID, filename)

	// Get file info from service
	response, err := h.service.GetFileInfo(userID, filename)
	if err != nil {
		logger.Errorf(ctx, "Failed to get file info: %v", err)
		if err.Error() == "file not found" {
			utilities.ErrorResponse(c, http.StatusNotFound, "File not found")
		} else {
			utilities.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	logger.Infof(ctx, "Successfully retrieved file info for user: %s", userID)
	utilities.SuccessResponse(c, response, "File info retrieved successfully")
}

package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"daybook-backend/middleware"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	MaxFileSize       = 10 << 20 // 10 MB
	MaxFilesPerUpload = 10
	UploadDir         = "./uploads"
)

var AllowedFileTypes = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".txt":  true,
	".csv":  true,
}

// FileUploadResponse represents the response for uploaded files
type FileUploadResponse struct {
	FileName     string `json:"fileName"`
	OriginalName string `json:"originalName"`
	FilePath     string `json:"filePath"`
	FileURL      string `json:"fileUrl"`
	FileSize     int64  `json:"fileSize"`
	MimeType     string `json:"mimeType"`
}

// UploadFiles handles multiple file uploads
func UploadFiles(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse multipart form with max memory of 32 MB
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "File too large or invalid form data")
		return
	}

	form := c.Request.MultipartForm
	files := form.File["files"]

	if len(files) == 0 {
		utilities.ErrorResponse(c, http.StatusBadRequest, "No files provided")
		return
	}

	if len(files) > MaxFilesPerUpload {
		utilities.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Maximum %d files allowed per upload", MaxFilesPerUpload))
		return
	}

	// Create uploads directory if it doesn't exist
	userUploadDir := filepath.Join(UploadDir, userID.String())
	if err := os.MkdirAll(userUploadDir, 0755); err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	var uploadedFiles []FileUploadResponse
	var errors []string

	for _, fileHeader := range files {
		// Validate file size
		if fileHeader.Size > MaxFileSize {
			errors = append(errors, fmt.Sprintf("File %s exceeds maximum size of 10MB", fileHeader.Filename))
			continue
		}

		// Validate file type
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if !AllowedFileTypes[ext] {
			errors = append(errors, fmt.Sprintf("File type %s not allowed for %s", ext, fileHeader.Filename))
			continue
		}

		// Open the uploaded file
		file, err := fileHeader.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to open file %s", fileHeader.Filename))
			continue
		}
		defer file.Close()

		// Generate unique filename
		uniqueFilename := generateUniqueFilename(fileHeader.Filename)
		filePath := filepath.Join(userUploadDir, uniqueFilename)

		// Create destination file
		dst, err := os.Create(filePath)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to save file %s", fileHeader.Filename))
			continue
		}
		defer dst.Close()

		// Copy file content
		if _, err := io.Copy(dst, file); err != nil {
			os.Remove(filePath) // Clean up on error
			errors = append(errors, fmt.Sprintf("Failed to write file %s", fileHeader.Filename))
			continue
		}

		// Build file URL
		fileURL := fmt.Sprintf("/api/v1/uploads/%s/%s", userID.String(), uniqueFilename)

		// Add to successful uploads
		uploadedFiles = append(uploadedFiles, FileUploadResponse{
			FileName:     uniqueFilename,
			OriginalName: fileHeader.Filename,
			FilePath:     filePath,
			FileURL:      fileURL,
			FileSize:     fileHeader.Size,
			MimeType:     fileHeader.Header.Get("Content-Type"),
		})
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
		utilities.ErrorResponse(c, http.StatusBadRequest, "No files were uploaded successfully")
		return
	}

	utilities.SuccessResponse(c, response, "Files uploaded successfully")
}

// UploadSingleFile handles single file upload
func UploadSingleFile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	// Validate file size
	if fileHeader.Size > MaxFileSize {
		utilities.ErrorResponse(c, http.StatusBadRequest, "File exceeds maximum size of 10MB")
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !AllowedFileTypes[ext] {
		utilities.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("File type %s not allowed", ext))
		return
	}

	// Create uploads directory if it doesn't exist
	userUploadDir := filepath.Join(UploadDir, userID.String())
	if err := os.MkdirAll(userUploadDir, 0755); err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	// Generate unique filename
	uniqueFilename := generateUniqueFilename(fileHeader.Filename)
	filePath := filepath.Join(userUploadDir, uniqueFilename)

	// Save file
	if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Build file URL
	fileURL := fmt.Sprintf("/api/v1/uploads/%s/%s", userID.String(), uniqueFilename)

	response := FileUploadResponse{
		FileName:     uniqueFilename,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileURL:      fileURL,
		FileSize:     fileHeader.Size,
		MimeType:     fileHeader.Header.Get("Content-Type"),
	}

	utilities.CreatedResponse(c, response, "File uploaded successfully")
}

// ServeUploadedFile serves the uploaded files
func ServeUploadedFile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	requestedUserID := c.Param("userId")
	filename := c.Param("filename")

	// Verify user can only access their own files
	if userID.String() != requestedUserID {
		utilities.ErrorResponse(c, http.StatusForbidden, "Access denied")
		return
	}

	filePath := filepath.Join(UploadDir, requestedUserID, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utilities.ErrorResponse(c, http.StatusNotFound, "File not found")
		return
	}

	c.File(filePath)
}

// DeleteFile deletes an uploaded file
func DeleteFile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	filename := c.Param("filename")

	filePath := filepath.Join(UploadDir, userID.String(), filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utilities.ErrorResponse(c, http.StatusNotFound, "File not found")
		return
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete file")
		return
	}

	utilities.SuccessResponse(c, gin.H{
		"filename": filename,
	}, "File deleted successfully")
}

// generateUniqueFilename generates a unique filename with timestamp and UUID
func generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	nameWithoutExt := strings.TrimSuffix(originalFilename, ext)

	// Sanitize filename
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, " ", "_")
	nameWithoutExt = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return '_'
	}, nameWithoutExt)

	timestamp := time.Now().Unix()
	uniqueID := uuid.New().String()[:8]

	return fmt.Sprintf("%s_%d_%s%s", nameWithoutExt, timestamp, uniqueID, ext)
}

// ValidateAndSanitizeFilename validates and sanitizes a filename
func ValidateAndSanitizeFilename(filename string) (string, error) {
	if filename == "" {
		return "", fmt.Errorf("filename cannot be empty")
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if !AllowedFileTypes[ext] {
		return "", fmt.Errorf("file type %s not allowed", ext)
	}

	// Remove any path traversal attempts
	filename = filepath.Base(filename)

	return filename, nil
}

// GetFileInfo returns information about an uploaded file
func GetFileInfo(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	filename := c.Param("filename")

	filePath := filepath.Join(UploadDir, userID.String(), filename)

	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		utilities.ErrorResponse(c, http.StatusNotFound, "File not found")
		return
	}

	// Open file to detect mime type
	file, err := os.Open(filePath)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to read file")
		return
	}
	defer file.Close()

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to read file")
		return
	}

	mimeType := http.DetectContentType(buffer)
	fileURL := fmt.Sprintf("/api/v1/uploads/%s/%s", userID.String(), filename)

	response := FileUploadResponse{
		FileName:     filename,
		OriginalName: filename,
		FilePath:     filePath,
		FileURL:      fileURL,
		FileSize:     fileInfo.Size(),
		MimeType:     mimeType,
	}

	utilities.SuccessResponse(c, response, "File info retrieved successfully")
}

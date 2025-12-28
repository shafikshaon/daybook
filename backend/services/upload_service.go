package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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

// UploadService handles file upload business logic
type UploadService interface {
	// UploadFile uploads a single file
	UploadFile(userID uuid.UUID, fileHeader *multipart.FileHeader) (*FileUploadResponse, error)

	// UploadFiles uploads multiple files
	UploadFiles(userID uuid.UUID, fileHeaders []*multipart.FileHeader) ([]FileUploadResponse, []string, error)

	// DeleteFile deletes a file
	DeleteFile(userID uuid.UUID, filename string) error

	// GetFileInfo retrieves file information
	GetFileInfo(userID uuid.UUID, filename string) (*FileUploadResponse, error)

	// GetFilePath builds the file path for serving
	GetFilePath(userID uuid.UUID, filename string) (string, error)
}

type uploadService struct {
}

// NewUploadService creates a new upload service
func NewUploadService() UploadService {
	return &uploadService{}
}

// UploadFile uploads a single file
func (s *uploadService) UploadFile(userID uuid.UUID, fileHeader *multipart.FileHeader) (*FileUploadResponse, error) {
	// Validate file size
	if fileHeader.Size > MaxFileSize {
		return nil, fmt.Errorf("file exceeds maximum size of 10MB")
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !AllowedFileTypes[ext] {
		return nil, fmt.Errorf("file type %s not allowed", ext)
	}

	// Create uploads directory
	userUploadDir := filepath.Join(UploadDir, userID.String())
	if err := os.MkdirAll(userUploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate unique filename
	uniqueFilename := generateUniqueFilename(fileHeader.Filename)
	filePath := filepath.Join(userUploadDir, uniqueFilename)

	// Open source file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(filePath) // Clean up on error
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	// Build file URL
	fileURL := fmt.Sprintf("/api/v1/uploads/%s/%s", userID.String(), uniqueFilename)

	return &FileUploadResponse{
		FileName:     uniqueFilename,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileURL:      fileURL,
		FileSize:     fileHeader.Size,
		MimeType:     fileHeader.Header.Get("Content-Type"),
	}, nil
}

// UploadFiles uploads multiple files
func (s *uploadService) UploadFiles(userID uuid.UUID, fileHeaders []*multipart.FileHeader) ([]FileUploadResponse, []string, error) {
	if len(fileHeaders) == 0 {
		return nil, nil, fmt.Errorf("no files provided")
	}

	if len(fileHeaders) > MaxFilesPerUpload {
		return nil, nil, fmt.Errorf("maximum %d files allowed per upload", MaxFilesPerUpload)
	}

	var uploadedFiles []FileUploadResponse
	var errors []string

	for _, fileHeader := range fileHeaders {
		response, err := s.UploadFile(userID, fileHeader)
		if err != nil {
			errors = append(errors, fmt.Sprintf("File %s: %v", fileHeader.Filename, err))
			continue
		}
		uploadedFiles = append(uploadedFiles, *response)
	}

	return uploadedFiles, errors, nil
}

// DeleteFile deletes a file
func (s *uploadService) DeleteFile(userID uuid.UUID, filename string) error {
	// Sanitize filename
	cleanFilename := filepath.Base(filename)
	if cleanFilename == "." || cleanFilename == ".." {
		return fmt.Errorf("invalid filename")
	}

	filePath := filepath.Join(UploadDir, userID.String(), cleanFilename)

	// Verify path is within user's directory
	userDir := filepath.Join(UploadDir, userID.String())
	if !filepath.HasPrefix(filePath, userDir) {
		return fmt.Errorf("access denied")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found")
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// GetFileInfo retrieves file information
func (s *uploadService) GetFileInfo(userID uuid.UUID, filename string) (*FileUploadResponse, error) {
	// Sanitize filename
	cleanFilename := filepath.Base(filename)
	filePath := filepath.Join(UploadDir, userID.String(), cleanFilename)

	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found")
	}

	// Open file to detect mime type
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	defer file.Close()

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	mimeType := http.DetectContentType(buffer)
	fileURL := fmt.Sprintf("/api/v1/uploads/%s/%s", userID.String(), filename)

	return &FileUploadResponse{
		FileName:     filename,
		OriginalName: filename,
		FilePath:     filePath,
		FileURL:      fileURL,
		FileSize:     fileInfo.Size(),
		MimeType:     mimeType,
	}, nil
}

// GetFilePath builds the file path for serving and validates access
func (s *uploadService) GetFilePath(userID uuid.UUID, filename string) (string, error) {
	// Sanitize filename
	cleanFilename := filepath.Base(filename)
	if cleanFilename == "." || cleanFilename == ".." {
		return "", fmt.Errorf("invalid filename")
	}

	filePath := filepath.Join(UploadDir, userID.String(), cleanFilename)

	// Verify path is within user's directory
	userDir := filepath.Join(UploadDir, userID.String())
	if !filepath.HasPrefix(filePath, userDir) {
		return "", fmt.Errorf("access denied")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	}

	return filePath, nil
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

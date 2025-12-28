package handlers

import (
	"net/http"

	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/services"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	service services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Signup creates a new user account
func (h *AuthHandler) Signup(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "Signup handler - Entry")

	var req models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf(ctx, "Invalid signup request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Processing signup request for username: %s, email: %s", req.Username, req.Email)

	response, err := h.service.Signup(ctx, &req)
	if err != nil {
		logger.Errorf(ctx, "Signup failed: %v", err)
		if err.Error() == "username or email already exists" {
			utilities.ErrorResponse(c, http.StatusConflict, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	logger.Infof(ctx, "User signup completed successfully for user: %s", response.User.ID)
	utilities.CreatedResponse(c, response, "User registered successfully")
}

// Login authenticates a user and returns a JWT token
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "Login handler - Entry")

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf(ctx, "Invalid login request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Processing login request for username/email: %s", req.Username)

	response, err := h.service.Login(ctx, &req)
	if err != nil {
		logger.Warnf(ctx, "Login failed: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	logger.Infof(ctx, "User login completed successfully for user: %s (ID: %s)", response.User.Username, response.User.ID)
	utilities.SuccessResponse(c, response, "Login successful")
}

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetProfile handler - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access to profile: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Fetching profile for user ID: %s", userID)

	user, err := h.service.GetProfile(ctx, userID)
	if err != nil {
		logger.Errorf(ctx, "User not found: %s, error: %v", userID, err)
		utilities.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	logger.Infof(ctx, "Profile retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, user, "Profile retrieved successfully")
}

// UpdateProfile updates the current user's profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateProfile handler - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access to update profile: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf(ctx, "Invalid update profile request: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Updating profile for user ID: %s", userID)

	user, err := h.service.UpdateProfile(ctx, userID, &req)
	if err != nil {
		logger.Errorf(ctx, "Failed to update profile: %v", err)
		if err.Error() == "email already in use" {
			utilities.ErrorResponse(c, http.StatusConflict, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	logger.Infof(ctx, "Profile updated successfully for user: %s", userID)
	utilities.SuccessResponse(c, user, "Profile updated successfully")
}

// ChangePassword changes the current user's password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ChangePassword handler - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access to change password: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf(ctx, "Invalid change password request: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Processing password change for user ID: %s", userID)

	if err := h.service.ChangePassword(ctx, userID, &req); err != nil {
		logger.Errorf(ctx, "Failed to change password: %v", err)
		if err.Error() == "current password is incorrect" {
			utilities.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		} else {
			utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	logger.Infof(ctx, "Password changed successfully for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Password changed successfully")
}

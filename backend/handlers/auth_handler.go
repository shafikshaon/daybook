package handlers

import (
	"net/http"
	"time"

	"daybook-backend/database"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
)

// Signup creates a new user account
func Signup(c *gin.Context) {
	var req models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		utilities.ErrorResponse(c, http.StatusConflict, "Username or email already exists")
		return
	}

	// Hash the password
	hashedPassword, err := utilities.HashPassword(req.Password)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create new user
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
		Role:     "user",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Create default settings for the user
	settings := models.Settings{
		UserID:         user.ID,
		Currency:       "USD",
		DarkMode:       false,
		DateFormat:     "MM/DD/YYYY",
		FirstDayOfWeek: 0,
		Language:       "en",
		Notifications: &models.Notifications{
			Push:          true,
			Email:         true,
			BudgetAlerts:  true,
			BillReminders: true,
		},
	}
	database.DB.Create(&settings)

	// Generate JWT token
	token, err := utilities.GenerateToken(&user)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  &user,
	}

	utilities.CreatedResponse(c, response, "User registered successfully")
}

// Login authenticates a user and returns a JWT token
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Find user by username
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	// Check password
	if err := utilities.CheckPassword(user.Password, req.Password); err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	database.DB.Save(&user)

	// Generate JWT token
	token, err := utilities.GenerateToken(&user)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  &user,
	}

	utilities.SuccessResponse(c, response, "Login successful")
}

// GetProfile returns the current user's profile
func GetProfile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utilities.SuccessResponse(c, user, "Profile retrieved successfully")
}

// UpdateProfile updates the current user's profile
func UpdateProfile(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	// Check if email is being changed and if it's already in use
	if req.Email != "" && req.Email != user.Email {
		var existingUser models.User
		if err := database.DB.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser).Error; err == nil {
			utilities.ErrorResponse(c, http.StatusConflict, "Email already in use")
			return
		}
		user.Email = req.Email
	}

	// Update profile fields
	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if err := database.DB.Save(&user).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	utilities.SuccessResponse(c, user, "Profile updated successfully")
}

// ChangePassword changes the current user's password
func ChangePassword(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	// Verify current password
	if err := utilities.CheckPassword(user.Password, req.CurrentPassword); err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Current password is incorrect")
		return
	}

	// Hash new password
	hashedPassword, err := utilities.HashPassword(req.NewPassword)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Update password
	user.Password = hashedPassword
	if err := database.DB.Save(&user).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update password")
		return
	}

	utilities.SuccessResponse(c, nil, "Password changed successfully")
}

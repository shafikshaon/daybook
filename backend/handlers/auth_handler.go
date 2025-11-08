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
)

// Signup creates a new user account
func Signup(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "Signup handler - Entry")

	var req models.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf(ctx, "Invalid signup request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Processing signup request for username: %s, email: %s", req.Username, req.Email)

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.WithContext(ctx).Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		logger.Warnf(ctx, "Signup failed - username or email already exists: %s", req.Username)
		utilities.ErrorResponse(c, http.StatusConflict, "Username or email already exists")
		return
	}

	// Hash the password
	logger.Debugf(ctx, "Hashing password for user: %s", req.Username)
	hashedPassword, err := utilities.HashPassword(req.Password)
	if err != nil {
		logger.Errorf(ctx, "Failed to hash password: %v", err)
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

	logger.Infof(ctx, "Creating new user: %s", req.Username)
	if err := database.DB.WithContext(ctx).Create(&user).Error; err != nil {
		logger.Errorf(ctx, "Failed to create user in database: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user")
		return
	}
	logger.Infof(ctx, "User created successfully with ID: %s", user.ID)

	// Create default settings for the user
	logger.Debugf(ctx, "Creating default settings for user: %s", user.ID)
	settings := models.Settings{
		UserID:         user.ID,
		Currency:       "BDT",
		DarkMode:       false,
		DateFormat:     "MM/DD/YYYY",
		FirstDayOfWeek: 0,
		Language:       "en",
		Notifications: &models.Notifications{
			Push:         true,
			Email:        true,
			BudgetAlerts: true,
		},
	}
	database.DB.WithContext(ctx).Create(&settings)
	logger.Debugf(ctx, "Default settings created for user: %s", user.ID)

	// Create default account types for the user
	logger.Debugf(ctx, "Creating default account types for user: %s", user.ID)
	if err := models.SeedDefaultAccountTypes(database.DB.WithContext(ctx), user.ID); err != nil {
		logger.Errorf(ctx, "Failed to create default account types: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create default account types")
		return
	}
	logger.Debugf(ctx, "Default account types created for user: %s", user.ID)

	// Generate JWT token
	logger.Debugf(ctx, "Generating JWT token for user: %s", user.ID)
	token, err := utilities.GenerateToken(&user)
	if err != nil {
		logger.Errorf(ctx, "Failed to generate JWT token: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  &user,
	}

	// Log signup activity
	utilities.LogAuthActivity(c, user.ID, models.ActionCreate, "User registered successfully")

	logger.Infof(ctx, "User signup completed successfully for user: %s", user.ID)
	utilities.CreatedResponse(c, response, "User registered successfully")
}

// Login authenticates a user and returns a JWT token
func Login(c *gin.Context) {
	ctx := middleware.GetContext(c)
	logger.Infof(ctx, "Login handler - Entry")

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf(ctx, "Invalid login request body: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Infof(ctx, "Processing login request for username/email: %s", req.Username)

	// Find user by username or email
	var user models.User
	if err := database.DB.WithContext(ctx).Where("LOWER(username) = LOWER(?) OR LOWER(email) = LOWER(?)",
		req.Username, req.Username).First(&user).Error; err != nil {
		logger.Warnf(ctx, "Login failed - user not found: %s", req.Username)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	logger.Debugf(ctx, "User found: %s (ID: %s)", user.Username, user.ID)

	// Check password
	if err := utilities.CheckPassword(user.Password, req.Password); err != nil {
		logger.Warnf(ctx, "Login failed - invalid password for user: %s", user.Username)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	logger.Debugf(ctx, "Password verified successfully for user: %s", user.Username)

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	database.DB.WithContext(ctx).Save(&user)
	logger.Debugf(ctx, "Last login updated for user: %s", user.ID)

	// Generate JWT token
	logger.Debugf(ctx, "Generating JWT token for user: %s", user.ID)
	token, err := utilities.GenerateToken(&user)
	if err != nil {
		logger.Errorf(ctx, "Failed to generate JWT token: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  &user,
	}

	// Log login activity
	utilities.LogAuthActivity(c, user.ID, models.ActionLogin, "User logged in successfully")

	logger.Infof(ctx, "User login completed successfully for user: %s (ID: %s)", user.Username, user.ID)
	utilities.SuccessResponse(c, response, "Login successful")
}

// GetProfile returns the current user's profile
func GetProfile(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetProfile handler - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "Unauthorized access to profile: %v", err)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "Fetching profile for user ID: %s", userID)

	var user models.User
	if err := database.DB.WithContext(ctx).First(&user, userID).Error; err != nil {
		logger.Errorf(ctx, "User not found in database: %s, error: %v", userID, err)
		utilities.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	logger.Infof(ctx, "Profile retrieved successfully for user: %s", userID)
	utilities.SuccessResponse(c, user, "Profile retrieved successfully")
}

// UpdateProfile updates the current user's profile
func UpdateProfile(c *gin.Context) {
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

	var user models.User
	if err := database.DB.WithContext(ctx).First(&user, userID).Error; err != nil {
		logger.Errorf(ctx, "User not found: %s, error: %v", userID, err)
		utilities.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	// Check if email is being changed and if it's already in use
	if req.Email != "" && req.Email != user.Email {
		logger.Debugf(ctx, "Email change requested from %s to %s", user.Email, req.Email)
		var existingUser models.User
		if err := database.DB.WithContext(ctx).Where("email = ? AND id != ?", req.Email, userID).First(&existingUser).Error; err == nil {
			logger.Warnf(ctx, "Email already in use: %s", req.Email)
			utilities.ErrorResponse(c, http.StatusConflict, "Email already in use")
			return
		}
		user.Email = req.Email
	}

	// Update profile fields
	if req.FullName != "" {
		logger.Debugf(ctx, "Updating full name from '%s' to '%s'", user.FullName, req.FullName)
		user.FullName = req.FullName
	}

	if err := database.DB.WithContext(ctx).Save(&user).Error; err != nil {
		logger.Errorf(ctx, "Failed to update profile: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	// Log profile update activity
	utilities.LogAuthActivity(c, userID, models.ActionUpdate, "User profile updated")

	logger.Infof(ctx, "Profile updated successfully for user: %s", userID)
	utilities.SuccessResponse(c, user, "Profile updated successfully")
}

// ChangePassword changes the current user's password
func ChangePassword(c *gin.Context) {
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

	var user models.User
	if err := database.DB.WithContext(ctx).First(&user, userID).Error; err != nil {
		logger.Errorf(ctx, "User not found: %s, error: %v", userID, err)
		utilities.ErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	// Verify current password
	logger.Debugf(ctx, "Verifying current password for user: %s", userID)
	if err := utilities.CheckPassword(user.Password, req.CurrentPassword); err != nil {
		logger.Warnf(ctx, "Current password verification failed for user: %s", userID)
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Current password is incorrect")
		return
	}

	// Hash new password
	logger.Debugf(ctx, "Hashing new password for user: %s", userID)
	hashedPassword, err := utilities.HashPassword(req.NewPassword)
	if err != nil {
		logger.Errorf(ctx, "Failed to hash new password: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Update password
	user.Password = hashedPassword
	if err := database.DB.WithContext(ctx).Save(&user).Error; err != nil {
		logger.Errorf(ctx, "Failed to update password in database: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update password")
		return
	}

	// Log password change activity
	utilities.LogAuthActivity(c, userID, models.ActionUpdate, "User password changed")

	logger.Infof(ctx, "Password changed successfully for user: %s", userID)
	utilities.SuccessResponse(c, nil, "Password changed successfully")
}

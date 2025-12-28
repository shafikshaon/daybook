package services

import (
	"context"
	"errors"
	"time"

	"daybook-backend/models"
	"daybook-backend/repository"
	"daybook-backend/utilities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthService handles authentication and user management business logic
type AuthService interface {
	// Signup creates a new user with default settings, account types, and categories
	Signup(ctx context.Context, req *models.SignupRequest) (*models.LoginResponse, error)

	// Login authenticates a user and returns a JWT token
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)

	// GetProfile retrieves a user's profile
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error)

	// UpdateProfile updates a user's profile
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *models.UpdateProfileRequest) (*models.User, error)

	// ChangePassword changes a user's password
	ChangePassword(ctx context.Context, userID uuid.UUID, req *models.ChangePasswordRequest) error
}

type authService struct {
	userRepo        repository.UserRepository
	settingsRepo    repository.SettingsRepository
	categoryRepo    repository.CategoryRepository
	accountTypeRepo repository.AccountTypeRepository
	txManager       repository.TransactionManager
	activityLogger  ActivityLogService
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo repository.UserRepository,
	settingsRepo repository.SettingsRepository,
	categoryRepo repository.CategoryRepository,
	accountTypeRepo repository.AccountTypeRepository,
	txManager repository.TransactionManager,
	activityLogger ActivityLogService,
) AuthService {
	return &authService{
		userRepo:        userRepo,
		settingsRepo:    settingsRepo,
		categoryRepo:    categoryRepo,
		accountTypeRepo: accountTypeRepo,
		txManager:       txManager,
		activityLogger:  activityLogger,
	}
}

// Signup creates a new user with default settings and data
func (s *authService) Signup(ctx context.Context, req *models.SignupRequest) (*models.LoginResponse, error) {
	// Check if user already exists
	exists, err := s.userRepo.ExistsByUsernameOrEmail(ctx, req.Username, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username or email already exists")
	}

	// Hash the password
	hashedPassword, err := utilities.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
		Role:     "user",
	}

	// Use transaction to create user and default data atomically
	err = s.txManager.WithTransaction(ctx, func(ctx context.Context, tx *gorm.DB) error {
		// Create user
		userRepoTx := s.userRepo.WithTx(tx)
		if err := userRepoTx.Create(ctx, &user); err != nil {
			return err
		}

		// Create default settings
		settings := models.Settings{
			UserID:         user.ID,
			Currency:       "BDT",
			DarkMode:       false,
			DateFormat:     "DD/MM/YYYY",
			FirstDayOfWeek: 0,
			Language:       "en",
			Notifications: &models.Notifications{
				Push:         true,
				Email:        true,
				BudgetAlerts: true,
			},
		}
		if err := tx.WithContext(ctx).Create(&settings).Error; err != nil {
			return err
		}

		// Create default account types
		if err := models.SeedDefaultAccountTypes(tx.WithContext(ctx), user.ID); err != nil {
			return err
		}

		// Create default categories
		defaultCategories := models.GetDefaultCategories(user.ID)
		if len(defaultCategories) > 0 {
			if err := tx.WithContext(ctx).Create(&defaultCategories).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := utilities.GenerateToken(&user)
	if err != nil {
		return nil, err
	}

	// Log signup activity
	s.activityLogger.LogAuthActivity(
		ctx,
		user.ID,
		models.ActionCreate,
		"User registered successfully",
	)

	return &models.LoginResponse{
		Token: token,
		User:  &user,
	}, nil
}

// Login authenticates a user
func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Find user by username or email
	user, err := s.userRepo.FindByUsernameOrEmail(ctx, req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Check password
	if err := utilities.CheckPassword(user.Password, req.Password); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := utilities.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	// Log login activity
	s.activityLogger.LogAuthActivity(
		ctx,
		user.ID,
		models.ActionLogin,
		"User logged in successfully",
	)

	return &models.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetProfile retrieves a user's profile
func (s *authService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	// Note: BaseRepository.FindByID expects (id, userID) but for User, we only have userID
	// So we'll use the db directly here
	var user models.User
	err := s.userRepo.Query(ctx, userID).First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateProfile updates a user's profile
func (s *authService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *models.UpdateProfileRequest) (*models.User, error) {
	// Get existing user
	user, err := s.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if email is being changed and if it's already in use
	if req.Email != "" && req.Email != user.Email {
		exists, err := s.userRepo.ExistsByEmail(ctx, req.Email, userID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email already in use")
		}
		user.Email = req.Email
	}

	// Update profile fields
	if req.FullName != "" {
		user.FullName = req.FullName
	}

	// Save updates
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Log profile update activity
	s.activityLogger.LogAuthActivity(
		ctx,
		userID,
		models.ActionUpdate,
		"User profile updated",
	)

	return user, nil
}

// ChangePassword changes a user's password
func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, req *models.ChangePasswordRequest) error {
	// Get existing user
	user, err := s.GetProfile(ctx, userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := utilities.CheckPassword(user.Password, req.CurrentPassword); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utilities.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Log password change activity
	s.activityLogger.LogAuthActivity(
		ctx,
		userID,
		models.ActionUpdate,
		"User password changed",
	)

	return nil
}

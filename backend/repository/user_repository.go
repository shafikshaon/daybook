package repository

import (
	"context"

	"daybook-backend/models"

	"gorm.io/gorm"
)

// UserRepository handles user database operations
type UserRepository interface {
	BaseRepository[models.User]

	// FindByUserID finds a user by their ID (users table has no user_id column)
	FindByUserID(ctx context.Context, userID uint) (*models.User, error)

	// FindByUsername finds a user by username
	FindByUsername(ctx context.Context, username string) (*models.User, error)

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*models.User, error)

	// FindByUsernameOrEmail finds a user by username or email
	FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*models.User, error)

	// ExistsByUsernameOrEmail checks if a user exists by username or email
	ExistsByUsernameOrEmail(ctx context.Context, username, email string) (bool, error)

	// ExistsByEmail checks if a user exists by email (excluding a specific user ID)
	ExistsByEmail(ctx context.Context, email string, excludeUserID uint) (bool, error)
}

type userRepository struct {
	*GormBaseRepository[models.User]
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		GormBaseRepository: NewGormBaseRepository[models.User](db),
	}
}

// FindByUserID finds a user by their ID
func (r *userRepository) FindByUserID(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username
func (r *userRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("LOWER(username) = LOWER(?)", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("LOWER(email) = LOWER(?)", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsernameOrEmail finds a user by username or email
func (r *userRepository) FindByUsernameOrEmail(ctx context.Context, usernameOrEmail string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("LOWER(username) = LOWER(?) OR LOWER(email) = LOWER(?)", usernameOrEmail, usernameOrEmail).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ExistsByUsernameOrEmail checks if a user exists by username or email
func (r *userRepository) ExistsByUsernameOrEmail(ctx context.Context, username, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("username = ? OR email = ?", username, email).
		Count(&count).Error
	return count > 0, err
}

// ExistsByEmail checks if a user exists by email (excluding a specific user ID)
func (r *userRepository) ExistsByEmail(ctx context.Context, email string, excludeUserID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ? AND id != ?", email, excludeUserID).
		Count(&count).Error
	return count > 0, err
}

package repository

import (
	"context"

	"gorm.io/gorm"
)

// BaseRepository defines common CRUD operations for all repositories
// T is the entity type (e.g., models.Account, models.Transaction)
type BaseRepository[T any] interface {
	// Create inserts a new entity into the database
	Create(ctx context.Context, entity *T) error

	// FindByID retrieves an entity by its ID and user ID (user-scoped)
	FindByID(ctx context.Context, id uint, userID uint) (*T, error)

	// FindAll retrieves all entities for a specific user
	FindAll(ctx context.Context, userID uint) ([]T, error)

	// Update saves changes to an existing entity
	Update(ctx context.Context, entity *T) error

	// Delete removes an entity by ID (user-scoped, soft delete)
	Delete(ctx context.Context, id uint, userID uint) error

	// WithTx returns a new repository instance using the provided transaction
	// This allows repositories to participate in transactions
	WithTx(tx *gorm.DB) BaseRepository[T]

	// Query returns a query builder scoped to the user
	// Useful for custom queries in specific repositories
	Query(ctx context.Context, userID uint) *gorm.DB
}

// GormBaseRepository implements BaseRepository using GORM
type GormBaseRepository[T any] struct {
	db *gorm.DB
}

// NewGormBaseRepository creates a new generic GORM repository
func NewGormBaseRepository[T any](db *gorm.DB) *GormBaseRepository[T] {
	return &GormBaseRepository[T]{db: db}
}

// Create inserts a new entity
func (r *GormBaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// FindByID retrieves an entity by ID and user ID
func (r *GormBaseRepository[T]) FindByID(ctx context.Context, id uint, userID uint) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindAll retrieves all entities for a user
func (r *GormBaseRepository[T]) FindAll(ctx context.Context, userID uint) ([]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&entities).Error
	return entities, err
}

// Update saves changes to an entity
func (r *GormBaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete soft-deletes an entity
func (r *GormBaseRepository[T]) Delete(ctx context.Context, id uint, userID uint) error {
	var entity T
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&entity).Error
}

// WithTx returns a new repository using the provided transaction
func (r *GormBaseRepository[T]) WithTx(tx *gorm.DB) BaseRepository[T] {
	return &GormBaseRepository[T]{db: tx}
}

// Query returns a query builder scoped to the user
func (r *GormBaseRepository[T]) Query(ctx context.Context, userID uint) *gorm.DB {
	return r.db.WithContext(ctx).Where("user_id = ?", userID)
}

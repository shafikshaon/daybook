package repository

import (
	"context"
	"time"

	"daybook-backend/models"

	"gorm.io/gorm"
)

// ReconciliationFilters represents query filters for reconciliations
type ReconciliationFilters struct {
	AccountID *uint
}

// ReconciliationStats represents reconciliation statistics
type ReconciliationStats struct {
	TotalReconciliations       int64      `json:"totalReconciliations"`
	CompletedReconciliations   int64      `json:"completedReconciliations"`
	PendingReconciliations     int64      `json:"pendingReconciliations"`
	DiscrepancyReconciliations int64      `json:"discrepancyReconciliations"`
	LastReconciliationDate     *time.Time `json:"lastReconciliationDate"`
	AverageDifference          float64    `json:"averageDifference"`
}

// ReconciliationRepository handles reconciliation data access
type ReconciliationRepository interface {
	BaseRepository[models.Reconciliation]

	// FindWithFilters retrieves reconciliations with optional filters
	FindWithFilters(ctx context.Context, userID uint, filters ReconciliationFilters) ([]models.Reconciliation, error)

	// FindByIDWithPreloads retrieves a reconciliation with all relationships
	FindByIDWithPreloads(ctx context.Context, id, userID uint) (*models.Reconciliation, error)

	// GetUnreconciledTransactions retrieves unreconciled transactions for an account
	GetUnreconciledTransactions(ctx context.Context, accountID, userID uint) ([]models.Transaction, error)

	// GetStats calculates reconciliation statistics for an account
	GetStats(ctx context.Context, accountID, userID uint) (*ReconciliationStats, error)

	// LinkTransactions links transactions to a reconciliation
	LinkTransactions(ctx context.Context, reconciliationID uint, transactionIDs []uint) error

	// UnlinkAllTransactions unlinks all transactions from a reconciliation
	UnlinkAllTransactions(ctx context.Context, reconciliationID uint) error

	// CheckColumnExists checks if a column exists in the transactions table
	CheckColumnExists(ctx context.Context, columnName string) bool
}

type reconciliationRepository struct {
	*GormBaseRepository[models.Reconciliation]
}

// NewReconciliationRepository creates a new reconciliation repository
func NewReconciliationRepository(db *gorm.DB) ReconciliationRepository {
	return &reconciliationRepository{
		GormBaseRepository: NewGormBaseRepository[models.Reconciliation](db),
	}
}

// FindWithFilters retrieves reconciliations with optional filters
func (r *reconciliationRepository) FindWithFilters(ctx context.Context, userID uint, filters ReconciliationFilters) ([]models.Reconciliation, error) {
	var reconciliations []models.Reconciliation

	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Account")

	if filters.AccountID != nil {
		query = query.Where("account_id = ?", *filters.AccountID)
	}

	err := query.Order("reconciliation_date DESC").Find(&reconciliations).Error
	return reconciliations, err
}

// FindByIDWithPreloads retrieves a reconciliation with all relationships
func (r *reconciliationRepository) FindByIDWithPreloads(ctx context.Context, id, userID uint) (*models.Reconciliation, error) {
	var reconciliation models.Reconciliation
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Preload("Account").
		Preload("Transactions.Transaction").
		First(&reconciliation).Error
	if err != nil {
		return nil, err
	}
	return &reconciliation, nil
}

// GetUnreconciledTransactions retrieves unreconciled transactions for an account
func (r *reconciliationRepository) GetUnreconciledTransactions(ctx context.Context, accountID, userID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := r.db.WithContext(ctx).Where("account_id = ? AND user_id = ?", accountID, userID)

	// Check if reconciled column exists
	if r.CheckColumnExists(ctx, "reconciled") {
		query = query.Where("reconciled = ? OR reconciled IS NULL", false)
	}

	err := query.Order("date DESC").Find(&transactions).Error
	return transactions, err
}

// GetStats calculates reconciliation statistics for an account
func (r *reconciliationRepository) GetStats(ctx context.Context, accountID, userID uint) (*ReconciliationStats, error) {
	stats := &ReconciliationStats{}

	// Total reconciliations
	r.db.WithContext(ctx).Model(&models.Reconciliation{}).
		Where("account_id = ? AND user_id = ?", accountID, userID).
		Count(&stats.TotalReconciliations)

	// Completed reconciliations
	r.db.WithContext(ctx).Model(&models.Reconciliation{}).
		Where("account_id = ? AND user_id = ? AND status = ?", accountID, userID, models.ReconciliationCompleted).
		Count(&stats.CompletedReconciliations)

	// Pending reconciliations
	r.db.WithContext(ctx).Model(&models.Reconciliation{}).
		Where("account_id = ? AND user_id = ? AND status = ?", accountID, userID, models.ReconciliationPending).
		Count(&stats.PendingReconciliations)

	// Discrepancy reconciliations
	r.db.WithContext(ctx).Model(&models.Reconciliation{}).
		Where("account_id = ? AND user_id = ? AND status = ?", accountID, userID, models.ReconciliationDiscrepancy).
		Count(&stats.DiscrepancyReconciliations)

	// Last reconciliation date
	var lastReconciliation models.Reconciliation
	if err := r.db.WithContext(ctx).
		Where("account_id = ? AND user_id = ?", accountID, userID).
		Order("reconciliation_date DESC").
		First(&lastReconciliation).Error; err == nil {
		stats.LastReconciliationDate = &lastReconciliation.ReconciliationDate
	}

	// Average difference
	var avgDiff struct {
		AvgDifference float64
	}
	r.db.WithContext(ctx).Model(&models.Reconciliation{}).
		Select("AVG(ABS(difference)) as avg_difference").
		Where("account_id = ? AND user_id = ?", accountID, userID).
		Scan(&avgDiff)
	stats.AverageDifference = avgDiff.AvgDifference

	return stats, nil
}

// LinkTransactions links transactions to a reconciliation
func (r *reconciliationRepository) LinkTransactions(ctx context.Context, reconciliationID uint, transactionIDs []uint) error {
	for _, transactionID := range transactionIDs {
		reconciliationTransaction := models.ReconciliationTransaction{
			ReconciliationID: reconciliationID,
			TransactionID:    transactionID,
		}

		if err := r.db.WithContext(ctx).Create(&reconciliationTransaction).Error; err != nil {
			return err
		}

		// Mark transaction as reconciled
		if err := r.db.WithContext(ctx).
			Model(&models.Transaction{}).
			Where("id = ?", transactionID).
			Updates(map[string]interface{}{
				"reconciled":        true,
				"reconciliation_id": reconciliationID,
			}).Error; err != nil {
			return err
		}
	}
	return nil
}

// UnlinkAllTransactions unlinks all transactions from a reconciliation
func (r *reconciliationRepository) UnlinkAllTransactions(ctx context.Context, reconciliationID uint) error {
	// Get all reconciliation transactions
	var reconciliationTransactions []models.ReconciliationTransaction
	if err := r.db.WithContext(ctx).
		Where("reconciliation_id = ?", reconciliationID).
		Find(&reconciliationTransactions).Error; err != nil {
		return err
	}

	// Unmark all transactions
	for _, rt := range reconciliationTransactions {
		r.db.WithContext(ctx).
			Model(&models.Transaction{}).
			Where("id = ?", rt.TransactionID).
			Updates(map[string]interface{}{
				"reconciled":        false,
				"reconciliation_id": nil,
			})
	}

	// Delete reconciliation transaction links
	return r.db.WithContext(ctx).
		Where("reconciliation_id = ?", reconciliationID).
		Delete(&models.ReconciliationTransaction{}).Error
}

// CheckColumnExists checks if a column exists in a table
func (r *reconciliationRepository) CheckColumnExists(ctx context.Context, columnName string) bool {
	var exists bool
	r.db.WithContext(ctx).Raw(
		"SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'transactions' AND column_name = ?)",
		columnName,
	).Scan(&exists)
	return exists
}

// WithTx returns a new repository instance using the provided transaction
func (r *reconciliationRepository) WithTx(tx *gorm.DB) BaseRepository[models.Reconciliation] {
	return &reconciliationRepository{
		GormBaseRepository: NewGormBaseRepository[models.Reconciliation](tx),
	}
}

package handlers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"daybook-backend/database"
	"daybook-backend/logger"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
)

// ListFixedDeposits returns all fixed deposits for the authenticated user
func ListFixedDeposits(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "ListFixedDeposits - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "ListFixedDeposits - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	logger.Debugf(ctx, "ListFixedDeposits - Fetching fixed deposits for user: %s", userID)

	query := database.DB.WithContext(ctx).Where("user_id = ?", userID)

	// Optional filter by withdrawn status
	if withdrawn := c.Query("withdrawn"); withdrawn != "" {
		logger.Debugf(ctx, "ListFixedDeposits - Filtering by withdrawn status: %s", withdrawn)
		query = query.Where("withdrawn = ?", withdrawn == "true")
	}

	var deposits []models.FixedDeposit
	if err := query.Order("maturity_date ASC").Find(&deposits).Error; err != nil {
		logger.Errorf(ctx, "ListFixedDeposits - Failed to fetch fixed deposits: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch fixed deposits")
		return
	}

	logger.Infof(ctx, "ListFixedDeposits - Successfully retrieved %d fixed deposits", len(deposits))
	utilities.SuccessResponse(c, deposits, "Fixed deposits retrieved successfully")
}

// GetFixedDeposit returns a specific fixed deposit by ID
func GetFixedDeposit(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "GetFixedDeposit - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "GetFixedDeposit - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	depositIDStr := c.Param("id")
	depositIDUint, err := strconv.ParseUint(depositIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "GetFixedDeposit - Invalid fixed deposit ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid fixed deposit ID")
		return
	}
	depositID := uint(depositIDUint)

	logger.Debugf(ctx, "GetFixedDeposit - Fetching fixed deposit: %s for user: %s", depositID, userID)

	var deposit models.FixedDeposit
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", depositID, userID).First(&deposit).Error; err != nil {
		logger.Warnf(ctx, "GetFixedDeposit - Fixed deposit not found: %s", depositID)
		utilities.ErrorResponse(c, http.StatusNotFound, "Fixed deposit not found")
		return
	}

	logger.Infof(ctx, "GetFixedDeposit - Successfully retrieved fixed deposit: %s", depositID)
	utilities.SuccessResponse(c, deposit, "Fixed deposit retrieved successfully")
}

// calculateMaturityAmount calculates the maturity amount based on compounding
func calculateMaturityAmount(principal float64, rate float64, tenureMonths int, compounding string) float64 {
	// Convert annual rate to decimal
	r := rate / 100.0

	// Convert tenure to years
	t := float64(tenureMonths) / 12.0

	var maturityAmount float64

	switch compounding {
	case "simple":
		// Simple Interest: A = P(1 + rt)
		maturityAmount = principal * (1 + r*t)

	case "daily":
		// Daily compounding: A = P(1 + r/365)^(365*t)
		n := 365.0
		maturityAmount = principal * math.Pow(1+r/n, n*t)

	case "monthly":
		// Monthly compounding: A = P(1 + r/12)^(12*t)
		n := 12.0
		maturityAmount = principal * math.Pow(1+r/n, n*t)

	case "quarterly":
		// Quarterly compounding: A = P(1 + r/4)^(4*t)
		n := 4.0
		maturityAmount = principal * math.Pow(1+r/n, n*t)

	case "semi-annually":
		// Semi-annual compounding: A = P(1 + r/2)^(2*t)
		n := 2.0
		maturityAmount = principal * math.Pow(1+r/n, n*t)

	case "annually":
		// Annual compounding: A = P(1 + r)^t
		maturityAmount = principal * math.Pow(1+r, t)

	default:
		// Default to monthly compounding
		n := 12.0
		maturityAmount = principal * math.Pow(1+r/n, n*t)
	}

	return maturityAmount
}

// CreateFixedDeposit creates a new fixed deposit
func CreateFixedDeposit(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "CreateFixedDeposit - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "CreateFixedDeposit - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var depositData struct {
		models.FixedDeposit
		AccountID uint `json:"accountId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&depositData); err != nil {
		logger.Warnf(ctx, "CreateFixedDeposit - Invalid request data: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "CreateFixedDeposit - Creating fixed deposit at %s for user: %s", depositData.Institution, userID)

	depositData.FixedDeposit.UserID = userID

	// Verify account belongs to user
	var account models.Account
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", depositData.AccountID, userID).First(&account).Error; err != nil {
		logger.Warnf(ctx, "CreateFixedDeposit - Invalid account ID: %s", depositData.AccountID)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	// Check sufficient balance
	if account.Balance < depositData.Principal {
		logger.Warnf(ctx, "CreateFixedDeposit - Insufficient account balance. Required: %.2f, Available: %.2f", depositData.Principal, account.Balance)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Insufficient account balance")
		return
	}

	// Calculate maturity date from start date and tenure
	depositData.MaturityDate = depositData.StartDate.AddDate(0, depositData.TenureMonths, 0)

	// Calculate maturity amount
	depositData.MaturityAmount = calculateMaturityAmount(
		depositData.Principal,
		depositData.InterestRate,
		depositData.TenureMonths,
		depositData.Compounding,
	)

	logger.Debugf(ctx, "CreateFixedDeposit - Calculated maturity amount: %.2f", depositData.MaturityAmount)

	// Start transaction
	tx := database.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create fixed deposit
	if err := tx.Create(&depositData.FixedDeposit).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "CreateFixedDeposit - Failed to create fixed deposit: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create fixed deposit")
		return
	}

	// Create transaction record
	transaction := models.Transaction{
		UserID:         userID,
		AccountID:      depositData.AccountID,
		Type:           "expense",
		Amount:         depositData.Principal,
		CategoryID:     0, // Use 0 for system transactions
		Date:           models.Date{Time: depositData.StartDate},
		Description:    "Fixed Deposit: " + depositData.Institution + " - " + depositData.AccountNumber,
		FixedDepositID: &depositData.FixedDeposit.ID,
		Tags:           []string{"fixed_deposit"},
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "CreateFixedDeposit - Failed to create transaction: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	// Update account balance (debit)
	account.Balance -= depositData.Principal
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "CreateFixedDeposit - Failed to update account balance: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account balance")
		return
	}

	tx.Commit()

	logger.Infof(ctx, "CreateFixedDeposit - Successfully created fixed deposit: %s", depositData.FixedDeposit.ID)

	// Log fixed deposit creation activity
	utilities.LogEntityActivity(c, userID, models.ActionCreate, "fixed_deposit",
		"FixedDeposit", depositData.FixedDeposit.ID, "Created fixed deposit at "+depositData.Institution, nil)

	result := map[string]interface{}{
		"fixedDeposit": depositData.FixedDeposit,
		"transaction":  transaction,
	}

	utilities.CreatedResponse(c, result, "Fixed deposit created successfully")
}

// UpdateFixedDeposit updates an existing fixed deposit
func UpdateFixedDeposit(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "UpdateFixedDeposit - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "UpdateFixedDeposit - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	depositIDStr := c.Param("id")
	depositIDUint, err := strconv.ParseUint(depositIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "UpdateFixedDeposit - Invalid fixed deposit ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid fixed deposit ID")
		return
	}
	depositID := uint(depositIDUint)

	logger.Debugf(ctx, "UpdateFixedDeposit - Updating fixed deposit: %s for user: %s", depositID, userID)

	var existingDeposit models.FixedDeposit
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", depositID, userID).First(&existingDeposit).Error; err != nil {
		logger.Warnf(ctx, "UpdateFixedDeposit - Fixed deposit not found: %s", depositID)
		utilities.ErrorResponse(c, http.StatusNotFound, "Fixed deposit not found")
		return
	}

	// Don't allow updates if already withdrawn
	if existingDeposit.Withdrawn {
		logger.Warnf(ctx, "UpdateFixedDeposit - Cannot update withdrawn fixed deposit: %s", depositID)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Cannot update withdrawn fixed deposit")
		return
	}

	var updateData models.FixedDeposit
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.Warnf(ctx, "UpdateFixedDeposit - Invalid request data: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Update allowed fields
	existingDeposit.Institution = updateData.Institution
	existingDeposit.AccountNumber = updateData.AccountNumber
	existingDeposit.Principal = updateData.Principal
	existingDeposit.InterestRate = updateData.InterestRate
	existingDeposit.TenureMonths = updateData.TenureMonths
	existingDeposit.Compounding = updateData.Compounding
	existingDeposit.StartDate = updateData.StartDate
	existingDeposit.AutoRenew = updateData.AutoRenew
	existingDeposit.Notes = updateData.Notes

	// Recalculate maturity date and amount
	existingDeposit.MaturityDate = existingDeposit.StartDate.AddDate(0, existingDeposit.TenureMonths, 0)
	existingDeposit.MaturityAmount = calculateMaturityAmount(
		existingDeposit.Principal,
		existingDeposit.InterestRate,
		existingDeposit.TenureMonths,
		existingDeposit.Compounding,
	)

	logger.Debugf(ctx, "UpdateFixedDeposit - Recalculated maturity amount: %.2f", existingDeposit.MaturityAmount)

	if err := database.DB.WithContext(ctx).Save(&existingDeposit).Error; err != nil {
		logger.Errorf(ctx, "UpdateFixedDeposit - Failed to update fixed deposit: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update fixed deposit")
		return
	}

	logger.Infof(ctx, "UpdateFixedDeposit - Successfully updated fixed deposit: %s", depositID)

	// Log fixed deposit update activity
	utilities.LogEntityActivity(c, userID, models.ActionUpdate, "fixed_deposit",
		"FixedDeposit", existingDeposit.ID, "Updated fixed deposit at "+existingDeposit.Institution, nil)

	utilities.SuccessResponse(c, existingDeposit, "Fixed deposit updated successfully")
}

// DeleteFixedDeposit deletes a fixed deposit
func DeleteFixedDeposit(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "DeleteFixedDeposit - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "DeleteFixedDeposit - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	depositIDStr := c.Param("id")
	depositIDUint, err := strconv.ParseUint(depositIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "DeleteFixedDeposit - Invalid fixed deposit ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid fixed deposit ID")
		return
	}
	depositID := uint(depositIDUint)

	logger.Debugf(ctx, "DeleteFixedDeposit - Deleting fixed deposit: %s for user: %s", depositID, userID)

	var deposit models.FixedDeposit
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", depositID, userID).First(&deposit).Error; err != nil {
		logger.Warnf(ctx, "DeleteFixedDeposit - Fixed deposit not found: %s", depositID)
		utilities.ErrorResponse(c, http.StatusNotFound, "Fixed deposit not found")
		return
	}

	// Soft delete
	if err := database.DB.WithContext(ctx).Delete(&deposit).Error; err != nil {
		logger.Errorf(ctx, "DeleteFixedDeposit - Failed to delete fixed deposit: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete fixed deposit")
		return
	}

	logger.Infof(ctx, "DeleteFixedDeposit - Successfully deleted fixed deposit: %s", depositID)

	// Log fixed deposit deletion activity
	utilities.LogEntityActivity(c, userID, models.ActionDelete, "fixed_deposit",
		"FixedDeposit", deposit.ID, "Deleted fixed deposit at "+deposit.Institution, nil)

	utilities.SuccessResponse(c, nil, "Fixed deposit deleted successfully")
}

// WithdrawFixedDeposit marks a fixed deposit as withdrawn
func WithdrawFixedDeposit(c *gin.Context) {
	ctx := middleware.GetContextWithUserID(c)
	logger.Infof(ctx, "WithdrawFixedDeposit - Entry")

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Warnf(ctx, "WithdrawFixedDeposit - Unauthorized access attempt")
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	depositIDStr := c.Param("id")
	depositIDUint, err := strconv.ParseUint(depositIDStr, 10, 32)
	if err != nil {
		logger.Warnf(ctx, "WithdrawFixedDeposit - Invalid fixed deposit ID: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid fixed deposit ID")
		return
	}
	depositID := uint(depositIDUint)

	var withdrawalData struct {
		AccountID            uint       `json:"accountId" binding:"required"`
		WithdrawnDate        *time.Time `json:"withdrawnDate"`
		ActualMaturityAmount float64    `json:"actualMaturityAmount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&withdrawalData); err != nil {
		logger.Warnf(ctx, "WithdrawFixedDeposit - Invalid request data: %v", err)
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.Debugf(ctx, "WithdrawFixedDeposit - Withdrawing fixed deposit: %s for user: %s", depositID, userID)

	var deposit models.FixedDeposit
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", depositID, userID).First(&deposit).Error; err != nil {
		logger.Warnf(ctx, "WithdrawFixedDeposit - Fixed deposit not found: %s", depositID)
		utilities.ErrorResponse(c, http.StatusNotFound, "Fixed deposit not found")
		return
	}

	// Check if already withdrawn
	if deposit.Withdrawn {
		logger.Warnf(ctx, "WithdrawFixedDeposit - Fixed deposit already withdrawn: %s", depositID)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Fixed deposit already withdrawn")
		return
	}

	// Verify account belongs to user
	var account models.Account
	if err := database.DB.WithContext(ctx).Where("id = ? AND user_id = ?", withdrawalData.AccountID, userID).First(&account).Error; err != nil {
		logger.Warnf(ctx, "WithdrawFixedDeposit - Invalid account ID: %s", withdrawalData.AccountID)
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	withdrawnDate := func() time.Time {
		if withdrawalData.WithdrawnDate != nil {
			return *withdrawalData.WithdrawnDate
		}
		return time.Now()
	}()

	// Start transaction
	tx := database.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Mark as withdrawn
	deposit.Withdrawn = true
	deposit.WithdrawnDate = &withdrawnDate
	deposit.ActualMaturityAmount = withdrawalData.ActualMaturityAmount

	logger.Debugf(ctx, "WithdrawFixedDeposit - Actual maturity amount: %.2f", withdrawalData.ActualMaturityAmount)

	if err := tx.Save(&deposit).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "WithdrawFixedDeposit - Failed to withdraw fixed deposit: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to withdraw fixed deposit")
		return
	}

	// Create transaction record (income as money returns to account)
	transaction := models.Transaction{
		UserID:         userID,
		AccountID:      withdrawalData.AccountID,
		Type:           "income",
		Amount:         withdrawalData.ActualMaturityAmount,
		CategoryID:     0, // Use 0 for system transactions
		Date:           models.Date{Time: withdrawnDate},
		Description:    "FD Maturity: " + deposit.Institution + " - " + deposit.AccountNumber,
		FixedDepositID: &depositID,
		Tags:           []string{"fixed_deposit", "maturity"},
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "WithdrawFixedDeposit - Failed to create transaction: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	// Update account balance (credit)
	account.Balance += withdrawalData.ActualMaturityAmount
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "WithdrawFixedDeposit - Failed to update account balance: %v", err)
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account balance")
		return
	}

	tx.Commit()

	// Calculate interest earned and penalty (if withdrawn early)
	interestEarned := deposit.ActualMaturityAmount - deposit.Principal
	isEarlyWithdrawal := deposit.WithdrawnDate.Before(deposit.MaturityDate)
	var penalty float64

	if isEarlyWithdrawal {
		// If withdrawn early, penalty is the difference between expected and actual
		expectedInterest := deposit.MaturityAmount - deposit.Principal
		penalty = expectedInterest - interestEarned
		logger.Debugf(ctx, "WithdrawFixedDeposit - Early withdrawal detected. Penalty: %.2f", penalty)
	}

	logger.Infof(ctx, "WithdrawFixedDeposit - Successfully withdrew fixed deposit: %s", depositID)

	// Log fixed deposit withdrawal activity
	utilities.LogEntityActivity(c, userID, models.ActionUpdate, "fixed_deposit",
		"FixedDeposit", deposit.ID, "Withdrew fixed deposit from "+deposit.Institution, nil)

	result := map[string]interface{}{
		"deposit":           deposit,
		"transaction":       transaction,
		"interestEarned":    interestEarned,
		"isEarlyWithdrawal": isEarlyWithdrawal,
		"penalty":           penalty,
	}

	utilities.SuccessResponse(c, result, "Fixed deposit withdrawn successfully")
}

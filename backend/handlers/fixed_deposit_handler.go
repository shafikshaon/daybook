package handlers

import (
	"math"
	"net/http"
	"time"

	"daybook-backend/database"
	"daybook-backend/middleware"
	"daybook-backend/models"
	"daybook-backend/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListFixedDeposits returns all fixed deposits for the authenticated user
func ListFixedDeposits(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	query := database.DB.Where("user_id = ?", userID)

	// Optional filter by withdrawn status
	if withdrawn := c.Query("withdrawn"); withdrawn != "" {
		query = query.Where("withdrawn = ?", withdrawn == "true")
	}

	var deposits []models.FixedDeposit
	if err := query.Order("maturity_date ASC").Find(&deposits).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch fixed deposits")
		return
	}

	utilities.SuccessResponse(c, deposits, "Fixed deposits retrieved successfully")
}

// GetFixedDeposit returns a specific fixed deposit by ID
func GetFixedDeposit(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	depositID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid fixed deposit ID")
		return
	}

	var deposit models.FixedDeposit
	if err := database.DB.Where("id = ? AND user_id = ?", depositID, userID).First(&deposit).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Fixed deposit not found")
		return
	}

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
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var depositData struct {
		models.FixedDeposit
		AccountID uuid.UUID `json:"accountId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&depositData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	depositData.FixedDeposit.UserID = userID

	// Verify account belongs to user
	var account models.Account
	if err := database.DB.Where("id = ? AND user_id = ?", depositData.AccountID, userID).First(&account).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID")
		return
	}

	// Check sufficient balance
	if account.Balance < depositData.Principal {
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

	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create fixed deposit
	if err := tx.Create(&depositData.FixedDeposit).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create fixed deposit")
		return
	}

	// Create transaction record
	transaction := models.Transaction{
		UserID:         userID,
		AccountID:      depositData.AccountID,
		Type:           "expense",
		Amount:         depositData.Principal,
		CategoryID:     "fixed_deposit_investment",
		Date:           depositData.StartDate,
		Description:    "Fixed Deposit: " + depositData.Institution + " - " + depositData.AccountNumber,
		FixedDepositID: &depositData.FixedDeposit.ID,
		Tags:           []string{"fixed_deposit"},
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	// Update account balance (debit)
	account.Balance -= depositData.Principal
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update account balance")
		return
	}

	tx.Commit()

	result := map[string]interface{}{
		"fixedDeposit": depositData.FixedDeposit,
		"transaction":  transaction,
	}

	utilities.CreatedResponse(c, result, "Fixed deposit created successfully")
}

// UpdateFixedDeposit updates an existing fixed deposit
func UpdateFixedDeposit(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	depositID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid fixed deposit ID")
		return
	}

	var existingDeposit models.FixedDeposit
	if err := database.DB.Where("id = ? AND user_id = ?", depositID, userID).First(&existingDeposit).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Fixed deposit not found")
		return
	}

	// Don't allow updates if already withdrawn
	if existingDeposit.Withdrawn {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Cannot update withdrawn fixed deposit")
		return
	}

	var updateData models.FixedDeposit
	if err := c.ShouldBindJSON(&updateData); err != nil {
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

	if err := database.DB.Save(&existingDeposit).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to update fixed deposit")
		return
	}

	utilities.SuccessResponse(c, existingDeposit, "Fixed deposit updated successfully")
}

// DeleteFixedDeposit deletes a fixed deposit
func DeleteFixedDeposit(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	depositID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid fixed deposit ID")
		return
	}

	var deposit models.FixedDeposit
	if err := database.DB.Where("id = ? AND user_id = ?", depositID, userID).First(&deposit).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Fixed deposit not found")
		return
	}

	// Soft delete
	if err := database.DB.Delete(&deposit).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete fixed deposit")
		return
	}

	utilities.SuccessResponse(c, nil, "Fixed deposit deleted successfully")
}

// WithdrawFixedDeposit marks a fixed deposit as withdrawn
func WithdrawFixedDeposit(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		utilities.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	depositID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Invalid fixed deposit ID")
		return
	}

	var withdrawalData struct {
		AccountID            uuid.UUID  `json:"accountId" binding:"required"`
		WithdrawnDate        *time.Time `json:"withdrawnDate"`
		ActualMaturityAmount float64    `json:"actualMaturityAmount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&withdrawalData); err != nil {
		utilities.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var deposit models.FixedDeposit
	if err := database.DB.Where("id = ? AND user_id = ?", depositID, userID).First(&deposit).Error; err != nil {
		utilities.ErrorResponse(c, http.StatusNotFound, "Fixed deposit not found")
		return
	}

	// Check if already withdrawn
	if deposit.Withdrawn {
		utilities.ErrorResponse(c, http.StatusBadRequest, "Fixed deposit already withdrawn")
		return
	}

	// Verify account belongs to user
	var account models.Account
	if err := database.DB.Where("id = ? AND user_id = ?", withdrawalData.AccountID, userID).First(&account).Error; err != nil {
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
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Mark as withdrawn
	deposit.Withdrawn = true
	deposit.WithdrawnDate = &withdrawnDate
	deposit.ActualMaturityAmount = withdrawalData.ActualMaturityAmount

	if err := tx.Save(&deposit).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to withdraw fixed deposit")
		return
	}

	// Create transaction record (income as money returns to account)
	transaction := models.Transaction{
		UserID:         userID,
		AccountID:      withdrawalData.AccountID,
		Type:           "income",
		Amount:         withdrawalData.ActualMaturityAmount,
		CategoryID:     "fixed_deposit_maturity",
		Date:           withdrawnDate,
		Description:    "FD Maturity: " + deposit.Institution + " - " + deposit.AccountNumber,
		FixedDepositID: &depositID,
		Tags:           []string{"fixed_deposit", "maturity"},
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		utilities.ErrorResponse(c, http.StatusInternalServerError, "Failed to create transaction")
		return
	}

	// Update account balance (credit)
	account.Balance += withdrawalData.ActualMaturityAmount
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
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
	}

	result := map[string]interface{}{
		"deposit":           deposit,
		"transaction":       transaction,
		"interestEarned":    interestEarned,
		"isEarlyWithdrawal": isEarlyWithdrawal,
		"penalty":           penalty,
	}

	utilities.SuccessResponse(c, result, "Fixed deposit withdrawn successfully")
}

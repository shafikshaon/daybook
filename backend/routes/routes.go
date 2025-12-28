package routes

import (
	"daybook-backend/container"
	"daybook-backend/handlers"
	"daybook-backend/logger"
	"daybook-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, c *container.Container) {
	// Add logger middleware globally to populate trace IDs
	router.Use(logger.ModifyContext)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Daybook API is running"})
	})

	// API v1 routes
	api := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := api.Group("/auth")
		{
			auth.POST("/signup", c.AuthHandler.Signup)
			auth.POST("/login", c.AuthHandler.Login)
		}

		// Protected routes (authentication required)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Auth routes
			authRoutes := protected.Group("/auth")
			{
				authRoutes.GET("/me", c.AuthHandler.GetProfile)
				authRoutes.PUT("/profile", c.AuthHandler.UpdateProfile)
				authRoutes.PUT("/change-password", c.AuthHandler.ChangePassword)
			}

			// Account routes
			accountRoutes := protected.Group("/accounts")
			{
				accountRoutes.GET("", c.AccountHandler.ListAccounts)
				accountRoutes.POST("", c.AccountHandler.CreateAccount)
				accountRoutes.GET("/:id", c.AccountHandler.GetAccount)
				accountRoutes.PUT("/:id", c.AccountHandler.UpdateAccount)
				accountRoutes.DELETE("/:id", c.AccountHandler.DeleteAccount)
				// NOTE: Direct balance updates removed - balances are updated automatically by transactions
				// See BALANCE_SYSTEM.md for dual-balance accounting system documentation

				// Reconciliation routes for accounts (must use same wildcard name)
				accountRoutes.GET("/:id/reconciliations", c.ReconciliationHandler.ListReconciliations)
				accountRoutes.GET("/:id/reconciliations/stats", c.ReconciliationHandler.GetReconciliationStats)
				accountRoutes.GET("/:id/unreconciled-transactions", c.ReconciliationHandler.GetUnreconciledTransactions)
			}

			// Account Type routes
			accountTypeRoutes := protected.Group("/account-types")
			{
				accountTypeRoutes.GET("", c.AccountTypeHandler.ListAccountTypes)
				accountTypeRoutes.GET("/:id", c.AccountTypeHandler.GetAccountType)
				accountTypeRoutes.POST("", c.AccountTypeHandler.CreateAccountType)
				accountTypeRoutes.PUT("/:id", c.AccountTypeHandler.UpdateAccountType)
				accountTypeRoutes.DELETE("/:id", c.AccountTypeHandler.DeleteAccountType)
			}

			// Transaction routes
			transactionRoutes := protected.Group("/transactions")
			{
				transactionRoutes.GET("", c.TransactionHandler.ListTransactions)
				transactionRoutes.GET("/stats", c.TransactionHandler.GetTransactionStats)
				transactionRoutes.GET("/:id", c.TransactionHandler.GetTransaction)
				transactionRoutes.POST("", c.TransactionHandler.CreateTransaction)
				transactionRoutes.POST("/bulk", c.TransactionHandler.BulkImportTransactions)
				transactionRoutes.PUT("/:id", c.TransactionHandler.UpdateTransaction)
				transactionRoutes.DELETE("/:id", c.TransactionHandler.DeleteTransaction)
			}

			// Recurring transaction routes
			recurringTransactionRoutes := protected.Group("/recurring-transactions")
			{
				recurringTransactionRoutes.GET("", c.RecurringTransactionHandler.ListRecurringTransactions)
				recurringTransactionRoutes.GET("/:id", c.RecurringTransactionHandler.GetRecurringTransaction)
				recurringTransactionRoutes.POST("", c.RecurringTransactionHandler.CreateRecurringTransaction)
				recurringTransactionRoutes.PUT("/:id", c.RecurringTransactionHandler.UpdateRecurringTransaction)
				recurringTransactionRoutes.DELETE("/:id", c.RecurringTransactionHandler.DeleteRecurringTransaction)
				recurringTransactionRoutes.POST("/process", c.RecurringTransactionHandler.ProcessRecurringTransactions)
			}
			// Credit card routes
			creditCardRoutes := protected.Group("/credit-cards")
			{
				creditCardRoutes.GET("", c.CreditCardHandler.ListCreditCards)
				creditCardRoutes.GET("/:id", c.CreditCardHandler.GetCreditCard)
				creditCardRoutes.POST("", c.CreditCardHandler.CreateCreditCard)
				creditCardRoutes.PUT("/:id", c.CreditCardHandler.UpdateCreditCard)
				creditCardRoutes.DELETE("/:id", c.CreditCardHandler.DeleteCreditCard)

				// Transaction routes
				creditCardRoutes.POST("/:id/transactions", c.CreditCardHandler.RecordCreditCardTransaction)
				creditCardRoutes.GET("/:id/transactions", c.CreditCardHandler.GetCreditCardTransactions)
				creditCardRoutes.DELETE("/:id/transactions/:transactionId", c.CreditCardHandler.DeleteCreditCardTransaction)

				// Payment routes
				creditCardRoutes.POST("/:id/payment", c.CreditCardHandler.RecordPayment)
				creditCardRoutes.GET("/:id/payments", c.CreditCardHandler.GetPayments)

				// Statement routes
				creditCardRoutes.GET("/:id/statements", c.CreditCardHandler.GetStatements)
			}

			// Statement routes
			protected.POST("/statements", c.CreditCardHandler.CreateStatement)

			// Reward routes
			rewardRoutes := protected.Group("/rewards")
			{
				rewardRoutes.GET("", c.CreditCardHandler.ListRewards)
				rewardRoutes.POST("", c.CreditCardHandler.RecordReward)
			}

			// OLD Investment routes - REMOVED
			// Replaced by unified Goals system at /goals
			// Investment, Portfolio, and Dividend functionality now available as Goal Holdings

			// Budget routes
			budgetRoutes := protected.Group("/budgets")
			{
				budgetRoutes.GET("", c.BudgetHandler.ListBudgets)
				budgetRoutes.GET("/:id", c.BudgetHandler.GetBudget)
				budgetRoutes.GET("/:id/progress", c.BudgetHandler.GetBudgetProgress)
				budgetRoutes.POST("", c.BudgetHandler.CreateBudget)
				budgetRoutes.PUT("/:id", c.BudgetHandler.UpdateBudget)
				budgetRoutes.DELETE("/:id", c.BudgetHandler.DeleteBudget)
			}

			// Category routes
			categoryRoutes := protected.Group("/categories")
			{
				categoryRoutes.GET("", c.CategoryHandler.ListCategories)
				categoryRoutes.GET("/icons", c.CategoryHandler.GetAvailableIcons)
				categoryRoutes.PUT("/reorder", c.CategoryHandler.ReorderCategories)
				categoryRoutes.GET("/:id", c.CategoryHandler.GetCategory)
				categoryRoutes.POST("", c.CategoryHandler.CreateCategory)
				categoryRoutes.PUT("/:id", c.CategoryHandler.UpdateCategory)
				categoryRoutes.DELETE("/:id", c.CategoryHandler.DeleteCategory)
			}

			// Goal routes (Unified savings, investments, and fixed deposits)
			goalRoutes := protected.Group("/goals")
			{
				goalRoutes.GET("", c.GoalHandler.ListGoals)
				goalRoutes.GET("/:id", c.GoalHandler.GetGoal)
				goalRoutes.POST("", c.GoalHandler.CreateGoal)
				goalRoutes.PUT("/:id", c.GoalHandler.UpdateGoal)
				goalRoutes.DELETE("/:id", c.GoalHandler.DeleteGoal)

				// Holdings management
				goalRoutes.POST("/:id/holdings", c.GoalHandler.AddHolding)
				goalRoutes.PUT("/holdings/:holdingId", c.GoalHandler.UpdateHolding)
				goalRoutes.DELETE("/holdings/:holdingId", c.GoalHandler.RemoveHolding)

				// Utility endpoints
				goalRoutes.GET("/holding-types", c.GoalHandler.GetHoldingTypes)
			}

			// Settings routes
			settingsRoutes := protected.Group("/settings")
			{
				settingsRoutes.GET("", c.SettingsHandler.GetSettings)
				settingsRoutes.PUT("", c.SettingsHandler.UpdateSettings)

				// Category management under settings
				settingsRoutes.GET("/categories", c.CategoryHandler.ListCategories)
				settingsRoutes.GET("/categories/icons", c.CategoryHandler.GetAvailableIcons)
				settingsRoutes.GET("/categories/:id", c.CategoryHandler.GetCategory)
				settingsRoutes.POST("/categories", c.CategoryHandler.CreateCategory)
				settingsRoutes.PUT("/categories/:id", c.CategoryHandler.UpdateCategory)
				settingsRoutes.DELETE("/categories/:id", c.CategoryHandler.DeleteCategory)
			}

			// Reconciliation routes
			reconciliationRoutes := protected.Group("/reconciliations")
			{
				reconciliationRoutes.GET("", c.ReconciliationHandler.ListReconciliations)
				reconciliationRoutes.GET("/:id", c.ReconciliationHandler.GetReconciliation)
				reconciliationRoutes.POST("", c.ReconciliationHandler.CreateReconciliation)
				reconciliationRoutes.PUT("/:id", c.ReconciliationHandler.UpdateReconciliation)
				reconciliationRoutes.DELETE("/:id", c.ReconciliationHandler.DeleteReconciliation)
			}

			// Debt routes
			debtRoutes := protected.Group("/debts")
			{
				debtRoutes.GET("", c.DebtHandler.ListDebts)
				debtRoutes.GET("/:id", c.DebtHandler.GetDebt)
				debtRoutes.POST("", c.DebtHandler.CreateDebt)
				debtRoutes.PUT("/:id", c.DebtHandler.UpdateDebt)
				debtRoutes.DELETE("/:id", c.DebtHandler.DeleteDebt)
				debtRoutes.POST("/:id/payments", c.DebtHandler.RecordDebtPayment)
				debtRoutes.GET("/:id/payments", c.DebtHandler.ListDebtPayments)
			}

			// Lend routes
			lendRoutes := protected.Group("/lends")
			{
				lendRoutes.GET("", c.LendHandler.ListLends)
				lendRoutes.GET("/:id", c.LendHandler.GetLend)
				lendRoutes.POST("", c.LendHandler.CreateLend)
				lendRoutes.PUT("/:id", c.LendHandler.UpdateLend)
				lendRoutes.DELETE("/:id", c.LendHandler.DeleteLend)
				lendRoutes.POST("/:id/payments", c.LendHandler.RecordLendPayment)
				lendRoutes.GET("/:id/payments", c.LendHandler.ListLendPayments)
			}

			// Assets tracking routes
			assetRoutes := protected.Group("/assets")
			{
				assetRoutes.GET("", c.AssetHandler.ListAssets)
				assetRoutes.GET("/stats", c.AssetHandler.GetAssetsStats)
				assetRoutes.GET("/:id", c.AssetHandler.GetAsset)
				assetRoutes.POST("", c.AssetHandler.CreateAsset)
				assetRoutes.PUT("/:id", c.AssetHandler.UpdateAsset)
				assetRoutes.DELETE("/:id", c.AssetHandler.DeleteAsset)

				// Service records
				assetRoutes.POST("/:id/services", c.AssetHandler.CreateServiceRecord)
				assetRoutes.GET("/:id/services", c.AssetHandler.ListServiceRecords)
				assetRoutes.DELETE("/:id/services/:serviceId", c.AssetHandler.DeleteServiceRecord)

				// Attachments
				assetRoutes.POST("/:id/attachments", c.AssetHandler.AddAttachment)
				assetRoutes.GET("/:id/attachments", c.AssetHandler.ListAttachments)
				assetRoutes.DELETE("/:id/attachments/:attachmentId", c.AssetHandler.DeleteAttachment)
			}

			// Report routes
			reportRoutes := protected.Group("/reports")
			{
				reportRoutes.GET("/dashboard", c.ReportHandler.GetDashboardSummary)
				reportRoutes.GET("/income-expense", c.ReportHandler.GetIncomeExpenseReport)
				reportRoutes.GET("/category-analysis", c.ReportHandler.GetCategoryAnalysis)
				reportRoutes.GET("/accounts", c.ReportHandler.GetAccountReport)
				reportRoutes.GET("/accounts/:id/history", c.ReportHandler.GetAccountBalanceHistory)
				reportRoutes.GET("/net-worth", c.ReportHandler.GetNetWorthReport)
				reportRoutes.GET("/budget", c.ReportHandler.GetBudgetReport)
				reportRoutes.GET("/cash-flow", c.ReportHandler.GetCashFlowReport)
				reportRoutes.GET("/monthly-summary", c.ReportHandler.GetMonthlySummary)
				reportRoutes.GET("/yearly-summary", c.ReportHandler.GetYearlySummary)
				reportRoutes.POST("/comparison", c.ReportHandler.GetPeriodComparison)
			}

			// File upload routes
			uploadRoutes := protected.Group("/uploads")
			{
				uploadRoutes.POST("", c.UploadHandler.UploadFiles)             // Multiple files
				uploadRoutes.POST("/single", c.UploadHandler.UploadSingleFile) // Single file
				uploadRoutes.GET("/:userId/:filename", c.UploadHandler.ServeUploadedFile)
				uploadRoutes.DELETE("/:filename", c.UploadHandler.DeleteFile)
				uploadRoutes.GET("/info/:filename", c.UploadHandler.GetFileInfo)
			}

			// Activity log routes
			activityRoutes := protected.Group("/activity-logs")
			{
				activityRoutes.GET("", handlers.ListActivityLogs)
				activityRoutes.GET("/summary", handlers.GetActivitySummary)
				activityRoutes.GET("/:id", handlers.GetActivityLog)
				activityRoutes.DELETE("/cleanup", handlers.DeleteOldActivityLogs)
				activityRoutes.POST("/backfill", handlers.BackfillActivityLogs)
			}
		}
	}
}

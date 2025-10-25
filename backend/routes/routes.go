package routes

import (
	"daybook-backend/handlers"
	"daybook-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
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
			auth.POST("/signup", handlers.Signup)
			auth.POST("/login", handlers.Login)
		}

		// Protected routes (authentication required)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Auth routes
			authRoutes := protected.Group("/auth")
			{
				authRoutes.GET("/me", handlers.GetProfile)
				authRoutes.PUT("/profile", handlers.UpdateProfile)
				authRoutes.PUT("/change-password", handlers.ChangePassword)
			}

			// Account routes
			accountRoutes := protected.Group("/accounts")
			{
				accountRoutes.GET("", handlers.ListAccounts)
				accountRoutes.GET("/:id", handlers.GetAccount)
				accountRoutes.POST("", handlers.CreateAccount)
				accountRoutes.PUT("/:id", handlers.UpdateAccount)
				accountRoutes.DELETE("/:id", handlers.DeleteAccount)
				// NOTE: Direct balance updates removed - balances are updated automatically by transactions
				// See BALANCE_SYSTEM.md for dual-balance accounting system documentation
			}

			// Account Type routes
			accountTypeRoutes := protected.Group("/account-types")
			{
				accountTypeRoutes.GET("", handlers.ListAccountTypes)
				accountTypeRoutes.GET("/:id", handlers.GetAccountType)
				accountTypeRoutes.POST("", handlers.CreateAccountType)
				accountTypeRoutes.PUT("/:id", handlers.UpdateAccountType)
				accountTypeRoutes.DELETE("/:id", handlers.DeleteAccountType)
			}

			// Transaction routes
			transactionRoutes := protected.Group("/transactions")
			{
				transactionRoutes.GET("", handlers.ListTransactions)
				transactionRoutes.GET("/stats", handlers.GetTransactionStats)
				transactionRoutes.GET("/:id", handlers.GetTransaction)
				transactionRoutes.POST("", handlers.CreateTransaction)
				transactionRoutes.POST("/bulk", handlers.BulkImportTransactions)
				transactionRoutes.PUT("/:id", handlers.UpdateTransaction)
				transactionRoutes.DELETE("/:id", handlers.DeleteTransaction)
			}

			// Credit card routes
			creditCardRoutes := protected.Group("/credit-cards")
			{
				creditCardRoutes.GET("", handlers.ListCreditCards)
				creditCardRoutes.GET("/:id", handlers.GetCreditCard)
				creditCardRoutes.POST("", handlers.CreateCreditCard)
				creditCardRoutes.PUT("/:id", handlers.UpdateCreditCard)
				creditCardRoutes.DELETE("/:id", handlers.DeleteCreditCard)

				// Transaction routes
				creditCardRoutes.POST("/:id/transactions", handlers.RecordCreditCardTransaction)
				creditCardRoutes.GET("/:id/transactions", handlers.GetCreditCardTransactions)
				creditCardRoutes.DELETE("/:id/transactions/:transactionId", handlers.DeleteCreditCardTransaction)

				// Payment routes
				creditCardRoutes.POST("/:id/payment", handlers.RecordPayment)
				creditCardRoutes.GET("/:id/payments", handlers.GetPayments)

				// Statement routes
				creditCardRoutes.GET("/:id/statements", handlers.GetStatements)
			}

			// Statement routes
			protected.POST("/statements", handlers.CreateStatement)

			// Reward routes
			rewardRoutes := protected.Group("/rewards")
			{
				rewardRoutes.GET("", handlers.ListRewards)
				rewardRoutes.POST("", handlers.RecordReward)
			}

			// Investment routes
			investmentRoutes := protected.Group("/investments")
			{
				investmentRoutes.GET("", handlers.ListInvestments)
				investmentRoutes.GET("/:id", handlers.GetInvestment)
				investmentRoutes.POST("", handlers.CreateInvestment)
				investmentRoutes.PUT("/:id", handlers.UpdateInvestment)
				investmentRoutes.DELETE("/:id", handlers.DeleteInvestment)
				investmentRoutes.POST("/:id/buy", handlers.BuyShares)
				investmentRoutes.POST("/:id/sell", handlers.SellShares)
			}

			// Portfolio routes
			portfolioRoutes := protected.Group("/portfolios")
			{
				portfolioRoutes.GET("", handlers.ListPortfolios)
				portfolioRoutes.POST("", handlers.CreatePortfolio)
			}

			// Dividend routes
			dividendRoutes := protected.Group("/dividends")
			{
				dividendRoutes.GET("", handlers.ListDividends)
				dividendRoutes.POST("", handlers.RecordDividend)
			}

			// Bill routes
			billRoutes := protected.Group("/bills")
			{
				billRoutes.GET("", handlers.ListBills)
				billRoutes.GET("/:id", handlers.GetBill)
				billRoutes.POST("", handlers.CreateBill)
				billRoutes.PUT("/:id", handlers.UpdateBill)
				billRoutes.DELETE("/:id", handlers.DeleteBill)
				billRoutes.POST("/:id/pay", handlers.PayBill)
			}

			// Bill payment routes
			protected.GET("/bill-payments", handlers.GetBillPayments)

			// Budget routes
			budgetRoutes := protected.Group("/budgets")
			{
				budgetRoutes.GET("", handlers.ListBudgets)
				budgetRoutes.GET("/:id", handlers.GetBudget)
				budgetRoutes.GET("/:id/progress", handlers.GetBudgetProgress)
				budgetRoutes.POST("", handlers.CreateBudget)
				budgetRoutes.PUT("/:id", handlers.UpdateBudget)
				budgetRoutes.DELETE("/:id", handlers.DeleteBudget)
			}

			// Savings goal routes
			savingsGoalRoutes := protected.Group("/savings-goals")
			{
				savingsGoalRoutes.GET("", handlers.ListSavingsGoals)
				savingsGoalRoutes.GET("/:id", handlers.GetSavingsGoal)
				savingsGoalRoutes.POST("", handlers.CreateSavingsGoal)
				savingsGoalRoutes.PUT("/:id", handlers.UpdateSavingsGoal)
				savingsGoalRoutes.DELETE("/:id", handlers.DeleteSavingsGoal)
				savingsGoalRoutes.POST("/:id/contribute", handlers.AddContribution)
				savingsGoalRoutes.POST("/:id/withdraw", handlers.WithdrawFromGoal)
			}

			// Automated rule routes
			automatedRuleRoutes := protected.Group("/automated-rules")
			{
				automatedRuleRoutes.GET("", handlers.ListAutomatedRules)
				automatedRuleRoutes.POST("", handlers.CreateAutomatedRule)
			}

			// Fixed deposit routes
			fixedDepositRoutes := protected.Group("/fixed-deposits")
			{
				fixedDepositRoutes.GET("", handlers.ListFixedDeposits)
				fixedDepositRoutes.GET("/:id", handlers.GetFixedDeposit)
				fixedDepositRoutes.POST("", handlers.CreateFixedDeposit)
				fixedDepositRoutes.PUT("/:id", handlers.UpdateFixedDeposit)
				fixedDepositRoutes.DELETE("/:id", handlers.DeleteFixedDeposit)
				fixedDepositRoutes.POST("/:id/withdraw", handlers.WithdrawFixedDeposit)
			}

			// Settings routes
			settingsRoutes := protected.Group("/settings")
			{
				settingsRoutes.GET("", handlers.GetSettings)
				settingsRoutes.PUT("", handlers.UpdateSettings)
			}

			// File upload routes
			uploadRoutes := protected.Group("/uploads")
			{
				uploadRoutes.POST("", handlers.UploadFiles)             // Multiple files
				uploadRoutes.POST("/single", handlers.UploadSingleFile) // Single file
				uploadRoutes.GET("/:userId/:filename", handlers.ServeUploadedFile)
				uploadRoutes.DELETE("/:filename", handlers.DeleteFile)
				uploadRoutes.GET("/info/:filename", handlers.GetFileInfo)
			}
		}
	}
}

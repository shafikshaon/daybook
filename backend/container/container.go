package container

import (
	"daybook-backend/config"
	"daybook-backend/handlers"
	"daybook-backend/repository"
	"daybook-backend/services"

	"gorm.io/gorm"
)

// Container holds all application dependencies
// This is the dependency injection container that wires everything together
type Container struct {
	// Infrastructure
	TxManager repository.TransactionManager

	// Repositories
	ActivityLogRepo          repository.ActivityLogRepository
	SettingsRepo             repository.SettingsRepository
	CategoryRepo             repository.CategoryRepository
	AccountTypeRepo          repository.AccountTypeRepository
	AccountRepo              repository.AccountRepository
	UserRepo                 repository.UserRepository
	BudgetRepo               repository.BudgetRepository
	ReconciliationRepo       repository.ReconciliationRepository
	DebtRepo                 repository.DebtRepository
	LendRepo                 repository.LendRepository
	AssetRepo                repository.AssetRepository
	GoalRepo                 repository.GoalRepository
	CreditCardRepo           repository.CreditCardRepository
	RecurringTransactionRepo repository.RecurringTransactionRepository
	TransactionRepo          repository.TransactionRepository
	ReportRepo               repository.ReportRepository
	BackupRepo               repository.BackupRepository
	// TODO: Add more repositories as we migrate

	// Services
	ActivityLogService          services.ActivityLogService
	SettingsService             services.SettingsService
	CategoryService             services.CategoryService
	AccountTypeService          services.AccountTypeService
	AccountService              services.AccountService
	AuthService                 services.AuthService
	BudgetService               services.BudgetService
	UploadService               services.UploadService
	ReconciliationService       services.ReconciliationService
	DebtService                 services.DebtService
	LendService                 services.LendService
	AssetService                services.AssetService
	GoalService                 services.GoalService
	CreditCardService           services.CreditCardService
	RecurringTransactionService services.RecurringTransactionService
	TransactionService          services.TransactionService
	ReportService               services.ReportService
	ExportService               services.ExportService
	BackupService               services.BackupService
	// TODO: Add more services as we migrate

	// Handlers
	SettingsHandler             *handlers.SettingsHandler
	CategoryHandler             *handlers.CategoryHandler
	AccountTypeHandler          *handlers.AccountTypeHandler
	AccountHandler              *handlers.AccountHandler
	AuthHandler                 *handlers.AuthHandler
	BudgetHandler               *handlers.BudgetHandler
	UploadHandler               *handlers.UploadHandler
	ReconciliationHandler       *handlers.ReconciliationHandler
	DebtHandler                 *handlers.DebtHandler
	LendHandler                 *handlers.LendHandler
	AssetHandler                *handlers.AssetHandler
	GoalHandler                 *handlers.GoalHandler
	CreditCardHandler           *handlers.CreditCardHandler
	RecurringTransactionHandler *handlers.RecurringTransactionHandler
	TransactionHandler          *handlers.TransactionHandler
	ReportHandler               *handlers.ReportHandler
	ExportHandler               *handlers.ExportHandler
	BackupHandler               *handlers.BackupHandler
	// TODO: Add more handlers as we migrate
}

// NewContainer creates and wires all dependencies
func NewContainer(db *gorm.DB, cfg *config.Config) *Container {
	c := &Container{}

	// Initialize infrastructure
	c.TxManager = repository.NewTransactionManager(db)

	// Initialize repositories
	c.ActivityLogRepo = repository.NewActivityLogRepository(db)
	c.SettingsRepo = repository.NewSettingsRepository(db)
	c.CategoryRepo = repository.NewCategoryRepository(db)
	c.AccountTypeRepo = repository.NewAccountTypeRepository(db)
	c.AccountRepo = repository.NewAccountRepository(db)
	c.UserRepo = repository.NewUserRepository(db)
	c.BudgetRepo = repository.NewBudgetRepository(db)
	c.ReconciliationRepo = repository.NewReconciliationRepository(db)
	c.DebtRepo = repository.NewDebtRepository(db)
	c.LendRepo = repository.NewLendRepository(db)
	c.AssetRepo = repository.NewAssetRepository(db)
	c.GoalRepo = repository.NewGoalRepository(db)
	c.CreditCardRepo = repository.NewCreditCardRepository(db)
	c.RecurringTransactionRepo = repository.NewRecurringTransactionRepository(db)
	c.TransactionRepo = repository.NewTransactionRepository(db)
	c.ReportRepo = repository.NewReportRepository(db)
	c.BackupRepo = repository.NewBackupRepository(db)
	// TODO: Add more repositories as we migrate

	// Initialize services
	c.ActivityLogService = services.NewActivityLogService(c.ActivityLogRepo)

	c.SettingsService = services.NewSettingsService(
		c.SettingsRepo,
		c.ActivityLogService,
	)

	c.CategoryService = services.NewCategoryService(
		c.CategoryRepo,
		c.ActivityLogService,
	)

	c.AccountTypeService = services.NewAccountTypeService(
		c.AccountTypeRepo,
		c.ActivityLogService,
	)

	c.AccountService = services.NewAccountService(
		c.AccountRepo,
		c.CategoryRepo,
		c.TxManager,
		c.ActivityLogService,
	)

	c.AuthService = services.NewAuthService(
		c.UserRepo,
		c.SettingsRepo,
		c.CategoryRepo,
		c.AccountTypeRepo,
		c.TxManager,
		c.ActivityLogService,
	)

	c.BudgetService = services.NewBudgetService(
		c.BudgetRepo,
		c.ActivityLogService,
	)

	c.UploadService = services.NewUploadService()

	c.ReconciliationService = services.NewReconciliationService(
		c.ReconciliationRepo,
		c.AccountRepo,
		c.TxManager,
		c.ActivityLogService,
	)

	c.DebtService = services.NewDebtService(
		c.DebtRepo,
		c.AccountRepo,
		c.TxManager,
		c.ActivityLogService,
	)

	c.LendService = services.NewLendService(
		c.LendRepo,
		c.AccountRepo,
		c.TxManager,
		c.ActivityLogService,
	)

	c.AssetService = services.NewAssetService(
		c.AssetRepo,
		c.ActivityLogService,
	)

	c.GoalService = services.NewGoalService(
		c.GoalRepo,
		c.AccountRepo,
		c.CategoryRepo,
		c.TxManager,
		c.ActivityLogService,
	)

	c.CreditCardService = services.NewCreditCardService(
		c.CreditCardRepo,
		c.AccountRepo,
		c.TxManager,
		c.ActivityLogService,
	)

	c.RecurringTransactionService = services.NewRecurringTransactionService(
		c.RecurringTransactionRepo,
		c.AccountRepo,
		c.CreditCardRepo,
		c.TxManager,
		c.ActivityLogService,
	)

	c.TransactionService = services.NewTransactionService(
		c.TransactionRepo,
		c.AccountRepo,
		c.CreditCardRepo,
		c.CategoryRepo,
		c.TxManager,
		c.ActivityLogService,
	)

	c.ReportService = services.NewReportService(c.ReportRepo)

	c.ExportService = services.NewExportService(
		c.TransactionRepo,
		c.AccountRepo,
		c.BudgetRepo,
		c.GoalRepo,
		c.CategoryRepo,
		c.AssetRepo,
		c.CreditCardRepo,
		c.DebtRepo,
		c.LendRepo,
		c.ActivityLogService,
	)

	c.BackupService = services.NewBackupService(
		c.BackupRepo,
		c.ActivityLogService,
		cfg,
	)
	// TODO: Add more services as we migrate

	// Initialize handlers
	c.SettingsHandler = handlers.NewSettingsHandler(c.SettingsService)
	c.CategoryHandler = handlers.NewCategoryHandler(c.CategoryService)
	c.AccountTypeHandler = handlers.NewAccountTypeHandler(c.AccountTypeService)
	c.AccountHandler = handlers.NewAccountHandler(c.AccountService)
	c.AuthHandler = handlers.NewAuthHandler(c.AuthService)
	c.BudgetHandler = handlers.NewBudgetHandler(c.BudgetService)
	c.UploadHandler = handlers.NewUploadHandler(c.UploadService)
	c.ReconciliationHandler = handlers.NewReconciliationHandler(c.ReconciliationService)
	c.DebtHandler = handlers.NewDebtHandler(c.DebtService)
	c.LendHandler = handlers.NewLendHandler(c.LendService)
	c.AssetHandler = handlers.NewAssetHandler(c.AssetService)
	c.GoalHandler = handlers.NewGoalHandler(c.GoalService)
	c.CreditCardHandler = handlers.NewCreditCardHandler(c.CreditCardService)
	c.RecurringTransactionHandler = handlers.NewRecurringTransactionHandler(c.RecurringTransactionService)
	c.TransactionHandler = handlers.NewTransactionHandler(c.TransactionService)
	c.ReportHandler = handlers.NewReportHandler(c.ReportService)
	c.ExportHandler = handlers.NewExportHandler(c.ExportService)
	c.BackupHandler = handlers.NewBackupHandler(c.BackupService)
	// TODO: Add more handlers as we migrate

	return c
}

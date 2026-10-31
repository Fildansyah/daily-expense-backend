package routes

import (
	"time"

	"expense.com/m/adapter/handlers"
	"expense.com/m/config"
	"expense.com/m/middlewares"
	authmodule "expense.com/m/module/auth"
	"expense.com/m/module/bill"
	"expense.com/m/module/budget"
	"expense.com/m/module/category"
	"expense.com/m/module/lender"
	"expense.com/m/module/report"
	"expense.com/m/module/transaction"
	"expense.com/m/module/transaction_type"
	"expense.com/m/module/user"
	"expense.com/m/module/wallet"
	"expense.com/m/shared/token"
	"expense.com/m/shared/tx"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRouter(app *fiber.App, db *gorm.DB, authCfg *config.AuthConfig, tokenManager *token.Manager) {
	//public
	publicGroup := app.Group("/public")

	userRepo := user.NewUserRepository(db)
	userSvc := user.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userSvc)

	refreshRepo := authmodule.NewRefreshTokenRepository(db)
	authSvc := authmodule.NewService(userRepo, refreshRepo, tokenManager, authCfg.RefreshTTL)
	authHandler := handlers.NewAuthHandler(authSvc, authCfg)

	// Rate limit in-memory — cukup untuk satu instance, tanpa Redis.
	loginLimiter := middlewares.RateLimit(10, time.Minute, "Terlalu banyak percobaan, coba lagi nanti")
	refreshLimiter := middlewares.RateLimit(60, time.Minute, "Terlalu banyak permintaan refresh, coba lagi nanti")

	RegisterAuthRoutes(publicGroup, authHandler, loginLimiter, refreshLimiter)
	RegisterUserRoutes(publicGroup, userHandler)

	//protected
	protectedGroup := app.Group("/api")
	protectedGroup.Use(middlewares.Authenticate(tokenManager))

	//user (protected) — session check untuk user yang sedang login
	protectedGroup.Get("/user/me", userHandler.GetCurrentUser)

	//transaction type
	transactionTypeRepo := transaction_type.NewTransactionTypeRepository(db)
	transactionTypeSvc := transaction_type.NewTransactionTypeService(transactionTypeRepo)
	transactionTypehandlr := handlers.NewTransactionTypeHandler(transactionTypeSvc)
	RegisterTransactionTypeRoutes(protectedGroup, transactionTypehandlr)

	txManager := tx.NewManager(db)

	//wallet
	walletRepo := wallet.NewWalletRepository(db)
	walletSvc := wallet.NewWalletService(walletRepo, txManager)
	walletHandler := handlers.NewWalletHandler(walletSvc)

	//category
	categoryRepo := category.NewCategoryRepository(db)
	categorySvc := category.NewCategoryService(categoryRepo, txManager)
	categoryHandler := handlers.NewCategoryHandler(categorySvc)
	RegisterCategoryRoutes(protectedGroup, categoryHandler)

	//transaction
	transactionRepo := transaction.NewTransactionRepository(db)
	transactionSvc := transaction.NewTransactionService(transactionRepo, walletRepo, txManager)
	transactionHandler := handlers.NewTransactionHandler(transactionSvc)
	RegisterTransactionRoutes(protectedGroup, transactionHandler)

	//budget
	budgetRepo := budget.NewBudgetRepository(db)
	budgetSvc := budget.NewBudgetService(budgetRepo, txManager, categorySvc, transactionSvc)
	budgetHandler := handlers.NewBudgetHandler(budgetSvc)
	RegisterBudgetRoutes(protectedGroup, budgetHandler)

	RegisterWalletRoutes(protectedGroup, walletHandler, transactionHandler)

	//bill
	billRepo := bill.NewBillRepository(db)
	billSvc := bill.NewBillService(billRepo, txManager, walletRepo, categoryRepo, transactionSvc)
	billHandler := handlers.NewBillHandler(billSvc)
	RegisterBillRoutes(protectedGroup, billHandler)

	//lender
	lenderRepo := lender.NewLenderRepository(db)
	lenderSvc := lender.NewLenderService(lenderRepo, txManager)
	lenderHandler := handlers.NewLenderHandler(lenderSvc)
	RegisterLenderRoutes(protectedGroup, lenderHandler)

	//report & dashboard
	reportSvc := report.NewService(transactionSvc, billSvc, walletSvc, categorySvc, budgetSvc)
	reportHandler := handlers.NewReportHandler(reportSvc)
	RegisterReportRoutes(protectedGroup, reportHandler)
}

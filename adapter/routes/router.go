package routes

import (
	"time"

	"expense.com/m/adapter/handlers"
	"expense.com/m/config"
	"expense.com/m/middlewares"
	authmodule "expense.com/m/module/auth"
	"expense.com/m/module/bill"
	"expense.com/m/module/category"
	"expense.com/m/module/lender"
	"expense.com/m/module/transaction_type"
	"expense.com/m/module/user"
	"expense.com/m/module/wallet"
	"expense.com/m/shared/token"
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

	//category
	categoryRepo := category.NewCategoryRepository(db)
	categorySvc := category.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categorySvc)
	RegisterCategoryRoutes(protectedGroup, categoryHandler)

	//wallet
	walletRepo := wallet.NewWalletRepository(db)
	walletSvc := wallet.NewWalletService(walletRepo)
	walletHandler := handlers.NewWalletHandler(walletSvc)
	RegisterWalletRoutes(protectedGroup, walletHandler)

	//bill
	billRepo := bill.NewBillRepository(db)
	billSvc := bill.NewBillService(billRepo)
	billHandler := handlers.NewBillHandler(billSvc)
	RegisterBillRoutes(protectedGroup, billHandler)

	//lender
	lenderRepo := lender.NewLenderRepository(db)
	lenderSvc := lender.NewLenderService(lenderRepo)
	lenderHandler := handlers.NewLenderHandler(lenderSvc)
	RegisterLenderRoutes(protectedGroup, lenderHandler)
}

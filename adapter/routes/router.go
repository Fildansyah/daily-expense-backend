package routes

import (
	"expense.com/m/adapter/handlers"
	"expense.com/m/middlewares"
	"expense.com/m/module/category"
	"expense.com/m/module/transaction_type"
	"expense.com/m/module/user"
	"expense.com/m/module/wallet"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRouter(app *fiber.App, db *gorm.DB) {
	//public
	publicGroup := app.Group("/public")
	//user
	userRepo := user.NewUserRepository(db)
	userSvc := user.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userSvc)
	RegisterUserRoutes(publicGroup, userHandler)

	//protected
	protectedGroup := app.Group("/api")
	protectedGroup.Use(middlewares.Authenticate())

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
}

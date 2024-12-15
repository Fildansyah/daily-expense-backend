package routes

import (
	"expense.com/m/adapter/handlers"
	"expense.com/m/middlewares"
	"expense.com/m/module/transaction_type"
	"expense.com/m/module/user"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRouter(app *fiber.App, db *gorm.DB) {

	publicGroup := app.Group("/public")
	//user
	userRepo := user.NewUserRepository(db)
	userSvc := user.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userSvc)
	RegisterUserRoutes(publicGroup, userHandler)

	protectedGroup := app.Group("/api")
	protectedGroup.Use(middlewares.Authenticate(db))

	//transaction type
	transactionTypeRepo := transaction_type.NewTransactionTypeRepository(db)
	transactionTypeSvc := transaction_type.NewTransactionTypeService(transactionTypeRepo)
	transactionTypehandlr := handlers.NewTransactionTypeHandler(transactionTypeSvc)
	RegisterTransactionTypeRoutes(protectedGroup, transactionTypehandlr)
}

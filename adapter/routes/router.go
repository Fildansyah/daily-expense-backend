package routes

import (
	"expense.com/m/adapter/handlers"
	"expense.com/m/module/transaction_type"
	"expense.com/m/module/user"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InitRouter(app *fiber.App, db *gorm.DB) {
	//user
	userRepo := user.NewUserRepository(db)
	userSvc := user.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userSvc)
	RegisterUserRoutes(app, userHandler)

	//transaction type
	transactionTypeRepo := transaction_type.NewTransactionTypeRepository(db)
	transactionTypeSvc := transaction_type.NewTransactionTypeService(transactionTypeRepo)
	transactionTypehandlr := handlers.NewTransactionTypeHandler(transactionTypeSvc)
	RegisterTransactionTypeRoutes(app, transactionTypehandlr)
}

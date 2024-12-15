package routes

import (
	"expense.com/m/adapter/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterTransactionTypeRoutes(app fiber.Router, handler *handlers.TransactionTypeHandler) {
	routes := app.Group("/transaction-type")

	routes.Post("", handler.CreateTransactionType)
	routes.Get(":id", handler.GetTransactionTypeByID)
	routes.Put("", handler.UpdateTransactionType)
	routes.Delete(":id", handler.DeleteTransactionType)
	routes.Get("", handler.GetAllTransactionTypes)
}

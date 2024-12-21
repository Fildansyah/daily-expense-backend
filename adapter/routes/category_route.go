package routes

import (
	"expense.com/m/adapter/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterCategoryRoutes(app fiber.Router, handler *handlers.CategoryHandler) {
	routes := app.Group("/category")

	routes.Post("", handler.CreateCategory)
	routes.Get("/transaction-type/:id", handler.FindCategoryByTransactionTypeID)
}

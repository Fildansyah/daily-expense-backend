package routes

import (
	"expense.com/m/adapter/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterTransactionRoutes(app fiber.Router, handler *handlers.TransactionHandler) {
	routes := app.Group("/transactions")

	routes.Get("", handler.List)
	routes.Post("", handler.Create)
	// static route harus didaftarkan sebelum /:id
	routes.Get("/summary", handler.Summary)
	routes.Post("/bulk", handler.Bulk)
	routes.Get("/:id", handler.Get)
	routes.Put("/:id", handler.Update)
	routes.Delete("/:id", handler.Delete)
}

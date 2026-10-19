package routes

import (
	"expense.com/m/adapter/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterLenderRoutes(app fiber.Router, handler *handlers.LenderHandler) {
	routes := app.Group("/lenders")

	routes.Get("", handler.GetLenders)
}

package routes

import (
	"expense.com/m/adapter/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterWalletRoutes(app fiber.Router, handler *handlers.WalletHandler) {
	routes := app.Group("/wallet")

	routes.Post("", handler.CreateWallet)
	routes.Get("/user/:id", handler.FindWalletByUserID)
}

package routes

import (
	"expense.com/m/adapter/handlers"
	"github.com/gofiber/fiber/v2"
)

// RegisterWalletRoutes mendaftarkan /wallet (kanonik, dipakai frontend lama)
// dan /wallets (alias plural) ke handler yang sama.
func RegisterWalletRoutes(app fiber.Router, handler *handlers.WalletHandler, transferHandler *handlers.TransactionHandler) {
	for _, prefix := range []string{"/wallet", "/wallets"} {
		routes := app.Group(prefix)

		routes.Get("", handler.GetMyWallets)
		routes.Post("", handler.CreateWallet)
		routes.Get("/summary", handler.GetSummary)
		routes.Post("/transfer", transferHandler.Transfer)
		routes.Get("/:id", handler.GetWallet)
		routes.Put("/:id", handler.UpdateWallet)
		routes.Delete("/:id", handler.DeleteWallet)
	}
}

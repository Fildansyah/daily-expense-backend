package routes

import (
	"expense.com/m/adapter/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterBillRoutes(app fiber.Router, handler *handlers.BillHandler) {
	routes := app.Group("/bills")

	routes.Get("", handler.GetBills)
	routes.Post("", handler.CreateBill)
	// /summary harus didaftarkan sebelum /:id agar "summary" tidak tertangkap sebagai :id
	routes.Get("/summary", handler.GetSummary)
	routes.Get("/:id", handler.GetBillByID)
	routes.Put("/:id", handler.UpdateBill)
	routes.Delete("/:id", handler.DeleteBill)
}

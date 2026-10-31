package routes

import (
	"expense.com/m/adapter/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterBudgetRoutes(app fiber.Router, handler *handlers.BudgetHandler) {
	routes := app.Group("/budgets")

	// static route sebelum /:id
	routes.Get("/summary", handler.GetSummary)
	routes.Post("/copy", handler.Copy)

	routes.Get("", handler.GetBudgets)
	routes.Post("", handler.CreateBudget)
	routes.Put("/:id", handler.UpdateBudget)
	routes.Delete("/:id", handler.DeleteBudget)
}

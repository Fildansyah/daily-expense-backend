package routes

import (
	"expense.com/m/adapter/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterReportRoutes(app fiber.Router, handler *handlers.ReportHandler) {
	app.Get("/dashboard", handler.Dashboard)
	app.Get("/calendar", handler.Calendar)

	reports := app.Group("/reports")
	reports.Get("/summary", handler.Summary)
	reports.Get("/cashflow", handler.Cashflow)
	reports.Get("/categories", handler.Categories)
	reports.Get("/trend", handler.Trend)
	reports.Get("/debt", handler.Debt)
	reports.Get("/networth", handler.NetWorth)
}

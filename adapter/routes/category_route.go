package routes

import (
	"expense.com/m/adapter/handlers"
	"github.com/gofiber/fiber/v2"
)

// RegisterCategoryRoutes mendaftarkan /category (kanonik) dan /categories (alias).
func RegisterCategoryRoutes(app fiber.Router, handler *handlers.CategoryHandler) {
	for _, prefix := range []string{"/category", "/categories"} {
		routes := app.Group(prefix)

		// static route sebelum /:id
		routes.Get("/tree", handler.GetTree)
		routes.Put("/reorder", handler.Reorder)
		routes.Post("/merge", handler.Merge)
		routes.Get("/transaction-type/:id", handler.FindCategoryByTransactionTypeID)

		routes.Get("", handler.GetCategories)
		routes.Post("", handler.CreateCategory)
		routes.Put("/:id", handler.UpdateCategory)
		routes.Delete("/:id", handler.DeleteCategory)
	}
}

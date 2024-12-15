package routes

import (
	"expense.com/m/adapter/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app fiber.Router, handler *handlers.UserHandler) {
	routes := app.Group("/user")

	routes.Post("", handler.CreateUser)
	routes.Post("/login", handler.LoginUser)
	routes.Get("/email/:email", handler.FindUserByEmail)
}

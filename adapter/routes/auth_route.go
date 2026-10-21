package routes

import (
	"expense.com/m/adapter/handlers"
	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app fiber.Router, handler *handlers.AuthHandler, loginLimiter, refreshLimiter fiber.Handler) {
	app.Post("/user/login", loginLimiter, handler.Login)

	authGroup := app.Group("/auth")
	authGroup.Post("/refresh", refreshLimiter, handler.Refresh)
	authGroup.Post("/logout", handler.Logout)
}

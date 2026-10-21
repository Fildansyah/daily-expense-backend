package middlewares

import (
	"strings"

	"expense.com/m/adapter/handlers/base_response"
	"expense.com/m/shared/token"
	"github.com/gofiber/fiber/v2"
)

func Authenticate(tm *token.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := strings.TrimSpace(c.Get("Authorization"))
		if authHeader == "" {
			return base_response.Error(c, fiber.StatusUnauthorized, "Token otorisasi tidak ditemukan")
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return base_response.Error(c, fiber.StatusUnauthorized, "Format token tidak valid")
		}

		claims, err := tm.ParseAccess(parts[1])
		if err != nil {
			return base_response.Error(c, fiber.StatusUnauthorized, "Token tidak valid atau sudah kedaluwarsa")
		}

		c.Locals("userID", claims.Subject)

		return c.Next()
	}
}

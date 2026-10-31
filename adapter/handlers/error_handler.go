package handlers

import (
	"errors"

	"expense.com/m/adapter/handlers/base_response"
	"expense.com/m/shared/apperror"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// respondError memetakan error domain ke status HTTP yang konsisten.
func respondError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, apperror.ErrValidation):
		return base_response.Error(c, fiber.StatusBadRequest, err.Error())
	case errors.Is(err, apperror.ErrConflict):
		return base_response.Error(c, fiber.StatusConflict, err.Error())
	case errors.Is(err, apperror.ErrNotFound), errors.Is(err, gorm.ErrRecordNotFound):
		return base_response.Error(c, fiber.StatusNotFound, err.Error())
	default:
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
}

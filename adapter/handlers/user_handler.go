package handlers

import (
	"errors"

	"expense.com/m/adapter/handlers/base_response"
	"expense.com/m/adapter/presenter"
	"expense.com/m/entity"
	"expense.com/m/module/user"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *user.UserService
}

func NewUserHandler(service *user.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req entity.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	result, err := h.service.CreateUser(&req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidInput):
			return base_response.Error(c, fiber.StatusBadRequest, err.Error())
		case errors.Is(err, entity.ErrEmailExists):
			return base_response.Error(c, fiber.StatusConflict, err.Error())
		default:
			return base_response.Error(c, fiber.StatusInternalServerError, "Terjadi kesalahan pada server")
		}
	}

	return base_response.Created(c, new(presenter.UserPresenter).FromDTO(result))
}

// GetCurrentUser mengembalikan user yang sedang login berdasarkan
// userID yang disetel oleh middleware Authenticate.
func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.GetCurrentUser(userID)
	if err != nil {
		return base_response.Error(c, fiber.StatusNotFound, err.Error())
	}
	if result == nil {
		return base_response.Error(c, fiber.StatusUnauthorized, "User tidak ditemukan")
	}

	return base_response.Success(c, new(presenter.UserPresenter).FromDTO(result))
}

package handlers

import (
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

func (h *UserHandler) FindUserByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	result, err := h.service.FindUserByEmail(email)
	if err != nil {
		return base_response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return base_response.Success(c, new(presenter.UserPresenter).FromDTO(result))
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var dto entity.UserDTO
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	result, err := h.service.CreateUser(&dto)
	if err != nil {
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return base_response.Success(c, new(presenter.UserPresenter).FromDTO(result))
}

func (h *UserHandler) LoginUser(c *fiber.Ctx) error {
	var dto entity.SignInRequest
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	result, err := h.service.LoginUser(&dto)
	if err != nil {
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return base_response.Success(c, result)
}

// GetCurrentUser mengembalikan user yang sedang login berdasarkan
// userID yang disetel oleh middleware Authenticate.
func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.GetCurrentUser(userID)
	if err != nil {
		return base_response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return base_response.Success(c, new(presenter.UserPresenter).FromDTO(result))
}

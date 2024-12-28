package handlers

import (
	"expense.com/m/adapter/handlers/base_response"
	"expense.com/m/adapter/presenter"
	"expense.com/m/entity"
	"expense.com/m/module/wallet"
	"github.com/gofiber/fiber/v2"
)

type WalletHandler struct {
	service *wallet.WalletService
}

func NewWalletHandler(service *wallet.WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

func (w *WalletHandler) CreateWallet(c *fiber.Ctx) error {
	var dto entity.WalletDTO
	authId, ok := c.Locals("userID").(string)

	if !ok {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	dto.UserID = authId

	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	result, err := w.service.CreateWallet(&dto)

	if err != nil {
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return base_response.Success(c, result)
}

func (w *WalletHandler) FindWalletByUserID(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := w.service.FindWalletByUserID(c.Context(), id)

	if err != nil {
		return base_response.Error(c, fiber.StatusNotFound, err.Error())
	}

	resp := make([]presenter.WalletPresenter, len(result))
	for i, item := range result {
		resp[i] = *new(presenter.WalletPresenter).FromDTO(item)
	}

	return base_response.Success(c, resp)
}

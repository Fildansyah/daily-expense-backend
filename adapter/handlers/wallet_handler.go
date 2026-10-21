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
	authID, ok := c.Locals("userID").(string)
	if !ok || authID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var dto entity.WalletDTO
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	// UserID selalu dari token, bukan dari body request.
	dto.UserID = authID

	result, err := w.service.CreateWallet(&dto)
	if err != nil {
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return base_response.Created(c, new(presenter.WalletPresenter).FromDTO(result))
}

// GetMyWallets hanya mengembalikan wallet milik user yang sedang login.
func (w *WalletHandler) GetMyWallets(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := w.service.FindWalletByUserID(c.Context(), userID)
	if err != nil {
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := make([]presenter.WalletPresenter, len(result))
	for i, item := range result {
		resp[i] = *new(presenter.WalletPresenter).FromDTO(item)
	}

	return base_response.Success(c, resp)
}

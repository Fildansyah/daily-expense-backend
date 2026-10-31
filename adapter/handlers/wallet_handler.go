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
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var dto entity.WalletDTO
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}
	dto.UserID = userID

	result, err := w.service.CreateWallet(c.Context(), &dto)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Created(c, new(presenter.WalletPresenter).FromDTO(result))
}

// GetMyWallets hanya mengembalikan wallet milik user yang sedang login.
func (w *WalletHandler) GetMyWallets(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := w.service.FindWalletByUserID(c.Context(), userID, c.QueryBool("includeArchived", false))
	if err != nil {
		return respondError(c, err)
	}

	resp := make([]presenter.WalletPresenter, len(result))
	for i, item := range result {
		resp[i] = *new(presenter.WalletPresenter).FromDTO(item)
	}

	return base_response.Success(c, resp)
}

func (w *WalletHandler) GetWallet(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := w.service.GetWallet(c.Context(), c.Params("id"), userID)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, new(presenter.WalletPresenter).FromDTO(result))
}

func (w *WalletHandler) UpdateWallet(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var dto entity.WalletDTO
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}
	dto.ID = c.Params("id")
	dto.UserID = userID

	result, err := w.service.UpdateWallet(c.Context(), &dto)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, new(presenter.WalletPresenter).FromDTO(result))
}

func (w *WalletHandler) DeleteWallet(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	if err := w.service.DeleteWallet(c.Context(), c.Params("id"), userID, c.QueryBool("archive", false)); err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, fiber.Map{"message": "Akun berhasil diproses"})
}

func (w *WalletHandler) GetSummary(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := w.service.GetSummary(c.Context(), userID)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, new(presenter.WalletSummaryPresenter).FromDTO(result))
}

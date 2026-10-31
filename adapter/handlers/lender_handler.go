package handlers

import (
	"expense.com/m/adapter/handlers/base_response"
	"expense.com/m/adapter/presenter"
	"expense.com/m/entity"
	"expense.com/m/module/lender"
	"github.com/gofiber/fiber/v2"
)

type LenderHandler struct {
	service *lender.LenderService
}

func NewLenderHandler(service *lender.LenderService) *LenderHandler {
	return &LenderHandler{service: service}
}

func (h *LenderHandler) GetLenders(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.GetLenders(c.Context(), userID)
	if err != nil {
		return respondError(c, err)
	}

	resp := make([]presenter.LenderPresenter, len(result))
	for i, item := range result {
		resp[i] = *new(presenter.LenderPresenter).FromDTO(item)
	}

	return base_response.Success(c, resp)
}

func (h *LenderHandler) CreateLender(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var payload entity.LenderPayload
	if err := c.BodyParser(&payload); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	result, err := h.service.CreateLender(c.Context(), userID, payload.Name)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Created(c, new(presenter.LenderPresenter).FromDTO(result))
}

func (h *LenderHandler) UpdateLender(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var payload entity.LenderPayload
	if err := c.BodyParser(&payload); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	result, err := h.service.UpdateLender(c.Context(), userID, c.Params("id"), payload.Name)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, new(presenter.LenderPresenter).FromDTO(result))
}

func (h *LenderHandler) DeleteLender(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	if err := h.service.DeleteLender(c.Context(), userID, c.Params("id")); err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, fiber.Map{"message": "Lender berhasil dihapus"})
}

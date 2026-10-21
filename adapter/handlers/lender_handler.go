package handlers

import (
	"expense.com/m/adapter/handlers/base_response"
	"expense.com/m/adapter/presenter"
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
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := make([]presenter.LenderPresenter, len(result))
	for i, item := range result {
		resp[i] = *new(presenter.LenderPresenter).FromDTO(item)
	}

	return base_response.Success(c, resp)
}

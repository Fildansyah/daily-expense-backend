package handlers

import (
	"expense.com/m/adapter/handlers/base_response"
	"expense.com/m/entity"
	"expense.com/m/module/budget"
	"github.com/gofiber/fiber/v2"
)

type BudgetHandler struct {
	service *budget.BudgetService
}

func NewBudgetHandler(service *budget.BudgetService) *BudgetHandler {
	return &BudgetHandler{service: service}
}

func (h *BudgetHandler) GetBudgets(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.GetBudgets(c.Context(), userID, c.Query("month"))
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, result)
}

func (h *BudgetHandler) GetSummary(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.GetSummary(c.Context(), userID, c.Query("month"))
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, result)
}

func (h *BudgetHandler) CreateBudget(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var payload entity.BudgetPayload
	if err := c.BodyParser(&payload); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	result, err := h.service.CreateBudget(c.Context(), userID, &payload)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Created(c, result)
}

func (h *BudgetHandler) UpdateBudget(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var payload entity.BudgetPayload
	if err := c.BodyParser(&payload); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	result, err := h.service.UpdateBudget(c.Context(), userID, c.Params("id"), payload.Amount)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, result)
}

func (h *BudgetHandler) DeleteBudget(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	if err := h.service.DeleteBudget(c.Context(), userID, c.Params("id")); err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, fiber.Map{"message": "Anggaran berhasil dihapus"})
}

func (h *BudgetHandler) Copy(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var payload entity.BudgetCopyPayload
	if err := c.BodyParser(&payload); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	skipExisting := true
	if payload.SkipExisting != nil {
		skipExisting = *payload.SkipExisting
	}

	created, err := h.service.Copy(c.Context(), userID, payload.FromMonth, payload.ToMonth, skipExisting)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, fiber.Map{
		"message": "Anggaran berhasil disalin",
		"created": created,
	})
}

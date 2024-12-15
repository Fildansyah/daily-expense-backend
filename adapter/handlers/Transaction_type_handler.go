package handlers

import (
	"expense.com/m/adapter/handlers/base_response"
	"expense.com/m/adapter/presenter"
	"expense.com/m/entity"
	"expense.com/m/module/transaction_type"
	"github.com/gofiber/fiber/v2"
)

type TransactionTypeHandler struct {
	service *transaction_type.TransactionTypeService
}

func NewTransactionTypeHandler(service *transaction_type.TransactionTypeService) *TransactionTypeHandler {
	return &TransactionTypeHandler{service: service}
}

func (h *TransactionTypeHandler) CreateTransactionType(c *fiber.Ctx) error {
	var dto entity.TransactionTypeDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	result, err := h.service.CreateTransactionType(&dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func (h *TransactionTypeHandler) GetTransactionTypeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.service.GetTransactionTypeByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// UpdateTransactionType handles updating an existing transaction type.
func (h *TransactionTypeHandler) UpdateTransactionType(c *fiber.Ctx) error {
	var dto entity.TransactionTypeDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.service.UpdateTransactionType(&dto); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Transaction type updated successfully"})
}

// DeleteTransactionType handles deleting a transaction type by ID.
func (h *TransactionTypeHandler) DeleteTransactionType(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteTransactionType(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Transaction type deleted successfully"})
}

// GetAllTransactionTypes handles fetching all transaction types.
func (h *TransactionTypeHandler) GetAllTransactionTypes(c *fiber.Ctx) error {
	result, err := h.service.GetAllTransactionTypes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	resp := make([]presenter.TransactionTypePresenter, len(result))
	for i, item := range result {
		resp[i] = *new(presenter.TransactionTypePresenter).FromDTO(item)
	}

	return base_response.Success(c, resp)
}

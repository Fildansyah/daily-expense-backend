package handlers

import (
	"errors"

	"expense.com/m/adapter/handlers/base_response"
	"expense.com/m/adapter/presenter"
	"expense.com/m/entity"
	"expense.com/m/module/transaction_type"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TransactionTypeHandler struct {
	service *transaction_type.TransactionTypeService
}

func NewTransactionTypeHandler(service *transaction_type.TransactionTypeService) *TransactionTypeHandler {
	return &TransactionTypeHandler{service: service}
}

func (h *TransactionTypeHandler) CreateTransactionType(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var dto entity.TransactionTypeDTO
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	dto.UserID = userID

	result, err := h.service.CreateTransactionType(&dto)
	if err != nil {
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return base_response.Created(c, new(presenter.TransactionTypePresenter).FromDTO(result))
}

func (h *TransactionTypeHandler) GetTransactionTypeByID(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.GetTransactionTypeByID(c.Params("id"), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return base_response.Error(c, fiber.StatusNotFound, "Transaction type tidak ditemukan")
		}
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return base_response.Success(c, new(presenter.TransactionTypePresenter).FromDTO(result))
}

func (h *TransactionTypeHandler) UpdateTransactionType(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var dto entity.TransactionTypeDTO
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}
	if dto.ID == "" {
		return base_response.Error(c, fiber.StatusBadRequest, "ID transaction type wajib diisi")
	}

	dto.UserID = userID

	if err := h.service.UpdateTransactionType(&dto); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return base_response.Error(c, fiber.StatusNotFound, "Transaction type tidak ditemukan")
		}
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return base_response.Success(c, fiber.Map{"message": "Transaction type berhasil diperbarui"})
}

func (h *TransactionTypeHandler) DeleteTransactionType(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	if err := h.service.DeleteTransactionType(c.Params("id"), userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return base_response.Error(c, fiber.StatusNotFound, "Transaction type tidak ditemukan")
		}
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return base_response.Success(c, fiber.Map{"message": "Transaction type berhasil dihapus"})
}

func (h *TransactionTypeHandler) GetAllTransactionTypes(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.GetAllTransactionTypes(userID)
	if err != nil {
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := make([]presenter.TransactionTypePresenter, len(result))
	for i, item := range result {
		resp[i] = *new(presenter.TransactionTypePresenter).FromDTO(item)
	}

	return base_response.Success(c, resp)
}

package handlers

import (
	"strconv"

	"expense.com/m/adapter/handlers/base_response"
	"expense.com/m/adapter/presenter"
	"expense.com/m/entity"
	"expense.com/m/module/transaction"
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	service *transaction.TransactionService
}

func NewTransactionHandler(service *transaction.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) Create(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var dto entity.TransactionDTO
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}
	dto.UserID = userID

	result, err := h.service.Create(c.Context(), &dto)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Created(c, new(presenter.TransactionPresenter).FromDTO(result))
}

func (h *TransactionHandler) List(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	filter := entity.TransactionFilter{
		From:       c.Query("from"),
		To:         c.Query("to"),
		Type:       c.Query("type"),
		CategoryId: c.Query("categoryId"),
		WalletId:   c.Query("walletId"),
		Q:          c.Query("q"),
		Page:       c.QueryInt("page", 1),
		Limit:      c.QueryInt("limit", 20),
		Sort:       c.Query("sort"),
	}
	if value := c.Query("minAmount"); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			filter.MinAmount = &parsed
		}
	}
	if value := c.Query("maxAmount"); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			filter.MaxAmount = &parsed
		}
	}
	items, total, err := h.service.List(c.Context(), userID, &filter)
	if err != nil {
		return respondError(c, err)
	}

	resp := make([]presenter.TransactionPresenter, len(items))
	for i, item := range items {
		resp[i] = *new(presenter.TransactionPresenter).FromDTO(item)
	}

	totalPages := 0
	if filter.Limit > 0 {
		totalPages = int((total + int64(filter.Limit) - 1) / int64(filter.Limit))
	}

	return base_response.Success(c, fiber.Map{
		"items": resp,
		"meta": fiber.Map{
			"page":       filter.Page,
			"limit":      filter.Limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

func (h *TransactionHandler) Get(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.Get(c.Context(), c.Params("id"), userID)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, new(presenter.TransactionPresenter).FromDTO(result))
}

func (h *TransactionHandler) Update(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var dto entity.TransactionDTO
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}
	dto.ID = c.Params("id")
	dto.UserID = userID

	result, err := h.service.Update(c.Context(), &dto)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, new(presenter.TransactionPresenter).FromDTO(result))
}

func (h *TransactionHandler) Delete(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	if err := h.service.Delete(c.Context(), c.Params("id"), userID); err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, fiber.Map{"message": "Transaksi berhasil dihapus"})
}

func (h *TransactionHandler) Bulk(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var dto entity.BulkTransactionDTO
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	if err := h.service.Bulk(c.Context(), userID, &dto); err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, fiber.Map{"message": "Transaksi berhasil diperbarui"})
}

func (h *TransactionHandler) Summary(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.Summary(c.Context(), userID, c.Query("from"), c.Query("to"))
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, result)
}

func (h *TransactionHandler) Transfer(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var dto entity.TransactionDTO
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}
	dto.UserID = userID

	result, err := h.service.Transfer(c.Context(), &dto)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Created(c, new(presenter.TransactionPresenter).FromDTO(result))
}

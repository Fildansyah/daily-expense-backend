package handlers

import (
	"expense.com/m/adapter/handlers/base_response"
	"expense.com/m/adapter/presenter"
	"expense.com/m/entity"
	"expense.com/m/module/bill"
	"github.com/gofiber/fiber/v2"
)

type BillHandler struct {
	service *bill.BillService
}

func NewBillHandler(service *bill.BillService) *BillHandler {
	return &BillHandler{service: service}
}

func (h *BillHandler) CreateBill(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var dto entity.BillDTO
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}
	dto.UserID = userID

	result, err := h.service.CreateBill(c.Context(), &dto)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Created(c, new(presenter.BillPresenter).FromDTO(result))
}

func (h *BillHandler) GetBills(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	filter := entity.BillFilter{
		From:     c.Query("from"),
		To:       c.Query("to"),
		Lender:   c.Query("lender"),
		Status:   c.Query("status"),
		WalletId: c.Query("walletId"),
		Q:        c.Query("q"),
		Page:     c.QueryInt("page", 1),
		Limit:    c.QueryInt("limit", 20),
		Sort:     c.Query("sort"),
	}

	items, total, err := h.service.GetBills(c.Context(), userID, &filter)
	if err != nil {
		return respondError(c, err)
	}

	resp := make([]presenter.BillPresenter, len(items))
	for i, item := range items {
		resp[i] = *new(presenter.BillPresenter).FromDTO(item)
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

func (h *BillHandler) GetBillByID(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.GetBillByID(c.Context(), c.Params("id"), userID)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, new(presenter.BillPresenter).FromDTO(result))
}

func (h *BillHandler) UpdateBill(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var dto entity.BillDTO
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}
	dto.ID = c.Params("id")
	dto.UserID = userID

	result, err := h.service.UpdateBill(c.Context(), &dto)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, new(presenter.BillPresenter).FromDTO(result))
}

func (h *BillHandler) DeleteBill(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	if err := h.service.DeleteBill(c.Context(), c.Params("id"), userID); err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, fiber.Map{"message": "Tagihan berhasil dihapus"})
}

func (h *BillHandler) GetSummary(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.GetBillSummary(c.Context(), userID, c.Query("from"), c.Query("to"))
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, new(presenter.BillSummaryPresenter).FromDTO(result))
}

func (h *BillHandler) GetPayments(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	payments, err := h.service.GetPayments(c.Context(), c.Params("id"), userID)
	if err != nil {
		return respondError(c, err)
	}

	resp := make([]presenter.BillPaymentPresenter, len(payments))
	for i, payment := range payments {
		resp[i] = *new(presenter.BillPaymentPresenter).FromDTO(payment)
	}

	return base_response.Success(c, resp)
}

func (h *BillHandler) AddPayment(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var dto entity.BillPaymentDTO
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	payment, updatedBill, transaction, err := h.service.AddPayment(c.Context(), c.Params("id"), userID, &dto)
	if err != nil {
		return respondError(c, err)
	}

	payload := fiber.Map{
		"payment": new(presenter.BillPaymentPresenter).FromDTO(payment),
		"bill":    new(presenter.BillPresenter).FromDTO(updatedBill),
	}
	if transaction != nil {
		payload["transaction"] = new(presenter.TransactionPresenter).FromDTO(transaction)
	}

	return base_response.Created(c, payload)
}

func (h *BillHandler) DeletePayment(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	if err := h.service.DeletePayment(c.Context(), c.Params("id"), c.Params("paymentId"), userID); err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, fiber.Map{"message": "Pembayaran berhasil dihapus"})
}

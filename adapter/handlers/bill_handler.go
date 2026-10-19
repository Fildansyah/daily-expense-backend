package handlers

import (
	"errors"
	"fmt"
	"time"

	"expense.com/m/adapter/handlers/base_response"
	"expense.com/m/adapter/presenter"
	"expense.com/m/entity"
	"expense.com/m/module/bill"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BillHandler struct {
	service *bill.BillService
}

func NewBillHandler(service *bill.BillService) *BillHandler {
	return &BillHandler{service: service}
}

func (h *BillHandler) validateInput(dto *entity.BillDTO) error {
	if dto.Lender == "" {
		return fmt.Errorf("lender is required")
	}
	if dto.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	if dto.PaidAmount < 0 {
		return fmt.Errorf("paidAmount cannot be negative")
	}
	if _, err := time.Parse("2006-01-02", dto.DueDate); err != nil {
		return fmt.Errorf("dueDate must be in yyyy-mm-dd format")
	}
	return nil
}

func (h *BillHandler) CreateBill(c *fiber.Ctx) error {
	var dto entity.BillDTO

	authId, ok := c.Locals("userID").(string)
	if !ok {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	dto.UserID = authId

	if err := h.validateInput(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	result, err := h.service.CreateBill(&dto)
	if err != nil {
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return base_response.Success(c, new(presenter.BillPresenter).FromDTO(result))
}

func (h *BillHandler) GetBills(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	if err := validateDateParam(c.Query("from"), "from"); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	if err := validateDateParam(c.Query("to"), "to"); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	filter := entity.BillFilter{
		From:   c.Query("from"),
		To:     c.Query("to"),
		Lender: c.Query("lender"),
		Status: c.Query("status"),
		Q:      c.Query("q"),
	}

	result, err := h.service.GetBills(c.Context(), userID, filter)
	if err != nil {
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := make([]presenter.BillPresenter, len(result))
	for i, item := range result {
		resp[i] = *new(presenter.BillPresenter).FromDTO(item)
	}

	return base_response.Success(c, resp)
}

func (h *BillHandler) GetBillByID(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.GetBillByID(c.Context(), c.Params("id"), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return base_response.Error(c, fiber.StatusNotFound, "Bill not found")
		}
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return base_response.Success(c, new(presenter.BillPresenter).FromDTO(result))
}

func (h *BillHandler) UpdateBill(c *fiber.Ctx) error {
	var dto entity.BillDTO

	userID, ok := c.Locals("userID").(string)
	if !ok {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	dto.ID = c.Params("id")
	dto.UserID = userID

	if err := h.validateInput(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	result, err := h.service.UpdateBill(&dto)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return base_response.Error(c, fiber.StatusNotFound, "Bill not found")
		}
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return base_response.Success(c, new(presenter.BillPresenter).FromDTO(result))
}

func (h *BillHandler) DeleteBill(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	if err := h.service.DeleteBill(c.Params("id"), userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return base_response.Error(c, fiber.StatusNotFound, "Bill not found")
		}
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return base_response.Success(c, fiber.Map{"message": "Bill deleted successfully"})
}

func (h *BillHandler) GetSummary(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	if err := validateDateParam(c.Query("from"), "from"); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	if err := validateDateParam(c.Query("to"), "to"); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	filter := entity.BillFilter{
		From: c.Query("from"),
		To:   c.Query("to"),
	}

	result, err := h.service.GetBillSummary(c.Context(), userID, filter)
	if err != nil {
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return base_response.Success(c, new(presenter.BillSummaryPresenter).FromDTO(result))
}

func validateDateParam(val, name string) error {
	if val == "" {
		return nil
	}
	if _, err := time.Parse("2006-01-02", val); err != nil {
		return fmt.Errorf("invalid %s date, expected yyyy-mm-dd", name)
	}
	return nil
}

package handlers

import (
	"expense.com/m/adapter/handlers/base_response"
	"expense.com/m/module/report"
	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	service *report.Service
}

func NewReportHandler(service *report.Service) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) Dashboard(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.Dashboard(c.Context(), userID, c.Query("from"), c.Query("to"))
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, result)
}

func (h *ReportHandler) Calendar(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.Calendar(c.Context(), userID, c.Query("from"), c.Query("to"))
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, result)
}

func (h *ReportHandler) Summary(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.Summary(c.Context(), userID, c.Query("from"), c.Query("to"), c.QueryBool("compare", false))
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, result)
}

func (h *ReportHandler) Cashflow(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.Cashflow(c.Context(), userID, c.Query("from"), c.Query("to"), c.Query("granularity"))
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, result)
}

func (h *ReportHandler) Categories(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.Categories(c.Context(), userID, c.Query("from"), c.Query("to"), c.Query("type"))
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, result)
}

func (h *ReportHandler) Trend(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.Trend(c.Context(), userID, c.Query("from"), c.Query("to"), c.QueryBool("compare", false))
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, result)
}

func (h *ReportHandler) Debt(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.Debt(c.Context(), userID, c.Query("from"), c.Query("to"))
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, result)
}

func (h *ReportHandler) NetWorth(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.NetWorth(c.Context(), userID, c.Query("from"), c.Query("to"))
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, result)
}

package handlers

import (
	"expense.com/m/adapter/handlers/base_response"
	"expense.com/m/adapter/presenter"
	"expense.com/m/entity"
	"expense.com/m/module/category"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	service *category.CategoryService
}

func NewCategoryHandler(service *category.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var dto entity.CategoryDTO
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	dto.UserID = userID

	result, err := h.service.CreateCategory(c.Context(), &dto)
	if err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return base_response.Created(c, new(presenter.CategoryPresenter).FromDTO(result))
}

func (h *CategoryHandler) FindCategoryByTransactionTypeID(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.FindCategoryByTransactionTypeID(c.Context(), c.Params("id"), userID)
	if err != nil {
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	resp := make([]presenter.CategoryPresenter, len(result))
	for i, item := range result {
		resp[i] = *new(presenter.CategoryPresenter).FromDTO(item)
	}

	return base_response.Success(c, resp)
}

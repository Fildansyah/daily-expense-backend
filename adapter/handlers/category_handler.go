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
	var dto entity.CategoryDTO
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	result, err := h.service.CreateCategory(&dto)

	if err != nil {
		return base_response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return base_response.Success(c, result)
}

func (h *CategoryHandler) FindCategoryByTransactionTypeID(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := h.service.FindCategoryByTransactionTypeID(c.Context(), id)

	if err != nil {
		return base_response.Error(c, fiber.StatusNotFound, err.Error())
	}

	resp := make([]presenter.CategoryPresenter, len(result))
	for i, item := range result {
		resp[i] = *new(presenter.CategoryPresenter).FromDTO(item)
	}

	return base_response.Success(c, resp)
}

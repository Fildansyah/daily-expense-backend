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
		return respondError(c, err)
	}

	return base_response.Created(c, new(presenter.CategoryPresenter).FromDTO(result))
}

func (h *CategoryHandler) GetCategories(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.GetTree(c.Context(), userID, c.Query("type"), c.QueryBool("includeArchived", false))
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, presenter.CategoryPresenterList(result))
}

func (h *CategoryHandler) GetTree(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.GetTree(c.Context(), userID, c.Query("type"), c.QueryBool("includeArchived", false))
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, presenter.CategoryPresenterList(result))
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var dto entity.CategoryDTO
	if err := c.BodyParser(&dto); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}
	dto.ID = c.Params("id")
	dto.UserID = userID

	result, err := h.service.UpdateCategory(c.Context(), &dto)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, new(presenter.CategoryPresenter).FromDTO(result))
}

func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	if err := h.service.DeleteCategory(c.Context(), c.Params("id"), userID, c.QueryBool("archive", false)); err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, fiber.Map{"message": "Kategori berhasil diproses"})
}

func (h *CategoryHandler) Reorder(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var req entity.CategoryReorderRequest
	if err := c.BodyParser(&req); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	if err := h.service.Reorder(c.Context(), userID, req.Items); err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, fiber.Map{"message": "Urutan kategori berhasil disimpan"})
}

type mergeCategoryRequest struct {
	SourceId string `json:"sourceId"`
	TargetId string `json:"targetId"`
}

func (h *CategoryHandler) Merge(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	var req mergeCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return base_response.Error(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	if err := h.service.Merge(c.Context(), userID, req.SourceId, req.TargetId); err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, fiber.Map{"message": "Kategori berhasil digabungkan"})
}

func (h *CategoryHandler) FindCategoryByTransactionTypeID(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return base_response.Error(c, fiber.StatusUnauthorized, "Invalid user ID")
	}

	result, err := h.service.FindCategoryByTransactionTypeID(c.Context(), c.Params("id"), userID)
	if err != nil {
		return respondError(c, err)
	}

	return base_response.Success(c, presenter.CategoryPresenterList(result))
}

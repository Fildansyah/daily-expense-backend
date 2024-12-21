package category

import (
	"context"

	"expense.com/m/entity"
)

type CategoryService struct {
	repo *CategoryRepository
}

func NewCategoryService(repo *CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) FindCategoryByTransactionTypeID(ctx context.Context, id string) ([]*entity.CategoryDTO, error) {
	return s.repo.FindCategoryByTransactionTypeID(ctx, id)
}

func (s *CategoryService) CreateCategory(input *entity.CategoryDTO) (*entity.CategoryDTO, error) {
	return s.repo.CreateCategory(input)
}

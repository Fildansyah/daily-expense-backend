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

func (s *CategoryService) FindCategoryByTransactionTypeID(ctx context.Context, transactionTypeID, userID string) ([]*entity.CategoryDTO, error) {
	return s.repo.FindCategoryByTransactionTypeID(ctx, transactionTypeID, userID)
}

func (s *CategoryService) CreateCategory(ctx context.Context, input *entity.CategoryDTO) (*entity.CategoryDTO, error) {
	return s.repo.CreateCategory(ctx, input)
}

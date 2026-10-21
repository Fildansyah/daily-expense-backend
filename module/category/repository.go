package category

import (
	"context"
	"errors"
	"fmt"

	"expense.com/m/entity"
	"expense.com/m/model"
	baserepository "expense.com/m/shared/base_repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	*baserepository.BaseRepository[model.Category]
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		BaseRepository: baserepository.NewBaseRepository[model.Category](db),
		db:             db,
	}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, input *entity.CategoryDTO) (*entity.CategoryDTO, error) {
	var transactionType model.TransactionType
	err := r.db.WithContext(ctx).
		Where("id = ? AND (user_id IS NULL OR user_id = ?)", input.TransactionTypeId, input.UserID).
		First(&transactionType).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("transaction type tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	modelEntity := input.ToModel()
	modelEntity.ID = uuid.NewString()

	if err := r.Create(modelEntity); err != nil {
		return nil, err
	}
	return input.FromModel(modelEntity), nil
}

func (r *CategoryRepository) FindCategoryByTransactionTypeID(ctx context.Context, transactionTypeID, userID string) ([]*entity.CategoryDTO, error) {
	var entities []model.Category
	err := r.db.WithContext(ctx).Model(&model.Category{}).
		Preload("TransactionTypes").
		Where("categories.transaction_type_id = ?", transactionTypeID).
		Where("categories.user_id IS NULL OR categories.user_id = ?", userID).
		Find(&entities).Error

	if err != nil {
		return nil, err
	}

	dtos := make([]*entity.CategoryDTO, len(entities))
	for i, modelEntity := range entities {
		dto := new(entity.CategoryDTO)
		dtos[i] = dto.FromModel(&modelEntity)
	}

	return dtos, nil
}

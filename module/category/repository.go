package category

import (
	"context"

	"expense.com/m/entity"
	"expense.com/m/model"
	"expense.com/m/module/base_repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	*base_repository.BaseRepository[model.Category]
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		BaseRepository: base_repository.NewBaseRepository[model.Category](db),
		db:             db,
	}
}

func (r *CategoryRepository) CreateCategory(input *entity.CategoryDTO) (*entity.CategoryDTO, error) {
	modelEntity := input.ToModel()
	modelEntity.ID = uuid.NewString()

	if err := r.Create(modelEntity); err != nil {
		return nil, err
	}
	return input.FromModel(modelEntity), nil
}

func (r *CategoryRepository) FindCategoryByTransactionTypeID(ctx context.Context, id string) ([]*entity.CategoryDTO, error) {
	var entities []model.Category
	err := r.db.WithContext(ctx).Model(&model.Category{}).
		Preload("TransactionTypes").
		Where("categories.transaction_type_id = ?", id).
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

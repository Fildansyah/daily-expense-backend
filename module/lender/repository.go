package lender

import (
	"context"

	"expense.com/m/entity"
	"expense.com/m/model"
	baserepository "expense.com/m/shared/base_repository"
	"gorm.io/gorm"
)

type LenderRepository struct {
	*baserepository.BaseRepository[model.Lender]
	db *gorm.DB
}

func NewLenderRepository(db *gorm.DB) *LenderRepository {
	return &LenderRepository{
		BaseRepository: baserepository.NewBaseRepository[model.Lender](db),
		db:             db,
	}
}

func (r *LenderRepository) FindAllLenders(ctx context.Context) ([]*entity.LenderDTO, error) {
	var entities []model.Lender
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&entities).Error; err != nil {
		return nil, err
	}

	dtos := make([]*entity.LenderDTO, len(entities))
	for i, modelEntity := range entities {
		dtos[i] = new(entity.LenderDTO).FromModel(&modelEntity)
	}

	return dtos, nil
}

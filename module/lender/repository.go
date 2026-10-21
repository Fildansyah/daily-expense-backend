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

// FindAllLenders mengembalikan lender global (user_id NULL) dan milik user.
func (r *LenderRepository) FindAllLenders(ctx context.Context, userID string) ([]*entity.LenderDTO, error) {
	var entities []model.Lender
	err := r.db.WithContext(ctx).
		Where("user_id IS NULL OR user_id = ?", userID).
		Order("name ASC").
		Find(&entities).Error
	if err != nil {
		return nil, err
	}

	dtos := make([]*entity.LenderDTO, len(entities))
	for i, modelEntity := range entities {
		dtos[i] = new(entity.LenderDTO).FromModel(&modelEntity)
	}

	return dtos, nil
}

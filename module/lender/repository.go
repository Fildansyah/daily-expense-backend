package lender

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"expense.com/m/entity"
	"expense.com/m/model"
	"expense.com/m/model/concern"
	"expense.com/m/shared/apperror"
	"expense.com/m/shared/tx"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LenderRepository struct {
	db *gorm.DB
}

func NewLenderRepository(db *gorm.DB) *LenderRepository {
	return &LenderRepository{db: db}
}

func (r *LenderRepository) FindAll(ctx context.Context, userID string) ([]*entity.LenderDTO, error) {
	var rows []model.Lender
	err := tx.DB(ctx, r.db).
		Where("user_id IS NULL OR user_id = ?", userID).
		Order("name ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]*entity.LenderDTO, len(rows))
	for i := range rows {
		items[i] = fromModel(&rows[i])
	}
	return items, nil
}

func (r *LenderRepository) FindByID(ctx context.Context, id, userID string) (*entity.LenderDTO, error) {
	return r.findByID(tx.DB(ctx, r.db), id, userID, false)
}

func (r *LenderRepository) FindOwned(ctx context.Context, id, userID string) (*entity.LenderDTO, error) {
	return r.findByID(tx.DB(ctx, r.db), id, userID, true)
}

func (r *LenderRepository) Create(ctx context.Context, input *entity.LenderDTO) (*entity.LenderDTO, error) {
	modelEntity := toModel(input)
	modelEntity.ID = uuid.NewString()

	if err := tx.DB(ctx, r.db).Create(modelEntity).Error; err != nil {
		return nil, err
	}
	return fromModel(modelEntity), nil
}

func (r *LenderRepository) Update(ctx context.Context, input *entity.LenderDTO) error {
	res := tx.DB(ctx, r.db).
		Model(&model.Lender{}).
		Where("id = ? AND user_id = ?", input.ID, input.UserID).
		Update("name", input.Name)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: lender tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *LenderRepository) Delete(ctx context.Context, id, userID string) error {
	res := tx.DB(ctx, r.db).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Lender{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: lender tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *LenderRepository) ExistsName(ctx context.Context, userID, name, excludeID string) (bool, error) {
	query := tx.DB(ctx, r.db).Model(&model.Lender{}).
		Where("LOWER(name) = ? AND (user_id IS NULL OR user_id = ?)", strings.ToLower(name), userID)
	if excludeID != "" {
		query = query.Where("id <> ?", excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *LenderRepository) CountBills(ctx context.Context, userID, lenderID string) (int64, error) {
	var count int64
	err := tx.DB(ctx, r.db).Model(&model.Bill{}).
		Where("user_id = ? AND lender_id = ?", userID, lenderID).
		Count(&count).Error
	return count, err
}

func (r *LenderRepository) findByID(db *gorm.DB, id, userID string, ownedOnly bool) (*entity.LenderDTO, error) {
	query := db.Where("id = ?", id)
	if ownedOnly {
		query = query.Where("user_id = ?", userID)
	} else {
		query = query.Where("user_id IS NULL OR user_id = ?", userID)
	}

	var modelEntity model.Lender
	if err := query.First(&modelEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: lender tidak ditemukan", apperror.ErrNotFound)
		}
		return nil, err
	}
	return fromModel(&modelEntity), nil
}

func toModel(input *entity.LenderDTO) *model.Lender {
	return &model.Lender{
		BaseFields: concern.BaseFields{
			ID:        input.ID,
			CreatedAt: input.CreatedAt,
		},
		UserID: input.UserID,
		Name:   input.Name,
	}
}

func fromModel(modelEntity *model.Lender) *entity.LenderDTO {
	if modelEntity == nil {
		return nil
	}

	return &entity.LenderDTO{
		CommonDTOFields: entity.CommonDTOFields{
			ID:        modelEntity.ID,
			CreatedAt: modelEntity.CreatedAt,
		},
		UserID: modelEntity.UserID,
		Name:   modelEntity.Name,
	}
}

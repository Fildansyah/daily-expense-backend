package budget

import (
	"context"
	"errors"
	"fmt"
	"time"

	"expense.com/m/entity"
	"expense.com/m/model"
	"expense.com/m/model/concern"
	"expense.com/m/shared/apperror"
	"expense.com/m/shared/tx"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	dateLayout  = "2006-01-02"
	monthLayout = "2006-01"
)

type BudgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) Create(ctx context.Context, input *entity.BudgetDTO) (*entity.BudgetDTO, error) {
	modelEntity := toModel(input)
	modelEntity.ID = uuid.NewString()

	if err := tx.DB(ctx, r.db).Create(modelEntity).Error; err != nil {
		return nil, err
	}
	return fromModel(modelEntity), nil
}

func (r *BudgetRepository) UpdateAmount(ctx context.Context, id, userID string, amount float64) error {
	res := tx.DB(ctx, r.db).
		Model(&model.Budget{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("amount", amount)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: anggaran tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *BudgetRepository) Delete(ctx context.Context, id, userID string) error {
	res := tx.DB(ctx, r.db).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Budget{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: anggaran tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *BudgetRepository) FindByID(ctx context.Context, id, userID string) (*entity.BudgetDTO, error) {
	var modelEntity model.Budget
	err := tx.DB(ctx, r.db).
		Where("id = ? AND user_id = ?", id, userID).
		First(&modelEntity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("%w: anggaran tidak ditemukan", apperror.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	return fromModel(&modelEntity), nil
}

func (r *BudgetRepository) ListByMonth(ctx context.Context, userID, month string) ([]*entity.BudgetDTO, error) {
	var rows []model.Budget
	err := tx.DB(ctx, r.db).
		Where("user_id = ? AND period_month = ?", userID, parseMonth(month)).
		Order("created_at ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]*entity.BudgetDTO, len(rows))
	for i := range rows {
		items[i] = fromModel(&rows[i])
	}
	return items, nil
}

func (r *BudgetRepository) Exists(ctx context.Context, userID, categoryID, month, excludeID string) (bool, error) {
	query := tx.DB(ctx, r.db).Model(&model.Budget{}).
		Where("user_id = ? AND category_id = ? AND period_month = ?", userID, categoryID, parseMonth(month))
	if excludeID != "" {
		query = query.Where("id <> ?", excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func toModel(input *entity.BudgetDTO) *model.Budget {
	return &model.Budget{
		BaseFields: concern.BaseFields{
			ID:        input.ID,
			CreatedAt: input.CreatedAt,
		},
		UserID:      input.UserID,
		CategoryId:  input.CategoryId,
		PeriodMonth: parseMonth(input.PeriodMonth),
		Amount:      input.Amount,
	}
}

func fromModel(modelEntity *model.Budget) *entity.BudgetDTO {
	if modelEntity == nil {
		return nil
	}

	return &entity.BudgetDTO{
		CommonDTOFields: entity.CommonDTOFields{
			ID:        modelEntity.ID,
			CreatedAt: modelEntity.CreatedAt,
		},
		UpdatedAt:   modelEntity.UpdatedAt,
		UserID:      modelEntity.UserID,
		CategoryId:  modelEntity.CategoryId,
		PeriodMonth: modelEntity.PeriodMonth.Format(monthLayout),
		Amount:      modelEntity.Amount,
	}
}

func parseMonth(value string) time.Time {
	parsed, err := time.Parse(monthLayout, value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

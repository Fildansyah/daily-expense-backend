package category

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

type TransactionTypeData struct {
	ID   string
	Type string
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, input *entity.CategoryDTO) (*entity.CategoryDTO, error) {
	modelEntity := toModel(input)
	modelEntity.ID = uuid.NewString()

	if err := tx.DB(ctx, r.db).Create(modelEntity).Error; err != nil {
		return nil, err
	}
	return fromModel(modelEntity), nil
}

func (r *CategoryRepository) UpdateFields(ctx context.Context, input *entity.CategoryDTO) error {
	res := tx.DB(ctx, r.db).
		Model(&model.Category{}).
		Where("id = ? AND user_id = ?", input.ID, input.UserID).
		Updates(map[string]interface{}{
			"name":        input.Name,
			"bg_color":    input.BgColor,
			"icon":        input.Icon,
			"parent_id":   input.ParentId,
			"sort_order":  input.SortOrder,
			"is_archived": input.IsArchived,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: kategori tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id, userID string) error {
	res := tx.DB(ctx, r.db).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Category{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: kategori tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *CategoryRepository) Archive(ctx context.Context, id, userID string) error {
	res := tx.DB(ctx, r.db).
		Model(&model.Category{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_archived", true)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("%w: kategori tidak ditemukan", apperror.ErrNotFound)
	}
	return nil
}

func (r *CategoryRepository) FindByID(ctx context.Context, id, userID string) (*entity.CategoryDTO, error) {
	return r.findByID(tx.DB(ctx, r.db), id, userID, false)
}

func (r *CategoryRepository) FindOwned(ctx context.Context, id, userID string) (*entity.CategoryDTO, error) {
	return r.findByID(tx.DB(ctx, r.db), id, userID, true)
}

func (r *CategoryRepository) FindVisible(ctx context.Context, userID, transactionType string, includeArchived bool) ([]*entity.CategoryDTO, error) {
	query := tx.DB(ctx, r.db).
		Model(&model.Category{}).
		Where("categories.user_id IS NULL OR categories.user_id = ?", userID)

	if !includeArchived {
		query = query.Where("categories.is_archived = FALSE")
	}
	if transactionType != "" {
		query = query.
			Joins("JOIN transaction_types AS tt ON tt.id = categories.transaction_type_id").
			Where("tt.type = ?", strings.ToUpper(transactionType))
	}

	var rows []model.Category
	if err := query.Order("categories.sort_order ASC, categories.name ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]*entity.CategoryDTO, len(rows))
	for i := range rows {
		items[i] = fromModel(&rows[i])
	}
	return items, nil
}

func (r *CategoryRepository) FindTransactionType(ctx context.Context, id, userID string) (*TransactionTypeData, error) {
	var row TransactionTypeData
	err := tx.DB(ctx, r.db).
		Model(&model.TransactionType{}).
		Select("id, type").
		Where("id = ? AND (user_id IS NULL OR user_id = ?)", id, userID).
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == "" {
		return nil, fmt.Errorf("%w: transaction type tidak ditemukan", apperror.ErrNotFound)
	}
	return &row, nil
}

func (r *CategoryRepository) FindTransactionTypeNames(ctx context.Context, userID string, ids []string) (map[string]string, error) {
	names := map[string]string{}
	if len(ids) == 0 {
		return names, nil
	}

	var rows []struct {
		ID   string
		Type string
	}
	err := tx.DB(ctx, r.db).
		Model(&model.TransactionType{}).
		Select("id, type").
		Where("id IN ? AND (user_id IS NULL OR user_id = ?)", ids, userID).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		names[row.ID] = row.Type
	}
	return names, nil
}

func (r *CategoryRepository) CountChildren(ctx context.Context, userID, categoryID string) (int64, error) {
	var count int64
	err := tx.DB(ctx, r.db).
		Model(&model.Category{}).
		Where("parent_id = ? AND user_id = ?", categoryID, userID).
		Count(&count).Error
	return count, err
}

func (r *CategoryRepository) CountTransactions(ctx context.Context, userID, categoryID string) (int64, error) {
	var count int64
	err := tx.DB(ctx, r.db).
		Model(&model.Transaction{}).
		Where("user_id = ? AND category_id = ?", userID, categoryID).
		Count(&count).Error
	return count, err
}

func (r *CategoryRepository) CountByCategory(ctx context.Context, userID string) (map[string]int64, error) {
	var rows []struct {
		CategoryId string
		Count      int64
	}

	err := tx.DB(ctx, r.db).
		Model(&model.Transaction{}).
		Select("category_id, COUNT(*) AS count").
		Where("user_id = ? AND category_id IS NOT NULL", userID).
		Group("category_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64, len(rows))
	for _, row := range rows {
		counts[row.CategoryId] = row.Count
	}
	return counts, nil
}

func (r *CategoryRepository) Reorder(ctx context.Context, userID string, items []entity.CategoryReorderItem) error {
	for _, item := range items {
		res := tx.DB(ctx, r.db).
			Model(&model.Category{}).
			Where("id = ? AND user_id = ?", item.ID, userID).
			Updates(map[string]interface{}{
				"parent_id":  item.ParentId,
				"sort_order": item.SortOrder,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("%w: kategori %s tidak ditemukan", apperror.ErrNotFound, item.ID)
		}
	}
	return nil
}

func (r *CategoryRepository) ReassignTransactions(ctx context.Context, userID, sourceID, targetID string) error {
	return tx.DB(ctx, r.db).
		Model(&model.Transaction{}).
		Where("user_id = ? AND category_id = ?", userID, sourceID).
		Update("category_id", targetID).Error
}

func (r *CategoryRepository) findByID(db *gorm.DB, id, userID string, ownedOnly bool) (*entity.CategoryDTO, error) {
	query := db.Where("id = ?", id)
	if ownedOnly {
		query = query.Where("user_id = ?", userID)
	} else {
		query = query.Where("user_id IS NULL OR user_id = ?", userID)
	}

	var modelEntity model.Category
	if err := query.First(&modelEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: kategori tidak ditemukan", apperror.ErrNotFound)
		}
		return nil, err
	}
	return fromModel(&modelEntity), nil
}

func toModel(input *entity.CategoryDTO) *model.Category {
	return &model.Category{
		BaseFields: concern.BaseFields{
			ID:        input.ID,
			CreatedAt: input.CreatedAt,
		},
		UserID:            input.UserID,
		Name:              model.CategoryNameEnum(input.Name),
		TransactionTypeId: input.TransactionTypeId,
		BgColor:           input.BgColor,
		Icon:              input.Icon,
		ParentId:          input.ParentId,
		SortOrder:         input.SortOrder,
		IsArchived:        input.IsArchived,
	}
}

func fromModel(modelEntity *model.Category) *entity.CategoryDTO {
	if modelEntity == nil {
		return nil
	}

	return &entity.CategoryDTO{
		CommonDTOFields: entity.CommonDTOFields{
			ID:        modelEntity.ID,
			CreatedAt: modelEntity.CreatedAt,
		},
		UpdatedAt:         modelEntity.UpdatedAt,
		UserID:            modelEntity.UserID,
		Name:              string(modelEntity.Name),
		TransactionTypeId: modelEntity.TransactionTypeId,
		BgColor:           modelEntity.BgColor,
		Icon:              modelEntity.Icon,
		ParentId:          modelEntity.ParentId,
		SortOrder:         modelEntity.SortOrder,
		IsArchived:        modelEntity.IsArchived,
	}
}

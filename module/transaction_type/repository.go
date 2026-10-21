package transaction_type

import (
	"expense.com/m/entity"
	"expense.com/m/model"
	baserepository "expense.com/m/shared/base_repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionTypeRepository struct {
	*baserepository.BaseRepository[model.TransactionType]
	db *gorm.DB
}

func NewTransactionTypeRepository(db *gorm.DB) *TransactionTypeRepository {
	return &TransactionTypeRepository{
		BaseRepository: baserepository.NewBaseRepository[model.TransactionType](db),
		db:             db,
	}
}

func (r *TransactionTypeRepository) CreateTransactionType(input *entity.TransactionTypeDTO) (*entity.TransactionTypeDTO, error) {
	modelEntity := input.ToModel()
	modelEntity.ID = uuid.NewString()

	if err := r.Create(modelEntity); err != nil {
		return nil, err
	}
	return input.FromModel(modelEntity), nil
}

func (r *TransactionTypeRepository) FindTransactionTypeByID(id, userID string) (*entity.TransactionTypeDTO, error) {
	var modelEntity model.TransactionType
	err := r.db.
		Where("id = ? AND (user_id IS NULL OR user_id = ?)", id, userID).
		First(&modelEntity).Error
	if err != nil {
		return nil, err
	}
	dto := new(entity.TransactionTypeDTO)
	return dto.FromModel(&modelEntity), nil
}

func (r *TransactionTypeRepository) UpdateTransactionType(input *entity.TransactionTypeDTO) error {
	result := r.db.Model(&model.TransactionType{}).
		Where("id = ? AND user_id = ?", input.ID, input.UserID).
		Update("type", input.Type)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *TransactionTypeRepository) DeleteTransactionType(id, userID string) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.TransactionType{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *TransactionTypeRepository) FindAllTransactionTypes(userID string) ([]*entity.TransactionTypeDTO, error) {
	var modelEntities []model.TransactionType
	err := r.db.
		Where("user_id IS NULL OR user_id = ?", userID).
		Order("created_at ASC").
		Find(&modelEntities).Error
	if err != nil {
		return nil, err
	}

	dtos := make([]*entity.TransactionTypeDTO, len(modelEntities))
	for i, modelEntity := range modelEntities {
		dto := new(entity.TransactionTypeDTO)
		dtos[i] = dto.FromModel(&modelEntity)
	}

	return dtos, nil
}

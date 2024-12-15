package transaction_type

import (
	"expense.com/m/entity"
	"expense.com/m/model"
	"expense.com/m/module/base_repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionTypeRepository struct {
	*base_repository.BaseRepository[model.TransactionType]
}

func NewTransactionTypeRepository(db *gorm.DB) *TransactionTypeRepository {
	return &TransactionTypeRepository{base_repository.NewBaseRepository[model.TransactionType](db)}
}

func (r *TransactionTypeRepository) CreateTransactionType(input *entity.TransactionTypeDTO) (*entity.TransactionTypeDTO, error) {
	modelEntity := input.ToModel()
	modelEntity.ID = uuid.NewString()

	if err := r.Create(modelEntity); err != nil {
		return nil, err
	}
	return input.FromModel(modelEntity), nil
}

func (r *TransactionTypeRepository) FindTransactionTypeByID(id string) (*entity.TransactionTypeDTO, error) {
	modelEntity, err := r.FindByID(id)
	if err != nil {
		return nil, err
	}
	dto := new(entity.TransactionTypeDTO)
	return dto.FromModel(modelEntity), nil
}

func (r *TransactionTypeRepository) UpdateTransactionType(input *entity.TransactionTypeDTO) error {
	modelEntity := input.ToModel()
	return r.Update(modelEntity)
}

func (r *TransactionTypeRepository) DeleteTransactionType(id string) error {
	return r.Delete(id)
}

func (r *TransactionTypeRepository) FindAllTransactionTypes() ([]*entity.TransactionTypeDTO, error) {
	modelEntities, err := r.FindAll()
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

package transaction_type

import (
	"expense.com/m/entity"
)

type TransactionTypeService struct {
	repo *TransactionTypeRepository
}

func NewTransactionTypeService(repo *TransactionTypeRepository) *TransactionTypeService {
	return &TransactionTypeService{repo: repo}
}

func (s *TransactionTypeService) CreateTransactionType(dto *entity.TransactionTypeDTO) (*entity.TransactionTypeDTO, error) {
	return s.repo.CreateTransactionType(dto)
}

func (s *TransactionTypeService) GetTransactionTypeByID(id string) (*entity.TransactionTypeDTO, error) {
	return s.repo.FindTransactionTypeByID(id)
}

func (s *TransactionTypeService) UpdateTransactionType(dto *entity.TransactionTypeDTO) error {
	return s.repo.UpdateTransactionType(dto)
}

func (s *TransactionTypeService) DeleteTransactionType(id string) error {
	return s.repo.DeleteTransactionType(id)
}

func (s *TransactionTypeService) GetAllTransactionTypes() ([]*entity.TransactionTypeDTO, error) {
	return s.repo.FindAllTransactionTypes()
}

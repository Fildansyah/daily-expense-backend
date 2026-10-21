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

func (s *TransactionTypeService) GetTransactionTypeByID(id, userID string) (*entity.TransactionTypeDTO, error) {
	return s.repo.FindTransactionTypeByID(id, userID)
}

func (s *TransactionTypeService) UpdateTransactionType(dto *entity.TransactionTypeDTO) error {
	return s.repo.UpdateTransactionType(dto)
}

func (s *TransactionTypeService) DeleteTransactionType(id, userID string) error {
	return s.repo.DeleteTransactionType(id, userID)
}

func (s *TransactionTypeService) GetAllTransactionTypes(userID string) ([]*entity.TransactionTypeDTO, error) {
	return s.repo.FindAllTransactionTypes(userID)
}

package presenter

import (
	"expense.com/m/entity"
	"expense.com/m/model"
	"expense.com/m/model/concern"
)

type TransactionTypePresenter struct {
	concern.BaseFields
	Type model.TransactionTypeEnum
}

func (v TransactionTypePresenter) FromDTO(dto *entity.TransactionTypeDTO) *TransactionTypePresenter {
	if dto == nil {
		return nil
	}

	return &TransactionTypePresenter{
		BaseFields: concern.BaseFields{
			ID:        dto.ID,
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		},
		Type: dto.Type,
	}
}

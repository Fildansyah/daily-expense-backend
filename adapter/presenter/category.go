package presenter

import (
	"expense.com/m/entity"
	"expense.com/m/model"
	"expense.com/m/model/concern"
)

type CategoryPresenter struct {
	concern.BaseFields
	Name              model.CategoryNameEnum
	TransactionTypeId string
	TransactionTypes  model.TransactionType
	Icon              string
	BgColor           string
}

func (v CategoryPresenter) FromDTO(dto *entity.CategoryDTO) *CategoryPresenter {
	if dto == nil {
		return nil
	}

	return &CategoryPresenter{
		BaseFields: concern.BaseFields{
			ID:        dto.ID,
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		},
		Name:              dto.Name,
		TransactionTypeId: dto.TransactionTypeId,
		TransactionTypes:  dto.TransactionTypes,
		Icon:              dto.Icon,
		BgColor:           dto.BgColor,
	}
}

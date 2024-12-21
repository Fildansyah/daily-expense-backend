package entity

import (
	"expense.com/m/model"
	"expense.com/m/model/concern"
)

type CategoryDTO struct {
	CommonDTOFields

	Name              model.CategoryNameEnum
	TransactionTypeId string
	TransactionTypes  model.TransactionType
	BgColor           string
	Icon              string
}

func (d *CategoryDTO) FromModel(m *model.Category) *CategoryDTO {
	if m == nil {
		return nil
	}

	return &CategoryDTO{
		CommonDTOFields: CommonDTOFields{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
		Name:              m.Name,
		TransactionTypeId: m.TransactionTypeId,
		TransactionTypes:  m.TransactionTypes,
		BgColor:           m.BgColor,
		Icon:              m.Icon,
	}
}

func (d *CategoryDTO) ToModel() *model.Category {
	m := &model.Category{
		BaseFields: concern.BaseFields{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		},
		Name:              d.Name,
		TransactionTypeId: d.TransactionTypeId,
		BgColor:           d.BgColor,
		Icon:              d.Icon,
	}

	return m
}

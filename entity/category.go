package entity

import (
	"expense.com/m/model"
	"expense.com/m/model/concern"
)

type CategoryDTO struct {
	CommonDTOFields

	UserID            string `json:"-"`
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
		},
		UserID:            m.UserID,
		Name:              m.Name,
		TransactionTypeId: m.TransactionTypeId,
		TransactionTypes:  m.TransactionTypes,
		BgColor:           m.BgColor,
		Icon:              m.Icon,
	}
}

func (d *CategoryDTO) ToModel() *model.Category {
	return &model.Category{
		BaseFields: concern.BaseFields{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
		},
		UserID:            d.UserID,
		Name:              d.Name,
		TransactionTypeId: d.TransactionTypeId,
		BgColor:           d.BgColor,
		Icon:              d.Icon,
	}
}

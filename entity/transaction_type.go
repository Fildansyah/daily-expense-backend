package entity

import (
	"expense.com/m/model"
	"expense.com/m/model/concern"
	"gorm.io/gorm"
)

type TransactionTypeDTO struct {
	CommonDTOFields

	Type model.TransactionTypeEnum
}

func (d *TransactionTypeDTO) FromModel(m *model.TransactionType) *TransactionTypeDTO {
	if m == nil {
		return nil
	}

	return &TransactionTypeDTO{
		CommonDTOFields: CommonDTOFields{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
		},
		Type: m.Type,
	}
}

func (d *TransactionTypeDTO) ToModel() *model.TransactionType {
	m := &model.TransactionType{
		BaseFields: concern.BaseFields{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
		},
		Type: d.Type,
	}

	if d.DeletedAt != nil {
		m.DeletedAt = &gorm.DeletedAt{Time: *d.DeletedAt, Valid: true}
	}

	return m
}

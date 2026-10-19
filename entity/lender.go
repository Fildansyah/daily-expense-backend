package entity

import (
	"expense.com/m/model"
	"expense.com/m/model/concern"
)

type LenderDTO struct {
	CommonDTOFields

	Name string `json:"name"`
}

func (d *LenderDTO) FromModel(m *model.Lender) *LenderDTO {
	return &LenderDTO{
		CommonDTOFields: CommonDTOFields{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
		},
		Name: m.Name,
	}
}

func (d *LenderDTO) ToModel() *model.Lender {
	return &model.Lender{
		BaseFields: concern.BaseFields{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
		},
		Name: d.Name,
	}
}

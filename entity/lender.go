package entity

import (
	"expense.com/m/model"
	"expense.com/m/model/concern"
)

type LenderDTO struct {
	CommonDTOFields

	UserID string `json:"-"`
	Name   string `json:"name"`
}

func (d *LenderDTO) FromModel(m *model.Lender) *LenderDTO {
	if m == nil {
		return nil
	}

	return &LenderDTO{
		CommonDTOFields: CommonDTOFields{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
		},
		UserID: m.UserID,
		Name:   m.Name,
	}
}

func (d *LenderDTO) ToModel() *model.Lender {
	return &model.Lender{
		BaseFields: concern.BaseFields{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
		},
		UserID: d.UserID,
		Name:   d.Name,
	}
}

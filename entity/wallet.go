package entity

import (
	"expense.com/m/model"
	"expense.com/m/model/concern"
)

type WalletDTO struct {
	CommonDTOFields

	UserID    string  `json:"-"`
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
	Currency  string  `json:"currency"`
	Type      string  `json:"type"`
	Color     string  `json:"color"`
	Icon      string  `json:"icon"`
	IsExclude bool    `json:"isExclude"`
}

func (d *WalletDTO) FromModel(m *model.Wallet) *WalletDTO {
	return &WalletDTO{
		CommonDTOFields: CommonDTOFields{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
		},
		UserID:    m.UserID,
		Name:      m.Name,
		Balance:   m.Balance,
		Currency:  m.Currency,
		Type:      m.Type,
		Color:     m.Color,
		Icon:      m.Icon,
		IsExclude: m.IsExclude,
	}
}

func (d *WalletDTO) ToModel() *model.Wallet {
	return &model.Wallet{
		BaseFields: concern.BaseFields{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
		},
		UserID:    d.UserID,
		Name:      d.Name,
		Balance:   d.Balance,
		Currency:  d.Currency,
		Type:      d.Type,
		Color:     d.Color,
		Icon:      d.Icon,
		IsExclude: d.IsExclude,
	}
}

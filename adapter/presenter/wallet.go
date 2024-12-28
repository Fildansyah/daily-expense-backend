package presenter

import (
	"expense.com/m/entity"
	"expense.com/m/model"
	"expense.com/m/model/concern"
)

type WalletPresenter struct {
	concern.BaseFields
	UserID    string
	Users     model.Users
	Name      string
	Balance   float64
	Currency  string
	Type      string
	Color     string
	Icon      string
	IsExclude bool
}

func (v WalletPresenter) FromDTO(dto *entity.WalletDTO) *WalletPresenter {
	return &WalletPresenter{
		BaseFields: concern.BaseFields{
			ID:        dto.ID,
			CreatedAt: dto.CreatedAt,
		},
		UserID:    dto.UserID,
		Users:     dto.Users,
		Name:      dto.Name,
		Balance:   dto.Balance,
		Currency:  dto.Currency,
		Type:      dto.Type,
		Color:     dto.Color,
		Icon:      dto.Icon,
		IsExclude: dto.IsExclude,
	}
}

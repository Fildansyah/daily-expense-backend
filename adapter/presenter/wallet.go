package presenter

import (
	"time"

	"expense.com/m/entity"
)

type WalletPresenter struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	Type      string    `json:"type"`
	Color     string    `json:"color"`
	Icon      string    `json:"icon"`
	IsExclude bool      `json:"isExclude"`
	CreatedAt time.Time `json:"createdAt"`
}

func (v WalletPresenter) FromDTO(dto *entity.WalletDTO) *WalletPresenter {
	if dto == nil {
		return nil
	}

	return &WalletPresenter{
		ID:        dto.ID,
		UserID:    dto.UserID,
		Name:      dto.Name,
		Balance:   dto.Balance,
		Currency:  dto.Currency,
		Type:      dto.Type,
		Color:     dto.Color,
		Icon:      dto.Icon,
		IsExclude: dto.IsExclude,
		CreatedAt: dto.CreatedAt,
	}
}

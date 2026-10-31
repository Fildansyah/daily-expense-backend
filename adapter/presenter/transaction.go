package presenter

import (
	"time"

	"expense.com/m/entity"
)

type TransactionCategoryPresenter struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Icon       string  `json:"icon"`
	BgColor    string  `json:"bgColor"`
	ParentId   *string `json:"parentId,omitempty"`
	ParentName string  `json:"parentName,omitempty"`
}

type TransactionWalletPresenter struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

type TransactionPresenter struct {
	ID                string                        `json:"id"`
	Type              string                        `json:"type"`
	Amount            float64                       `json:"amount"`
	Memo              string                        `json:"memo"`
	Note              string                        `json:"note"`
	Date              string                        `json:"date"`
	TransactionTypeId string                        `json:"transactionTypeId"`
	Wallet            *TransactionWalletPresenter   `json:"wallet,omitempty"`
	TransferWallet    *TransactionWalletPresenter   `json:"transferWallet,omitempty"`
	Category          *TransactionCategoryPresenter `json:"category,omitempty"`
	CreatedAt         time.Time                     `json:"createdAt"`
	UpdatedAt         time.Time                     `json:"updatedAt"`
}

func (v TransactionPresenter) FromDTO(dto *entity.TransactionDTO) *TransactionPresenter {
	if dto == nil {
		return nil
	}

	return &TransactionPresenter{
		ID:                dto.ID,
		Type:              string(dto.Type),
		Amount:            dto.Amount,
		Memo:              dto.Memo,
		Note:              dto.Note,
		Date:              dto.Date,
		TransactionTypeId: dto.TransactionTypeId,
		Wallet:            walletBrief(dto.Wallet),
		TransferWallet:    walletBrief(dto.TransferWallet),
		Category:          categoryBrief(dto.Category),
		CreatedAt:         dto.CreatedAt,
		UpdatedAt:         dto.UpdatedAt,
	}
}

func walletBrief(dto *entity.WalletBriefDTO) *TransactionWalletPresenter {
	if dto == nil {
		return nil
	}
	return &TransactionWalletPresenter{
		ID:    dto.ID,
		Name:  dto.Name,
		Type:  dto.Type,
		Color: dto.Color,
		Icon:  dto.Icon,
	}
}

func categoryBrief(dto *entity.CategoryBriefDTO) *TransactionCategoryPresenter {
	if dto == nil {
		return nil
	}
	return &TransactionCategoryPresenter{
		ID:         dto.ID,
		Name:       string(dto.Name),
		Icon:       dto.Icon,
		BgColor:    dto.BgColor,
		ParentId:   dto.ParentId,
		ParentName: dto.ParentName,
	}
}

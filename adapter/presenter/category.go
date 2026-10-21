package presenter

import (
	"time"

	"expense.com/m/entity"
	"expense.com/m/model"
)

type CategoryTransactionTypePresenter struct {
	ID   string                    `json:"id"`
	Type model.TransactionTypeEnum `json:"type"`
}

type CategoryPresenter struct {
	ID                string                            `json:"id"`
	CreatedAt         time.Time                         `json:"createdAt"`
	Name              model.CategoryNameEnum            `json:"name"`
	TransactionTypeId string                            `json:"transactionTypeId"`
	TransactionTypes  *CategoryTransactionTypePresenter `json:"transactionTypes,omitempty"`
	Icon              string                            `json:"icon"`
	BgColor           string                            `json:"bgColor"`
}

func (v CategoryPresenter) FromDTO(dto *entity.CategoryDTO) *CategoryPresenter {
	if dto == nil {
		return nil
	}

	var transactionTypes *CategoryTransactionTypePresenter
	if dto.TransactionTypes.ID != "" {
		transactionTypes = &CategoryTransactionTypePresenter{
			ID:   dto.TransactionTypes.ID,
			Type: dto.TransactionTypes.Type,
		}
	}

	return &CategoryPresenter{
		ID:                dto.ID,
		CreatedAt:         dto.CreatedAt,
		Name:              dto.Name,
		TransactionTypeId: dto.TransactionTypeId,
		TransactionTypes:  transactionTypes,
		Icon:              dto.Icon,
		BgColor:           dto.BgColor,
	}
}

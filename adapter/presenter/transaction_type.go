package presenter

import (
	"time"

	"expense.com/m/entity"
	"expense.com/m/model"
)

type TransactionTypePresenter struct {
	ID        string                    `json:"id"`
	CreatedAt time.Time                 `json:"createdAt"`
	Type      model.TransactionTypeEnum `json:"type"`
}

func (v TransactionTypePresenter) FromDTO(dto *entity.TransactionTypeDTO) *TransactionTypePresenter {
	if dto == nil {
		return nil
	}

	return &TransactionTypePresenter{
		ID:        dto.ID,
		CreatedAt: dto.CreatedAt,
		Type:      dto.Type,
	}
}

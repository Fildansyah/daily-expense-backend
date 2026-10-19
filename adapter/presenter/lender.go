package presenter

import "expense.com/m/entity"

type LenderPresenter struct {
	ID     string `json:"id"`
	Lender string `json:"lender"`
}

func (v LenderPresenter) FromDTO(dto *entity.LenderDTO) *LenderPresenter {
	if dto == nil {
		return nil
	}

	return &LenderPresenter{
		ID:     dto.ID,
		Lender: dto.Name,
	}
}

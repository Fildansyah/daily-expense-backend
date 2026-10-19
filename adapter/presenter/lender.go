package presenter

import "expense.com/m/entity"

type LenderPresenter struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (v LenderPresenter) FromDTO(dto *entity.LenderDTO) *LenderPresenter {
	if dto == nil {
		return nil
	}

	return &LenderPresenter{
		ID:   dto.ID,
		Name: dto.Name,
	}
}

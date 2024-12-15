package presenter

import (
	"expense.com/m/entity"
	"expense.com/m/model/concern"
)

type UserPresenter struct {
	concern.BaseFields
	FirstName string
	LastName  string
	Email     string
}

func (v UserPresenter) FromDTO(dto *entity.UserDTO) *UserPresenter {
	if dto == nil {
		return nil
	}

	return &UserPresenter{
		BaseFields: concern.BaseFields{
			ID:        dto.ID,
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		},
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
	}
}

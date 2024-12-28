package entity

import (
	"expense.com/m/model"
	"expense.com/m/model/concern"
	"gorm.io/gorm"
)

type UserDTO struct {
	CommonDTOFields
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (d *UserDTO) FromModel(m *model.Users) *UserDTO {
	if m == nil {
		return nil
	}

	return &UserDTO{
		CommonDTOFields: CommonDTOFields{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Email:     m.Email,
		Password:  m.Password,
	}
}

func (d *UserDTO) ToModel() *model.Users {
	m := &model.Users{
		BaseFields: concern.BaseFields{
			ID:        d.ID,
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		},
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Email:     d.Email,
		Password:  d.Password,
	}

	if d.DeletedAt != nil {
		m.DeletedAt = &gorm.DeletedAt{Time: *d.DeletedAt, Valid: true}
	}

	return m
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

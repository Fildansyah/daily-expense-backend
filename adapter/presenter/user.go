package presenter

import (
	"time"

	"expense.com/m/entity"
)

// UserPresenter hanya mengekspos field yang aman untuk client
// (tanpa password).
type UserPresenter struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
}

func (v UserPresenter) FromDTO(dto *entity.UserDTO) *UserPresenter {
	if dto == nil {
		return nil
	}

	return &UserPresenter{
		ID:        dto.ID,
		CreatedAt: dto.CreatedAt,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
	}
}

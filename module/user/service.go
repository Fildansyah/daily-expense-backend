package user

import (
	"expense.com/m/entity"
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) FindUserByEmail(email string) (*entity.UserDTO, error) {
	return s.repo.FindUserByEmail(email)
}

func (s *UserService) GetCurrentUser(id string) (*entity.UserDTO, error) {
	return s.repo.FindUserByID(id)
}

func (s *UserService) CreateUser(input *entity.RegisterRequest) (*entity.UserDTO, error) {
	input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}

	return s.repo.CreateUser(&entity.UserDTO{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password,
	})
}

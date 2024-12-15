package user

import "expense.com/m/entity"

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) FindUserByEmail(email string) (*entity.UserDTO, error) {
	return s.repo.FindUserByEmail(email)
}

func (s *UserService) CreateUser(input *entity.UserDTO) (*entity.UserDTO, error) {
	return s.repo.CreateUser(input)
}

func (s *UserService) LoginUser(input *entity.SignInRequest) (*entity.SignInResponse, error) {
	return s.repo.LoginUser(input)
}

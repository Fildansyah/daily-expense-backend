package user

import (
	"fmt"

	"expense.com/m/entity"
	"expense.com/m/middlewares"
	"expense.com/m/model"
	"expense.com/m/module/base_repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	*base_repository.BaseRepository[model.Users]
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: base_repository.NewBaseRepository[model.Users](db),
		db:             db,
	}
}

func (r *UserRepository) FindUserByEmail(email string) (*entity.UserDTO, error) {
	var user model.Users
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	dto := new(entity.UserDTO)
	return dto.FromModel(&user), nil
}

func (r *UserRepository) CreateUser(input *entity.UserDTO) (*entity.UserDTO, error) {

	existingEmail, _ := r.FindUserByEmail(input.Email)

	if existingEmail != nil {
		return nil, fmt.Errorf("user with email %s already exists", input.Email)
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	modelEntity := input.ToModel()
	modelEntity.ID = uuid.NewString()
	modelEntity.Password = string(hashedPass)

	if err := r.Create(modelEntity); err != nil {
		return nil, err
	}
	return input.FromModel(modelEntity), nil
}

func (r *UserRepository) LoginUser(input *entity.SignInRequest) (*entity.SignInResponse, error) {

	user, err := r.FindUserByEmail(input.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user with email %s not found", input.Email)
	}

	if user.Password == "" {
		return nil, fmt.Errorf("user password is empty")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, err
	}

	token, err := middlewares.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %v", err)
	}

	return &entity.SignInResponse{Token: token}, nil
}

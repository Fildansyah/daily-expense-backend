package user

import (
	"errors"

	"expense.com/m/entity"
	"expense.com/m/model"
	baserepository "expense.com/m/shared/base_repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	*baserepository.BaseRepository[model.Users]
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db:             db,
		BaseRepository: baserepository.NewBaseRepository[model.Users](db),
	}
}

// FindUserByEmail mengembalikan (nil, nil) jika user tidak ditemukan agar
// pemanggil bisa membedakan "tidak ada" dari kegagalan database.
func (r *UserRepository) FindUserByEmail(email string) (*entity.UserDTO, error) {
	var user model.Users
	err := r.db.Where("LOWER(email) = ?", entity.NormalizeEmail(email)).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	dto := new(entity.UserDTO)
	return dto.FromModel(&user), nil
}

func (r *UserRepository) FindUserByID(id string) (*entity.UserDTO, error) {
	var user model.Users
	err := r.db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	dto := new(entity.UserDTO)
	return dto.FromModel(&user), nil
}

func (r *UserRepository) CreateUser(input *entity.UserDTO) (*entity.UserDTO, error) {
	input.Email = entity.NormalizeEmail(input.Email)

	existing, err := r.FindUserByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, entity.ErrEmailExists
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	modelEntity := input.ToModel()
	modelEntity.ID = uuid.NewString()
	modelEntity.Password = string(hashedPass)

	if err := r.Create(modelEntity); err != nil {
		// Unique constraint bisa menang balapan dengan request lain.
		if current, findErr := r.FindUserByEmail(input.Email); findErr == nil && current != nil {
			return nil, entity.ErrEmailExists
		}
		return nil, err
	}

	return input.FromModel(modelEntity), nil
}

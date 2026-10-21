package entity

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"expense.com/m/model"
	"expense.com/m/model/concern"
	"gorm.io/gorm"
)

var (
	ErrEmailExists  = errors.New("email sudah terdaftar")
	ErrInvalidInput = errors.New("input tidak valid")
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

type RegisterRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (r *RegisterRequest) Normalize() {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Email = NormalizeEmail(r.Email)
}

func (r *RegisterRequest) Validate() error {
	if r.FirstName == "" || r.LastName == "" {
		return fmt.Errorf("%w: nama depan dan belakang wajib diisi", ErrInvalidInput)
	}

	addr, err := mail.ParseAddress(r.Email)
	if err != nil || addr.Address != r.Email {
		return fmt.Errorf("%w: format email tidak valid", ErrInvalidInput)
	}

	if len([]rune(r.Password)) < 8 {
		return fmt.Errorf("%w: kata sandi minimal 8 karakter", ErrInvalidInput)
	}

	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

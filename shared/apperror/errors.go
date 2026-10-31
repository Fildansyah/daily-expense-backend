package apperror

import "errors"

var (
	ErrValidation = errors.New("validasi gagal")
	ErrNotFound   = errors.New("data tidak ditemukan")
	ErrConflict   = errors.New("konflik data")
)

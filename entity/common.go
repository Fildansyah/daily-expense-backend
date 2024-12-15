package entity

import "time"

type CommonDTOFields struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

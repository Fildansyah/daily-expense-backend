package entity

import "time"

type CommonDTOFields struct {
	ID        string
	CreatedAt time.Time
	DeletedAt *time.Time
}

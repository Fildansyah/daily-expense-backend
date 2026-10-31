package entity

import "time"

type CommonDTOFields struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

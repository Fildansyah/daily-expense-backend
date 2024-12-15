package concern

import (
	"time"

	"gorm.io/gorm"
)

type BaseFields struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *gorm.DeletedAt
}

package model

import "time"

type RefreshToken struct {
	ID         string    `gorm:"primaryKey;type:uuid"`
	UserID     string    `gorm:"type:uuid;not null;index"`
	FamilyID   string    `gorm:"type:uuid;not null;index"`
	TokenHash  string    `gorm:"type:char(64);not null;uniqueIndex"`
	ExpiresAt  time.Time `gorm:"not null"`
	CreatedAt  time.Time
	RevokedAt  *time.Time
	ReplacedBy *string `gorm:"type:uuid"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	FullName     string    `gorm:"not null"`
	Role         string    `gorm:"default:user"`
	IsActive     bool      `gorm:"default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

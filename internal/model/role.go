package model

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Code      string    `gorm:"uniqueIndex;not null"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

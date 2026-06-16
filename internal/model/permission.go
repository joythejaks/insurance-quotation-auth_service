package model

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Code        string    `gorm:"uniqueIndex;not null"`
	Name        string    `gorm:"not null"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

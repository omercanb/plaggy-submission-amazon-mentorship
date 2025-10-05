package database

import (
	"time"

	"gorm.io/gorm"
)

type Instructor struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Name      string `gorm:"size:255;not null"`
	Surname   string `gorm:"size:255;not null"`
	Email     string `gorm:"size:255;not null;uniqueIndex"`
	Password  string `gorm:"size:255;not null"`
}

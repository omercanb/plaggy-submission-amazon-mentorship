package database

import (
	"time"

	"gorm.io/gorm"
)

type Assignment struct {
	ID          uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
	Title       string    `gorm:"size:255;not null"`
	DueDate     time.Time `gorm:"not null"`
	ClassroomID uint      `gorm:"not null;index"`
	Classroom   Classroom `gorm:"foreignKey:ClassroomID"`
}

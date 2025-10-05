package database

import (
	"time"

	"gorm.io/gorm"
)

type Classroom struct {
	ID           uint `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
	Title        string     `gorm:"size:255;not null;uniqueIndex"`
	InstructorID uint       `gorm:"not null;index"`
	Instructor   Instructor `gorm:"foreignKey:InstructorID"`
}

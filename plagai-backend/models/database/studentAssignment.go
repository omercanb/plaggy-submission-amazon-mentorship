package database

import (
	"time"

	"gorm.io/gorm"
)

type StudentAssignment struct {
	ID uint `gorm:"primaryKey"`
	//Creation timestamp also implies submission timestamp.
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
	StudentID    uint       `gorm:"not null;index"`
	Student      Student    `gorm:"foreignKey:StudentID"`
	AssignmentID uint       `gorm:"not null;index"`
	Assignment   Assignment `gorm:"foreignKey:AssignmentID"`
}

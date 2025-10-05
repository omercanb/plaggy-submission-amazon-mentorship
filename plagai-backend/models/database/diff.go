package database

import (
	"time"

	"gorm.io/gorm"
)

type Diff struct {
	ID                  uint `gorm:"primaryKey"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt
	StudentAssignmentID uint              `gorm:"not null;index"`
	StudentAssignment   StudentAssignment `gorm:"foreignKey:StudentAssignmentID"`
	FilePath            string            `gorm:"not null"`
	DiffData            string            `gorm:"not null"`
}

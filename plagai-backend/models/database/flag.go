package database

import (
	"time"

	"gorm.io/gorm"
)

type Flag struct {
	ID                  uint `gorm:"primaryKey"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt
	Text                string            `gorm:"not null"`
	DiffID              uint              `gorm:"not null;index"`
	Diff                Diff              `gorm:"foreignKey:DiffID"`
	Severity            uint              `gorm:"not null"`
	StudentAssignmentID uint              `gorm:"not null;index"`
	StudentAssignment   StudentAssignment `gorm:"foreignKey:StudentAssignmentID"`
}

package repository

import (
	"errors"
	"log"

	"github.com/plagai/plagai-backend/models/database"
	"github.com/plagai/plagai-backend/models/domain"
	"gorm.io/gorm"
)

var ErrAssignmentNotFound = errors.New("assignment does not belong to classroom")

type AssignmentRepo interface {
	GetAssignmentsForClassroomID(classroomID uint) ([]domain.Assignment, error)
}

type assignmentRepo struct {
	db *gorm.DB
}

func NewAssignmentRepo(db *gorm.DB) AssignmentRepo {
	return &assignmentRepo{db: db}
}

func (repo *assignmentRepo) GetAssignmentsForClassroomID(classroomID uint) ([]domain.Assignment, error) {
	var assignmentsFromDB []database.Assignment
	if err := repo.db.
		Where("classroom_id = ?", classroomID).
		Find(&assignmentsFromDB).Error; err != nil {
		log.Println("assignment query failed:", err)
		return []domain.Assignment{}, ErrAssignmentNotFound
	}
	assignments := make([]domain.Assignment, len(assignmentsFromDB))
	for i, dbAssignment := range assignmentsFromDB {
		assignments[i] = domain.Assignment{
			ID:          dbAssignment.ID,
			Title:       dbAssignment.Title,
			DueDate:     dbAssignment.DueDate,
			ClassroomID: dbAssignment.ClassroomID,
		}
	}
	return assignments, nil
}

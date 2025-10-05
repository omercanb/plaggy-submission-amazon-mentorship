package repository

import (
	"errors"
	"fmt"

	"github.com/plagai/plagai-backend/models/database"
	"github.com/plagai/plagai-backend/models/domain"
	"gorm.io/gorm"
)

var ErrStudentAssignmentNotFound = errors.New("error querying for student assignment")

type StudentAssignmentRepo interface {
	GetStudentAssignments(studentID uint) []domain.StudentAssignment
	NewStudentAssignment(studentID uint, assignmentID uint) (domain.StudentAssignment, error)
}

type studentAssignmentRepo struct {
	db *gorm.DB
}

func NewStudentAssignmentRepo(db *gorm.DB) StudentAssignmentRepo {
	return &studentAssignmentRepo{db: db}
}

func (repo *studentAssignmentRepo) NewStudentAssignment(studentID uint, assignmentID uint) (domain.StudentAssignment, error) {
	studentAssignment := database.StudentAssignment{
		StudentID:    studentID,
		AssignmentID: assignmentID,
	}

	res := repo.db.Create(&studentAssignment)
	if res.Error != nil {
		return domain.StudentAssignment{}, res.Error
	}
	if res.RowsAffected == 0 {
		return domain.StudentAssignment{}, fmt.Errorf("insert failed")
	}

	return domain.StudentAssignment{
		ID:           studentAssignment.ID,
		StudentID:    studentAssignment.StudentID,
		AssignmentID: studentAssignment.AssignmentID,
	}, nil
}

func (repo *studentAssignmentRepo) GetStudentAssignments(studentID uint) []domain.StudentAssignment {
	dbAssignments := []database.StudentAssignment{}
	res := repo.db.Where("student_id = ?", studentID).Find(&dbAssignments)
	if res.Error != nil {
		return []domain.StudentAssignment{}
	}
	assignments := make([]domain.StudentAssignment, len(dbAssignments))
	for i, v := range dbAssignments {
		assignments[i] = domain.StudentAssignment{
			ID:             v.ID,
			SubmissionTime: v.CreatedAt,
			StudentID:      v.StudentID,
			AssignmentID:   v.AssignmentID,
		}
	}
	return assignments
}

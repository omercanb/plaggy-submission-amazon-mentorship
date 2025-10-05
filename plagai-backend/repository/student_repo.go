package repository

import (
	"errors"

	"github.com/plagai/plagai-backend/models/database"
	"github.com/plagai/plagai-backend/models/domain"
	"gorm.io/gorm"
)

var ErrStudentNotFound = errors.New("student not found")

type StudentRepository interface {
	FindByEmail(email string) (domain.Student, error)
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (repo *studentRepository) FindByEmail(email string) (domain.Student, error) {
	student := database.Student{}
	res := repo.db.First(&student, "email = ?", email)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// No matching record
		return domain.Student{}, ErrStudentNotFound
	}
	if res.Error != nil {
		// DB-level error (connection, query, etc.)
		return domain.Student{}, res.Error
	}

	return domain.Student{
		ID:          student.ID,
		Name:        student.Name,
		Email:       student.Email,
		ClassroomID: student.ClassroomID,
	}, nil
}

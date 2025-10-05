package repository

import (
	"errors"
	"log"

	"github.com/plagai/plagai-backend/models/database"
	"github.com/plagai/plagai-backend/models/domain"
	"gorm.io/gorm"
)

// A better error message would be nice
var ErrClassroomNotFound = errors.New("classroom not found")

type ClassroomRepo interface {
	GetClassroomByID(id uint) (domain.Classroom, error)
}

type classroomRepo struct {
	db *gorm.DB
}

func NewClassroomRepo(db *gorm.DB) ClassroomRepo {
	return &classroomRepo{db: db}
}

func (repo *classroomRepo) GetClassroomByID(id uint) (domain.Classroom, error) {
	var classroomFromDB database.Classroom
	if dbErr := repo.db.Where("id = ?", id).First(&classroomFromDB).Error; dbErr != nil {
		log.Println("classroom query failed:", dbErr)
		return domain.Classroom{}, ErrClassroomNotFound
	}
	classroom := domain.Classroom{
		ID:           classroomFromDB.ID,
		Title:        classroomFromDB.Title,
		InstructorID: classroomFromDB.InstructorID,
	}
	return classroom, nil
}

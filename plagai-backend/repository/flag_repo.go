package repository

import (
	"errors"
	"time"

	"github.com/plagai/plagai-backend/models/database"
	"github.com/plagai/plagai-backend/models/domain"
	"gorm.io/gorm"
)

var (
	ErrUnauthorized = errors.New("unauthorized access to classroom")
	ErrDatabase     = errors.New("database error")
)

type FlagRepository interface {
	AddFlag(flag *domain.Flag, studentAssignmentID uint) error
	FindByID(id uint) (*domain.Flag, error)
	FindByStudentAssignmentID(studentAssignmentID uint) ([]domain.Flag, error)
	Delete(id uint) error
}

type flagRepository struct {
	db *gorm.DB
}

func NewFlagRepository(db *gorm.DB) FlagRepository {
	return &flagRepository{db: db}
}

func (r *flagRepository) AddFlag(flag *domain.Flag, studentAssignmentID uint) error {
	dbFlag := toDBFlag(flag, studentAssignmentID)
	return r.db.Create(&dbFlag).Error
}

func (r *flagRepository) FindByID(id uint) (*domain.Flag, error) {
	var dbFlag database.Flag
	if err := r.db.First(&dbFlag, id).Error; err != nil {
		return nil, err
	}
	return toDomainFlag(&dbFlag), nil
}

func (r *flagRepository) FindByStudentAssignmentID(studentAssignmentID uint) ([]domain.Flag, error) {
	var dbFlags []database.Flag
	if err := r.db.Where("student_assignment_id = ?", studentAssignmentID).Find(&dbFlags).Error; err != nil {
		return nil, err
	}
	flags := make([]domain.Flag, len(dbFlags))
	for i, dbFlag := range dbFlags {
		flags[i] = *toDomainFlag(&dbFlag)
	}
	return flags, nil
}

func (r *flagRepository) Delete(id uint) error {
	return r.db.Delete(&database.Flag{}, id).Error
}

/*
func (r *flagRepository) GetFlagsForClassroom(section uint, homeworkNum uint, instructorEmail string, limit int) ([]domain.Flag, error) {
	var classroom database.Classroom
	if err := r.db.First(&classroom, section).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAssignmentNotFound
		}
		return nil, ErrDatabase
	}

	var assignment database.Assignment
	if err := r.db.Where("id = ? AND classroom_id = ?", homeworkNum, classroom.ID).
		First(&assignment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInstructorNotFound
		}
		return nil, ErrDatabase
	}

	var inst database.Instructor
	if err := r.db.Where("email = ?", instructorEmail).First(&inst).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInstructorNotFound
		}
		return nil, ErrDatabase
	}

	if classroom.InstructorID != inst.ID {
		return nil, ErrUnauthorized
	}

	type row struct {
		DiffID        uint
		DiffData      string
		DiffCreatedAt time.Time
		FilePath      string
		StudentEmail  string
		AssignmentID  uint
	}

	var rows []row
	q := r.db.Table("diffs").
		Select(`
			diffs.id         AS diff_id,
			diffs.diff_data  AS diff_data,
			diffs.created_at AS diff_created_at,
			diffs.file_path  AS file_path,
			students.email   AS student_email,
			assignments.id   AS assignment_id
		`).
		Joins(`JOIN student_assignments sa ON sa.id = diffs.student_assignment_id`).
		Joins(`JOIN students            ON students.id = sa.student_id`).
		Joins(`JOIN assignments         ON assignments.id = sa.assignment_id`).
		Where(`assignments.id = ? AND assignments.classroom_id = ?`, assignment.ID, classroom.ID).
		Order(`diffs.created_at ASC`).
		Limit(limit)

	if err := q.Find(&rows).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrDatabase
	}

	flags := make([]domain.Flag, 0, len(rows))
	for _, row := range rows {
		flags = append(flags, domain.Flag{
			ID:              row.DiffID,
			PatchText:       row.DiffData,
			FlagExplanation: "Explanations not implemented yet",
			Severity:        0,
			CreatedAt:       row.DiffCreatedAt,
			FilePath:        row.FilePath,
			StudentEmail:    row.StudentEmail,
			AssignmentID:    strconv.Itoa(int(row.AssignmentID)),
		})
	}

	return flags, nil
}
*/

func toDBFlag(flag *domain.Flag, studentAssignmentID uint) database.Flag {
	return database.Flag{
		ID:   0, // GORM will auto-generate
		Text: flag.FlagExplanation,
		Diff: database.Diff{
			StudentAssignmentID: studentAssignmentID,
			FilePath:            flag.Diff.FilePath,
			DiffData:            flag.Diff.PatchText,
		},
		Severity:            uint(flag.Severity),
		StudentAssignmentID: studentAssignmentID,
		CreatedAt:           time.Now(),
	}
}

func toDomainFlag(dbFlag *database.Flag) *domain.Flag {
	return &domain.Flag{
		ID: dbFlag.ID,
		Diff: domain.Diff{
			FilePath:  dbFlag.Diff.FilePath,
			PatchText: dbFlag.Diff.DiffData,
			Timestamp: dbFlag.Diff.CreatedAt,
		},
		FlagExplanation: dbFlag.Text,
		Severity:        int(dbFlag.Severity),
	}
}

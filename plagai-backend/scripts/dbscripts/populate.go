package dbscripts

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/plagai/plagai-backend/models/database"
	"github.com/plagai/plagai-backend/security"
	"gorm.io/gorm"
)

type rawSubmission struct {
	PatchText string `json:"patch_text"`
	Timestamp int64  `json:"timestamp"` // Unix ms
	FilePath  string `json:"file_path"` // optional; may be absent
}

func HashPasswords(db *gorm.DB) error {
		targets := []string{"test@example.com", "brt@crt.com"}

	var rows []database.Instructor
	if err := db.
		Where("email IN ? AND password IS NOT NULL AND password <> ''", targets).
		Find(&rows).Error; err != nil {
		return err
	}

	for _, inst := range rows {
		if strings.HasPrefix(inst.Password, "$argon2id$") {
			continue // already migrated
		}
		// Client sends SHA-256(password) going forward; we migrate to match that scheme:
		digestHex := security.Sha256Hex(inst.Password)
		phc, err := security.HashFromClientDigestHex(digestHex, security.DefaultParams)
		if err != nil {
			return err
		}
		if err := db.Model(&database.Instructor{}).
			Where("id = ?", inst.ID).
			Update("password", phc).Error; err != nil {
			return err
		}
		log.Printf("[hashpw] migrated %s (id=%d)", inst.Email, inst.ID)
	}
	return nil
}

func Populate(db *gorm.DB) {
	instructors := []database.Instructor{
		{
			Name:     "Test",
			Surname:  "Example",
			Email:    "test@example.com",
			Password: "password123",
		},
		{
			Name:     "brt",
			Surname:  "crt",
			Email:    "brt@crt.com",
			Password: "password",
		},
	}
	db.CreateInBatches(instructors, len(instructors))
	var instFromDB []database.Instructor
	db.Where("").Find(&instFromDB)
	classrooms := []database.Classroom{
		{
			Title:        "CS 101-2",
			InstructorID: instFromDB[0].ID,
		},
		{
			Title:        "CS 101-3",
			InstructorID: instFromDB[0].ID,
		},
		{
			Title:        "CS 473-1",
			InstructorID: instFromDB[1].ID,
		},
	}
	db.CreateInBatches(classrooms, len(classrooms))
	var classroomsFromDB []database.Classroom
	db.Where("").Find(&classroomsFromDB)
	homeworks := []database.Assignment{
		{
			Title:       "Algo I",
			DueDate:     time.Now().Add(7 * 24 * time.Hour),
			ClassroomID: classroomsFromDB[2].ID,
		},
		{
			Title:       "Algo II",
			DueDate:     time.Now().Add(7 * 24 * time.Hour),
			ClassroomID: classroomsFromDB[2].ID,
		},
		{
			Title:       "Arrays",
			DueDate:     time.Now().Add(14 * 24 * time.Hour),
			ClassroomID: classroomsFromDB[1].ID,
		},
		{
			Title:       "Arrays",
			DueDate:     time.Now().Add(14 * 24 * time.Hour),
			ClassroomID: classroomsFromDB[0].ID,
		},
	}
	db.CreateInBatches(homeworks, len(homeworks))

	studentsToCreate := []database.Student{
		{
			Name:        "crt",
			Surname:     "brt",
			Email:       "crt@brt.com",
			ClassroomID: classroomsFromDB[2].ID,
		},
		{
			Name:        "crt",
			Surname:     "crt",
			Email:       "crt@crt.com",
			ClassroomID: classroomsFromDB[0].ID,
		},
		{
			Name:        "brt",
			Surname:     "brt",
			Email:       "brt@brt.com",
			ClassroomID: classroomsFromDB[1].ID,
		},
	}
	db.CreateInBatches(studentsToCreate, 2)

	fileData, err := os.ReadFile("./scripts/mockdata/filesubmission.json")
	if err != nil {
		log.Println(`{"status":"ERROR","message":"failed to read data file"}`)
		return
	}

	var raws []rawSubmission
	if err := json.Unmarshal(fileData, &raws); err != nil {
		log.Printf(`{"status":"ERROR","message":"failed to parse submissions json: %v"}`, err)
		return
	}
	if len(raws) == 0 {
		log.Println(`{"status":"WARN","message":"no submissions in json; skipping diff seeding"}`)
		return
	}

	var assignments []database.Assignment
	if err := db.Find(&assignments).Error; err != nil {
		log.Printf(`{"status":"ERROR","message":"failed to load assignments: %v"}`, err)
		return
	}
	if len(assignments) == 0 {
		log.Println(`{"status":"WARN","message":"no assignments found; skipping studentassignment/diff seeding"}`)
		return
	}

	var students []database.Student
	if err := db.Find(&students).Error; err != nil {
		log.Printf(`{"status":"ERROR","message":"failed to load students: %v"}`, err)
		return
	}
	if len(students) == 0 {
		log.Println(`{"status":"WARN","message":"no students found; skipping studentassignment/diff seeding"}`)
		return
	}

	studentsByClass := make(map[uint][]database.Student)
	for _, s := range students {
		studentsByClass[s.ClassroomID] = append(studentsByClass[s.ClassroomID], s)
	}

	var sasToCreate []database.StudentAssignment
	for _, asg := range assignments {
		ss := studentsByClass[asg.ClassroomID]
		for _, stu := range ss {
			sasToCreate = append(sasToCreate, database.StudentAssignment{
				StudentID:    stu.ID,
				AssignmentID: asg.ID,
			})
		}
	}
	if len(sasToCreate) == 0 {
		log.Println(`{"status":"WARN","message":"no student-assignment pairs; skipping diffs"}`)
		return
	}
	if err := db.CreateInBatches(&sasToCreate, 100).Error; err != nil {
		log.Printf(`{"status":"ERROR","message":"failed to create studentassignments: %v"}`, err)
		return
	}

	const diffsPerSA = 5
	var diffsToCreate []database.Diff
	cursor := 0

	for saIdx := range sasToCreate {
		for i := 0; i < diffsPerSA; i++ {
			raw := raws[cursor%len(raws)]
			cursor++

			decoded, decErr := url.QueryUnescape(raw.PatchText)
			if decErr != nil {
				decoded = raw.PatchText
			}

			path := raw.FilePath
			if path == "" {
				path = "src/file_" + time.Now().Format("20060102") + "_" + strconv.Itoa(saIdx) + "_" + strconv.Itoa(i) + ".go"
			}

			createdAt := time.UnixMilli(raw.Timestamp)
			if raw.Timestamp == 0 {
				createdAt = time.Now()
			}

			diffsToCreate = append(diffsToCreate, database.Diff{
				StudentAssignmentID: sasToCreate[saIdx].ID,
				FilePath:            path,
				DiffData:            decoded,
				CreatedAt:           createdAt,
				UpdatedAt:           createdAt,
			})
		}
	}

	if err := db.CreateInBatches(&diffsToCreate, 200).Error; err != nil {
		log.Printf(`{"status":"ERROR","message":"failed to create diffs: %v"}`, err)
		return
	}

	log.Printf(`{"status":"OK","message":"seeded %d studentassignments and %d diffs"}`, len(sasToCreate), len(diffsToCreate))
}

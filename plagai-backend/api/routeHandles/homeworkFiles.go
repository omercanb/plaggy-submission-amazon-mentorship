package routeHandles

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/plagai/plagai-backend/core"
	"github.com/plagai/plagai-backend/middleware"
	"github.com/plagai/plagai-backend/models"
	"github.com/plagai/plagai-backend/models/database"
	"gorm.io/gorm"
)

type fileDTO struct {
	FilePath    string     `json:"filePath"`
	PatchCount  int64      `json:"patchCount"`
	FlagCount   int64      `json:"flagCount"`
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
}

// Send files associated with a homework for a specific student
func (h *Handler) ListStudentFiles(w http.ResponseWriter, r *http.Request) {
	sectionStr := r.URL.Query().Get("section")
	homeworkStr := r.URL.Query().Get("homework")
	studentEmail := r.URL.Query().Get("student")
	flaggedOnly := r.URL.Query().Get("flaggedOnly") == "1"

	sectionID, err := strconv.Atoi(sectionStr)
	if err != nil || sectionID <= 0 || studentEmail == "" {
		http.Error(w, `{"status":"ERROR","message":"invalid params"}`, http.StatusBadRequest)
		return
	}
	homeworkID, err := strconv.Atoi(homeworkStr)
	if err != nil || homeworkID <= 0 {
		http.Error(w, `{"status":"ERROR","message":"invalid 'homework'"}`, http.StatusBadRequest)
		return
	}

	var classroom database.Classroom
	if err := h.DB.First(&classroom, sectionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"status":"ERROR","message":"section not found"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"status":"ERROR","message":"db error loading section"}`, http.StatusInternalServerError)
		return
	}
	var assignment database.Assignment
	if err := h.DB.Where("id = ? AND classroom_id = ?", homeworkID, classroom.ID).
		First(&assignment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"status":"ERROR","message":"homework not in section"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"status":"ERROR","message":"db error loading homework"}`, http.StatusInternalServerError)
		return
	}

	claims := middleware.Claims{}
	core.ConvertToken(r.Header.Get("Authorization"), &claims)
	var inst database.Instructor
	if err := h.DB.Where("email = ?", claims.Email).First(&inst).Error; err != nil || inst.ID == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(models.Response[string]{
			Data: "Error", Status: "Unauthorized",
			Message: "Only instructors are allowed", Error: "Unauthorized",
		})
		return
	}
	if classroom.InstructorID != inst.ID {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(models.Response[string]{
			Data: "Error", Status: "Unauthorized",
			Message: "You don't have access to this section", Error: "Unauthorized",
		})
		return
	}

	var rows []fileDTO
	q := h.DB.
		Table("student_assignments sa").
		Select(`
			diffs.file_path AS file_path,
			COUNT(diffs.id) AS patch_count,
			COUNT(flags.id) AS flag_count,
			MAX(diffs.created_at) AS last_updated
		`).
		Joins(`JOIN students s ON s.id = sa.student_id`).
		Joins(`JOIN diffs ON diffs.student_assignment_id = sa.id`).
		Joins(`LEFT JOIN flags ON flags.diff_id = diffs.id AND flags.student_assignment_id = sa.id`).
		Where(`sa.assignment_id = ? AND s.email = ? AND s.classroom_id = ?`, assignment.ID, studentEmail, classroom.ID).
		Group(`diffs.file_path`).
		Order(`last_updated DESC NULLS LAST`)

	if flaggedOnly {
		q = q.Having("COUNT(flags.id) > 0")
	}

	if err := q.Scan(&rows).Error; err != nil && err != gorm.ErrRecordNotFound {
		http.Error(w, `{"status":"ERROR","message":"db error loading files"}`, http.StatusInternalServerError)
		return
	}

	resp := models.Response[[]fileDTO]{Data: rows, Status: "OK"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

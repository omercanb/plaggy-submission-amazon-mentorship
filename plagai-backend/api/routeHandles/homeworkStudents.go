package routeHandles

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/plagai/plagai-backend/core"
	"github.com/plagai/plagai-backend/middleware"
	"github.com/plagai/plagai-backend/models"
	"github.com/plagai/plagai-backend/models/database"
	"gorm.io/gorm"
)

type studentDTO struct {
	ID         uint   `json:"id"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	PatchCount int64  `json:"patchCount"`
	FlagCount  int64  `json:"flagCount"`
}

// Send students in a specific section 
func (h *Handler) ListHomeworkStudents(w http.ResponseWriter, r *http.Request) {
	sectionStr := r.URL.Query().Get("section")
	homeworkStr := r.URL.Query().Get("homework")
	withActivity := r.URL.Query().Get("withActivity") == "1"

	sectionID, err := strconv.Atoi(sectionStr)
	if err != nil || sectionID <= 0 {
		http.Error(w, `{"status":"ERROR","message":"invalid 'section'"}`, http.StatusBadRequest)
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

	var rows []studentDTO
	if err := h.DB.
		Table("students AS s").
		Select(`
			s.id, s.email, s.name, s.surname,
			(
				SELECT COUNT(*)
				FROM student_assignments sa
				JOIN diffs d ON d.student_assignment_id = sa.id
				WHERE sa.student_id = s.id AND sa.assignment_id = ?
			) AS patch_count,
			(
				SELECT COUNT(*)
				FROM student_assignments sa2
				JOIN flags f ON f.student_assignment_id = sa2.id
				WHERE sa2.student_id = s.id AND sa2.assignment_id = ?
			) AS flag_count
		`, assignment.ID, assignment.ID).
		Where("s.classroom_id = ?", classroom.ID).
		Order("s.surname ASC, s.name ASC").
		Scan(&rows).Error; err != nil {
		http.Error(w, `{"status":"ERROR","message":"db error loading students"}`, http.StatusInternalServerError)
		return
	}

	if withActivity {
		filtered := rows[:0]
		for _, r := range rows {
			if r.PatchCount > 0 || r.FlagCount > 0 {
				filtered = append(filtered, r)
			}
		}
		rows = filtered
	}

	resp := models.Response[[]studentDTO]{Data: rows, Status: "OK"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

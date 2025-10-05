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

func (h *Handler) SendDetections(w http.ResponseWriter, r *http.Request) {
	// Section and homework based retrieval of detections.
	sectionStr := r.URL.Query().Get("section")
	if sectionStr == "" {
		http.Error(w, `{"status":"ERROR","message":"missing 'section' query param"}`, http.StatusBadRequest)
		return
	}

	sectionNum, err := strconv.Atoi(sectionStr)
	if err != nil || sectionNum < 0 {
		http.Error(w, `{"status":"ERROR","message":"'section' must be a non-negative integer"}`, http.StatusBadRequest)
		return
	}

	homeworkStr := r.URL.Query().Get("homework")
	if homeworkStr == "" {
		http.Error(w, `{"status":"ERROR","message":"missing 'homework' query param"}`, http.StatusBadRequest)
		return
	}

	homeworkNum, err := strconv.Atoi(homeworkStr)
	if err != nil || homeworkNum < 0 {
		http.Error(w, `{"status":"ERROR","message":"'homework' must be a non-negative integer"}`, http.StatusBadRequest)
		return
	}

	// Change this limit once most (I fucking hope so) flags have some sort of
	// patch data in its associated diff. ~brtcrt
	limit := 200
	if limStr := r.URL.Query().Get("limit"); limStr != "" {
		if v, err := strconv.Atoi(limStr); err == nil && v > 0 { // In fact, since we are using offset based pagination
			limit = v // this whole thing can just be hardcoded but we can keep
		} // the url option for now ~brtcrt
	}
	if limit > 500 {
		limit = 500
	}
	page := 1
	if pStr := r.URL.Query().Get("page"); pStr != "" {
		if v, err := strconv.Atoi(pStr); err == nil && v > 0 {
			page = v
		}
	}
	offset := (page - 1) * limit

	var classroom database.Classroom
	if err := h.DB.First(&classroom, sectionNum).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"status":"ERROR","message":"section not found"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"status":"ERROR","message":"database error while loading section"}`, http.StatusInternalServerError)
		return
	}

	var assignment database.Assignment
	if err := h.DB.Where("id = ? AND classroom_id = ?", homeworkNum, classroom.ID).
		First(&assignment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"status":"ERROR","message":"homework does not belong to the specified section"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"status":"ERROR","message":"database error while validating homework"}`, http.StatusInternalServerError)
		return
	}

	claims := middleware.Claims{}
	core.ConvertToken(r.Header.Get("Authorization"), &claims)

	var inst database.Instructor
	if err := h.DB.Where("email = ?", claims.Email).First(&inst).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(models.Response[string]{
				Data:    "Error",
				Status:  "Unauthorized",
				Message: "Only instructors are allowed",
				Error:   "Unauthorized",
			})
			return
		}
		http.Error(w, `{"status":"ERROR","message":"database error while checking instructor"}`, http.StatusInternalServerError)
		return
	}

	if classroom.InstructorID != inst.ID {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(models.Response[string]{
			Data:    "Error",
			Status:  "Unauthorized",
			Message: "You don't have access to this section",
			Error:   "Unauthorized",
		})
		return
	}

	whereSQL := `
		assignments.id = ? AND assignments.classroom_id = ? AND
		diffs.diff_data IS NOT NULL AND TRIM(diffs.diff_data) <> ''
	`
	// Above is for filtering only the ones with patch data, below is anything. Ideally we should use the one below change later ~brtcrt
	// whereSQL := `assignments.id = ? AND assignments.classroom_id = ?`
	whereArgs := []any{assignment.ID, classroom.ID}

	var total int64
	if err := h.DB.Model(&database.Flag{}).
		Joins(`JOIN student_assignments sa ON sa.id = flags.student_assignment_id`).
		Joins(`JOIN diffs ON diffs.id = flags.diff_id AND diffs.student_assignment_id = sa.id`).
		Joins(`JOIN assignments ON assignments.id = sa.assignment_id`).
		Where(whereSQL, whereArgs...).
		Distinct("flags.id").
		Count(&total).Error; err != nil {
		http.Error(w, `{"status":"ERROR","message":"database error while counting flags"}`, http.StatusInternalServerError)
		return
	}

	type row struct {
		FlagID        uint
		Text          string
		FlagCreatedAt time.Time
		FlagSeverity  uint
		StudentEmail  string
		AssignmentID  uint
		DiffFilePath  string
		DiffPatchData string
	}

	var rows []row
	q := h.DB.Model(&database.Flag{}).
		Select(`
			flags.id          AS flag_id,
			flags.text        AS text,
			flags.created_at  AS flag_created_at,
			flags.severity    AS flag_severity,
			students.email    AS student_email,
			assignments.id    AS assignment_id,
			diffs.file_path   AS diff_file_path,
			diffs.diff_data   AS diff_patch_data
		`).
		Joins(`JOIN student_assignments sa ON sa.id = flags.student_assignment_id`).
		Joins(`JOIN diffs ON diffs.id = flags.diff_id AND diffs.student_assignment_id = sa.id`).
		Joins(`JOIN students    ON students.id = sa.student_id`).
		Joins(`JOIN assignments ON assignments.id = sa.assignment_id`).
		Where(whereSQL, whereArgs...).
		Order("flags.severity DESC, flags.created_at DESC").
		Limit(limit).
		Offset(offset)

	if err := q.Find(&rows).Error; err != nil && err != gorm.ErrRecordNotFound {
		http.Error(w, `{"status":"ERROR","message":"database error while fetching flags"}`, http.StatusInternalServerError)
		return
	}

	dets := make([]models.Detection, 0, len(rows))
	for _, r := range rows {
		dets = append(dets, models.Detection{
			ID:         strconv.Itoa(int(r.FlagID)),
			CreatedBy:  r.StudentEmail,
			Content:    r.Text,
			CreatedAt:  r.FlagCreatedAt,
			HomeworkID: strconv.Itoa(int(r.AssignmentID)),
			Severity:   r.FlagSeverity,
			FilePath:   r.DiffFilePath,
			DiffData:   r.DiffPatchData,
		})
	}

	hasMore := (int64(page)*int64(limit) < total)
	w.Header().Set("X-Total-Count", strconv.FormatInt(total, 10))
	w.Header().Set("X-Page", strconv.Itoa(page))
	w.Header().Set("X-Limit", strconv.Itoa(limit))
	if hasMore {
		w.Header().Set("X-Has-More", "true")
		w.Header().Set("X-Next-Page", strconv.Itoa(page+1))
	} else {
		w.Header().Set("X-Has-More", "false")
	}

	resp := models.Response[[]models.Detection]{
		Data:   dets,
		Status: "OK",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

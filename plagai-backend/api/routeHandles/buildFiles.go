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
	"github.com/plagai/plagai-backend/models/domain"
	"github.com/plagai/plagai-backend/service"
	"gorm.io/gorm"
)

type patchItem struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	PatchText string    `json:"patchText"`
}

// I can't be bothered to make this into its of seperate model fuck me ~brtcrt
type fileDiffsPayload struct {
	FilePath  string      `json:"filePath"`
	Student   string      `json:"student"`
	FinalText string      `json:"finalText"`
	Patches   []patchItem `json:"patches"`
}

func (h *Handler) BuildFile(w http.ResponseWriter, r *http.Request) {
	sectionStr := r.URL.Query().Get("section")
	homeworkStr := r.URL.Query().Get("homework")
	studentEmail := r.URL.Query().Get("student")
	filePath := r.URL.Query().Get("file")
	if sectionStr == "" || homeworkStr == "" || studentEmail == "" || filePath == "" {
		http.Error(w, `{"status":"ERROR","message":"missing required query params: section, homework, student, file"}`, http.StatusBadRequest)
		return
	}
	sectionID, err := strconv.Atoi(sectionStr)
	if err != nil || sectionID < 0 {
		http.Error(w, `{"status":"ERROR","message":"'section' must be a non-negative integer"}`, http.StatusBadRequest)
		return
	}
	homeworkID, err := strconv.Atoi(homeworkStr)
	if err != nil || homeworkID < 0 {
		http.Error(w, `{"status":"ERROR","message":"'homework' must be a non-negative integer"}`, http.StatusBadRequest)
		return
	}
	flaggedOnly := r.URL.Query().Get("flaggedOnly") == "1"

	limit := 50
	if limStr := r.URL.Query().Get("limit"); limStr != "" {
		if v, err := strconv.Atoi(limStr); err == nil && v > 0 {
			limit = v
		}
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
	if err := h.DB.First(&classroom, sectionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"status":"ERROR","message":"section not found"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"status":"ERROR","message":"database error while loading section"}`, http.StatusInternalServerError)
		return
	}
	var assignment database.Assignment
	if err := h.DB.Where("id = ? AND classroom_id = ?", homeworkID, classroom.ID).
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
				Data:    "Error", Status: "Unauthorized",
				Message: "Only instructors are allowed", Error: "Unauthorized",
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
			Data:    "Error", Status: "Unauthorized",
			Message: "You don't have access to this section", Error: "Unauthorized",
		})
		return
	}

	var student database.Student
	if err := h.DB.Where("email = ? AND classroom_id = ?", studentEmail, classroom.ID).
		First(&student).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"status":"ERROR","message":"student not found in this section"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"status":"ERROR","message":"database error while loading student"}`, http.StatusInternalServerError)
		return
	}

	var sa database.StudentAssignment
	if err := h.DB.Where("student_id = ? AND assignment_id = ?", student.ID, assignment.ID).
		First(&sa).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, `{"status":"ERROR","message":"no student-assignment record for this student & homework"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"status":"ERROR","message":"database error while loading student assignment"}`, http.StatusInternalServerError)
		return
	}

	whereSQL := "diffs.student_assignment_id = ? AND diffs.file_path = ? AND diffs.diff_data IS NOT NULL AND TRIM(diffs.diff_data) <> ''"
	whereArgs := []any{sa.ID, filePath}

	if flaggedOnly {
		whereSQL += " AND EXISTS (SELECT 1 FROM flags f WHERE f.diff_id = diffs.id AND f.student_assignment_id = diffs.student_assignment_id)"
	}

	var total int64
	if err := h.DB.Table("diffs").
		Where(whereSQL, whereArgs...).
		Count(&total).Error; err != nil {
		http.Error(w, `{"status":"ERROR","message":"database error while counting patches"}`, http.StatusInternalServerError)
		return
	}

	type row struct {
		ID        uint
		CreatedAt time.Time
		PatchText string
	}
	var pageRows []row
	if err := h.DB.Table("diffs").
		Select("diffs.id AS id, diffs.created_at AS created_at, diffs.diff_data AS patch_text").
		Where(whereSQL, whereArgs...).
		Order("diffs.created_at DESC, diffs.id DESC").
		Limit(limit).
		Offset(offset).
		Find(&pageRows).Error; err != nil && err != gorm.ErrRecordNotFound {
		http.Error(w, `{"status":"ERROR","message":"database error while fetching patches"}`, http.StatusInternalServerError)
		return
	}

	patches := make([]patchItem, 0, len(pageRows))
	for _, r := range pageRows {
		patches = append(patches, patchItem{
			ID:        r.ID,
			CreatedAt: r.CreatedAt,
			PatchText: r.PatchText,
		})
	}

	type fullRow struct {
		PatchText string
	}
	var allRows []fullRow
	if err := h.DB.Table("diffs").
		Select("diffs.diff_data AS patch_text").
		Where(whereSQL, whereArgs...).
		Order("diffs.created_at ASC, diffs.id ASC").
		Find(&allRows).Error; err != nil && err != gorm.ErrRecordNotFound {
		http.Error(w, `{"status":"ERROR","message":"database error while loading full patch history"}`, http.StatusInternalServerError)
		return
	}
	domainPatches := make([]domain.Diff, 0, len(allRows))
	for _, r := range allRows {
		domainPatches = append(domainPatches, domain.Diff{
			FilePath:  filePath,
			PatchText: r.PatchText,
		})
	}
	finalText, buildErr := service.BuildFileFromPatchesAndStartText("", domainPatches)
	if buildErr != nil {
		finalText = ""
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

	payload := fileDiffsPayload{
		FilePath:  filePath,
		Student:   studentEmail,
		FinalText: finalText,
		Patches:   patches,
	}
	resp := models.Response[fileDiffsPayload]{Data: payload, Status: "OK"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

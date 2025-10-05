package routeHandles

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/plagai/plagai-backend/core"
	"github.com/plagai/plagai-backend/middleware"
	"github.com/plagai/plagai-backend/models"
	"github.com/plagai/plagai-backend/models/database"
	"gorm.io/gorm"
)

// This will be modified later when db is used.
func (h *Handler) SendHomeworkDetails(w http.ResponseWriter, r *http.Request) {
	homeworkStr := r.URL.Query().Get("homework")
	if homeworkStr == "" {
		http.Error(w, `{"status":"ERROR","message":"missing 'homework' query param"}`, http.StatusBadRequest)
		return
	}

	homeworkNum, err := strconv.Atoi(homeworkStr)
	if err != nil || homeworkNum <= 0 {
		http.Error(w, `{"status":"ERROR","message":"'homework' out of range"}`, http.StatusBadRequest)
		return
	}

	var hwFromDB database.Assignment
	h.DB.Where("id = ?", (uint)(homeworkNum)).First(&hwFromDB)

	response := models.Response[models.Homework]{
		Data: models.Homework{
			ID:         strconv.FormatUint(uint64(hwFromDB.ID), 10),
			Title:      hwFromDB.Title,
			AssignedAt: hwFromDB.CreatedAt,
			DueDate:    hwFromDB.DueDate,
		},
		Status: "OK",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) SendHomeworks(w http.ResponseWriter, r *http.Request) {
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

	fileData, err := os.ReadFile("./scripts/mockdata/filesubmission.json")
	if err != nil {
		http.Error(w, `{"status":"ERROR","message":"failed to read data file"}`, http.StatusInternalServerError)
		return
	}

	var subs []models.PatchWithTimestamp
	if err := json.Unmarshal(fileData, &subs); err != nil {
		http.Error(w, `{"status":"ERROR","message":"failed to parse data file"}`, http.StatusInternalServerError)
		return
	}

	claims := middleware.Claims{}
	core.ConvertToken(r.Header.Get("Authorization"), &claims)

	var userFromDB database.Instructor
	h.DB.Where("Email LIKE ?", claims.Email).First(&userFromDB)

	log.Println("Email from DB ", userFromDB.Email)

	var classroomsFromDB []database.Classroom
	h.DB.Where("instructor_id = ?", userFromDB.ID).Find(&classroomsFromDB)
	log.Println("Classroom from DB ", classroomsFromDB)
	match := false
	var classroom database.Classroom
	for i := 0; i < len(classroomsFromDB); i++ {
		if classroomsFromDB[i].ID == (uint)(sectionNum) {
			match = true
			classroom = classroomsFromDB[i]
			break
		}
	}
	if !match {
		errMsg := models.Response[string]{
			Data:    "Error",
			Message: "Unauthorized. You don't have access to this section",
			Error:   "Unauthorized. You don't have access to this section",
			Status:  "Unauthorized",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	var homeworksFromDB []database.Assignment
	if err := h.DB.
		Where("classroom_id = ?", classroom.ID).
		Find(&homeworksFromDB).Error; err != nil {
		log.Println("assignment query failed:", err)
		http.Error(w, `{"status":"ERROR","message":"failed to load homeworks"}`, http.StatusInternalServerError)
		return
	}
	log.Printf("Loaded %d homeworks for classroom %d\n", len(homeworksFromDB), classroom.ID)
	sendHws := make([]models.Homework, len(homeworksFromDB))
	for i, c := range homeworksFromDB {
		sendHws[i].ID = strconv.FormatUint(uint64(c.ID), 10)
		sendHws[i].Title = c.Title
		sendHws[i].AssignedAt = c.CreatedAt
		sendHws[i].DueDate = c.DueDate
	}
	response := models.Response[[]models.Homework]{
		Data:   sendHws,
		Status: "OK",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

type createHomeworksReq struct {
	SectionIDs []uint `json:"sectionIDs"`
	Title      string `json:"title"`
	DueAt      string `json:"dueAt"`
}

type homeworkDTO struct {
	ID        uint      `json:"id"`
	SectionID uint      `json:"sectionID"`
	Title     string    `json:"title"`
	DueAt     time.Time `json:"dueAt"`
}

func (h *Handler) CreateHomework(w http.ResponseWriter, r *http.Request) {
	var req createHomeworksReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Error:", err.Error())
		http.Error(w, `{"status":"ERROR","message":"invalid JSON body"}`, http.StatusBadRequest)
		return
	}
	log.Println("Req: ", req)
	if len(req.SectionIDs) == 0 {
		http.Error(w, `{"status":"ERROR","message":"sectionIDs must be a non-empty array"}`, http.StatusBadRequest)
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		http.Error(w, `{"status":"ERROR","message":"title is required"}`, http.StatusBadRequest)
		return
	}
	dueAt, err := time.Parse(time.RFC3339, req.DueAt)
	if err != nil {
		http.Error(w, `{"status":"ERROR","message":"dueAt must be an RFC3339 timestamp (from toISOString())"}`, http.StatusBadRequest)
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
				Message: "Only instructors can create homeworks",
				Error:   "Unauthorized",
			})
			return
		}
		http.Error(w, `{"status":"ERROR","message":"database error while loading instructor"}`, http.StatusInternalServerError)
		return
	}

	uniqueSectionIDs := make([]uint, 0, len(req.SectionIDs))
	seen := make(map[uint]struct{}, len(req.SectionIDs))
	for _, id := range req.SectionIDs {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			uniqueSectionIDs = append(uniqueSectionIDs, id)
		}
	}

	var classrooms []database.Classroom
	if err := h.DB.Where("id IN ? AND instructor_id = ?", uniqueSectionIDs, inst.ID).
		Find(&classrooms).Error; err != nil {
		http.Error(w, `{"status":"ERROR","message":"database error while loading sections"}`, http.StatusInternalServerError)
		return
	}
	if len(classrooms) != len(uniqueSectionIDs) {
		okSet := make(map[uint]struct{}, len(classrooms))
		for _, c := range classrooms {
			okSet[c.ID] = struct{}{}
		}
		var missing []string
		for _, id := range uniqueSectionIDs {
			if _, ok := okSet[id]; !ok {
				missing = append(missing, strconv.FormatUint(uint64(id), 10))
			}
		}
		http.Error(w, `{"status":"ERROR","message":"you do not own all requested sections: `+strings.Join(missing, ",")+`"}`, http.StatusUnauthorized)
		return
	}

	var created []database.Assignment
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		for _, c := range classrooms {
			asg := database.Assignment{
				Title:       req.Title,
				DueDate:     dueAt,
				ClassroomID: c.ID,
			}
			if err := tx.Create(&asg).Error; err != nil {
				return err
			}
			created = append(created, asg)
			var studs []database.Student
			if err := tx.Where("classroom_id = ?", c.ID).Find(&studs).Error; err != nil {
				return err
			}
			if len(studs) > 0 {
				sas := make([]database.StudentAssignment, 0, len(studs))
				for _, s := range studs {
					sas = append(sas, database.StudentAssignment{
						StudentID:    s.ID,
						AssignmentID: asg.ID,
					})
				}
				if err := tx.Create(&sas).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		http.Error(w, `{"status":"ERROR","message":"failed to create homeworks"}`, http.StatusInternalServerError)
		return
	}

	// Build response DTOs
	out := make([]homeworkDTO, 0, len(created))
	for _, a := range created {
		out = append(out, homeworkDTO{
			ID:        a.ID,
			SectionID: a.ClassroomID,
			Title:     a.Title,
			DueAt:     a.DueDate,
		})
	}

	resp := models.Response[[]homeworkDTO]{
		Data:   out,
		Status: "OK",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

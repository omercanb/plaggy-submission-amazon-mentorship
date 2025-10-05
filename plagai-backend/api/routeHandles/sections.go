package routeHandles

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/plagai/plagai-backend/core"
	"github.com/plagai/plagai-backend/middleware"
	"github.com/plagai/plagai-backend/models"
	"github.com/plagai/plagai-backend/models/database"
)

func (h *Handler) SendSectionDetails(w http.ResponseWriter, r *http.Request) {
	sectionStr := r.URL.Query().Get("section")
	if sectionStr == "" {
		http.Error(w, `{"status":"ERROR","message":"missing 'section' query param"}`, http.StatusBadRequest)
		return
	}

	sectionNum, err := strconv.Atoi(sectionStr)
	if err != nil {
		http.Error(w, `{"status":"ERROR","message":"'section' conversion error"}`, http.StatusBadRequest)
		return
	}
	var classroomFromDB database.Classroom
	if dbErr := h.DB.Where("id = ?", (uint)(sectionNum)).First(&classroomFromDB).Error; dbErr != nil {
		log.Println("classroom query failed:", dbErr)
		http.Error(w, `{"status":"ERROR","message":"failed to load classrooms"}`, http.StatusInternalServerError)
		return
	}

	response := models.Response[models.Section]{
		Data:   models.Section{ID: (string)(classroomFromDB.ID), Name: classroomFromDB.Title},
		Status: "OK",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) SendSections(w http.ResponseWriter, r *http.Request) {
	// TODO this is just mock data, eventually connect to db. ~brtcrt
	// Also this is just to test if different users can have different sections
	// assigned to them. When we connect this to a db, probably just hold a list of
	// ids of the sections in the user and query again from the sections table.
	// Honestly though, I really don't give a shit how you structure the db.
	claims := middleware.Claims{}
	// No need to validate this, middleware should handle it.
	core.ConvertToken(r.Header.Get("Authorization"), &claims)

	var userFromDB database.Instructor
	if err := h.DB.
		Where("LOWER(email) = LOWER(?)", strings.TrimSpace(claims.Email)).
		First(&userFromDB).Error; err != nil {
		// Not found or other DB error
		log.Println("instructor lookup failed:", err)
		http.Error(w, `{"status":"ERROR","message":"instructor not found"}`, http.StatusNotFound)
		return
	}

	log.Println("Email from DB:", userFromDB.Email, "ID:", userFromDB.ID)

	var classroomsFromDB []database.Classroom
	if err := h.DB.
		Where("instructor_id = ?", userFromDB.ID).
		Find(&classroomsFromDB).Error; err != nil {
		log.Println("classroom query failed:", err)
		http.Error(w, `{"status":"ERROR","message":"failed to load classrooms"}`, http.StatusInternalServerError)
		return
	}
	log.Printf("Loaded %d classrooms for instructor %d\n", len(classroomsFromDB), userFromDB.ID)

	sendSections := make([]models.Section, len(classroomsFromDB))
	for i, c := range classroomsFromDB {
		sendSections[i].ID = strconv.FormatUint(uint64(c.ID), 10)
		sendSections[i].Name = c.Title
		log.Println("Section id:", c.ID, "title:", c.Title)
	}

	response := models.Response[[]models.Section]{
		Data:   sendSections,
		Status: "OK",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

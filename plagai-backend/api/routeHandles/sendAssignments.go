package routeHandles

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/plagai/plagai-backend/api"
	"github.com/plagai/plagai-backend/models"
	"github.com/plagai/plagai-backend/repository"
)

func (h *Handler) SendAssignments(w http.ResponseWriter, r *http.Request) {
	claims, err := api.GetClaimsFromAuthorization(r)
	if err != nil {
		switch err {
		case api.ErrMissingAuthHeader:
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		case api.ErrInvalidToken:
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		log.Println("auth error:", err)
		return
	}

	studentRepo := repository.NewStudentRepository(h.DB)

	student, err := studentRepo.FindByEmail(claims.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("No student found with email: %v", claims.Email), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	classroomRepo := repository.NewClassroomRepo(h.DB)

	classroom, err := classroomRepo.GetClassroomByID(student.ClassroomID)
	if err != nil {
		http.Error(w, fmt.Sprintf("No classroom found for student id: %v", student.ClassroomID), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	assignmentRepo := repository.NewAssignmentRepo(h.DB)

	assignments, err := assignmentRepo.GetAssignmentsForClassroomID(classroom.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Assignments couldn't be found for classroom with id: %v", classroom.ID), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	assignmentsDto := make([]models.Homework, len(assignments))
	for i, assignment := range assignments {
		assignmentsDto[i] = models.Homework{
			ID:         strconv.FormatUint(uint64(assignment.ID), 10),
			Title:      assignment.Title,
			AssignedAt: assignment.AssignedAt,
			DueDate:    assignment.DueDate,
		}
	}

	assignmentJson, err := json.Marshal(assignmentsDto)
	if err != nil {
		http.Error(w, "Failed to marshal events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(assignmentJson)
}

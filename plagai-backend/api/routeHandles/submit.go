package routeHandles

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/plagai/plagai-backend/api"
	"github.com/plagai/plagai-backend/flagging"
	"github.com/plagai/plagai-backend/models"
	"github.com/plagai/plagai-backend/models/database"
	"github.com/plagai/plagai-backend/models/domain"
	"github.com/plagai/plagai-backend/repository"
)

func (h *Handler) SubmitHandler(w http.ResponseWriter, r *http.Request) {
	// Parse JWT from Authorization header
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

	var submission models.Submission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	edits := submission.Edits

	email := claims.Email
	studentRepo := repository.NewStudentRepository(h.DB)
	student, err := studentRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrStudentNotFound) {
			log.Printf("no student found with email: %v for submission: %v", email, err)
		} else {
			log.Println(err)
		}
		http.Error(w, "No student found", http.StatusNotFound)
		return
	}

	studentAssignmentRepo := repository.NewStudentAssignmentRepo(h.DB)
	assignments := studentAssignmentRepo.GetStudentAssignments(student.ID)
	studentAssignmentToSubmitTo := domain.StudentAssignment{}
	assignmentFound := false
	for _, a := range assignments {
		if a.AssignmentID == submission.AssignmentId {
			assignmentFound = true
			studentAssignmentToSubmitTo = a
			break
		}
	}
	if assignmentFound == false {
		studentAssignmentToSubmitTo, err = studentAssignmentRepo.NewStudentAssignment(student.ID, submission.AssignmentId)
		if err != nil {
			log.Println(err)
			return
		}
	}

	var editEventsForDB []models.DBEditEvent
	for _, editDTO := range edits {
		editEventsForDB = append(editEventsForDB, models.DBEditEvent{
			PatchText: editDTO.Patch,
			Timestamp: editDTO.Timestamp.UnixMilli(),
			FilePath:  editDTO.FilePath,
		})
	}

	var diffsToCreate []database.Diff
	for _, event := range editEventsForDB {
		createdAt := time.UnixMilli(event.Timestamp)
		if event.Timestamp == 0 {
			createdAt = time.Now()
		}

		diffsToCreate = append(diffsToCreate, database.Diff{
			StudentAssignmentID: studentAssignmentToSubmitTo.ID,
			FilePath:            event.FilePath,
			DiffData:            event.PatchText,
			CreatedAt:           createdAt,
			UpdatedAt:           createdAt,
		})
	}

	if err := h.DB.CreateInBatches(&diffsToCreate, 200).Error; err != nil {
		log.Printf(`{"status":"ERROR","message":"failed to create diffs: %v"}`, err)
		return
	}

	log.Printf("Recieved %d edit events and added to db", len(diffsToCreate))

	ruleEngine := flagging.GetDefaultFlaggingEngine()
	flags := ruleEngine.FlagAssignment(edits)

	flagRepo := repository.NewFlagRepository(h.DB)

	for _, flag := range flags {
		err = flagRepo.AddFlag(&flag, studentAssignmentToSubmitTo.ID)
		if err != nil {
			log.Printf("failed to add flag for assignment %d, flag text: %q: %v",
				submission.AssignmentId, flag.FlagExplanation, err)
		}
	}
	fmt.Printf("Received %d edit events from %s:\n", len(edits), claims.Email)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "received"})
}

/*
func mapEventsToDiffs(events []models.EditEventDTO) []domain.Diff {
	diffs := make([]domain.Diff, 0)
	for _, event := range events {
		if event.EventType != models.APIEventModified {
			continue
		}
		diffs = append(diffs, domain.Diff{
			PatchText: event.Patch,
			Timestamp: event.Timestamp,
		})
	}
	return diffs
}
*/

package api

import (
	"aiplag-agent/cli/models"
	"aiplag-agent/common/api/dtomodels"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	//BackendBaseURL = "http://localhost:8080"
	BackendBaseURL     = "https://plaggy.xyz"
	SubmissionEndpoint = BackendBaseURL + "/api/v1/submit"
	AssignmentEndpoint = BackendBaseURL + "/api/v1/assignments"
)

func FetchAssignments(studentEmail string, token string) ([]models.Assignment, error) {
	req, err := http.NewRequest("GET", AssignmentEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, ServerError
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ServerError
	}

	currentAssignmentsDto := []dtomodels.Assignment{}
	if err := json.NewDecoder(resp.Body).Decode(&currentAssignmentsDto); err != nil {
		return nil, fmt.Errorf("failed to unmarshal returned assignments: %w", err)
	}

	currentAssignments := []models.Assignment{}
	for _, assignmentDto := range currentAssignmentsDto {
		assignmentId, err := strconv.ParseUint(assignmentDto.ID, 10, 64)
		if err != nil {
			log.Printf("Error when parsing assignment id: %s\n", assignmentDto.ID)
			continue
		}

		currentAssignments = append(currentAssignments, models.Assignment{
			ID:         uint(assignmentId),
			Title:      assignmentDto.Title,
			DueDate:    assignmentDto.DueDate,
			AssignedAt: assignmentDto.AssignedAt,
		})
	}

	return currentAssignments, nil
}

package api

import (
	"aiplag-agent/common/api/dtomodels"
	"aiplag-agent/common/db"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
)

var ServerError = errors.New("server is offline")

// SubmitEdits posts all edit events for a given assignment to a specified URL.
// The token currently includes the email encoded inside
func SubmitEdits(assignmentID uint, eh *db.EditHistoryStore, path string, token string) error {
	// Get events from the store
	internalAssignmentID, err := eh.GetAssignmentIDByFullPath(path)
	events, err := eh.GetEventsByAssignment(internalAssignmentID)
	if err != nil {
		return fmt.Errorf("failed to get events: %w", err)
	}

	// Convert to API models
	editEvents := dtomodels.ConvertEditEvents(events)
	submission := dtomodels.Submission{
		AssignmentId: assignmentID,
		Edits:        editEvents,
	}

	// Convert paths to relative filepaths to hide the students full path the the homework direcotry
	for i := range submission.Edits {
		relativeFilepath, err := filepath.Rel(path, submission.Edits[i].FilePath)
		if err == nil {
			submission.Edits[i].FilePath = relativeFilepath
		}
	}

	data, err := json.Marshal(submission)
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	req, err := http.NewRequest("POST", SubmissionEndpoint, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ServerError
	}
	return nil
}

package dtomodels

import (
	"aiplag-agent/daemon/models"
	"time"
)

// EditEventType represents event type in the API
type EditEventType string

const (
	APIEventAdded    EditEventType = "added"
	APIEventModified EditEventType = "modified"
	APIEventDeleted  EditEventType = "deleted"
	APIEventRenamed  EditEventType = "renamed"
)

// EditEvent is the JSON representation sent over HTTP
type EditEvent struct {
	ID           int           `json:"id"`
	AssignmentID int           `json:"assignment_id"`
	FilePath     string        `json:"file_path"`
	EventType    EditEventType `json:"event_type"`
	Patch        string        `json:"patch,omitempty"`
	Timestamp    time.Time     `json:"timestamp"`
}

// ConvertEditEvent maps internal EditEvent to APIEditEvent
func ConvertEditEvent(e models.EditEvent) EditEvent {
	var apiType EditEventType
	switch e.EventType {
	case models.EventAdded:
		apiType = APIEventAdded
	case models.EventModified:
		apiType = APIEventModified
	case models.EventDeleted:
		apiType = APIEventDeleted
	case models.EventRenamed:
		apiType = APIEventRenamed
	}

	return EditEvent{
		ID:           e.ID,
		AssignmentID: e.AssignmentID,
		FilePath:     e.FilePath,
		EventType:    apiType,
		Patch:        e.Patch,
		Timestamp:    e.Timestamp,
	}
}

// ConvertEditEvents converts a slice of internal EditEvents to APIEditEvents
func ConvertEditEvents(events []models.EditEvent) []EditEvent {
	apiEvents := make([]EditEvent, len(events))
	for i, e := range events {
		apiEvents[i] = ConvertEditEvent(e)
	}
	return apiEvents
}

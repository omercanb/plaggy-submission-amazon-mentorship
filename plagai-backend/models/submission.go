package models

import "time"

type PatchWithTimestamp struct {
	PatchText string `json:"patch_text"`
	Timestamp int64  `json:"timestamp"`
}

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

type Submission struct {
	AssignmentId uint        `json:"assignmentID"`
	Edits        []EditEvent `json:"edits"`
}

type DBEditEvent struct {
	PatchText string
	Timestamp int64
	FilePath  string
}

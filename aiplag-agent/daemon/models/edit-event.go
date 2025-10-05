package models

import (
	"fmt"
	"time"
)

// EditEvent stores information about a single edit event.
type EditEvent struct {
	ID           int
	AssignmentID int
	Patch        string
	FilePath     string
	EventType    EditEventType
	Timestamp    time.Time
}

// EditEventType represents the type of edit event applied to a file.
type EditEventType string

const (
	EventAdded    EditEventType = "added"
	EventModified EditEventType = "modified"
	EventDeleted  EditEventType = "deleted"
	EventRenamed  EditEventType = "renamed"
)

// StringToEditEventType converts a string to an EditEventType constant.
func StringToEditEventType(s string) (EditEventType, error) {
	switch s {
	case string(EventAdded):
		return EventAdded, nil
	case string(EventModified):
		return EventModified, nil
	case string(EventDeleted):
		return EventDeleted, nil
	case string(EventRenamed):
		return EventRenamed, nil
	default:
		return "", fmt.Errorf("invalid EditEventType: %q", s)
	}
}

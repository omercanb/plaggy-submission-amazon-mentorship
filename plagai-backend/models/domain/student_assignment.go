package domain

import (
	"time"
)

type StudentAssignment struct {
	ID             uint
	SubmissionTime time.Time
	StudentID      uint
	AssignmentID   uint
}
